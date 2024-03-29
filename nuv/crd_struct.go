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
	"io/ioutil"
	"os"
	"path/filepath"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

//go:embed embed/config.yaml
var configYaml []byte

type WhiskSpec struct {
	Components ComponentsS `json:"components" yaml:"components"`
	OpenWhisk  OpenWhiskS  `json:"openwhisk" yaml:"openwhisk"`
	Nuvolaris  *NuvolarisS `json:"nuvolaris,omitempty" yaml:"nuvolaris, omitempty"`
	CouchDb    CouchDbS    `json:"couchdb" yaml:"couchdb"`
	MongoDb    MongoDbS    `json:"mongodb" yaml:"mongodb"`
	Kafka      KafkaS      `json:"kafka" yaml:"kafka"`
	S3         S3S         `json:"s3" yaml:"s3"`
	Scheduler  SchedulerS  `json:"scheduler,omitempty" yaml:"scheduler, omitempty"`
}

type ComponentsS struct {
	Openwhisk bool `json:"openwhisk" yaml:"openwhisk"`
	Invoker   bool `json:"invoker" yaml:"invoker"`
	CouchDb   bool `json:"couchdb" yaml:"couchdb"`
	Kafka     bool `json:"kafka" yaml:"kafka"`
	MongoDb   bool `json:"mongodb" yaml:"mongodb"`
	Redis     bool `json:"redis" yaml:"redis"`
	Cron      bool `json:"cron" yaml:"cron"`
	S3Bucket  bool `json:"s3bucket" yaml:"redis"`
}

type OpenWhiskS struct {
	Namespaces NamespacesS `json:"namespaces" yaml:"namespaces"`
	Limits     LimitsS     `json:"limits, omitempty" yaml:"limits, omitempty"`
}

type NamespacesS struct {
	WhiskSystem string `json:"whisk-system" yaml:"whisk-system"`
	Nuvolaris   string `json:"nuvolaris" yaml:"nuvolaris"`
}

type LimitsS struct {
	LimitActions  LimitActionsS  `json:"actions" yaml:"actions"`
	LimitTriggers LimitTriggersS `json:"triggers" yaml:"triggers"`
}

type LimitActionsS struct {
	SequenceMaxLength string `json:"sequence-maxLength" yaml:"sequence-maxLength"`
	InvokesPerMinute  string `json:"invokes-perMinute" yaml:"invokes-perMinute"`
	InvokesConcurrent string `json:"invokes-concurrent" yaml:"invokes-concurrent"`
}

type LimitTriggersS struct {
	FiresPerMinute string `json:"fires-perMinute" yaml:"fires-perMinute"`
}

type CouchDbS struct {
	//Host       string `json:"host" yaml:"host"`
	VolumeSize int    `json:"volume-size" yaml:"volume-size"`
	Admin      AdminS `json:"admin" yaml:"admin"`
	Controller AdminS `json:"controller" yaml:"controller"`
	Invoker    AdminS `json:"invoker" yaml:"invoker"`
}

type AdminS struct {
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
}

type NuvolarisS struct {
	ApiHost string `json:"apihost" yaml:"apihost"`
}

type MongoDbS struct {
	Host       string `json:"host" yaml:"host"`
	VolumeSize int    `json:"volume-size" yaml:"volume-size"`
	Admin      AdminS `json:"admin" yaml:"admin"`
	Nuvolaris  AdminS `json:"nuvolaris" yaml:"nuvolaris"`
}

type KafkaS struct {
	Host       string `json:"host" yaml:"host"`
	VolumeSize int    `json:"volume-size" yaml:"volume-size"`
}

type S3S struct {
	VolumeSize int    `json:"volume-size" yaml:"volume-size"`
	Id         string `json:"id" yaml:"id"`
	Key        string `json:"key" yaml:"key"`
	Region     string `json:"region" yaml:"region"`
}

type SchedulerS struct {
	Schedule string `json:"schedule" yaml:"schedule"`
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
		Nuvolaris:  in.Spec.Nuvolaris,
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

func configureCrd(apiHost string) error {
	var result WhiskSpec
	yaml.Unmarshal(configYaml, &result)
	result.OpenWhisk.Namespaces.WhiskSystem = keygen(alphanum, 64)
	result.OpenWhisk.Namespaces.Nuvolaris = keygen(alphanum, 64)
	if apiHost == "auto" {
		result.Nuvolaris = nil
	} else {
		result.Nuvolaris.ApiHost = apiHost
	}
	result.CouchDb.Admin.Password = GenerateRandomSeq(alphanum, 8)
	result.CouchDb.Controller.Password = GenerateRandomSeq(alphanum, 8)
	result.CouchDb.Invoker.Password = GenerateRandomSeq(alphanum, 8)
	result.MongoDb.Admin.Password = GenerateRandomSeq(alphanum, 8)
	result.S3.Id = generateAwsAccessKeyId()
	result.S3.Key = generateAwsSecretAccessKey()
	content, _ := yaml.Marshal(result)
	_, err := WriteFileToNuvolarisConfigDir("config.yaml", content)
	nuvolarisHome, err := GetOrCreateNuvolarisConfigDir()
	if err != nil {
		return err
	}
	path := filepath.Join(nuvolarisHome, "config.yaml")
	fmt.Println("Nuvolaris configuration written to", path)
	fmt.Println("please edit this configuration file if you need to change parameters")
	return err
}

func updateApihostInConfig(apiHost string) error {
	var result WhiskSpec
	content, err := ReadFileFromNuvolarisConfigDir("config.yaml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(content, &result)
	if err != nil {
		return err
	}
	if result.Nuvolaris == nil {
		result.Nuvolaris = &NuvolarisS{}
	}
	result.Nuvolaris.ApiHost = apiHost
	content, err = yaml.Marshal(result)
	if err != nil {
		return err
	}
	_, err = WriteFileToNuvolarisConfigDir("config.yaml", content)
	return err
}

func readOrCreateCrdConfig(apiHost string) (*WhiskSpec, error) {
	var result WhiskSpec
	nuvHomedir, err := GetOrCreateNuvolarisConfigDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(nuvHomedir, "config.yaml")
	if _, err := os.Stat(path); err != nil {
		err = configureCrd(apiHost)
		if err != nil {
			return nil, err
		}
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(content, &result)
	return &result, err
}
