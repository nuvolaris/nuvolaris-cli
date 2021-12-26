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
package commands

import (
	"fmt"
	"net/http"
	"os"

	"github.com/apache/openwhisk-client-go/whisk"
)

type WskCmd struct {
}

func (wsk *WskCmd) Run() error {
	runWskApiInteractionSample()
	return nil
}

func runWskApiInteractionSample() {
	fmt.Println("Doing something with wsk...")

	client, err := whisk.NewClient(http.DefaultClient, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	options := &whisk.ActionListOptions{
		Limit: 10,
		Skip:  0,
	}

	fmt.Printf("Retrieving actions list...  \n")

	actions, resp, err := client.Actions.List("", options)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Println("Returned with status: ", resp.Status)
	fmt.Printf("Returned actions: \n %+v", actions)
}