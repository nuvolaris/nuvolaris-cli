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

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
)

const (
	CRDKind       string = "Whisk"
	CRDPlural     string = "whisks"
	CRDSingular   string = "whisk"
	CRDShortName  string = "wsk"
	CRDGroup      string = "nuvolaris.org"
	CRDVersion    string = "v1"
	FullCRDName   string = CRDPlural + "." + CRDGroup
	namespace     string = "nuvolaris"
	wskObjectName string = "controller"
	apiVersion    string = "nuvolaris.org/v1"
)

var preserveUnknownFields bool = true

func configureCRD() *apiextensions.CustomResourceDefinition {

	whiskCrd := apiextensions.CustomResourceDefinition{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      FullCRDName,
			Namespace: namespace,
		},
		Status: apiextensions.CustomResourceDefinitionStatus{
			StoredVersions: []string{CRDVersion},
		},
		Spec: apiextensions.CustomResourceDefinitionSpec{
			Scope: apiextensions.NamespaceScoped,
			Group: CRDGroup,
			Names: apiextensions.CustomResourceDefinitionNames{
				Kind:       CRDKind,
				Plural:     CRDPlural,
				Singular:   CRDSingular,
				ShortNames: []string{CRDShortName},
			},
			Versions: []apiextensions.CustomResourceDefinitionVersion{
				{
					Name:    CRDVersion,
					Served:  true,
					Storage: true,
					Subresources: &apiextensions.CustomResourceSubresources{
						Status: &apiextensions.CustomResourceSubresourceStatus{},
					},
					Schema: &apiextensions.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
							Type: "object",
							Properties: map[string]apiextensions.JSONSchemaProps{
								"spec": {
									Type: "object",
									Properties: map[string]apiextensions.JSONSchemaProps{
										"debug": {Type: "boolean"},
										"couchdb": {
											Type: "object",
											Properties: map[string]apiextensions.JSONSchemaProps{
												"whisk_admin": {Type: "string"},
											},
										},
										"mongodb": {
											Type: "object",
											Properties: map[string]apiextensions.JSONSchemaProps{
												"whisk_admin": {Type: "string"},
											},
										},
										"bucket": {Type: "string"},
										"openwhisk": {
											Type: "object",
											Properties: map[string]apiextensions.JSONSchemaProps{
												"whisk.system": {Type: "string"},
												"nuvolaris":    {Type: "string"},
											},
										},
									},
								},
								"status": {
									Type:                   "object",
									XPreserveUnknownFields: &preserveUnknownFields,
								},
							},
						},
					},
					AdditionalPrinterColumns: []apiextensions.CustomResourceColumnDefinition{
						{
							Name:        "Debug",
							Type:        "string",
							Priority:    0,
							JSONPath:    ".spec.debug",
							Description: "Debugging enabled",
						},
						{
							Name:        "Message",
							Type:        "string",
							Priority:    0,
							JSONPath:    ".status.whisk_create.message",
							Description: "As returned from the handler (sometimes)",
						},
					},
				},
			},
		},
	}
	return &whiskCrd
}

func (c *KubeClient) deployCRD() error {
	_, err := c.apiextclientset.ApiextensionsV1().CustomResourceDefinitions().Get(c.ctx, FullCRDName, metaV1.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			crd := configureCRD()
			_, err := c.apiextclientset.ApiextensionsV1().CustomResourceDefinitions().Create(c.ctx, crd, metaV1.CreateOptions{})
			if err != nil {
				return err
			}
			fmt.Println("✓ Custom resource definition for openwhisk created")
			return nil
		}
		return err
	}
	fmt.Println("custom resource definition for whisk already exists...skipping")
	return nil
}

type WhiskSpec struct {
	Debug   bool      `json:"debug"`
	Couchdb WskAdmin  `json:"couchdb"`
	Mongodb WskAdmin  `json:"mongodb"`
	Bucket  string    `json:"bucket"`
	OW      OpenWhisk `json:"openwhisk"`
}

type WskAdmin struct {
	WhiskAdmin string `json:"whisk_admin"`
}
type OpenWhisk struct {
	WhiskSystem string `json:"whisk.system"`
	Nuvolaris   string `json:"nuvolaris"`
}

type Whisk struct {
	metaV1.TypeMeta   `json:",inline"`
	metaV1.ObjectMeta `json:"metadata"`
	Spec              WhiskSpec `json:"spec"`
}

type WhiskList struct {
	metaV1.TypeMeta `json:",inline"`
	metaV1.ListMeta `json:"metadata,omitempty"`

	Items []Whisk `json:"items"`
}

func (in *Whisk) DeepCopyInto(out *Whisk) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = WhiskSpec{
		Debug:   in.Spec.Debug,
		Couchdb: in.Spec.Couchdb,
		Mongodb: in.Spec.Mongodb,
		Bucket:  in.Spec.Bucket,
		OW:      in.Spec.OW,
	}

}

func (in *Whisk) DeepCopy() *Whisk {
	if in == nil {
		return nil
	}
	out := new(Whisk)
	in.DeepCopyInto(out)
	return out
}

func (in *Whisk) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *WhiskList) DeepCopyObject() runtime.Object {
	out := WhiskList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]Whisk, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}

	return &out
}

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

func createWhiskOperatorObject(c *KubeClient) error {
	whisk := &Whisk{
		TypeMeta: metaV1.TypeMeta{
			Kind:       CRDKind,
			APIVersion: apiVersion,
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      wskObjectName,
			Namespace: namespace,
		},
		Spec: WhiskSpec{
			Debug: false,
			Couchdb: WskAdmin{
				WhiskAdmin: GenerateRandomSeq(alphanum, 8),
			},
			Mongodb: WskAdmin{
				WhiskAdmin: GenerateRandomSeq(alphanum, 8),
			},
			Bucket: awsKeygen(),
			OW: OpenWhisk{
				WhiskSystem: keygen(alphanum, 32),
				Nuvolaris:   keygen(alphanum, 32),
			},
		},
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
			fmt.Println("✓ Openwhisk operator created")
			return getWhisk(client)
		}
		return err
	}
	fmt.Println("openwhisk operator already exists...skipping")
	return nil
}
