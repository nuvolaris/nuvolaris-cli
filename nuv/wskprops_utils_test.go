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
	"os"
	"testing"
)

func TestWriteWskPropsFile(t *testing.T) {
	content := []byte("APIHOST=http://localhost:3232\nKEY_TO_REPLACE=valuetoreplace")
	WriteFileToNuvolarisConfigDir(".wskprops", content)
	keyValuePairToAdd := wskPropsKeyValue{
		wskPropsKey:   "NEW_KEY",
		wskPropsValue: "somevalue",
	}
	err := writeWskPropsFile(keyValuePairToAdd)
	if err != nil {
		t.Errorf(err.Error())
	}
	propsMap, err := readWskPropsAsMap()
	if err != nil {
		t.Errorf(err.Error())
	}
	assert.Equal(t, propsMap["APIHOST"], "http://localhost:3232")
	assert.Equal(t, propsMap["KEY_TO_REPLACE"], "valuetoreplace")
	assert.Equal(t, propsMap["NEW_KEY"], "somevalue")
	assert.Equal(t, len(propsMap), 3)

	keyValuePairToReplace := wskPropsKeyValue{
		wskPropsKey:   "KEY_TO_REPLACE",
		wskPropsValue: "replacedvalue",
	}
	err = writeWskPropsFile(keyValuePairToReplace)
	propsMap, err = readWskPropsAsMap()
	if err != nil {
		t.Errorf(err.Error())
	}
	assert.Equal(t, propsMap["KEY_TO_REPLACE"], "replacedvalue")
	assert.Equal(t, len(propsMap), 3)

	wskPropsPath, err := getWhiskPropsPath()
	if err != nil {
		t.Errorf(err.Error())
	}
	if _, err := os.Stat(wskPropsPath); err == nil {
		os.Remove(wskPropsPath)
	}
	err = writeWskPropsFile(keyValuePairToAdd)
	if err != nil {
		t.Errorf(err.Error())
	}
	propsMap, err = readWskPropsAsMap()
	assert.Equal(t, propsMap["NEW_KEY"], "somevalue")
	assert.Equal(t, len(propsMap), 1)

}

func TestGetOrCreateWhiskPropsFile(t *testing.T) {
	content := []byte("APIHOST=http://localhost:3232\n")
	WriteFileToNuvolarisConfigDir(".wskprops", content)
	readContent, err := getOrCreateWhiskPropsFile()
	if err != nil {
		t.Errorf(err.Error())
	}
	assert.Equal(t, readContent, content)
	wskPropsPath, err := getWhiskPropsPath()
	if err != nil {
		t.Errorf(err.Error())
	}
	if _, err := os.Stat(wskPropsPath); err == nil {
		os.Remove(wskPropsPath)
	}
	readContent, err = getOrCreateWhiskPropsFile()
	assert.Equal(t, readContent, []byte(nil))

}
func TestReadWskPropsAsMap(t *testing.T) {
	WriteFileToNuvolarisConfigDir(".wskprops", []byte("APIHOST=http://localhost:3232\nAUTH=23bc46b1-71f6-4ed5-8c54-816aa4f8c502:123zO3xZCLrMN6v2BKK1dXYFpXlPkccOFqm12CdAsMgRU4VrNZ9lyGVCGuMDGIwP\n"))
	keyvalues, err := readWskPropsAsMap()
	if err != nil {
		t.Errorf(err.Error())
	}
	assert.Equal(t, keyvalues["APIHOST"], "http://localhost:3232")
	assert.Equal(t, keyvalues["AUTH"], "23bc46b1-71f6-4ed5-8c54-816aa4f8c502:123zO3xZCLrMN6v2BKK1dXYFpXlPkccOFqm12CdAsMgRU4VrNZ9lyGVCGuMDGIwP")
	assert.Equal(t, len(keyvalues), 2)
}

func TestWskPropsAsEnvVar(t *testing.T) {
	WriteFileToNuvolarisConfigDir(".wskprops", []byte("APIHOST=http://localhost:3232\nAUTH=23bc46b1-71f6-4ed5-8c54-816aa4f8c502:123zO3xZCLrMN6v2BKK1dXYFpXlPkccOFqm12CdAsMgRU4VrNZ9lyGVCGuMDGIwP\n"))
	err := setWskPropsAsEnvVariable()
	if err != nil {
		t.Errorf(err.Error())
	}

	assert.Equal(t, os.Getenv("APIHOST"), "http://localhost:3232")
	assert.Equal(t, os.Getenv("AUTH"), "23bc46b1-71f6-4ed5-8c54-816aa4f8c502:123zO3xZCLrMN6v2BKK1dXYFpXlPkccOFqm12CdAsMgRU4VrNZ9lyGVCGuMDGIwP")
}

func TestFlattenWskPropsMap(t *testing.T) {
	annotations := make(map[string]string)
	annotations["APIHOST"] = "somehost"
	annotations["AUTH"] = "someauth"
	annotations["MONGODB_PASS"] = "somepass"
	wskProps := flattenWskPropsMap(annotations)
	assert.Contains(t, wskProps, wskPropsKeyValue{
		wskPropsKey:   "APIHOST",
		wskPropsValue: "somehost",
	})
	assert.Contains(t, wskProps, wskPropsKeyValue{
		wskPropsKey:   "AUTH",
		wskPropsValue: "someauth",
	})
	assert.Contains(t, wskProps, wskPropsKeyValue{
		wskPropsKey:   "MONGODB_PASS",
		wskPropsValue: "somepass",
	})
}
