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
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_createNuvolarisHomeDir(t *testing.T) {

	realOsStatFunc := osStatFunc
	realOsMkdirDirFunc := osMkdirFunc
	realOsIsNotExistFunc := osIsNotExistFunc

	defer func() {
		osMkdirFunc = realOsMkdirDirFunc
		osIsNotExistFunc = realOsIsNotExistFunc
		osStatFunc = realOsStatFunc
	}()

	//case: dir does not exist yet
	osStatFunc = func(path string) (fs.FileInfo, error) {
		return nil, os.ErrNotExist
	}

	osMkdirFunc = func(name string, perm fs.FileMode) error {
		return nil //dir created
	}

	osIsNotExistFunc = func(err error) bool {
		return true
	}

	out, err := createNuvolarisConfigDirIfNotExists("/home/userdir")
	assert.Equal(t, err, nil, "")
	assert.Equal(t, out, "/home/userdir/.nuvolaris", "")

	//case: dir already exists
	osStatFunc = func(path string) (fs.FileInfo, error) {
		return nil, nil
	}
	out, err = createNuvolarisConfigDirIfNotExists("/home/userdir")
	assert.Equal(t, err, nil, "")
	assert.Equal(t, out, "/home/userdir/.nuvolaris", "")

}

func Test_rewriteKindConfigFile(t *testing.T) {
	realOsStatFunc := osStatFunc
	realOsRemoveFunc := osRemoveFunc
	realOsWriteFileFunc := osWriteFileFunc

	defer func() {
		osStatFunc = realOsStatFunc
		osRemoveFunc = realOsRemoveFunc
		osWriteFileFunc = realOsWriteFileFunc

	}()

	//case: config file does not exist yet
	osStatFunc = func(path string) (fs.FileInfo, error) {
		return nil, os.ErrNotExist
	}

	osWriteFileFunc = func(name string, data []byte, perm fs.FileMode) error {
		return nil
	}

	out, err := rewriteKindConfigFile("/home/userdir/.nuvolaris")
	assert.Equal(t, err, nil, "")
	assert.Equal(t, out, "/home/userdir/.nuvolaris/kind.yaml", "")

	//case: config file already exists
	osRemoveFunc = func(name string) error {
		return nil
	}
	out, err = rewriteKindConfigFile("/home/userdir/.nuvolaris")
	assert.Equal(t, err, nil, "")
	assert.Equal(t, out, "/home/userdir/.nuvolaris/kind.yaml", "")

}

func Test_clusterAlreadyRunning(t *testing.T) {
	DryRunPush("")
	out, err := clusterAlreadyRunning(true)
	assert.Equal(t, out, false, "")
	assert.Equal(t, err, nil, "")

	DryRunPush("nuvolaris")
	out, err = clusterAlreadyRunning(true)
	assert.Equal(t, out, true, "")
	assert.Equal(t, err, nil, "")

	DryRunPush("kind nuvolaris keep adding cluster names")
	out, err = clusterAlreadyRunning(true)
	assert.Equal(t, out, true, "")
	assert.Equal(t, err, nil, "")

	DryRunPush("kind")
	out, err = clusterAlreadyRunning(true)
	assert.Equal(t, out, false, "")
	assert.Equal(t, err, nil, "")
	//output
	//kind get clusters
}

func Test_startCluster(t *testing.T) {
	err := startCluster(true, "./embed/kind.yaml")
	assert.Equal(t, err, nil, "")
	//output
	//kind create cluster --wait=1m --config=./embed/kind.yaml
}

func Test_createCluster(t *testing.T) {
	err := createCluster(true)
	assert.Equal(t, err, nil, "")
	//output
	//kind get clusters
	//starting nuvolaris kind cluster...hang tight
	//kind create cluster --wait=1m --config=/home/nuvolaris/.nuvolaris/kind.yaml
	//nuvolaris kind cluster started. Have a nice day! 👋
}

func Test_destroyCluster(t *testing.T) {
	err := destroyCluster(true)
	assert.Equal(t, err, nil, "")
	//kind get clusters
	//kind cluster nuvolaris not found. Skipping...
}
