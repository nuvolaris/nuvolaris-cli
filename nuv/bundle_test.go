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
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func Test_bundleRun(t *testing.T) {
	hd, _ := GetHomeDir()
	targetFile := filepath.Join(hd, "test-bundle-output.zip")
	bundleCmd := BundleCmd{Path: "./test-bundle/test0", Target: targetFile}
	bundleCmd.Run()

	exists := fileExists(targetFile)
	assert.True(t, exists)
	err := os.Remove(targetFile)
	assert.NoError(t, err)
}

func Test_bundleFolderDoesNotExist(t *testing.T) {
	hd, _ := GetHomeDir()
	targetFile := filepath.Join(hd, "test-bundle-output.zip")
	bundleCmd := BundleCmd{Path: "./test-bundle/fake", Target: targetFile}
	err := bundleCmd.Run()
	assert.ErrorContains(t, err, "folder './test-bundle/fake' not found!")
}

func Test_indexFileDoesNotExist(t *testing.T) {
	hd, _ := GetHomeDir()
	targetFile := filepath.Join(hd, "test-bundle-output.zip")
	bundleCmd := BundleCmd{Path: "./test-bundle/test1", Target: targetFile}
	err := bundleCmd.Run()
	assert.ErrorContains(t, err, "folder './test-bundle/test1' does not contain an index.html file")
}

func Test_indexJsFileExist(t *testing.T) {
	hd, _ := GetHomeDir()
	targetFile := filepath.Join(hd, "test-bundle-output.zip")
	bundleCmd := BundleCmd{Path: "./test-bundle/test2", Target: targetFile}
	err := bundleCmd.Run()
	assert.ErrorContains(t, err, "folder './test-bundle/test2' contains an index.js file. Bundle structure not valid.")
}

func Test_PackageJsonFileExist(t *testing.T) {
	hd, _ := GetHomeDir()
	targetFile := filepath.Join(hd, "test-bundle-output.zip")
	bundleCmd := BundleCmd{Path: "./test-bundle/test3", Target: targetFile}
	err := bundleCmd.Run()
	assert.ErrorContains(t, err, "folder './test-bundle/test3' contains a package.json file. Bundle structure not valid.")
}

func Test_NonValidBundleTarget(t *testing.T) {
	hd, _ := GetHomeDir()
	targetFile := filepath.Join(hd, "test-bundle-output.tar")
	bundleCmd := BundleCmd{Path: "./test-bundle/test0", Target: targetFile}
	err := bundleCmd.Run()
	assert.ErrorContains(t, err, "target '"+targetFile+"' is not valid! Please use .zip extension.")
}
