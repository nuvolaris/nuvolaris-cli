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
	"testing"

	"github.com/stretchr/testify/assert"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCreateServiceAccount(t *testing.T) {
	testclient.clientset = fake.NewSimpleClientset(nspace)
	err := testclient.createServiceAccount()
	if err != nil {
		t.Errorf(err.Error())
	}
	account, _ := testclient.clientset.CoreV1().ServiceAccounts("nuvolaris").Get(testclient.ctx, "nuvolaris-operator", meta_v1.GetOptions{})
	assert.Equal(t, account.Name, "nuvolaris-operator")
	assert.Equal(t, account.Namespace, "nuvolaris")
	assert.Equal(t, account.Labels, map[string]string{"app": "nuvolaris-operator"})
}

func TestSkipCreateServiceAccountIfAlreadyExists(t *testing.T) {
	testclient.clientset = fake.NewSimpleClientset(nspace, service_account)
	err := testclient.createServiceAccount()
	if err != nil {
		t.Errorf(err.Error())
	}
}

func Example_createServiceAccount() {
	testclient.clientset = fake.NewSimpleClientset(nspace)
	testclient.createServiceAccount()
	// Output:
	// ✓ Service account created
}

func Example_skipCreationServiceAccount() {
	testclient.clientset = fake.NewSimpleClientset(nspace, service_account)
	testclient.createServiceAccount()
	// Output:
	// service account already created...skipping
}

func TestSkipOperatorPodIfRunning(t *testing.T) {
	testclient.clientset = fake.NewSimpleClientset(nspace, operator_pod)
	err := testclient.createOperatorPod()
	if err != nil {
		t.Errorf(err.Error())
	}
}

func Example_skipOperatorPod() {
	testclient.clientset = fake.NewSimpleClientset(nspace, operator_pod)
	testclient.createOperatorPod()
	// Output:
	// nuvolaris operator pod already running...skipping
}
