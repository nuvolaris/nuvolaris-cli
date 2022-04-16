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
	"time"
)

type UninstallCmd struct {
	Args []string `optional:"" name:"args" help:"uninstall nuvolaris"`
}

func (uninstallCmd *UninstallCmd) Run() error {
	return uninstallNuvolaris()
}

type UninstallPipeline struct {
	kubeClient     *KubeClient
	currentContext string
	err            error
}

type uninstallStep func(sp *UninstallPipeline)

func (up *UninstallPipeline) step(f uninstallStep) {
	if up.err != nil {
		return
	}
	f(up)
	time.Sleep(2 * time.Second)
}

func uninstallNuvolaris() error {
	up := UninstallPipeline{}
	up.step(getCurrentK8sConfig)
	up.step(initK8sClient)
	up.step(resetNuv)
	return up.err
}

func getCurrentK8sConfig(up *UninstallPipeline) {
	k8sConfig, _ := getK8sConfig()
	up.currentContext = k8sConfig.CurrentContext
}

func initK8sClient(up *UninstallPipeline) {
	up.kubeClient, up.err = initClients(up.currentContext)
}

func resetNuv(up *UninstallPipeline) {
	up.err = up.kubeClient.cleanup()
}
