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
	"os"
	"path/filepath"
)

var homeDir, _ = GetHomeDir()

func Example_devClusterWrongAction() {

	config := KindConfig{}
	config.manageKindCluster(NewLogger(), "delete")
	// Output:
	// did you mean nuv devcluster create/destroy?
}

func Example_devClusterAlreadyRunning() {
	config := KindConfig{
		nuvolarisClusterName: "nuvolaris",
		kind: func(...string) error {
			fmt.Println("nuvolaris")
			return nil
		},
	}
	config.manageKindCluster(NewLogger(), "create")
	// Output:
	// nuvolaris kind cluster is already running...skipping
}

func Example_multipleDevClustersRunningWithNuvolaris() {
	config := KindConfig{
		nuvolarisClusterName: "nuvolaris",
		kind: func(...string) error {
			fmt.Println("other")
			fmt.Println("nuvolaris")
			return nil
		},
	}
	config.manageKindCluster(NewLogger(), "create")
	// Output:
	// nuvolaris kind cluster is already running...skipping
}

func Example_preflightChecksNok() {
	config := KindConfig{
		nuvolarisClusterName: "nuvolaris",
		homedir:              homeDir,
		kind: func(...string) error {
			return nil
		},
		preflightChecks: func(*Logger, string) error {
			return fmt.Errorf("docker is not running")
		},
	}
	config.manageKindCluster(NewLogger(), "create")
	// Output:
	// Running Preflight checks...
}

func Example_successfullClusterStartFromScratch() {

	config := KindConfig{
		homedir:              homeDir,
		kindYaml:             kind_yaml,
		nuvolarisClusterName: "nuvolaris",
		nuvolarisConfigDir:   ".nuvolaris",
		kindConfigFile:       "kind.yaml",
		fullConfigPath:       "",
		preflightChecks: func(*Logger, string) error {
			return nil
		},
		kind: func(...string) error {
			return nil
		},
	}
	fullPath := filepath.Join(config.homedir, config.nuvolarisConfigDir)
	os.RemoveAll(fullPath)

	config.manageKindCluster(NewLogger(), "create")
	// Output:
	// Running Preflight checks...
	// Preflight checks passed!
	// nuvolaris config dir created
	// kind.yaml written
	// Starting nuvolaris kind cluster... hang tight
	// Nuvolaris kind cluster started. Have a nice day! 👋
}

func Example_destroyRunningCluster() {

	config := KindConfig{
		nuvolarisClusterName: "nuvolaris",
		kind: func(...string) error {
			fmt.Println("nuvolaris")
			return nil
		},
	}

	config.manageKindCluster(NewLogger(), "destroy")
	// Output:
	// nuvolaris
	// kind cluster nuvolaris destroyed
}

func Example_destroyClusterNotRunning() {

	config := KindConfig{
		nuvolarisClusterName: "nuvolaris",
		kind: func(...string) error {
			return nil
		},
	}

	config.manageKindCluster(NewLogger(), "destroy")
	// Output:
	// kind cluster nuvolaris not found...skipping
}
