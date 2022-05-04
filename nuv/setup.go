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

type SetupCmd struct {
	Devcluster bool   `help:"start dev kind k8s cluster" xor:"devcluster-or-uninstall-or-context"`
	Configure  bool   `help:"generate configuration file"`
	ImageTag   string `default:"${image_tag}" help:"nuvolaris operator docker image tag to deploy"`
	Uninstall  string `help:"uninstall nuvolaris from given context" xor:"devcluster-or-uninstall-or-context"`
	Context    string `help:"set kubernetes context to install nuvolaris" xor:"devcluster-or-uninstall-or-context"`
	Apihost    string `help:"set kubernetes host IP"`
}

func (setupCmd *SetupCmd) Run(logger *Logger) error {
	return setupNuvolaris(logger, setupCmd)
}
