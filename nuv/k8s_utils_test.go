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
	"github.com/stretchr/testify/assert"
	fakeclient "k8s.io/client-go/kubernetes/fake"
	"testing"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var configmap = &coreV1.ConfigMap{
	ObjectMeta: metaV1.ObjectMeta{
		Name:      "config",
		Namespace: "nuvolaris",
		Annotations: map[string]string{
			"nuvolaris-apihost":        "https://localhost:3232",
			"nuvolaris-auth":           "23bc46b1-71f6-4ed5-8c54-816aa4f8c502:123zO3xZCLrMN6v2BKK1dXYFpXlPkccOFqm12CdAsMgRU4VrNZ9lyGVCGuMDGIwP",
			"nuvolaris-mongo-db-pass":  "338fe176-4856-4b21-adae-1fe7a8a9a4c9:123zO3xZCLrMN6v2BKK1dXYFpXlPkccOFqm12CdAsMgRU4VrNZ9lyGVCGuMDGIwP",
			"nuvolaris-couchdb":        "19fb8b3b-8a34-403b-a120-4205a7749e97:123zO3xZCLrMN6v2BKK1dXYFpXlPkccOFqm12CdAsMgRU4VrNZ9lyGVCGuMDGIwP",
			"not-nuvolaris-annotation": "something",
		},
	},
}

func TestReadClusterConfig(t *testing.T) {
	testclient.clientset = fakeclient.NewSimpleClientset(nspace, configmap)
	_, err := testclient.clientset.CoreV1().ConfigMaps(testclient.namespace).Get(testclient.ctx, "config", metaV1.GetOptions{})
	if err != nil {
		t.Errorf(err.Error())
	}
	annotations, err := readClusterConfig(&testclient, "config")
	if err != nil {
		t.Errorf(err.Error())
	}
	assert.Equal(t, annotations["APIHOST"], "http://localhost:3232")
	assert.Equal(t, annotations["AUTH"], "23bc46b1-71f6-4ed5-8c54-816aa4f8c502:123zO3xZCLrMN6v2BKK1dXYFpXlPkccOFqm12CdAsMgRU4VrNZ9lyGVCGuMDGIwP")
	assert.Equal(t, annotations["MONGO_DB_PASS"], "338fe176-4856-4b21-adae-1fe7a8a9a4c9:123zO3xZCLrMN6v2BKK1dXYFpXlPkccOFqm12CdAsMgRU4VrNZ9lyGVCGuMDGIwP")
	assert.Equal(t, annotations["COUCHDB"], "19fb8b3b-8a34-403b-a120-4205a7749e97:123zO3xZCLrMN6v2BKK1dXYFpXlPkccOFqm12CdAsMgRU4VrNZ9lyGVCGuMDGIwP")
	assert.Empty(t, annotations["not-nuvolaris-annotation"])
	assert.Empty(t, annotations["NOT_NUVOLARIS-ANNOTATION"])
	assert.Empty(t, annotations["NUVOLARIS-APIHOST"])
	assert.Empty(t, annotations["nuvolaris-apihost"])
}

func TestWriteConfigToWskProps(t *testing.T) {
	testclient.clientset = fakeclient.NewSimpleClientset(nspace, configmap)
	err := writeConfigToWskProps(&testclient, "config")
	if err != nil {
		t.Errorf(err.Error())
	}
	content, err := ReadFileFromNuvolarisConfigDir(WskPropsFilename)
	wskProps := string(content)
	assert.Contains(t, wskProps, "APIHOST=http://localhost:3232\n")
	assert.NotContains(t, wskProps, "nuvolaris-apihost")
	assert.Contains(t, wskProps, "AUTH=23bc46b1-71f6-4ed5-8c54-816aa4f8c502:123zO3xZCLrMN6v2BKK1dXYFpXlPkccOFqm12CdAsMgRU4VrNZ9lyGVCGuMDGIwP\n")
	assert.NotContains(t, wskProps, "nuvolaris-auth=23bc46b1-71f6-4ed5-8c54-816aa4f8c502:123zO3xZCLrMN6v2BKK1dXYFpXlPkccOFqm12CdAsMgRU4VrNZ9lyGVCGuMDGIwP")
}
