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
package preflight

import (
	"fmt"

	. "github.com/nuvolaris/nuvolaris-cli/nuv/globals"
	utils "github.com/nuvolaris/nuvolaris-cli/nuv/utils"
)

func Example_ensureDockerVersion() {
	utils.DryRunPush("19.03.5", "10.03.5", MinDockerVersion, "!no docker")

	p := PreflightChecksPipeline{dryRun: true}
	p.step(ensureDockerVersion)
	fmt.Println(p.err)

	p = PreflightChecksPipeline{dryRun: true}
	p.step(ensureDockerVersion)
	fmt.Println(p.err)

	p = PreflightChecksPipeline{dryRun: true}
	p.step(ensureDockerVersion)
	fmt.Println(p.err)

	p = PreflightChecksPipeline{dryRun: true}
	p.step(ensureDockerVersion)
	fmt.Println(p.err)
	// Output:
	// docker version --format {{.Server.Version}}
	// <nil>
	// docker version --format {{.Server.Version}}
	// installed docker version 10.3.5 is no longer supported
	// docker version --format {{.Server.Version}}
	// <nil>
	// docker version --format {{.Server.Version}}
	// no docker
}

func Example_isInHomePath() {
	p := PreflightChecksPipeline{dir: "/home/nuvolaris"}
	p.step(isInHomePath)
	fmt.Println(p.err)

	p = PreflightChecksPipeline{dir: "/var/run"}
	p.step(isInHomePath)
	fmt.Println(p.err)

	p = PreflightChecksPipeline{dir: ""}
	p.step(isInHomePath)
	fmt.Println(p.err)
	// Output:
	// <nil>
	// work directory /var/run should be below your home directory /home/nuvolaris;
	// this is required to be accessible by Docker
	// <nil>
}

func Example_checkDockerMemory() {
	p := PreflightChecksPipeline{dockerData: "\nTotal Memory: 11GiB\n"}
	p.step(checkDockerMemory)
	fmt.Println(p.err)

	p = PreflightChecksPipeline{dockerData: "\nTotal Memory: 3GiB\n"}
	p.step(checkDockerMemory)
	fmt.Println(p.err)
	// Output:
	// <nil>
	// nuv needs 4GB memory allocatable on docker
}
