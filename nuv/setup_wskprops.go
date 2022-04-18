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

type WskPropsPipeline struct {
	kubeClient *KubeClient
	k8sContext string
	err        error
	logger     *Logger
}

type wskSetupStep func(sp *WskPropsPipeline)

func (wsp *WskPropsPipeline) wStep(f wskSetupStep) {
	if wsp.err != nil {
		return
	}
	f(wsp)
}

func setupWskProps(logger *Logger, cmd *WskPropsCmd) error {
	wsp := WskPropsPipeline{
		logger: logger,
	}
	if len(cmd.Context) == 0 {
		config, err := getK8sConfig()
		if err != nil {
			return err
		}
		wsp.k8sContext = config.CurrentContext
	} else {
		wsp.k8sContext = cmd.Context
	}

	wsp.wStep(assertClusterConfig)
	wsp.wStep(writeNuvClusterConfig)

	return wsp.err
}

func assertClusterConfig(wsp *WskPropsPipeline) {
	wsp.kubeClient, wsp.err = initClients(wsp.k8sContext)
}

func writeNuvClusterConfig(wsp *WskPropsPipeline) {
	wsp.err = waitForApihostSet(wsp.kubeClient, NuvolarisConfigmapName)
	wsp.err = writeConfigToWskProps(wsp.kubeClient, NuvolarisConfigmapName)
}
