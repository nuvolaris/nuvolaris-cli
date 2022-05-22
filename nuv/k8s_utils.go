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
	"path/filepath"
	"strings"
	"time"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/clientcmd"
)

const apihostAnnotation = "apihost"
const nuvAnnotationPrefix = "nuvolaris-"

func isPodRunning(c *KubeClient, podName string) wait.ConditionFunc {
	return func() (bool, error) {
		fmt.Printf(".")

		pod, err := getPod(c, podName)
		if err != nil {
			return false, err
		}

		switch pod.Status.Phase {
		case coreV1.PodPending:
			return false, nil
		case coreV1.PodRunning:
			return true, nil
		case coreV1.PodFailed, coreV1.PodSucceeded:
			return false, fmt.Errorf("pod cannot start...aborting")
		}
		return false, nil
	}
}

func isPodCompleted(c *KubeClient, podName string) wait.ConditionFunc {
	return func() (bool, error) {
		fmt.Printf(".")

		pod, err := getPod(c, podName)
		if err != nil {
			return false, err
		}

		switch pod.Status.Phase {
		case coreV1.PodPending, coreV1.PodRunning:
			return false, nil
		case coreV1.PodSucceeded:
			return true, nil
		case coreV1.PodFailed:
			return false, fmt.Errorf("pod cannot start...aborting")
		}
		return false, nil
	}
}

func isPodCreated(c *KubeClient, podName string) wait.ConditionFunc {
	return func() (bool, error) {
		fmt.Printf(".")
		_, err := getPod(c, podName)
		if err != nil {
			return false, nil
		}
		return true, nil
	}
}

func isNamespaceTerminated(c *KubeClient, namespace string) wait.ConditionFunc {
	return func() (bool, error) {
		fmt.Printf(".")

		_, err := getNamespace(c, namespace)
		if err != nil {
			return true, err
		}
		return false, nil
	}
}

func isApihostSet(c *KubeClient, configmap string) wait.ConditionFunc {
	return func() (bool, error) {
		fmt.Printf(".")

		cm, err := getConfigmap(c, configmap)
		if err != nil {
			return false, err
		}

		if cm.Annotations == nil {
			return false, fmt.Errorf("no annotations found")
		}

		host := cm.Annotations[apihostAnnotation]
		if host == "https://pending" || host == "" {
			return false, nil
		} else {
			return true, nil
		}
	}
}
func isConfigmapReady(c *KubeClient, configmap string) wait.ConditionFunc {
	return func() (bool, error) {
		fmt.Printf(".")
		_, err := getConfigmap(c, configmap)
		if err != nil {
			return false, err
		}
		return true, nil
	}
}

func readClusterConfig(c *KubeClient, configmap string) (map[string]string, error) {
	waitForConfigmapReady(c, configmap)
	cm, err := getConfigmap(c, configmap)
	if err != nil {
		return nil, err
	}
	wskPropsEntries := make(map[string]string)
	for k, v := range cm.Annotations {
		if strings.HasPrefix(k, nuvAnnotationPrefix) {
			key := strings.TrimPrefix(k, nuvAnnotationPrefix)
			key = strings.ToUpper(key)
			key = strings.ReplaceAll(key, "-", "_")
			wskPropsEntries[key] = v
		}
	}
	apihost := cm.Annotations["apihost"]
	//TODO remove temporary workaround to replace https with http
	wskPropsEntries["APIHOST"] = strings.ReplaceAll(apihost, "https", "http")
	updateApihostInConfig(wskPropsEntries["APIHOST"])
	return wskPropsEntries, nil
}

func writeConfigToWskProps(c *KubeClient, configmapName string) error {
	wskPropsMap, err := readClusterConfig(c, configmapName)
	if err != nil {
		return err
	}
	wskPropsEntries := flattenWskPropsMap(wskPropsMap)
	return writeWskPropsFile(wskPropsEntries...)
}

func getK8sConfig() clientcmd.ClientConfig {
	kubeConfig := getKubeconfigPath()
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfig},
		&clientcmd.ConfigOverrides{})
}

func getKubeconfigPath() string {
	if home, _ := GetHomeDir(); home != "" {
		return filepath.Join(home, ".kube", "config")
	}
	return ""
}

func getPod(c *KubeClient, podName string) (*coreV1.Pod, error) {
	return c.clientset.CoreV1().Pods(c.namespace).Get(c.ctx, podName, metaV1.GetOptions{})
}

func getNamespace(c *KubeClient, namespace string) (*coreV1.Namespace, error) {
	return c.clientset.CoreV1().Namespaces().Get(c.ctx, namespace, metaV1.GetOptions{})
}

func getConfigmap(c *KubeClient, configmapName string) (*coreV1.ConfigMap, error) {
	return c.clientset.CoreV1().ConfigMaps(c.namespace).Get(c.ctx, configmapName, metaV1.GetOptions{})
}

func waitForPod(c *KubeClient, podName string) error {
	return waitFor(c, isPodCreated, podName)
}

func waitForPodRunning(c *KubeClient, podName string) error {
	return waitFor(c, isPodRunning, podName)
}

func waitForPodCompleted(c *KubeClient, podName string) error {
	return waitFor(c, isPodCompleted, podName)
}

func waitForNamespaceToBeTerminated(c *KubeClient, namespace string) error {
	return waitFor(c, isNamespaceTerminated, namespace)
}

func waitForConfigmapReady(c *KubeClient, configmap string) error {
	return waitFor(c, isConfigmapReady, configmap)
}

func waitForApihostSet(c *KubeClient, configmap string) error {
	return waitFor(c, isApihostSet, configmap)
}

func waitFor(c *KubeClient, f checkCondition, resourceName string) error {
	return wait.PollImmediate(time.Second, time.Duration(TimeoutInSec)*time.Second, f(c, resourceName))
}

type checkCondition func(c *KubeClient, resourceName string) wait.ConditionFunc
