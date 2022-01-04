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
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alecthomas/units"
	"github.com/coreos/go-semver/semver"
	"github.com/mitchellh/go-homedir"

	log "github.com/sirupsen/logrus"
)

// Preflight perform preflight checks
func Preflight(skipDockerVersion bool, dir string) (string, error) {
	info, err := DockerInfo(false)
	if err != nil {
		return "", err
	}
	err = CheckDockerMemory(info)
	if err != nil {
		return "", err
	}
	if !skipDockerVersion {
		err = EnsureDockerVersion(false)
		if err != nil {
			return "", err
		}
	}
	err = IsInHomePath(dir)
	if err != nil {
		return "", err
	}
	return info, nil
}

func EnsureDockerVersion(dryRun bool) error {
	version, err := DockerVersion(dryRun)
	if err != nil {
		return err
	}
	vA := semver.New(MinDockerVersion)
	vB := semver.New(strings.TrimSpace(version))
	if vB.Compare(*vA) == -1 {
		return fmt.Errorf("installed docker version %s is no longer supported", vB)
	}
	return nil
}

func IsInHomePath(dir string) error {
	// do not check if the directory is empty
	if dir == "" {
		return nil
	}
	homePath, err := homedir.Dir()
	if err != nil {
		return err
	}
	dir, err = filepath.Abs(dir)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(dir, homePath) {
		return fmt.Errorf("work directory %s should be below your home directory %s;\nthis is required to be accessible by Docker", dir, homePath)
	}
	return nil
}

// CheckDockerMemory checks docker memory
func CheckDockerMemory(info string) error {
	var search = regexp.MustCompile(`Total Memory: (.*)`)
	result := search.FindString(string(info))
	if result == "" {
		return fmt.Errorf("docker is not running")
	}
	mem := strings.Split(result, ":")
	memory := strings.TrimSpace(mem[1])
	n, err := units.ParseStrictBytes(memory)
	if err != nil {
		return err
	}
	log.Debug("mem:", n)
	//fmt.Println(n)
	if n <= int64(MinDockerMem) {
		return fmt.Errorf("nuv needs 4GB memory allocatable on docker")
	}
	return nil

}
