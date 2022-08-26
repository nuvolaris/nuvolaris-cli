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
	"os"
	"path/filepath"
	"strings"
)

func writeWskPropsFile(keyValues ...wskPropsKeyValue) error {
	currentContent, err := readWskPropsAsMap()
	if err != nil {
		return err
	}
	for _, keyValue := range keyValues {
		currentContent[keyValue.wskPropsKey] = keyValue.wskPropsValue
	}

	var sb strings.Builder
	for k, v := range currentContent {
		sb.WriteString(k + "=" + v + "\n")
	}
	_, err = WriteFileToNuvolarisConfigDir(WskPropsFilename, []byte(sb.String()))
	if err != nil {
		return err
	}
	return setWskEnvVariable(true)
}

func setWskEnvVariable(check bool) error {
	_, ok := os.LookupEnv("WSK_CONFIG_FILE")
	if !ok {
		path, err := getWhiskPropsPath()
		if err != nil {
			return err
		}
		if check {
			_, err = os.Stat(path)
			if os.IsNotExist(err) {
				return fmt.Errorf(".wskprops file not found. Run nuv setup")
			}
		}
		os.Setenv("WSK_CONFIG_FILE", path)
	}
	return nil
}

func getWhiskPropsPath() (string, error) {
	path, err := GetOrCreateNuvolarisConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(path, WskPropsFilename), nil
}

func getOrCreateWhiskPropsFile() ([]byte, error) {
	path, err := getWhiskPropsPath()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(path); err == nil {
		return os.ReadFile(path)
	}
	return nil, os.WriteFile(path, nil, 0600)
}

func readWskPropsAsMap() (map[string]string, error) {
	content, err := getOrCreateWhiskPropsFile()
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	if content != nil {
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			if strings.Contains(line, "=") {
				keyValue := strings.Split(line, "=")
				m[keyValue[0]] = keyValue[1]
			}
		}
	}
	return m, nil
}

func flattenWskPropsMap(annotations map[string]string) []wskPropsKeyValue {
	var wskPropsEntries []wskPropsKeyValue
	for k, v := range annotations {
		wskPropsEntries = append(wskPropsEntries, wskPropsKeyValue{
			wskPropsKey:   k,
			wskPropsValue: v,
		})
	}
	return wskPropsEntries
}

func setWskPropsAsEnvVariable() error {
	propMap, err := readWskPropsAsMap()
	if err != nil {
		fmt.Errorf(err.Error())
		return err
	}

	for k, v := range propMap {
		os.Setenv(k, v)
	}

	return nil
}

type wskPropsKeyValue struct {
	wskPropsKey   string
	wskPropsValue string
}
