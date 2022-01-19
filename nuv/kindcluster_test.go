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
	"os"
	"testing"
)

func Test_manageKindCluster(t *testing.T) {

	realCreateCluster := createCluster
	realDestroyCluster := destroyCluster

	defer func() {
		createCluster = realCreateCluster
		destroyCluster = realDestroyCluster
	}()

	tests := []struct {
		name           string
		action         string
		createCluster  func() error
		destroyCluster func() error
		expectedError  error
		expectedOutput string
	}{
		{
			name:          "successfully creating cluster",
			action:        "create",
			createCluster: func() error { return nil },
			expectedError: nil,
		},
		{
			name:           "successfully destroying cluster",
			action:         "destroy",
			destroyCluster: func() error { return nil },
			expectedError:  nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			createCluster = test.createCluster
			destroyCluster = test.destroyCluster
			err := manageKindCluster(test.action)
			if err != test.expectedError {
				t.Errorf("Expected %v, got %v", test.expectedError, err)
			}
		})
	}

}

func Test_createNuvolarisHomeDir(t *testing.T) {

	realOsStat := osStat
	realOsMkdir := osMkdir
	realOsIsNotExist := osIsNotExist

	defer func() {
		osMkdir = realOsMkdir
		osIsNotExist = realOsIsNotExist
		osStat = realOsStat
	}()

	tests := []struct {
		name         string
		homedir      string
		expectedErr  error
		expectedDir  string
		osStat       func(string) (fs.FileInfo, error)
		osMkdir      func(string, fs.FileMode) error
		osIsNotExist func(error) bool
	}{
		{
			name:        ".nuvolaris dir does not exist yet",
			homedir:     "/home/userdir",
			expectedDir: "/home/userdir/.nuvolaris",
			expectedErr: nil,
			osStat: func(string) (fs.FileInfo, error) {
				return nil, os.ErrNotExist
			},
			osMkdir: func(string, fs.FileMode) error {
				return nil
			},
			osIsNotExist: func(error) bool {
				return true
			},
		},
		{
			name:        ".nuvolaris dir already exists",
			homedir:     "/home/userdir",
			expectedDir: "/home/userdir/.nuvolaris",
			expectedErr: nil,
			osStat: func(string) (fs.FileInfo, error) {
				return nil, nil
			},
			osMkdir: func(string, fs.FileMode) error {
				return nil
			},
			osIsNotExist: func(error) bool {
				return true
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			osStat = test.osStat
			osMkdir = test.osMkdir
			osIsNotExist = test.osIsNotExist
			result, err := createNuvolarisConfigDirIfNotExists(test.homedir)
			if result != test.expectedDir {
				t.Errorf("Expected %s, got %s", test.expectedDir, result)
			}
			if err != test.expectedErr {
				t.Errorf("Expected %v, got %v", test.expectedErr, err)
			}
		})
	}

}

func Test_rewriteKindConfigFile(t *testing.T) {
	realOsStatFunc := osStat
	realOsRemoveFunc := osRemove
	realOsWriteFileFunc := osWriteFile

	defer func() {
		osStat = realOsStatFunc
		osRemove = realOsRemoveFunc
		osWriteFile = realOsWriteFileFunc

	}()

	tests := []struct {
		name           string
		path           string
		expectedErr    error
		expectedResult string
		osStat         func(string) (fs.FileInfo, error)
		osWriteFile    func(string, []byte, fs.FileMode) error
		osRemove       func(string) error
	}{
		{
			name:           "config file does not exist yet",
			path:           "/home/userdir/.nuvolaris",
			expectedResult: "/home/userdir/.nuvolaris/kind.yaml",
			expectedErr:    nil,
			osStat: func(string) (fs.FileInfo, error) {
				return nil, os.ErrNotExist
			},
			osWriteFile: func(string, []byte, fs.FileMode) error {
				return nil
			},
			osRemove: func(string) error {
				return nil
			},
		},
		{
			name:           "config file already exists",
			path:           "/home/userdir/.nuvolaris",
			expectedResult: "/home/userdir/.nuvolaris/kind.yaml",
			expectedErr:    nil,
			osStat: func(string) (fs.FileInfo, error) {
				return nil, nil
			},
			osWriteFile: func(string, []byte, fs.FileMode) error {
				return nil
			},
			osRemove: func(string) error {
				return nil
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			osStat = test.osStat
			osWriteFile = test.osWriteFile
			osRemove = test.osRemove
			result, err := rewriteKindConfigFile(test.path)
			if result != test.expectedResult {
				t.Errorf("Expected %s, got %s", test.expectedResult, result)
			}
			if err != test.expectedErr {
				t.Errorf("Expected %v, got %v", test.expectedErr, err)
			}
		})
	}
}

func Test_clusterAlreadyRunning(t *testing.T) {
	realKind := kind

	defer func() {
		kind = realKind
	}()

	tests := []struct {
		name           string
		result         string
		expectedResult bool
		kind           func(...string) error
	}{
		{
			name:           "no running clusters",
			expectedResult: false,
			kind: func(...string) error {
				fmt.Println("")
				return nil
			},
		},
		{
			name:           "nuvolaris cluster running",
			expectedResult: true,
			kind: func(...string) error {
				fmt.Println("nuvolaris")
				return nil
			},
		},
		{
			name:           "nuvolaris and other clusters running",
			expectedResult: true,
			kind: func(...string) error {
				fmt.Println("kind nuvolaris keep adding cluster names")
				return nil
			},
		},
		{
			name:           "other cluster running but not nuvolaris",
			expectedResult: false,
			kind: func(...string) error {
				fmt.Println("kind")
				return nil
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			kind = test.kind
			result, _ := clusterAlreadyRunning()
			if result != test.expectedResult {
				t.Errorf("Expected %t, got %t", test.expectedResult, result)
			}
		})
	}
}
