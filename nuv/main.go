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
	"github.com/alecthomas/kong"
)

// CLIVersion holds the current version, to be set by the build with
//  go build -ldflags "-X main.CLIVersion=<version>"
var CLIVersion = "latest"

// ImageTag holds the version of the Docker image used for the nuvolaris
// operator used in setup
var ImageTag = "0.2.0-trinity.22052321"

func main() {
	cli := CLI{}
	logger := NewLogger()

	ctx := kong.Parse(&cli,
		kong.Name(Name),
		kong.Description(Description),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact:             true,
			NoExpandSubcommands: true,
		}),
		kong.Vars{
			"version":   CLIVersion,
			"image_tag": ImageTag,
		},
		kong.Bind(logger),
	)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
