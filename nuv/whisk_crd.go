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
	"fmt"
	"strings"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
)

const (
	CRDKind       string = "Whisk"
	CRDPlural     string = "whisks"
	CRDGroup      string = "nuvolaris.org"
	CRDVersion    string = "v1"
	namespace     string = NuvolarisNamespace
	wskObjectName string = "controller"
	apiVersion    string = "nuvolaris.org/v1"
)

var preserveUnknownFields = true

var SchemeGroupVersion = schema.GroupVersion{Group: CRDGroup, Version: CRDVersion}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Whisk{},
		&WhiskList{},
	)
	metaV1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

func restClient(cfg *rest.Config) (*rest.RESTClient, error) {
	scheme := runtime.NewScheme()
	SchemeBuilder := runtime.NewSchemeBuilder(addKnownTypes)
	if err := SchemeBuilder.AddToScheme(scheme); err != nil {
		return nil, err
	}
	config := *cfg
	config.GroupVersion = &SchemeGroupVersion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.NewCodecFactory(scheme)
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func createWhisk(obj *Whisk, c *rest.RESTClient) error {
	result := c.Post().Namespace(namespace).Resource(CRDPlural).Body(obj).Do(context.Background())
	if result.Error() != nil {
		return fmt.Errorf(result.Error().Error())
	}
	return nil
}

func getWhisk(c *rest.RESTClient) error {
	_, err := c.Get().Namespace(namespace).Resource(CRDPlural).
		Name(wskObjectName).DoRaw(context.Background())
	return err
}

func createWhiskOperatorObject(c *KubeClient, apiHost string) error {
	spec, err := readOrCreateCrdConfig(apiHost)
	if err != nil {
		return err
	}
	whisk := &Whisk{
		TypeMeta: metaV1.TypeMeta{
			Kind:       CRDKind,
			APIVersion: apiVersion,
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      wskObjectName,
			Namespace: namespace,
		},
		Spec: *spec,
	}
	client, err := restClient(c.cfg)
	if err != nil {
		return err
	}
	err = getWhisk(client)
	if err != nil {
		if strings.Contains(err.Error(), "not find") {
			err = createWhisk(whisk, client)
			if err != nil {
				return err
			}
			fmt.Println("✓ Openwhisk operator started")

			//TODO remove temporary workaround
			writeWskPropsFile(wskPropsKeyValue{
				wskPropsKey:   "AUTH",
				wskPropsValue: spec.OpenWhisk.Namespaces.Nuvolaris,
			})
			return nil
		}
		return err
	}
	fmt.Println("openwhisk operator already running...skipping")
	return nil
}
