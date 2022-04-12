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
	"context"
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	"strings"

	coreV1 "k8s.io/api/core/v1"
	extclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// KubeClient represents the wrapper of Kubernetes API clients
type KubeClient struct {
	clientset       kubernetes.Interface
	apiextclientset extclientset.Interface
	namespace       string
	ctx             context.Context
	cfg             *rest.Config
}

func initClients(logger *Logger, createDevcluster bool, k8sContext string) (*KubeClient, error) {

	if createDevcluster {
		fmt.Println("Starting devcluster...")
		cfg, err := configKind()
		if err != nil {
			return nil, err
		}
		err = cfg.manageKindCluster(logger, "create")
		if err != nil {
			return nil, err
		}
	}

	kubeconfig := flag.String("kubeconfig", getKubeconfigPath(), "")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("looks like no cluster is running. Run nuv devcluster create or nuv setup --devcluster")
	}

	err = assertNuvolarisContext(k8sContext)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %s", err)
	}

	apics, err := extclientset.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create apiextensions client: %s", err)
	}

	return &KubeClient{
		clientset:       clientset,
		apiextclientset: apics,
		namespace:       "nuvolaris",
		ctx:             context.Background(),
		cfg:             config,
	}, nil
}

func assertNuvolarisContext(k8sContext string) error {
	config, err := getK8sConfig()
	if err != nil {
		return err
	}

	var nuvolarisContext string

	for context := range config.Contexts {
		if context == k8sContext {
			nuvolarisContext = context
			break
		}
	}

	if nuvolarisContext == "" {
		return fmt.Errorf("context not found")
	}

	config.CurrentContext = nuvolarisContext
	err = clientcmd.ModifyConfig(clientcmd.NewDefaultPathOptions(), config, true)
	if err != nil {
		return fmt.Errorf("error ModifyConfig: %w", err)
	}

	fmt.Println("✓ Current context set to", nuvolarisContext)
	return nil
}

func (c *KubeClient) getNuvolarisNamespace() (*coreV1.Namespace, error) {
	ns, err := c.clientset.CoreV1().Namespaces().Get(c.ctx, c.namespace, metaV1.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, nil
		}
		return nil, err
	}
	return ns, nil
}

func (c *KubeClient) createNuvolarisNamespace() error {
	ns, err := c.getNuvolarisNamespace()
	if err != nil {
		return err
	}
	if ns == nil {
		namespace := &coreV1.Namespace{
			ObjectMeta: metaV1.ObjectMeta{
				Name: c.namespace,
			},
		}
		_, err := c.clientset.CoreV1().Namespaces().Create(c.ctx, namespace, metaV1.CreateOptions{})
		if err != nil {
			fmt.Println("failed creation of namespace nuvolaris")
			return err
		}
		fmt.Println("✓ Namespace nuvolaris created")
		return nil
	}
	fmt.Println("namespace nuvolaris already exists...skipping")
	return nil
}

func (c *KubeClient) cleanup() error {

	_, err := c.clientset.CoreV1().Namespaces().Get(c.ctx, c.namespace, metaV1.GetOptions{})
	if err != nil {
		fmt.Println("nuvolaris namespace not found. Nothing to do.")
		return nil
	}

	//manually remove wsk controller
	//to avoid namespace staying forever in Terminating state
	//to find out what resources are preventing deletion of namespace, run
	//kubectl api-resources --verbs=list --namespaced -o name | xargs -n 1 kubectl get -n nuvolaris

	client, err := restClient(c.cfg)
	if err != nil {
		return err
	}

	patch := []byte(`{"metadata":{"finalizers":[]}}`)
	err = client.Patch(types.MergePatchType).Namespace(c.namespace).Resource(CRDPlural).Name(wskObjectName).Body(patch).Do(c.ctx).Error()
	if err != nil {
		return err
	}
	err = client.Delete().Namespace(c.namespace).Resource(CRDPlural).Name(wskObjectName).Do(c.ctx).Error()
	if err != nil {
		return err
	}

	err = c.clientset.CoreV1().Namespaces().Delete(c.ctx, c.namespace, metaV1.DeleteOptions{})
	if err != nil {
		return err
	}

	fmt.Println("waiting for nuvolaris namespace to be terminated...a little patience please")
	waitForNamespaceToBeTerminated(c, c.namespace)
	fmt.Println("nuvolaris uninstalled.")
	return nil
}
