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
	"regexp"

	"github.com/alecthomas/kong"
)

// CLIVersion holds the current version, to be set by the build with
//  go build -ldflags "-X main.CLIVersion=<version>"
var CLIVersion = "vX.Y.Z-milestone.build"

// ImageTag holds the version of the Docker image used for the nuvolaris
// operator used in setup
var defaultOperatorTag = "0.2.1-trinity.22070510"
var defaultOperatorImage = "ghcr.io/nuvolaris/nuvolaris-operator"

func main() {


	cli := CLI{}
	logger := NewLogger()

	operatorImage := os.Getenv("OPERATOR_IMAGE")
	if operatorImage == "" {
		operatorImage = defaultOperatorImage
	}

	operatorTag := os.Getenv("OPERATOR_TAG")
	if operatorTag == "" {
		operatorTag = defaultOperatorTag
	}

	// parse timestamp
	version := CLIVersion
	re, _ := regexp.Compile("v(.*)-([a-z]+).([0-9]+)")
	res := re.FindStringSubmatch(CLIVersion)
	if len(res) > 0 {
		version = fmt.Sprintf("nuv: v%s (milestone: %s, build: %s)", res[1], res[2], res[3])
	}

	ctx := kong.Parse(&cli,
		kong.Name(Name),
		kong.Description(Description),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact:             true,
			NoExpandSubcommands: true,
		}),
		kong.Vars{
			"version":        version + "\nnuvolaris-operator: " + operatorImage + ":" + operatorTag,
			"operator_image": operatorImage,
			"operator_tag":   operatorTag,
		},
		kong.Bind(logger),
	)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}

var Debugging = os.Getenv("DEBUG") != ""

func debug(format string, args ...any) {
	if Debugging {
		fmt.Printf(">>> "+format, args...)
	}
}
