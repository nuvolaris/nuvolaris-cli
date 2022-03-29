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

	coreV1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const operatorName = "nuvolaris-operator"
const operatorBinding = "nuvolaris-operator-crb"

var serviceAccount = &coreV1.ServiceAccount{
	ObjectMeta: metaV1.ObjectMeta{
		Name:      operatorName,
		Namespace: namespace,
		Labels:    map[string]string{"app": operatorName},
	},
}

var clusterRoleBinding = &rbacV1.ClusterRoleBinding{
	ObjectMeta: metaV1.ObjectMeta{
		Name:      operatorBinding,
		Namespace: namespace,
		Labels:    map[string]string{"app": operatorName},
	},
	Subjects: []rbacV1.Subject{
		{
			Kind:      "ServiceAccount",
			Name:      operatorName,
			Namespace: namespace,
		},
	},
	RoleRef: rbacV1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "ClusterRole",
		Name:     "cluster-admin",
	},
}

func configOperatorPod(operatorDockerImage string) *coreV1.Pod {
	return &coreV1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      operatorName,
			Namespace: namespace,
		},
		Spec: coreV1.PodSpec{
			Containers: []coreV1.Container{
				{
					Name:  operatorName,
					Image: operatorDockerImage,
				},
			},
			ServiceAccountName: operatorName,
		},
	}
}

func (c *KubeClient) createServiceAccount() error {

	_, err := c.clientset.CoreV1().ServiceAccounts(c.namespace).Get(c.ctx, operatorName, metaV1.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			_, err := c.clientset.CoreV1().ServiceAccounts(c.namespace).Create(c.ctx, serviceAccount, metaV1.CreateOptions{})
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
	_, err := c.clientset.RbacV1().ClusterRoleBindings().Get(c.ctx, operatorBinding, metaV1.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			_, err := c.clientset.RbacV1().ClusterRoleBindings().Create(c.ctx, clusterRoleBinding, metaV1.CreateOptions{})
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

func (c *KubeClient) createOperatorPod(dockerImg string) error {
	_, err := getPod(c, operatorName)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			fmt.Println("Deploying nuvolaris operator image " + dockerImg)
			_, err := c.clientset.CoreV1().Pods(c.namespace).Create(c.ctx, configOperatorPod(dockerImg), metaV1.CreateOptions{})
			if err != nil {
				return err
			}
			fmt.Println("Waiting for nuvolaris operator pod...hang tight")
			err = waitForPodRunning(c, operatorName)
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
