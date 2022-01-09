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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHomedir(t *testing.T) {
	realHomeDirFunc := homeDirFunc
	defer func() {
		homeDirFunc = realHomeDirFunc
	}()
	homeDirFunc = func() (string, error) {
		return "/home/userdir", nil
	}

	home, err := GetHomedir()
	assert.Equal(t, home, "/home/userdir", "")
	assert.Equal(t, err, nil, "")

	homeDirFunc = func() (string, error) {
		return "", fmt.Errorf("some error returned from homedir")
	}

	home, err = GetHomedir()
	assert.Equal(t, home, "", "")
	assert.Equal(t, err.Error(), "some error returned from homedir", "")

}
