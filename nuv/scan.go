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
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type ScanCmd struct {
	// Path string `arg:"" optional:"" help:"Path to scan." type:"path"`
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
	root, err := processDir(aferoFs, path, "packages")
	if err != nil {
		return ProjectTree{}, err
	}
	return root, nil
}

func processDir(aferoFs afero.Fs, path string, dir string) (ProjectTree, error) {

	pt := ProjectTree{name: dir}
	var folders []*ProjectTree
	var files []*ProjectFile

	filename := filepath.Join(path, dir)
	children, err := afero.ReadDir(aferoFs, filename)

	if err != nil {
		return ProjectTree{}, err
	}

	for _, info := range children {
		if info.IsDir() {
			childPT, err := processDir(aferoFs, filename, info.Name())
			if err != nil {
				return pt, err
			}
			childPT.parent = &pt
			folders = append(folders, &childPT)
		} else {
			files = append(files, &ProjectFile{name: info.Name(), parent: &pt})
		}
	}

	pt.folders = folders
	pt.files = files
	return pt, nil
}

type ProjectTree struct {
	name    string
	parent  *ProjectTree
	folders []*ProjectTree
	files   []*ProjectFile
}

type ProjectFile struct {
	name   string
	parent *ProjectTree
}
