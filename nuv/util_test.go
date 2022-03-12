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
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleSysErr() {
	SysErr("/bin/echo 1 2 3")
	SysErr("/bin/echo 3", "4", "5")
	SysErr("@sh -c", "echo foo >/tmp/foo")
	outFoo, _ := SysErr("cat /tmp/foo")
	fmt.Print(outFoo)
	SysErr("@sh -c", "echo bar >/tmp/bar")
	outBar, _ := SysErr("@cat /tmp/bar")
	fmt.Print(outBar)
	_, err := SysErr("false")
	fmt.Println("ERR", err)
	_, err = SysErr("donotexist")
	fmt.Println("ERR", err)
	// Output:
	// 1 2 3
	// 3 4 5
	// foo
	// foo
	// bar
	// ERR exit status 1
	// ERR exec: "donotexist": executable file not found in $PATH
}

func ExampleDryRunSysErr() {
	DryRunPush("first", "second", "!third")
	out, err := DryRunSysErr("dummy")
	fmt.Println(1, out, err)
	out, err = DryRunSysErr("dummy", "alpha", "beta")
	fmt.Println(2, out, err)
	out, err = DryRunSysErr("dummy")
	fmt.Println(3, "out", out, "err", err)
	out, err = DryRunSysErr("dummy")
	fmt.Println(4, "out", out, "err", err)
	// Output:
	// dummy
	// 1 first <nil>
	// dummy alpha beta
	// 2 second <nil>
	// dummy
	// 3 out  err third
	// dummy
	// 4 out  err <nil>
}

func TestV4UUIDFormat(t *testing.T) {
	uuid := GenerateUUID()
	assert.Equal(t, len(uuid), 36, "")
	r := regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$")
	match, err := regexp.MatchString(r.String(), uuid)
	assert.True(t, match)
	assert.Nil(t, err)
}

func TestGenerateRandSequence(t *testing.T) {
	random := GenerateRandomSeq(alphanum, 32)
	assert.Equal(t, len(random), 32)
	assert.NotContains(t, random, " ")
	random = GenerateRandomSeq(alphanum, 64)
	assert.Equal(t, len(random), 64)
	assert.NotContains(t, random, "@")
}

func TestKeygen(t *testing.T) {
	key := keygen(alphanum, 32)
	assert.Equal(t, len(key), 69)
	assert.Contains(t, key, ":")
}

func TestAwsAccessKeyId(t *testing.T) {
	keyId := generateAwsAccessKeyId()
	assert.Equal(t, len(keyId), 20)
	assert.NotContains(t, keyId, ":")
	assert.NotContains(t, keyId, "9")
	assert.True(t, strings.HasPrefix(keyId, "AKIA"))
}

func TestAwsSecretAccessKey(t *testing.T) {
	keyId := generateAwsSecretAccessKey()
	assert.Equal(t, len(keyId), 40)
	assert.NotContains(t, keyId, ":")
	assert.NotContains(t, keyId, "@")
}
