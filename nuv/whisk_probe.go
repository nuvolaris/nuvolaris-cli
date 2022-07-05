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
	_ "embed"
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

type WskProbe struct {
	wsk func([]string, ...string) error
}

//go:embed embed/hello.js
var helloContent []byte

func readinessProbe(c *KubeClient) error {
	fmt.Print("Reading Nuvolaris cluster config...")
	err := waitForApihostSet(c, NuvolarisConfigmapName)
	if err != nil {
		return err
	}

	writeConfigToWskProps(c, NuvolarisConfigmapName)

	wskProbe := WskProbe{wsk: Wsk}

	//fmt.Println("\nNow, waiting for openwhisk pods to start...waiting is the hardest part 💚")

	pods := []string{"controller", "couchdb", "redis"}
	for _, pod := range pods {
		fmt.Printf("\nWaiting for %s Running...", pod)
		waitForPod(c, pod+"-0")
		err = waitForPodRunning(c, pod+"-0")
		if err != nil {
			return err
		}
	}

	runtimes := []string{"nodejs", "python", "go"}
	for _, runtime := range runtimes {
		fmt.Printf("\nWaiting for Runtime %s Ready...", runtime)
		podName := "preload-runtime-" + runtime
		debug("preloading " + podName)
		err = waitForPodCompleted(c, podName)
		if err != nil {
			return err
		}
	}

	fmt.Println("\n✓ Openwhisk running")

	fmt.Println("Creating an action.\n(It can take a few minutes please be patient)")
	path, err := WriteFileToNuvolarisConfigDir("hello.js", helloContent)
	if err != nil {
		return err
	}
	err = wskProbe.waitFor(TimeoutInSec, wskProbe.isActionCreated(path))
	if err != nil {
		return err
	}

	fmt.Println("✓ Openwhisk action successfully created")

	fmt.Println("Invoking action...")

	err = wskProbe.wsk([]string{"wsk", "action"}, "invoke", "hello", "-r")
	if err != nil {
		return err
	}

	fmt.Println("✓ Openwhisk action successfully invoked. Done.")
	fmt.Println("  You are all set! Thanks for using Nuvolaris 😊")
	return nil
}

func (probe *WskProbe) isOpenWhiskDeployed() wait.ConditionFunc {
	return func() (bool, error) {
		fmt.Printf(".")

		err := probe.wsk([]string{"wsk", "namespace"}, "get")
		if err != nil {
			return false, nil
		}
		return true, nil
	}
}

func (probe *WskProbe) isActionCreated(pathToHello string) wait.ConditionFunc {
	return func() (bool, error) {
		fmt.Printf(".")

		err := probe.wsk([]string{"wsk", "action"}, "create", "hello", pathToHello)
		if err != nil {
			if strings.Contains(err.Error(), "resource already exists") {
				fmt.Println("Openwhisk action already created...skipping")
				return true, nil
			}
			return false, nil
		}
		return true, nil
	}
}

func (probe *WskProbe) waitFor(timeoutSec int, function wait.ConditionFunc) error {
	return wait.PollImmediate(time.Second, time.Duration(timeoutSec)*time.Second, function)
}
