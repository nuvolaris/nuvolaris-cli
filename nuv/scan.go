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
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

type ScanCmd struct {
	Path string `arg:"" optional:"" default:"./" help:"Path to scan." type:"path"`
}

func (s *ScanCmd) Run() error {
	fsys := os.DirFS(s.Path)

	taskfile, err := generateTaskfile(fsys)
	if err != nil {
		return err
	}

	// Save to ~/.nuvolaris/nuvolaris.yml
	_, err = WriteFileToNuvolarisConfigDir("nuvolaris.yml", []byte(taskfile))
	if err != nil {
		return err
	}
	return nil
}

func generateTaskfile(fsys fs.FS) (string, error) {

	// 1. Check that ScanFolder is present and accessible
	b, err := packagesFolderExists(fsys)
	if !b {
		// packages folder not found, stop here
		return "", fmt.Errorf("folder '%s' not found! Cannot scan project :(", ScanFolder)
	}
	if err != nil {
		log.Error("Error reading packages folder!") // TODO: improve feedback to user...
		log.Debug(err)
		return "", err
	}

	// 2. Visit the ScanFolder and parse the contents into a tree object
	projectTree, err := visitScanFolder(fsys)
	if err != nil {
		return "", err
	}

	// 3. Turn the tree into a list of tasks for the Taskfile
	tasks := parseProjectTree(&projectTree)

	// 4. Merge the tasks into yaml format for a Taskfile
	taskfile := mergeIntoYaml(tasks)

	return taskfile, nil
}

// 1.
func packagesFolderExists(fsys fs.FS) (bool, error) {
	_, err := fs.Stat(fsys, ScanFolder)
	if os.IsNotExist(err) {
		return false, err
	}
	return true, err
}

// 2.
func visitScanFolder(fsys fs.FS) (ScanTree, error) {
	root, err := processDir(fsys, "", ScanFolder, true)
	if err != nil {
		return ScanTree{}, err
	}
	return root, nil
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

type ScanTree struct {
	name     string
	path     string
	packages []*ScanTree

	mfActions []*Action
	sfActions []*Action
}

type Action struct {
	name    string
	runtime string
	path    string
}

func processDir(fsys fs.FS, parentPath string, dir string, rootLevel bool) (ScanTree, error) {
	pt := ScanTree{name: dir}
	var folders []*ScanTree
	var mfActions []*Action
	var sfActions []*Action

	dirPath := filepath.Join(parentPath, dir)
	children, err := fs.ReadDir(fsys, dirPath)

	if err != nil { // TODO: ReadDir returns the entries it was able to read before the error. Parse what was read anyway?
		return ScanTree{}, err
	}

	for _, info := range children {
		if info.IsDir() {
			if rootLevel {
				// root level: folders == packages and continue walk
				childPT, err := processDir(fsys, dirPath, info.Name(), false)
				if err != nil {
					return pt, err
				}
				folders = append(folders, &childPT)
			} else {
				// inner level: folders = multi file actions and stop
				mfPath := filepath.Join(dirPath, info.Name())
				runtime, err := findMfaRuntime(fsys, mfPath)
				if err != nil {
					return pt, err
				}
				mfActions = append(mfActions, &Action{name: info.Name(), runtime: runtime, path: mfPath})

			}
		} else {
			ext := filepath.Ext(info.Name())
			if extRuntimes[ext] == "" {
				return pt, fmt.Errorf("no supported runtime found for file %s", info.Name())
			}
			actionName := strings.TrimSuffix(info.Name(), ext) // remove extension from filename
			sfActions = append(sfActions, &Action{name: actionName, runtime: ext, path: filepath.Join(dirPath, info.Name())})
		}
	}

	pt.path = dirPath
	pt.packages = folders
	pt.mfActions = mfActions
	pt.sfActions = sfActions
	return pt, nil
}

func findMfaRuntime(fsys fs.FS, mfPath string) (string, error) {
	found, err := searchRuntime(fsys, mfPath, jsRuntime, "package.json")
	if found {
		return jsRuntime, err
	}
	found, err = searchRuntime(fsys, mfPath, pyRuntime, "requirements.txt")
	if found {
		return pyRuntime, err
	}
	found, err = searchRuntime(fsys, mfPath, javaRuntime, "pom.xml")
	if found {
		return javaRuntime, err
	}
	found, err = searchRuntime(fsys, mfPath, goRuntime, "go.mod")
	if found {
		return goRuntime, err
	}

	if err != nil {
		return "", err
	}
	return "", fmt.Errorf("no supported runtime found")
}

func searchRuntime(fsys fs.FS, mfPath, ext, file string) (bool, error) {
	if _, err := fs.Stat(fsys, filepath.Join(mfPath, file)); err == nil {
		return true, nil

	} else if errors.Is(err, os.ErrNotExist) {

		pattern := filepath.Join(mfPath, "*"+ext)
		matches, err := fs.Glob(fsys, pattern)
		if err != nil {
			return false, err
		}
		return len(matches) > 0, nil

	} else {
		return false, err
	}
}

// 3.
func parseProjectTree(projectRoot *ScanTree) []string {

	// First level commands: actions from files
	rootTasks := parseRootSingleFileActions(projectRoot)

	subTasks := parseSubFolders(projectRoot)

	return append(rootTasks, subTasks...)
}

func parseRootSingleFileActions(projectRoot *ScanTree) []string {
	childTasks := make([]string, len(projectRoot.sfActions))
	for i, sfAction := range projectRoot.sfActions {
		cmd := actionUpdate("", sfAction.name, sfAction.path, extRuntimes[sfAction.runtime])
		childTasks[i] = cmd
	}
	return childTasks
}

var wg sync.WaitGroup

func parseSubFolders(projectRoot *ScanTree) []string {
	tasks := make([]string, 0)
	var subTasks []string

	for _, subf := range projectRoot.packages {
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

func parseSingleFileActions(taskQueue chan string, parent *ScanTree) {
	defer wg.Done()

	wskPkg := parent.name + "/"
	for _, sfAction := range parent.sfActions {
		cmd := actionUpdate(wskPkg, sfAction.name, sfAction.path, extRuntimes[sfAction.runtime])
		taskQueue <- cmd
	}

}

func parseMultiFileActions(taskQueue chan string, parent *ScanTree) {
	defer wg.Done()

	wskPkg := parent.name + "/"
	for _, mfAction := range parent.mfActions {
		packCmd := fmt.Sprintf("nuv pack -r %s/%s.zip %s/*", mfAction.path, mfAction.name, mfAction.path)
		packPath := fmt.Sprintf("%s/%s.zip", mfAction.path, mfAction.name)
		cmd := actionUpdate(wskPkg, mfAction.name, packPath, extRuntimes[mfAction.runtime])
		taskQueue <- packCmd
		taskQueue <- cmd
	}
}

func actionUpdate(pkg, actionName, filepath, runtime string) string {
	return fmt.Sprintf("nuv wsk action update %s%s %s %s", pkg, actionName, filepath, runtime)
}
func packageUpdate(pkgName string) string {
	return fmt.Sprintf("nuv wsk package update %s", pkgName)
}

// 4.
func mergeIntoYaml(tasks []string) string {
	taskfile := "version: 3\n\ntasks:\n  default:\n    cmds:"

	for _, t := range tasks {
		taskfile = fmt.Sprintf("%s\n      - %s", taskfile, t)
	}

	taskfile = fmt.Sprintf("%s\n", taskfile)
	return taskfile
}
