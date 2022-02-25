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
)

func writeWskPropertiesFile() error {
	content := []byte("AUTH=" + wsk_auth + "\nAPIHOST=" + wsk_apihost)
	path, err := writeFileToNuvolarisHomedir(".wskprops", content)
	if err != nil {
		return err
	}
	fmt.Println("✓ .wskprops file written")
	os.Setenv("WSK_CONFIG_FILE", path)
	fmt.Println("✓ WSK_CONFIG_FILE env variable set to " + os.Getenv("WSK_CONFIG_FILE"))
	return nil
}

func writeFileToNuvolarisHomedir(filename string, content []byte) (string, error) {
	homedir, err := GetHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(homedir, ".nuvolaris", filename)
	if _, err := os.Stat(path); err == nil {
		os.Remove(path)
	}

	if err := os.WriteFile(path, content, 0600); err != nil {
		return "", err
	}
	return path, nil
}

func getWskPropsFilePath() (string, error) {
	homedir, err := GetHomeDir()
	if err != nil {
		return "", err
	}
	var path = filepath.Join(homedir, ".nuvolaris", ".wskprops")
	fileDoesNotExistException := checkFileExists(err, path)
	return path, fileDoesNotExistException
}

func checkFileExists(err error, path string) error {
	_, err = os.Stat(path)
	ex := fmt.Errorf("file: %s does not exist", path)
	fmt.Println(ex)
	return ex
}
