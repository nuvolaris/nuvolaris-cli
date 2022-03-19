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
	"encoding/json"
	"fmt"
	"io/fs"
	"os"

	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Cmd struct {
	Mb      mb      `cmd:"" help:"creates an S3 bucket"`
	List    ls      `cmd:"" help:"lists S3 objects and common prefixes under a prefix or all S3 buckets"`
	Put     put     `cmd:"" help:"puts a local file in a S3 bucket"`
	Secrets secrets `cmd:"" help:"sets secrets for S3 buckets"`
}
type mb struct {
	BucketName string `arg:"" type:"string"`
}

func (c *mb) Run() error {
	_, err := newS3session()
	if err != nil {
		return err
	}
	fmt.Printf("Creating bucket %q...\n", c.BucketName)
	return nil
}

type ls struct{}

func (c *ls) Run() error {
	return nil
}

type put struct{}

func (c *put) Run() error {
	return nil
}

type secrets struct{}

func (c *secrets) Run() error {
	return nil
}

func newS3session() (*s3.S3, error) {
	path, err := GetOrCreateNuvolarisConfigDir()
	if err != nil {
		return nil, err
	}
	fsys := os.DirFS(path)
	_, err = readS3Secrets(fsys)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

type s3SecretsJSON struct {
	Id     string
	Key    string
	Region string
}

func readS3Secrets(fsys fs.FS) (s3SecretsJSON, error) {
	content, err := fs.ReadFile(fsys, "secrets.json")
	if err != nil {
		return s3SecretsJSON{}, fmt.Errorf("unable to read s3 secrets. Did you set them with nuv s3 secrets?")
	}
	var secrets s3SecretsJSON
	err = json.Unmarshal(content, &secrets)
	if err != nil {
		return s3SecretsJSON{}, err
	}
	return secrets, nil
}
