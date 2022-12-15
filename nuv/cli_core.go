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

//go:build !fast
// +build !fast

package main

import (
	"github.com/alecthomas/kong"
)

type CLI struct {
	Version    kong.VersionFlag `short:"v" cmd:"" passthrough:"" help:"show nuvolaris version"`
	Action     ActionCmd        `aliases:"a" cmd:"" passthrough:"" help:"work with actions"`
	Package    PackageCmd       `aliases:"p" cmd:"" passthrough:"" help:"work with packages"`
	Activation ActivationCmd    `aliases:"ac" cmd:"" passthrough:"" help:"work with actions"`
	Rule       RuleCmd          `aliases:"ru" cmd:"" passthrough:"" help:"work with rules"`
	Trigger    TriggerCmd       `aliases:"tr" cmd:"" passthrough:"" help:"work with triggers"`
	Project    ProjectCmd       `aliases:"pr" cmd:"" passthrough:"" help:"work with projects"`

	// utils
	Invoke InvokeCmd `aliases:"i" cmd:"" help:"invoke an action and return the result"`
	Url    UrlCmd    `cmd:"" passthrough:"" help:"show url of an action"`
	Logs   LogsCmd   `aliases:"l" cmd:"" passthrough:"" help:"show activation logs"`
	Result ResultCmd `aliases:"r"  cmd:"" passthrough:"" help:"show activation results"`
	Poll   PollCmd   `aliases:"po" cmd:"" help:"poll activations"`
	Bundle BundleCmd `aliases:"bu" cmd:"" help:"creates a web application bundle"`

	// Setup
	Setup      SetupCmd      `cmd:"" help:"setup nuvolaris"`
	Auth       AuthCmd       `cmd:"" help:"configure authentication"`
	Devcluster DevClusterCmd `cmd:"" help:"create or destroy kind k8s cluster"`

	// work in progress
	Scan ScanCmd `cmd:"" help:"scan subcommand" hidden:""`
	S3   S3Cmd   `cmd:"" name:"s3" help:"s3 subcommand" hidden:""`
	Wsk  WskCmd  `cmd:"" passthrough:"" help:"legacy wsk subcommand"`

	// not to be seen by users
	Task TaskCmd `cmd:"" passthrough:"" help:"task subcommand" hidden:""`
	Kind KindCmd `cmd:"" passthrough:"" help:"kind subcommand" hidden:""`
}
