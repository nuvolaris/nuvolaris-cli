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

type WhiskSpec struct {
	Components ComponentsS `json:"components"`
	OpenWhisk  OpenWhiskS  `json:"openwhisk"`
	CouchDb    CouchDbS    `json:"couchdb"`
	MongoDb    MongoDbS    `json:"mongodb"`
	Kafka      KafkaS      `json:"kafka"`
	S3         S3S         `json:"s3"`
}

type ComponentsS struct {
	Openwhisk bool `json:"openwhisk"`
	Invoker   bool `json:"invoker"`
	CouchDb   bool `json:"couchdb"`
	Kafka     bool `json:"kafka"`
	MongoDb   bool `json:"mongodb"`
	Redis     bool `json:"redis"`
	S3Bucket  bool `json:"s3bucket"`
}

type OpenWhiskS struct {
	Namespaces NamespacesS `json:"namespaces"`
}

type NamespacesS struct {
	WhiskSystem string `json:"whisk-system"`
	Nuvolaris   string `json:"nuvolaris"`
}

type CouchDbS struct {
	Host       string `json:"host"`
	VolumeSize int    `json:"volume-size"`
	Admin      AdminS `json:"admin"`
	Controller AdminS `json:"controller"`
	Invoker    AdminS `json:"invoker"`
}

type AdminS struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type MongoDbS struct {
	Host       string `json:"host"`
	VolumeSize int    `json:"volume-size"`
	Admin      AdminS `json:"admin"`
}

type KafkaS struct {
	Host       string `json:"host"`
	VolumeSize int    `json:"volume-size"`
}

type S3S struct {
	VolumeSize int    `json:"volume-size"`
	Id         string `json:"id"`
	Key        string `json:"key"`
	Region     string `json:"region"`
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
		Components: in.Spec.Components,
		OpenWhisk:  in.Spec.OpenWhisk,
		CouchDb:    in.Spec.CouchDb,
		MongoDb:    in.Spec.MongoDb,
		Kafka:      in.Spec.Kafka,
		S3:         in.Spec.S3,
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
	authKey := keygen(alphanum, 64)
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
			Components: ComponentsS{
				Openwhisk: true,
				Invoker:   false,
				CouchDb:   true,
				Kafka:     false,
				MongoDb:   false,
				Redis:     false,
				S3Bucket:  false,
			},
			OpenWhisk: OpenWhiskS{
				Namespaces: NamespacesS{
					WhiskSystem: keygen(alphanum, 64),
					Nuvolaris:   authKey,
				},
			},
			CouchDb: CouchDbS{
				Host:       "couchdb",
				VolumeSize: 10,
				Admin: AdminS{
					User:     "whisk_admin",
					Password: GenerateRandomSeq(alphanum, 8),
				},
				Controller: AdminS{
					User:     "invoker_admin",
					Password: GenerateRandomSeq(alphanum, 8),
				},
				Invoker: AdminS{
					User:     "controller_admin",
					Password: GenerateRandomSeq(alphanum, 8),
				},
			},
			MongoDb: MongoDbS{
				Host:       "mongodb",
				VolumeSize: 10,
				Admin: AdminS{
					User:     "admin",
					Password: GenerateRandomSeq(alphanum, 8),
				},
			},
			Kafka: KafkaS{
				Host:       "kafka",
				VolumeSize: 10,
			},
			S3: S3S{
				VolumeSize: 10,
				Id:         generateAwsAccessKeyId(),
				Key:        generateAwsSecretAccessKey(),
				Region:     "eu-central-1",
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
			fmt.Println("✓ Openwhisk operator started")
			//TODO remove temporary workaround
			writeWskPropsFile(wskPropsKeyValue{
				wskPropsKey:   "AUTH",
				wskPropsValue: authKey,
			})
			return nil
		}
		return err
	}
	fmt.Println("openwhisk operator already running...skipping")
	return nil
}
