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

type KindConfig struct {
	homedir              string
	kindYaml             []byte
	nuvolarisClusterName string
	nuvolarisConfigDir   string
	kindConfigFile       string
	fullConfigPath       string
	preflightChecks      func(*Logger, string) error
	kind                 func(...string) error
}

//go:embed embed/kind.yaml
var kindYaml []byte

func configKind() (*KindConfig, error) {

	homeDir, err := GetHomeDir()
	if err != nil {
		return nil, err
	}

	config := KindConfig{
		homedir:              homeDir,
		kindYaml:             kindYaml,
		nuvolarisClusterName: "nuvolaris",
		kindConfigFile:       "kind.yaml",
		fullConfigPath:       "",
		preflightChecks:      RunPreflightChecks,
		kind:                 Kind,
	}
	return &config, nil
}

func (config *KindConfig) manageKindCluster(logger *Logger, action string) error {
	if action == "create" {
		if err := config.createCluster(logger); err != nil {
			return err
		}
		return nil
	}
	if action == "destroy" {
		removeConfigYaml()
		if err := config.destroyCluster(); err != nil {
			return err
		}
		return nil
	}
	fmt.Println("subcommand not available")
	return nil
}

func (config *KindConfig) createCluster(logger *Logger) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("error in create cluster: %w", err)
		}
	}()

	clusterIsRunning, err := config.clusterAlreadyRunning()
	if err != nil {
		return err
	}
	if clusterIsRunning {
		logger.Info("nuvolaris kind cluster is already running...skipping")
		return nil
	}

	logger.Info("Running Preflight checks...")
	if err = config.preflightChecks(logger, config.homedir); err != nil {
		return err
	}
	logger.Info("Preflight checks passed!")

	_, err = GetOrCreateNuvolarisConfigDir()
	if err != nil {
		return err
	}

	fullConfigPath, err := config.rewriteKindConfigFile()
	if err != nil {
		return err
	}

	config.fullConfigPath = fullConfigPath

	logger.Info("Starting nuvolaris kind cluster... hang tight")
	if err = config.startCluster(); err != nil {
		return err
	}

	logger.Info("Nuvolaris kind cluster started. Have a nice day! 👋")
	return nil
}

func (config *KindConfig) destroyCluster() error {
	clusterIsRunning, err := config.clusterAlreadyRunning()
	if err != nil {
		return err
	}
	if clusterIsRunning {
		if err := config.stopCluster(); err != nil {
			return err
		}
		fmt.Println("kind cluster nuvolaris destroyed")
	} else {
		fmt.Println("kind cluster nuvolaris not found...skipping")
	}
	return nil
}

func removeConfigYaml() {
	// remove configuation
	homeDir, _ := GetHomeDir()
	configYaml := filepath.Join(homeDir, ".nuvolaris", "config.yaml")

	if _, err := os.Stat(configYaml); err == nil {
		err = os.Remove(configYaml)
		if err == nil {
			fmt.Printf("%s removed \n", configYaml)
		} else {
			fmt.Printf("cannot remove %s - please remove it manually\n", configYaml)
		}
	}
}

func (config *KindConfig) clusterAlreadyRunning() (bool, error) {
	//capture cmd output
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := config.kind("get", "clusters")

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	if err != nil {
		return false, err
	}
	if strings.Contains(string(out), config.nuvolarisClusterName) {
		return true, nil
	} else {
		return false, nil
	}
}

func (config *KindConfig) rewriteKindConfigFile() (string, error) {
	nuvHomedir, err := GetOrCreateNuvolarisConfigDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(nuvHomedir, config.kindConfigFile)
	if _, err := os.Stat(path); err == nil {
		os.Remove(path)
	}

	// set the path for the data dir
	dataDir := filepath.Join(config.homedir, ".nuvolaris_data")
	// here docker is remote to we cannot know the remote home and we use /tmp
	if os.Getenv("DOCKER_HOST") != "" {
		dataDir = "/tmp/nuvolaris_data"
	}
	replacedConfigYaml := strings.ReplaceAll(string(config.kindYaml), "$NUV_DATA_DIR", dataDir)
	if err := os.WriteFile(path, []byte(replacedConfigYaml), 0600); err != nil {
		return "", err
	}
	fmt.Println(config.kindConfigFile + " written")
	return path, nil
}

func (config *KindConfig) startCluster() error {
	if err := config.kind("create", "cluster", "--wait=5m", "--config="+config.fullConfigPath); err != nil {
		return err
	}
	return nil
}

func (config *KindConfig) stopCluster() error {
	if err := config.kind("delete", "cluster", "--name="+config.nuvolarisClusterName); err != nil {
		return err
	}
	return nil
}
