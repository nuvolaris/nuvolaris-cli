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
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
)

var dryRunBuf = []string{}

// DryRunPush saves dummy results for dry run execution
func DryRunPush(buf ...string) {
	dryRunBuf = buf
}

// DryRunPop returns a value from the buffer of dry run results
// returns an empty string if the buffer is empty
func DryRunPop(buf ...string) string {
	res := ""
	if len(dryRunBuf) > 0 {
		res = dryRunBuf[0]
		dryRunBuf = dryRunBuf[1:]
	}
	return res
}

// SysErr executes a command in a convenient way:
// it splits the paramenter in arguments if separated by spaces,
// then accepts multiple arguments;
// logs errors in stderr and prints output in stdout;
// also returns output as a string, or an error if there is an error
// If the command starts with "@" do not print the output.
func SysErr(cli string, args ...string) (string, error) {
	return sysErr(false, cli, args...)
}

// DryRunSysErr performs a dry run of SysErr
// in this case it always prints the command
func DryRunSysErr(cli string, args ...string) (string, error) {
	return sysErr(true, cli, args...)
}

func sysErr(dryRun bool, cli string, args ...string) (string, error) {
	re := regexp.MustCompile(`[\r\t\n\f ]+`)
	a := strings.Split(re.ReplaceAllString(cli, " "), " ")
	params := args
	//fmt.Println(params)
	if len(a) > 1 {
		params = append(a[1:], args...)
	}
	exe := strings.TrimPrefix(a[0], "@")
	silent := strings.HasPrefix(a[0], "@")
	if dryRun {
		if len(params) > 0 {
			fmt.Printf("%s %s\n", exe, strings.Join(params, " "))
		} else {
			fmt.Println(exe)
		}
		res := DryRunPop()
		if strings.HasPrefix(res, "!") {
			return "", errors.New(res[1:])
		}
		return res, nil
	}

	cmd := exec.Command(exe, params...)
	out, err := cmd.CombinedOutput()
	res := string(out)
	if err != nil {
		return "", err
	}
	if !silent {
		fmt.Print(res)
	}
	return res, nil
}

func ExecutingInContainer() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	fsys := os.DirFS("/")
	// if .dockerenv exists and is a regular file
	if info, err := fs.Stat(fsys, ".dockerenv"); os.IsNotExist(err) || !info.Mode().IsRegular() {
		return false
	}

	// and if docker-host.sock exists and is a socket
	if info, err := fs.Stat(fsys, "var/run/docker-host.sock"); os.IsNotExist(err) || info.Mode().Type() != fs.ModeSocket {
		return false
	}

	return true
}

func DockerHostKind() error {
	if os.Args[1] == "create" || os.Args[1] == "delete" {
		appendKubeConfig()
	}
	os.Args = append([]string{"env", "DOCKER_HOST=unix:///var/run/docker-host.sock"}, os.Args...)
	cmd := exec.Command("sudo", os.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func DockerHostEmpty() bool {
	return len(os.Getenv("DOCKER_HOST")) == 0
}

func appendKubeConfig() {
	homedir, _ := GetHomeDir()
	kc := filepath.Join(homedir, ".kube/config")
	os.Args = append(os.Args, "--kubeconfig="+kc)
}

func GenerateUUID() string {
	id, _ := uuid.NewRandom()
	return id.String()
}

var alphanum = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var awsAlphanum = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ234567")
var awsBase64 = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+/")

func GenerateRandomSeq(pattern []rune, length int) string {
	rand.Seed(time.Now().UnixNano())
	seq := make([]rune, length)
	for i := range seq {
		seq[i] = pattern[rand.Intn(len(pattern))]
	}
	return string(seq)
}

func keygen(pattern []rune, seqLength int) string {
	return GenerateUUID() + ":" + GenerateRandomSeq(pattern, seqLength)
}

func generateAwsAccessKeyId() string {
	return "AKIA" + GenerateRandomSeq(awsAlphanum, 16)
}

func generateAwsSecretAccessKey() string {
	return GenerateRandomSeq(awsBase64, 40)
}

func awsKeygen() string {
	return generateAwsAccessKeyId() + ":" + GenerateRandomSeq(awsBase64, 40)
}

// FileExists reports whether the named file exists as a boolean
func fileExists(name string) bool {
	if fi, err := os.Stat(name); err == nil {
		if fi.Mode().IsRegular() {
			return true
		}
	}
	return false
}

// DirExists reports whether the dir exists as a boolean
func dirExists(name string) bool {
	if fi, err := os.Stat(name); err == nil {
		if fi.Mode().IsDir() {
			return true
		}
	}
	return false
}
