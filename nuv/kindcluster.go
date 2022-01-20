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
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	nuvolarisConfigDir   = ".nuvolaris"
	nuvolarisClusterName = "nuvolaris"
	kindConfigFile       = "kind.yaml"
)

//go:embed embed/kind.yaml
var KindYaml []byte

//monkey patching functions for unit tests
var osMkdir = os.Mkdir
var osStat = os.Stat
var osIsNotExist = os.IsNotExist
var osWriteFile = os.WriteFile
var osRemove = os.Remove
var kind = Kind

var manageKindCluster = func(action string) error {

	switch action {
	case "create":
		if err := createCluster(); err != nil {
			return err
		}
	case "destroy":
		if err := destroyCluster(); err != nil {
			return err
		}
	default:
		fmt.Println("did you mean nuv devcluster create/destroy?")
	}
	return nil
}

var createCluster = func() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("error in create cluster: %w", err)
		}
	}()

	clusterIsRunning, err := clusterAlreadyRunning()
	if err != nil {
		return err
	}
	if clusterIsRunning {
		fmt.Println("nuvolaris kind cluster is already running. Skipping...")
		return nil
	}
	homedir, err := GetHomeDir()
	if err != nil {
		return err
	}

	fmt.Println("running preflight checks")
	if err = RunPreflightChecks(homedir); err != nil {
		return err
	}
	fmt.Println("preflight checks ok")

	configDir, err := createNuvolarisConfigDirIfNotExists(homedir)
	if err != nil {
		return err
	}

	filePath, err := rewriteKindConfigFile(configDir)
	if err != nil {
		return err
	}

	fmt.Println("starting nuvolaris kind cluster...hang tight")
	if err = startCluster(filePath); err != nil {
		return err
	}

	fmt.Println("nuvolaris kind cluster started. Have a nice day! 👋")
	return nil
}

var destroyCluster = func() error {
	clusterIsRunning, err := clusterAlreadyRunning()
	if err != nil {
		return err
	}
	if clusterIsRunning {
		if err := stopCluster(); err != nil {
			return err
		}
		fmt.Println("kind cluster nuvolaris destroyed")
	} else {
		fmt.Println("kind cluster nuvolaris not found. Skipping...")
	}
	return nil
}

var clusterAlreadyRunning = func() (bool, error) {
	//capture cmd output
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := kind("get", "clusters")

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	if err != nil {
		return false, err
	}
	if strings.Contains(string(out), nuvolarisClusterName) {
		return true, nil
	} else {
		return false, nil
	}
}

var createNuvolarisConfigDirIfNotExists = func(homedir string) (string, error) {
	fullPath := filepath.Join(homedir, nuvolarisConfigDir)
	_, err := osStat(fullPath)
	if osIsNotExist(err) {
		if err := osMkdir(fullPath, 0777); err != nil {
			return "", err
		}
	}
	return fullPath, nil
}

var rewriteKindConfigFile = func(configDir string) (string, error) {
	path := filepath.Join(configDir, kindConfigFile)
	if _, err := osStat(path); err == nil {
		osRemove(path)
	}
	if err := osWriteFile(path, KindYaml, 0600); err != nil {
		return "", err
	}
	return path, nil
}

var startCluster = func(configFile string) error {
	if err := kind("create", "cluster", "--wait=1m", "--config="+configFile); err != nil {
		return err
	}
	return nil

}

var stopCluster = func() error {
	if err := kind("delete", "cluster", "--name="+nuvolarisClusterName); err != nil {
		return err
	}
	return nil
}
