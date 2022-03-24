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
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ensureDockerVersion(t *testing.T) {
	DryRunPush("19.03.5", "10.03.5", MinDockerVersion, "!no docker")

	p := PreflightChecksPipeline{dryRun: true, logger: NewLogger()}
	p.step(ensureDockerVersion)
	assert.NoError(t, p.err)

	p = PreflightChecksPipeline{dryRun: true, logger: NewLogger()}
	p.step(ensureDockerVersion)
	assert.Error(t, p.err)

	p = PreflightChecksPipeline{dryRun: true, logger: NewLogger()}
	p.step(ensureDockerVersion)
	assert.NoError(t, p.err)

	p = PreflightChecksPipeline{dryRun: true, logger: NewLogger()}
	p.step(ensureDockerVersion)
	assert.Error(t, p.err)
}

func Test_isInHomePath(t *testing.T) {
	homedir, _ := GetHomeDir()
	p := PreflightChecksPipeline{dir: homedir, logger: NewLogger()}
	p.step(isInHomePath)
	assert.NoError(t, p.err)

	p = PreflightChecksPipeline{dir: "/var/run", logger: NewLogger()}
	p.step(isInHomePath)
	assert.Error(t, p.err)

	p = PreflightChecksPipeline{dir: "", logger: NewLogger()}
	p.step(isInHomePath)
	assert.NoError(t, p.err)
}

func Test_checkDockerMemory(t *testing.T) {
	p := PreflightChecksPipeline{dockerData: "\nTotal Memory: 11GiB\n", logger: NewLogger()}
	p.step(checkDockerMemory)
	assert.NoError(t, p.err)

	p = PreflightChecksPipeline{dockerData: "\nTotal Memory: 3GiB\n", logger: NewLogger()}
	p.step(checkDockerMemory)
	assert.Error(t, p.err)
}
