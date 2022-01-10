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
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type ScanCmd struct {
	// Path string `arg:"" optional:"" help:"Path to scan." type:"path"`
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
	name   string
	path   string
	parent *ProjectTree
}

func (s *ScanCmd) Run() error {
	fs := afero.NewOsFs()
	b, err := checkPackagesFolder(fs, "./") // TODO path
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
	return nil
}

func checkPackagesFolder(fs afero.Fs, path string) (bool, error) {
	dir := "packages"
	filename := filepath.Join(path, dir)
	b, err := afero.DirExists(fs, filename)
	return b, err
}

func scanPackagesFolder(aferoFs afero.Fs, path string) (ProjectTree, error) {
	root, err := processDir(aferoFs, path, "packages", true)
	if err != nil {
		return ProjectTree{}, err
	}
	return root, nil
}

func processDir(aferoFs afero.Fs, parentPath string, dir string, rootLevel bool) (ProjectTree, error) {
	pt := ProjectTree{name: dir}
	var folders []*ProjectTree
	var mfActions []*ProjectTreeAction
	var sfActions []*ProjectTreeAction

	dirPath := filepath.Join(parentPath, dir)
	children, err := afero.ReadDir(aferoFs, dirPath)

	if err != nil {
		return ProjectTree{}, err
	}

	for _, info := range children {
		if info.IsDir() {
			if rootLevel {
				// root level: folders = packages and continue walk
				childPT, err := processDir(aferoFs, dirPath, info.Name(), false)
				if err != nil {
					return pt, err
				}
				childPT.parent = &pt
				folders = append(folders, &childPT)
			} else {
				// inner level: folders = multi file actions and stop
				mfActions = appendNewAction(mfActions, info, dirPath, &pt)
			}
		} else {
			sfActions = appendNewAction(sfActions, info, dirPath, &pt)
		}
	}

	pt.path = dirPath
	pt.folders = folders
	pt.mfActions = mfActions
	pt.sfActions = sfActions
	return pt, nil
}

func appendNewAction(actions []*ProjectTreeAction, info fs.FileInfo, dirPath string, parent *ProjectTree) []*ProjectTreeAction {
	ext := path.Ext(info.Name())
	actionName := strings.TrimSuffix(info.Name(), ext) // remove extension from filename
	return append(actions, &ProjectTreeAction{name: actionName, path: filepath.Join(dirPath, info.Name()), parent: parent})
}

var wg sync.WaitGroup

func parseProjectTree(projectRoot *ProjectTree) TaskTree {
	taskRoot := TaskTree{}

	// First level commands: actions from files
	wg.Add(1)
	parseSingleFileActions(projectRoot, &taskRoot)

	// First level commands: packages from folders
	parseSubFolders(projectRoot, &taskRoot)

	wg.Wait()
	return taskRoot
}

func parseSubFolders(projectRoot *ProjectTree, taskRoot *TaskTree) {
	var tasks []*TaskTree
	for _, subf := range projectRoot.folders {
		t := TaskTree{parent: taskRoot, command: packageUpdate(subf.name)}

		// Second level commands: single file actions from subfolders
		wg.Add(1)
		go parseSingleFileActions(subf, &t)

		// Second level commands: multi file actions from subfolders
		// TODO: Multi file actions heuristics

		tasks = append(tasks, &t)
	}
	taskRoot.tasks = append(taskRoot.tasks, tasks...)
}

var extRuntimes = map[string]string{
	".go":   "--kind go:default",
	".java": "--kind java:default",
	".js":   "--kind nodejs:default",
	".py":   "--kind python:default",
}

func parseSingleFileActions(parent *ProjectTree, taskNode *TaskTree) {
	defer wg.Done()

	var tasks []*TaskTree

	wskPkg := ""
	if parent.parent != nil {
		wskPkg = parent.name + "/"
	}

	for _, file := range parent.sfActions {
		runtime := extRuntimes[filepath.Ext(file.path)]
		cmd := actionUpdate(wskPkg, file.name, file.path, runtime)
		t := TaskTree{parent: taskNode, command: cmd}
		tasks = append(tasks, &t)
	}

	taskNode.tasks = append(taskNode.tasks, tasks...)
}

func actionUpdate(pkg, actionName, filepath, runtime string) string {
	return fmt.Sprintf("wsk action update %s%s %s %s", pkg, actionName, filepath, runtime)
}
func packageUpdate(pkgName string) string {
	return fmt.Sprintf("wsk package update %s", pkgName)
}

type TaskTree struct {
	parent  *TaskTree
	tasks   []*TaskTree
	command string
}
