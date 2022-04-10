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

type SetupPipeline struct {
	kubeClient          *KubeClient
	createDevcluster    bool
	k8sContext          string
	operatorDockerImage string
	err                 error
	logger              *Logger
}

type setupStep func(sp *SetupPipeline)

func (sp *SetupPipeline) step(f setupStep) {
	if sp.err != nil {
		return
	}
	f(sp)
	time.Sleep(2 * time.Second)
}

func setupNuvolaris(logger *Logger, cmd *SetupCmd) error {
	imgTag := cmd.ImageTag

	sp := SetupPipeline{
		operatorDockerImage: "ghcr.io/nuvolaris/nuvolaris-operator:" + imgTag,
		logger:              logger,
	}

	sp.createDevcluster = cmd.Devcluster
	sp.k8sContext = cmd.Context

	sp.step(assertNuvolarisClusterConfig)

	if cmd.Reset {
		sp.step(resetNuvolaris)
	} else {
		sp.step(createNuvolarisNamespace)
		sp.step(deployWhiskCrd)
		sp.step(deployServiceAccount)
		sp.step(deployClusterRoleBinding)
		sp.step(runNuvolarisOperatorPod)
		sp.step(deployOperatorObject)
		sp.step(waitForOpenWhiskReady)
	}
	return sp.err
}

func assertNuvolarisClusterConfig(sp *SetupPipeline) {
	sp.kubeClient, sp.err = initClients(sp.logger, sp.createDevcluster, sp.k8sContext)
}

func createNuvolarisNamespace(sp *SetupPipeline) {
	sp.err = sp.kubeClient.createNuvolarisNamespace()
}

func deployWhiskCrd(sp *SetupPipeline) {
	sp.err = sp.kubeClient.deployCRD()
}

func deployServiceAccount(sp *SetupPipeline) {
	sp.err = sp.kubeClient.createServiceAccount()
}

func deployClusterRoleBinding(sp *SetupPipeline) {
	sp.err = sp.kubeClient.createClusterRoleBinding()
}

func runNuvolarisOperatorPod(sp *SetupPipeline) {
	sp.err = sp.kubeClient.createOperatorPod(sp.operatorDockerImage)
}

func deployOperatorObject(sp *SetupPipeline) {
	sp.err = createWhiskOperatorObject(sp.kubeClient)
}

func waitForOpenWhiskReady(sp *SetupPipeline) {
	sp.err = readinessProbe(sp.kubeClient)
}

func resetNuvolaris(sp *SetupPipeline) {
	sp.err = sp.kubeClient.cleanup()
}
