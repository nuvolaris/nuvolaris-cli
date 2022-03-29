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
)

func TestConfigureCRD(t *testing.T) {
	whiskCrd := configureCRD()
	assert.Equal(t, whiskCrd.Name, "whisks.nuvolaris.org")
	assert.Equal(t, whiskCrd.Namespace, "nuvolaris")
	assert.Equal(t, whiskCrd.Spec.Names.Kind, "Whisk")
	assert.Equal(t, whiskCrd.Spec.Names.Singular, "whisk")
	assert.Equal(t, whiskCrd.Spec.Names.Plural, "whisks")
	assert.Equal(t, whiskCrd.Spec.Names.ShortNames, []string{"wsk"})
	assert.Equal(t, whiskCrd.Spec.Group, "nuvolaris.org")

	openApiSpecProperties := whiskCrd.Spec.Versions[0].Schema.OpenAPIV3Schema.Properties["spec"].Properties
	assert.NotEmpty(t, openApiSpecProperties["debug"])
	assert.NotEmpty(t, openApiSpecProperties["couchdb"])
	assert.NotEmpty(t, openApiSpecProperties["couchdb"].Properties["whisk_admin"])
	assert.NotEmpty(t, openApiSpecProperties["mongodb"])
	assert.NotEmpty(t, openApiSpecProperties["mongodb"].Properties["whisk_admin"])
	assert.NotEmpty(t, openApiSpecProperties["bucket"])
	assert.NotEmpty(t, openApiSpecProperties["openwhisk"])
	assert.NotEmpty(t, openApiSpecProperties["openwhisk"].Properties["whisk.system"])
	assert.NotEmpty(t, openApiSpecProperties["openwhisk"].Properties["nuvolaris"])
}
