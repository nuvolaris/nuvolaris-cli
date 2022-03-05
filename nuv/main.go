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

// CLI_VERSION holds the current version, to be set by the build with
//  go build -ldflags "-X main.CLI_VERSION=<version>"
var CLI_VERSION string = "latest"

type CLI struct {
	Deploy     DeployCmd        `cmd:"" help:"deploy a nuvolaris cluster"`
	Destroy    DestroyCmd       `cmd:"" help:"destroy a nuvolaris cluster"`
	Wsk        WskCmd           `cmd:"" passthrough:"" help:"wsk subcommand"`
	Task       TaskCmd          `cmd:"" help:"task subcommand"`
	Kind       KindCmd          `cmd:"" help:"kind subcommand"`
	Devcluster DevClusterCmd    `cmd:"" help:"create or destroy kind k8s cluster"`
	Setup      SetupCmd         `cmd:"" help:"setup nuvolaris"`
	Version    kong.VersionFlag `short:"v" help:"print nuvolaris version"`
}

func main() {
	cli := CLI{}
	ctx := kong.Parse(&cli,
		kong.Name(Name),
		kong.Description(Description),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact:             true,
			NoExpandSubcommands: false,
		}),
		kong.Vars{"version": CLI_VERSION},
	)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
