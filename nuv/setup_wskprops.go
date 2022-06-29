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
	"strings"

	"sigs.k8s.io/yaml"
)

func setupWskProps(cmd *AuthCmd) error {
	var auth, apihost string

	var errMessage = "nuvolaris setup config file not found. Please setup nuvolaris or specify both --apihost and --auth"

	if cmd.Apihost == "" || cmd.Auth == "" {
		config, err := ReadFileFromNuvolarisConfigDir("config.yaml")
		if err != nil {
			return fmt.Errorf(errMessage)
		}
		var result WhiskSpec
		yaml.Unmarshal(config, &result)
		auth = result.OpenWhisk.Namespaces.Nuvolaris
		if result.Nuvolaris == nil {
			result.Nuvolaris = &NuvolarisS{}
		}
		apihost = result.Nuvolaris.ApiHost
		if auth == "" || apihost == "" {
			return fmt.Errorf(errMessage)
		}
	} else {
		auth = cmd.Auth
		apihost = cmd.Apihost
	}
	writeWskPropsFile(wskPropsKeyValue{
		wskPropsKey:   "AUTH",
		wskPropsValue: auth,
	})

	// TODO: temporary workaround adding http: all the time
	if !strings.HasPrefix(apihost, "http://") {
		apihost = "http://" + apihost
	}
	apihost = strings.ReplaceAll(apihost, "https:", "http:")
	writeWskPropsFile(wskPropsKeyValue{
		wskPropsKey:   "APIHOST",
		wskPropsValue: apihost,
	})
	if cmd.Show {
		fmt.Printf("Configure authentication:\nnuv auth --apihost %s --auth %s\n", apihost, auth)
	} else {
		fmt.Printf("Autentication ready.\n")
	}
	return nil
}
