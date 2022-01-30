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

	log "github.com/sirupsen/logrus"
)

type PreflightChecksPipeline struct {
	dryRun            bool
	skipDockerVersion bool
	dir               string
	dockerData        string
	err               error
}

type checkStep func(pd *PreflightChecksPipeline)

func (p *PreflightChecksPipeline) step(f checkStep) {
	if p.err != nil {
		return
	}
	f(p)
}

// RunPreflightChecks performs preflight checks
// checks docker version, available memory and dir paths
func RunPreflightChecks(dir string) error {

	// Preflight Checks pipeline
	// TODO: keep skipDockerVersion and dryRun?
	pp := PreflightChecksPipeline{skipDockerVersion: false, dryRun: false, dir: dir}

	pp.step(extractDockerInfo)
	pp.step(checkDockerMemory)
	pp.step(ensureDockerVersion)
	pp.step(isInHomePath)

	return pp.err
}

func extractDockerInfo(p *PreflightChecksPipeline) {
	p.dockerData, p.err = dockerInfo(p.dryRun)
}

func checkDockerMemory(p *PreflightChecksPipeline) {
	var search = regexp.MustCompile(`Total Memory: (.*)`)
	result := search.FindString(string(p.dockerData))
	if result == "" {
		p.err = fmt.Errorf("docker is not running")
		return
	}
	fmt.Println("docker is running...")
	mem := strings.Split(result, ":")
	memory := strings.TrimSpace(mem[1])
	n, err := units.ParseStrictBytes(memory)
	if err != nil {
		p.err = err
		return
	}
	log.Debug("mem:", n)
	if n <= int64(MinDockerMem) {
		p.err = fmt.Errorf("nuv needs 4GB memory allocatable on docker")
		return
	}
	fmt.Println("enough memory to allocate...")
}

func ensureDockerVersion(p *PreflightChecksPipeline) {
	if p.skipDockerVersion {
		return
	}
	version, err := dockerVersion(p.dryRun)
	if err != nil {
		p.err = err
		return
	}
	vA := semver.New(MinDockerVersion)
	vB := semver.New(strings.TrimSpace(version))
	if vB.Compare(*vA) == -1 {
		p.err = fmt.Errorf("installed docker version %s is no longer supported", vB)
		return
	}
	fmt.Printf("installed docker version %s ok...\n", vB)
}

func isInHomePath(p *PreflightChecksPipeline) {
	// do not check if the directory is empty
	if p.dir == "" {
		return
	}
	homePath, err := GetHomeDir()
	if err != nil {
		p.err = err
		return
	}
	dir, err := filepath.Abs(p.dir)
	if err != nil {
		p.err = err
		return
	}
	if !strings.HasPrefix(dir, homePath) {
		p.err = fmt.Errorf("work directory %s should be below your home directory;\nthis is required to be accessible by Docker", dir)
		return
	}
	fmt.Println("dir tree ok...")
}
