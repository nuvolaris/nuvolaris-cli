// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

type ScanCmd struct {
	Path string `arg:"" optional:"" help:"Path to scan." type:"path"`
	Save string `optional:"" short:"s" help:"Path to save output nuvolaris.yml." type:"path"`
}

func (s *ScanCmd) Run() error {
	b, err := packagesFolderExists(s.Path)
	if !b {
		// packages folder not found, stop here
		var errMsg string
		if s.Path == "" {
			errMsg = "folder 'packages' not found! Cannot scan project :("
		} else {
			errMsg = fmt.Sprintf("folder 'packages' in %s not found! Cannot scan project :(", s.Path)
		}
		return fmt.Errorf("%s", errMsg)
	}
	if err != nil {
		log.Error("Error reading packages folder!") // TODO: improve feedback to user...
		log.Debug(err)
		return err
	}

	projectTree, err := scanPackagesFolder(s.Path)
	if err != nil {
		return err
	}

	tasks := parseProjectTree(&projectTree)

	mergeIntoYaml(tasks, s.Save)

	return nil
}

func mergeIntoYaml(tasks []string, savePath string) {
	taskfile := "version: 3\n\ntasks:\n  default:\n    cmds:"

	for _, t := range tasks {
		taskfile = fmt.Sprintf("%s\n      - %s", taskfile, t)
	}

	taskfile = fmt.Sprintf("%s\n", taskfile)
	err := os.WriteFile(filepath.Join(savePath, "nuvolaris.yml"), []byte(taskfile), 0700)
	if err != nil {
		log.Fatal(err)
	}
}

const goRuntime = ".go"
const javaRuntime = ".java"
const jsRuntime = ".js"
const pyRuntime = ".py"

var extRuntimes = map[string]string{
	goRuntime:   "--kind go:default",
	javaRuntime: "--kind java:default",
	jsRuntime:   "--kind nodejs:default",
	pyRuntime:   "--kind python:default",
}

type ProjectTree struct {
	name    string
	path    string
	folders []*ProjectTree

	mfActions []*ProjectTreeAction
	sfActions []*ProjectTreeAction
}

type ProjectTreeAction struct {
	name    string
	runtime string
	path    string
}

func packagesFolderExists(path string) (bool, error) {
	dir := "packages"
	filename := filepath.Join(path, dir)
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false, err
	}
	return true, err
}

func scanPackagesFolder(path string) (ProjectTree, error) {
	root, err := processDir(path, "packages", true)
	if err != nil {
		return ProjectTree{}, err
	}
	return root, nil
}

func processDir(parentPath string, dir string, rootLevel bool) (ProjectTree, error) {
	pt := ProjectTree{name: dir}
	var folders []*ProjectTree
	var mfActions []*ProjectTreeAction
	var sfActions []*ProjectTreeAction

	dirPath := filepath.Join(parentPath, dir)
	children, err := os.ReadDir(dirPath)

	if err != nil { // TODO ReadDir returns the entries it was able to read before the error. Parse what was read anyway?
		return ProjectTree{}, err
	}

	for _, info := range children {
		if info.IsDir() {
			if rootLevel {
				// root level: folders = packages and continue walk
				childPT, err := processDir(dirPath, info.Name(), false)
				if err != nil {
					return pt, err
				}
				folders = append(folders, &childPT)
			} else {
				// inner level: folders = multi file actions and stop
				mfPath := filepath.Join(dirPath, info.Name())
				runtime, err := findMfaRuntime(mfPath)
				if err != nil {
					return pt, err
				}
				mfActions = append(mfActions, &ProjectTreeAction{name: info.Name(), runtime: runtime, path: mfPath})

			}
		} else {
			ext := path.Ext(info.Name())
			if extRuntimes[ext] == "" {
				return pt, fmt.Errorf("no supported runtime found for file %s", info.Name())
			}
			actionName := strings.TrimSuffix(info.Name(), ext) // remove extension from filename
			sfActions = append(sfActions, &ProjectTreeAction{name: actionName, runtime: ext, path: filepath.Join(dirPath, info.Name())})
		}
	}

	pt.path = dirPath
	pt.folders = folders
	pt.mfActions = mfActions
	pt.sfActions = sfActions
	return pt, nil
}

func findMfaRuntime(mfPath string) (string, error) {
	found, err := searchRuntime(mfPath, jsRuntime, "package.json")
	if found {
		return jsRuntime, err
	}
	found, err = searchRuntime(mfPath, pyRuntime, "requirements.txt")
	if found {
		return pyRuntime, err
	}
	found, err = searchRuntime(mfPath, javaRuntime, "pom.xml")
	if found {
		return javaRuntime, err
	}
	found, err = searchRuntime(mfPath, goRuntime, "go.mod")
	if found {
		return goRuntime, err
	}

	if err != nil {
		return "", err
	}
	return "", fmt.Errorf("no supported runtime found")
}

func searchRuntime(mfPath, ext, file string) (bool, error) {
	if _, err := os.Stat(path.Join(mfPath, file)); err == nil {
		return true, nil

	} else if errors.Is(err, os.ErrNotExist) {

		pattern := fmt.Sprintf("%s/*%s", mfPath, ext)
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return false, err
		}
		return len(matches) > 0, nil

	} else {
		return false, err
	}
}

func parseProjectTree(projectRoot *ProjectTree) []string {

	// First level commands: actions from files
	rootTasks := parseRootSingleFileActions(projectRoot)

	subTasks := parseSubFolders(projectRoot)

	return append(rootTasks, subTasks...)
}

func parseRootSingleFileActions(projectRoot *ProjectTree) []string {
	childTasks := make([]string, len(projectRoot.sfActions))
	for i, sfAction := range projectRoot.sfActions {
		cmd := actionUpdate("", sfAction.name, sfAction.path, extRuntimes[sfAction.runtime])
		childTasks[i] = cmd
	}
	return childTasks
}

var wg sync.WaitGroup

func parseSubFolders(projectRoot *ProjectTree) []string {
	tasks := make([]string, 0)
	var subTasks []string

	for _, subf := range projectRoot.folders {
		taskQueue := make(chan string, len(subf.sfActions)+(len(subf.mfActions)*2))

		// First level commands: packages from folders
		t := packageUpdate(subf.name)

		// Second level commands: single file actions from subfolders
		wg.Add(1)
		go parseSingleFileActions(taskQueue, subf)

		// Second level commands: multi file actions from subfolders
		wg.Add(1)
		go parseMultiFileActions(taskQueue, subf)

		wg.Wait()
		close(taskQueue)
		subTasks = appendTasks(taskQueue)

		tasks = append(tasks, t)
		tasks = append(tasks, subTasks...)
	}
	return tasks
}

func appendTasks(taskQueue chan string) []string {
	tasks := make([]string, 0)
	for subTask := range taskQueue {
		tasks = append(tasks, subTask)
	}
	return tasks
}

func parseSingleFileActions(taskQueue chan string, parent *ProjectTree) {
	defer wg.Done()

	wskPkg := parent.name + "/"
	for _, sfAction := range parent.sfActions {
		cmd := actionUpdate(wskPkg, sfAction.name, sfAction.path, extRuntimes[sfAction.runtime])
		taskQueue <- cmd
	}

}

func parseMultiFileActions(taskQueue chan string, parent *ProjectTree) {
	defer wg.Done()

	wskPkg := parent.name + "/"

	for _, mfAction := range parent.mfActions {
		zipCmd := fmt.Sprintf("zip -r %s/%s.zip %s/*", mfAction.path, mfAction.name, mfAction.path)
		zipPath := fmt.Sprintf("%s/%s.zip", mfAction.path, mfAction.name)
		cmd := actionUpdate(wskPkg, mfAction.name, zipPath, extRuntimes[mfAction.runtime])
		taskQueue <- zipCmd
		taskQueue <- cmd
	}
}

func actionUpdate(pkg, actionName, filepath, runtime string) string {
	return fmt.Sprintf("wsk action update %s%s %s %s", pkg, actionName, filepath, runtime)
}
func packageUpdate(pkgName string) string {
	return fmt.Sprintf("wsk package update %s", pkgName)
}
