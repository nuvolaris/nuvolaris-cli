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
	wsk func(...string) error
}

func readinessProbe(c *KubeClient) error {
	fmt.Println("Reading cluster config...")
	apihost, err := readAnnotationFromConfigmap(c, c.namespace, "config", "apihost")
	if err != nil {
		return err
	}
	writeWskPropertiesFile(apihost)
	fmt.Println("wsk properties file written...")

	wsk_probe := WskProbe{wsk: Wsk}

	fmt.Println("Waiting for openwhisk pod to complete...waiting is the hardest part 💚")
	err = waitForPodCompleted(c, "wsk-prewarm-nodejs14", TimeoutInSec)
	if err != nil {
		return err
	}
	fmt.Println("✓ Openwhisk running")

	fmt.Println("Creating an action...")
	hello_content := []byte("function main(args) { return { \"body\":\"hello from Nuvolaris\"} }")
	path, err := WriteFileToNuvolarisConfigDir("hello.js", hello_content)
	if err != nil {
		return err
	}
	err = wsk_probe.waitFor(TimeoutInSec, wsk_probe.isActionCreated(path))
	if err != nil {
		return err
	}

	fmt.Println("✓ Openwhisk action succesfully created")

	fmt.Println("Invoking action...")

	err = wsk_probe.wsk("action", "invoke", "hello", "-r")
	if err != nil {
		return err
	}

	fmt.Println("✓ Openwhisk action succesfully invoked. Done.")
	fmt.Println("  You are all set! Thanks for using Nuvolaris 😊")
	return nil
}

func (probe *WskProbe) isOpenWhiskDeployed() wait.ConditionFunc {
	return func() (bool, error) {
		fmt.Printf(".")

		err := probe.wsk("namespace", "get")
		if err != nil {
			return false, nil
		}
		return true, nil
	}
}

func (probe *WskProbe) isActionCreated(path_to_hello string) wait.ConditionFunc {
	return func() (bool, error) {
		fmt.Printf(".")

		err := probe.wsk("action", "create", "hello", path_to_hello)
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

func (probe *WskProbe) waitFor(timeout_sec int, function wait.ConditionFunc) error {
	return wait.PollImmediate(time.Second, time.Duration(timeout_sec)*time.Second, function)
}
