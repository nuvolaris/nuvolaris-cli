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

func getWhiskSpec() (WhiskSpec, error) {
	config, err := ReadFileFromNuvolarisConfigDir("config.yaml")

	if err != nil {
		return WhiskSpec{}, err
	}
	var result WhiskSpec
	err = yaml.Unmarshal(config, &result)
	return result, err
}

func setupWskProps(cmd *AuthCmd) error {
	var auth, apihost, redisurl, mongodburl, mdb_username, mdb_password string
	result, err := getWhiskSpec()
	var errMessage = "nuvolaris setup config file not found. Please setup nuvolaris or specify both --apihost and --auth"

	if cmd.Apihost == "" || cmd.Auth == "" {
		if err != nil {
			return fmt.Errorf(errMessage)
		}

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

	// ADD REDIS URL
	if cmd.Redis == "" && result.Components.Redis {
		redisurl = "redis://redis"
	} else {
		redisurl = cmd.Redis
	}

	if redisurl != "" {
		fmt.Println("Adding REDIS_URI to whisk properties")
		writeWskPropsFile(wskPropsKeyValue{
			wskPropsKey:   "REDIS_URI",
			wskPropsValue: redisurl,
		})
	}

	// ADD MONGODB IF IT IS THE CASE
	// TODO if enabled the url could be retrieved issuing a kubectl command similar to kubectl get secret -n nuvolaris -ojson nuvolaris-mongodb-nuvolaris-nuvolaris
	if cmd.Mongodb == "" && result.Components.MongoDb {
		mdb_username = result.MongoDb.Nuvolaris.User
		mdb_password = result.MongoDb.Nuvolaris.Password
		mongodburl = "mongodb://" + mdb_username + ":" + mdb_password + "@nuvolaris-mongodb-0.nuvolaris-mongodb-svc.nuvolaris.svc.cluster.local:27017/nuvolaris?replicaSet=nuvolaris-mongodb&ssl=false"
	} else {
		mongodburl = cmd.Mongodb
	}

	if mongodburl != "" {
		fmt.Println("Adding MONGODB_URI to whisk properties")
		writeWskPropsFile(wskPropsKeyValue{
			wskPropsKey:   "MONGODB_URI",
			wskPropsValue: mongodburl,
		})
	}

	if cmd.Show {
		fmt.Printf("nuv auth --apihost %s --auth %s\n", apihost, auth)
	} else {
		fmt.Printf("Autentication ready.\n")
	}
	return nil
}
