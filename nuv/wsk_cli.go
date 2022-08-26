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

package main

// wsk cmd
type WskCmd struct {
	Args []string `arg:"" name:"args"`
}

func (wsk *WskCmd) Run() error {
	return Wsk([]string{"wsk"}, wsk.Args...)
}

// action cmd
type ActionCmd struct {
	Args []string `arg:"" name:"args"`
}

func (wsk *ActionCmd) Run() error {
	return Wsk([]string{"wsk", "action"}, wsk.Args...)
}

// action cmd
type ActivationCmd struct {
	Args []string `arg:"" name:"args"`
}

func (wsk *ActivationCmd) Run() error {
	return Wsk([]string{"wsk", "activation"}, wsk.Args...)
}

// package cmd
type PackageCmd struct {
	Args []string `arg:"" name:"args"`
}

func (wsk *PackageCmd) Run() error {
	return Wsk([]string{"wsk", "package"}, wsk.Args...)
}

// rule cmd
type RuleCmd struct {
	Args []string `arg:"" name:"args"`
}

func (wsk *RuleCmd) Run() error {
	return Wsk([]string{"wsk", "rule"}, wsk.Args...)
}

// trigger cmd
type TriggerCmd struct {
	Args []string `arg:"" name:"args"`
}

func (wsk *TriggerCmd) Run() error {
	return Wsk([]string{"wsk", "trigger"}, wsk.Args...)
}

// project cmd
type ProjectCmd struct {
	Args []string `arg:"" name:"args"`
}

func (wsk *ProjectCmd) Run() error {
	setWskPropsAsEnvVariable()
	return Wsk([]string{"wsk", "project"}, wsk.Args...)
}

// invoke cmd
type InvokeCmd struct {
	Name   string            `arg:""`
	Params map[string]string `arg:"" optional:""`
	File   string            `short:"f" default:""`
}

func (wsk *InvokeCmd) Run() error {
	cmd := []string{"wsk", "action", "invoke", "-r"}
	args := []string{wsk.Name}
	for k, v := range wsk.Params {
		args = append(args, "-p", k, v)
	}
	if wsk.File != "" {
		args = append(args, "-P", wsk.File)
	}
	return Wsk(cmd, args...)
}

// log cmd
type LogsCmd struct {
	Args []string `arg:"" optional:""`
}

func (wsk *LogsCmd) Run() error {
	if len(wsk.Args) == 0 {
		return Wsk([]string{"wsk", "activation", "logs"}, "--last")
	}
	return Wsk([]string{"wsk", "activation", "logs"}, wsk.Args...)

}

// result cmd
type ResultCmd struct {
	Args []string `arg:"" optional:""`
}

func (wsk *ResultCmd) Run() error {
	if len(wsk.Args) == 0 {
		return Wsk([]string{"wsk", "activation", "result", "--last"})
	}
	return Wsk([]string{"wsk", "activation", "result"}, wsk.Args...)
}

// result cmd
type UrlCmd struct {
	Args []string `arg:"" name:"args"`
}

func (wsk *UrlCmd) Run() error {
	return Wsk([]string{"wsk", "action", "get", "--url"}, wsk.Args...)
}

// result cmd
type PollCmd struct {
}

func (wsk *PollCmd) Run() error {
	return Wsk([]string{"wsk", "activation", "poll"})
}

// auth cmd
type AuthCmd struct {
	Apihost string `help:"manually specify apihost" type:"string"`
	Auth    string `help:"manually specify authorization key" type:"string"`
	Redis   string `help:"manually specify redis url" type:"string"`
	Mongodb string `help:"manually specify mongodb url" type:"string"`
	Show    bool   `help:"show current auth"`
}

func (s *AuthCmd) Run(logger *Logger) error {
	return setupWskProps(s)
}
