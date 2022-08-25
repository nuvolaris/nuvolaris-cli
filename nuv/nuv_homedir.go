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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

// GetHomeDir detects the user's home directory in cross-compilation environments
var GetHomeDir = func() (string, error) {
	homedir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return homedir, nil
}

// GetOrCreateNuvolarisConfigDir creates .nuvolaris dir under user's homedir, if not already there
func GetOrCreateNuvolarisConfigDir() (string, error) {
	homedir, err := GetHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(homedir, ".nuvolaris")
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0777); err != nil {
			fmt.Println(err)
			return "", err
		}
		fmt.Println("nuvolaris config dir created")
	}
	return path, nil
}

// WriteFileToNuvolarisConfigDir writes file to .nuvolaris dir
func WriteFileToNuvolarisConfigDir(filename string, content []byte) (string, error) {
	nuvHomedir, err := GetOrCreateNuvolarisConfigDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(nuvHomedir, filename)
	if _, err := os.Stat(path); err == nil {
		os.Remove(path)
	}

	if err := os.WriteFile(path, content, 0600); err != nil {
		return "", err
	}
	return path, nil
}

// ReadFileFromNuvolarisConfigDir reads file from .nuvolaris dir
func ReadFileFromNuvolarisConfigDir(filename string) ([]byte, error) {
	nuvHomedir, err := GetOrCreateNuvolarisConfigDir()

	if err != nil {
		return nil, err
	}

	path := filepath.Join(nuvHomedir, filename)
	if _, err := os.Stat(path); err != nil {
		fmt.Println("File reading error", err)
		return nil, err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}
