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

	"github.com/apache/openwhisk-cli/commands"
	"github.com/apache/openwhisk-cli/wski18n"
	"github.com/apache/openwhisk-client-go/whisk"
	goi18n "github.com/nicksnyder/go-i18n/i18n"
)

type WskCmd struct {
	Args []string `arg:"" name:"args"`
}

var cliDebug = os.Getenv("WSK_CLI_DEBUG") // Useful for tracing init() code

var T goi18n.TranslateFunc

func init() {
	if len(cliDebug) > 0 {
		whisk.SetDebug(true)
	}

	T = wski18n.T

	// Rest of CLI uses the Properties struct, so set the build time there
	commands.Properties.CLIVersion = CLIVersion
}

func (w *WskCmd) BeforeApply() error {
	_, ok := os.LookupEnv("WSK_CONFIG_FILE")
	path, err := getWhiskPropsPath()
	if err != nil {
		return err
	}
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf(".wskprops file not found. Run nuv setup")
	}

	if !ok {
		os.Setenv("WSK_CONFIG_FILE", path)
	}
	return nil
}

func getWhiskPropsPath() (string, error) {
	path, err := GetOrCreateNuvolarisConfigDir()
	if err != nil {
		return "", err
	}
	wpath := filepath.Join(path, ".wskprops")
	return wpath, nil
}

func Wsk(args ...string) error {
	os.Args = append([]string{"wsk"}, args...)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			fmt.Println(T("Application exited unexpectedly"))
		}
	}()
	return commands.Execute()
}

func (wsk *WskCmd) Run() error {
	return Wsk(wsk.Args...)
}

const wskAuth = "23bc46b1-71f6-4ed5-8c54-816aa4f8c502:123zO3xZCLrMN6v2BKK1dXYFpXlPkccOFqm12CdAsMgRU4VrNZ9lyGVCGuMDGIwP"

func writeWskPropertiesFile(apihost string) error {
	content := []byte("AUTH=" + wskAuth + "\nAPIHOST=" + apihost)
	path, err := WriteFileToNuvolarisConfigDir(".wskprops", content)
	if err != nil {
		return err
	}
	os.Setenv("WSK_CONFIG_FILE", path)
	return nil
}

func getWhiskPropsPath() (string, error) {
	path, err := GetOrCreateNuvolarisConfigDir()
	if err != nil {
		return "", err
	}
	wpath := filepath.Join(path, ".wskprops")
	return wpath, nil
}
