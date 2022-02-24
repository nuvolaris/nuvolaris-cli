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
	"fmt"
	"strings"
	"time"

	core_v1 "k8s.io/api/core/v1"
	rbac_v1 "k8s.io/api/rbac/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const operator_name = "nuvolaris-operator"
const operator_image = "ghcr.io/nuvolaris/nuvolaris-operator:neo-22.0207.21"
const operator_binding = "nuvolaris-operator-crb"

var service_account = &core_v1.ServiceAccount{
	ObjectMeta: meta_v1.ObjectMeta{
		Name:      operator_name,
		Namespace: namespace,
		Labels:    map[string]string{"app": operator_name},
	},
}

var cluster_role_binding = &rbac_v1.ClusterRoleBinding{
	ObjectMeta: meta_v1.ObjectMeta{
		Name:      operator_binding,
		Namespace: namespace,
		Labels:    map[string]string{"app": operator_name},
	},
	Subjects: []rbac_v1.Subject{
		{
			Kind:      "ServiceAccount",
			Name:      operator_name,
			Namespace: namespace,
		},
	},
	RoleRef: rbac_v1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "ClusterRole",
		Name:     "cluster-admin",
	},
}
var operator_pod = &core_v1.Pod{
	ObjectMeta: meta_v1.ObjectMeta{
		Name:      operator_name,
		Namespace: namespace,
	},
	Spec: core_v1.PodSpec{
		Containers: []core_v1.Container{
			{
				Name:  operator_name,
				Image: operator_image,
			},
		},
		ServiceAccountName: operator_name,
	},
}

func (c *KubeClient) createServiceAccount() error {

	_, err := c.clientset.CoreV1().ServiceAccounts(c.namespace).Get(c.ctx, operator_name, meta_v1.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			_, err := c.clientset.CoreV1().ServiceAccounts(c.namespace).Create(c.ctx, service_account, meta_v1.CreateOptions{})
			if err != nil {
				return err
			}
			fmt.Println("✓ Service account created")
			return nil
		}
		return err
	}
	fmt.Println("service account already created...skipping")
	return nil
}

func (c *KubeClient) createClusterRoleBinding() error {
	_, err := c.clientset.RbacV1().ClusterRoleBindings().Get(c.ctx, operator_binding, meta_v1.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			_, err := c.clientset.RbacV1().ClusterRoleBindings().Create(c.ctx, cluster_role_binding, meta_v1.CreateOptions{})
			if err != nil {
				return err
			}
			fmt.Println("✓ Cluster role binding created")
			return nil
		}
		return err
	}
	fmt.Println("cluster role binding already created...skipping")
	return nil
}

func (c *KubeClient) createOperatorPod() error {
	_, err := getOperatorPod(c)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			_, err := c.clientset.CoreV1().Pods(c.namespace).Create(c.ctx, operator_pod, meta_v1.CreateOptions{})
			if err != nil {
				return err
			}
			fmt.Println("Waiting for nuvolaris operator pod...hang tight")
			err = waitForPodRunning(c, TimeoutInSec)
			if err != nil {
				return err
			}

			fmt.Println("")
			fmt.Println("✓ Nuvolaris operator pod running")
			return nil
		}
		return err
	}
	fmt.Println("nuvolaris operator pod already running...skipping")
	return nil
}

func getOperatorPod(c *KubeClient) (*core_v1.Pod, error) {
	return c.clientset.CoreV1().Pods(c.namespace).Get(c.ctx, operator_name, meta_v1.GetOptions{})
}

func isPodRunning(c *KubeClient) wait.ConditionFunc {
	return func() (bool, error) {
		fmt.Printf(".")

		pod, err := getOperatorPod(c)
		if err != nil {
			return false, err
		}

		switch pod.Status.Phase {
		case core_v1.PodPending:
			return false, nil
		case core_v1.PodRunning:
			return true, nil
		case core_v1.PodFailed, core_v1.PodSucceeded, core_v1.PodUnknown:
			return false, fmt.Errorf("nuvolaris-operator pod cannot start...aborting")
		}
		return false, nil
	}
}

func waitForPodRunning(c *KubeClient, timeout_sec int) error {
	return wait.PollImmediate(time.Second, time.Duration(timeout_sec)*time.Second, isPodRunning(c))
}
