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

//functions declared as vars for mocking in unit tests
var osMkdirFunc = os.Mkdir
var osStatFunc = os.Stat
var osIsNotExistFunc = os.IsNotExist
var osWriteFileFunc = os.WriteFile
var osRemoveFunc = os.Remove
var createClusterFunc = createCluster
var destroyClusterFunc = destroyCluster
var KindFunc = Kind

func manageKindCluster(action string) error {

	switch action {
	case "create":
		if err := createClusterFunc(false); err != nil {
			return err
		}
	case "destroy":
		if err := destroyClusterFunc(false); err != nil {
			return err
		}
	default:
		fmt.Println("did you mean nuv devcluster create/destroy?")
	}
	return nil
}

func createCluster(dryRun bool) (err error) {
	defer func() {
		if err != nil {
			fmt.Errorf("error in create cluster: %w", err)
		}
	}()

	clusterIsRunning, err := clusterAlreadyRunning(dryRun)
	if err != nil {
		return err
	}
	if clusterIsRunning {
		fmt.Println("nuvolaris kind cluster is already running. Skipping...")
		return nil
	}
	homedir, err := GetHomedir()
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

func destroyCluster(dryRun bool) error {
	clusterIsRunning, err := clusterAlreadyRunning(dryRun)
	if err != nil {
		return err
	}
	if clusterIsRunning {
		err := stopCluster()
		if err != nil {
			return err
		}
		fmt.Println("kind cluster nuvolaris destroyed")
	} else {
		fmt.Println("kind cluster nuvolaris not found. Skipping...")
	}
	return nil
}

func clusterAlreadyRunning(dryRun bool) (bool, error) {
	out, err := sysErr(dryRun, "@nuv kind get clusters")

	if err != nil {
		return false, err
	}

	if strings.Contains(out, nuvolarisClusterName) {
		return true, nil
	} else {
		return false, nil
	}
}

func createNuvolarisConfigDirIfNotExists(homedir string) (string, error) {
	fullPath := filepath.Join(homedir, nuvolarisConfigDir)
	_, err := osStatFunc(fullPath)
	if osIsNotExistFunc(err) {
		if err := osMkdirFunc(fullPath, 0777); err != nil {
			return "", err
		}
	}
	return fullPath, nil
}

func rewriteKindConfigFile(configDir string) (string, error) {
	path := filepath.Join(configDir, kindConfigFile)
	if _, err := osStatFunc(path); err == nil {
		osRemoveFunc(path)
	}
	if err := osWriteFileFunc(path, KindYaml, 0600); err != nil {
		return "", err
	}
	return path, nil
}

func startCluster(configFile string) error {
	if err := KindFunc("create", "cluster", "--wait=1m", "--config="+configFile); err != nil {
		return err
	}
	return nil

}

func stopCluster() error {
	if err := KindFunc("delete", "cluster", "--name="+nuvolarisClusterName); err != nil {
		return err
	}
	return nil
}
