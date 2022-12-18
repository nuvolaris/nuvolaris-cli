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
	"archive/zip"
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//go:embed embed-bundle/index.js
var idx []byte

//go:embed embed-bundle/package.json
var pkg []byte

type BundleCmd struct {
	Path   string `arg:"" help:"Path containing the web application bundle to assemble." type:"path"`
	Target string `arg:"" optional:"" help:"Name of of the output bundle" type:"path"`
}

func (s *BundleCmd) Run() error {
	err := validateBundleStructure(s.Path)
	if err != nil {
		return err
	}

	targetFile, err := getTargetOutput(filepath.Base(s.Path), s.Target)
	if err != nil {
		return err
	}

	fmt.Printf("Creatin zipfile %s scanning folder '%s'\n", targetFile, s.Path)
	err = ZipWriter(s.Path, targetFile)
	return err
}

func validateBundleStructure(basePath string) error {
	if !dirExists(basePath) {
		return fmt.Errorf("folder '%s' not found! Bundle requires an existing folder containing a valid web application source code", basePath)
	}

	fileToCheck := filepath.Join(basePath, "index.html")
	if !fileExists(fileToCheck) {
		return fmt.Errorf("folder '%s' does not contain an index.html file. Bundle structure not valid.", basePath)
	}

	fileToCheck = filepath.Join(basePath, "index.js")
	if fileExists(fileToCheck) {
		return fmt.Errorf("folder '%s' contains an index.js file. Bundle structure not valid.", basePath)
	}

	fileToCheck = filepath.Join(basePath, "package.json")
	if fileExists(fileToCheck) {
		return fmt.Errorf("folder '%s' contains a package.json file. Bundle structure not valid.", basePath)
	}

	return nil
}

// Calculates target zip filename
func getTargetOutput(basePath, target string) (string, error) {
	if target != "" && !strings.HasSuffix(target, ".zip") {
		return "", fmt.Errorf("target '%s' is not valid! Please use .zip extension.", target)
	}

	if target != "" && strings.HasSuffix(target, ".zip") {
		return target, nil
	}

	//We need to use the input folder name
	target = "./" + basePath + ".zip"
	return target, nil
}

// Zip the content cretaing the output file
func ZipWriter(baseFolder, outputFile string) error {

	// Get a Buffer to Write To
	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add folder file to the archive.
	addFiles(w, baseFolder, "")
	// Add files to expose the bundle as OpwnWhisk actions
	addContent(w, idx, "", "index.js")
	addContent(w, pkg, "", "package.json")

	if err != nil {
		fmt.Println(err)
		return err
	}

	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// Add Files recursively to the Output Folder
func addFiles(w *zip.Writer, basePath, baseInZip string) {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		filename := filepath.Join(basePath, file.Name())
		if !file.IsDir() {
			addFile(w, baseInZip, filename, file.Name())
		} else if file.IsDir() {
			fmt.Printf("Processing folder %s \n", filename)
			// Recurse
			newBase := filepath.Join(basePath, file.Name(), "/")
			addFiles(w, newBase, baseInZip+file.Name()+"/")
		}
	}
}

// Add a single file to the Output folder
func addFile(w *zip.Writer, baseInZip, filename, zipName string) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	// Add some files to the archive.
	f, err := w.Create(baseInZip + zipName)
	if err != nil {
		fmt.Println(err)
	}
	_, err = f.Write(dat)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Added file %s \n", filename)
}

// Add a single file to the Output folder
func addContent(w *zip.Writer, content []byte, baseInZip, zipName string) {
	// Add the content to the archive.
	f, err := w.Create(baseInZip + zipName)
	if err != nil {
		fmt.Println(err)
	}
	_, err = f.Write(content)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Added file %s \n", zipName)
}
