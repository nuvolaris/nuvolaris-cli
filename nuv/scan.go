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
	// Path string `arg:"" optional:"" help:"Path to scan." type:"path"`
}

func (s *ScanCmd) Run() error {
	b, err := packagesFolderExists("./") // TODO path
	if !b {
		// packages folder not found, stop here
		log.Error("Folder 'packages' not found! Cannot scan project :(")
		return nil
	}
	if err != nil {
		log.Error("Error reading packages folder!") // TODO: improve feedback to user...
		log.Debug(err)
		return err
	}

	// TODO: ignore non code files (files without supported runtimes)
	// projectTree, err := scanPackagesFolder(fs, "./")
	// if err != nil {
	// 	return err
	// }

	// taskTree := parseProjectTree(&projectTree)

	// for _, st := range taskTree.tasks {
	// 	log.Info(st.command)
	// 	for _, stt := range st.tasks {
	// 		log.Info(stt.command)
	// 	}
	// }
	return nil
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
	parent  *ProjectTree
	folders []*ProjectTree

	mfActions []*ProjectTreeAction
	sfActions []*ProjectTreeAction
}

type ProjectTreeAction struct {
	name    string
	runtime string
	path    string
	parent  *ProjectTree
}

type TaskTree struct {
	parent  *TaskTree
	tasks   []*TaskTree
	command string
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
				childPT.parent = &pt
				folders = append(folders, &childPT)
			} else {
				// inner level: folders = multi file actions and stop
				mfPath := filepath.Join(dirPath, info.Name())
				runtime, err := findMfaRuntime(mfPath)
				if err != nil {
					return pt, err
				}
				mfActions = append(mfActions, &ProjectTreeAction{name: info.Name(), runtime: runtime, path: mfPath, parent: &pt})

			}
		} else {
			ext := path.Ext(info.Name())
			actionName := strings.TrimSuffix(info.Name(), ext) // remove extension from filename
			sfActions = append(sfActions, &ProjectTreeAction{name: actionName, runtime: ext, path: filepath.Join(dirPath, info.Name()), parent: &pt})
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

func parseProjectTree(projectRoot *ProjectTree) TaskTree {
	taskRoot := TaskTree{}

	// First level commands: actions from files
	parseRootSingleFileActions(projectRoot, &taskRoot)

	parseSubFolders(projectRoot, &taskRoot)

	return taskRoot
}

func parseRootSingleFileActions(projectRoot *ProjectTree, taskRoot *TaskTree) {
	childTasks := make([]*TaskTree, len(projectRoot.sfActions))
	for i, sfAction := range projectRoot.sfActions {
		cmd := actionUpdate("", sfAction.name, sfAction.path, extRuntimes[sfAction.runtime])
		childTasks[i] = &TaskTree{command: cmd}
	}
	taskRoot.tasks = childTasks
}

var wg sync.WaitGroup

func parseSubFolders(projectRoot *ProjectTree, taskRoot *TaskTree) {
	tasks := make([]*TaskTree, len(projectRoot.folders))

	for i, subf := range projectRoot.folders {
		taskQueue := make(chan *TaskTree, len(subf.sfActions)+len(subf.mfActions))

		// First level commands: packages from folders
		t := TaskTree{parent: taskRoot, command: packageUpdate(subf.name)}

		// Second level commands: single file actions from subfolders
		wg.Add(1)
		go parseSingleFileActions(taskQueue, subf, &t)

		// Second level commands: multi file actions from subfolders
		wg.Add(1)
		go parseMultiFileActions(taskQueue, subf, &t)

		wg.Wait()
		close(taskQueue)
		appendTasks(taskQueue, &t)

		tasks[i] = &t
	}
	taskRoot.tasks = append(taskRoot.tasks, tasks...)
}

func appendTasks(taskQueue chan *TaskTree, taskNode *TaskTree) {
	for subTask := range taskQueue {
		taskNode.tasks = append(taskNode.tasks, subTask)
	}

}

func parseSingleFileActions(taskQueue chan *TaskTree, parent *ProjectTree, taskNode *TaskTree) {
	defer wg.Done()

	wskPkg := parent.name + "/"

	for _, sfAction := range parent.sfActions {
		cmd := actionUpdate(wskPkg, sfAction.name, sfAction.path, extRuntimes[sfAction.runtime])
		t := TaskTree{parent: taskNode, command: cmd}
		taskQueue <- &t
	}

}

func parseMultiFileActions(taskQueue chan *TaskTree, parent *ProjectTree, taskNode *TaskTree) {
	defer wg.Done()

	wskPkg := parent.name + "/"

	for _, mfAction := range parent.mfActions {
		zipCmd := fmt.Sprintf("zip -r %s/%s.zip %s/*", mfAction.path, mfAction.name, mfAction.path)
		zipPath := fmt.Sprintf("%s/%s.zip", mfAction.path, mfAction.name)
		cmd := actionUpdate(wskPkg, mfAction.name, zipPath, extRuntimes[mfAction.runtime])
		t := TaskTree{parent: taskNode, command: fmt.Sprintf("%s\n%s", zipCmd, cmd)}
		taskQueue <- &t
	}
}

func actionUpdate(pkg, actionName, filepath, runtime string) string {
	return fmt.Sprintf("wsk action update %s%s %s %s", pkg, actionName, filepath, runtime)
}
func packageUpdate(pkgName string) string {
	return fmt.Sprintf("wsk package update %s", pkgName)
}
