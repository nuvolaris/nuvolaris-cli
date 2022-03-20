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
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

func runS3(f func(s3iface.S3API, string) error, bucket string) error {
	session, err := newS3session()
	if err != nil {
		return err
	}
	return f(session, bucket)
}

type S3Cmd struct {
	Mb      mb      `cmd:"" help:"creates an S3 bucket"`
	List    ls      `cmd:"" help:"lists S3 objects and common prefixes under a prefix or all S3 buckets"`
	Put     put     `cmd:"" help:"puts a local file in a S3 bucket"`
	Secrets secrets `cmd:"" help:"sets secrets for S3 buckets"`
}
type mb struct {
	BucketName string `arg:"" type:"string" help:"the name of the bucket to create"`
}

func (c *mb) Run() error {
	return runS3(createBucket, c.BucketName)
}

type ls struct {
	BucketName string `arg:"" type:"string" help:"the name of the bucket to list"`
}

func (c *ls) Run() error {
	return runS3(listBucketContent, c.BucketName)
}

type put struct {
	BucketName string `arg:"" type:"string" help:"the name of the bucket to use"`
	File       string `arg:"" type:"path" help:"the file to put in the bucket"`
}

func (c *put) Run() error {
	content, err := os.ReadFile(c.File)
	if err != nil {
		return err
	}
	putFile := preparePut(c.File, string(content))
	return runS3(putFile, c.BucketName)
}

type secrets struct{}

func (c *secrets) Run() error {
	return nil
}

const s3endpoint = "a-dummy-endpoint"

func newS3session() (s3iface.S3API, error) {
	path, err := GetOrCreateNuvolarisConfigDir()
	if err != nil {
		return nil, err
	}
	fsys := os.DirFS(path)
	secrets, err := readS3Secrets(fsys)
	if err != nil {
		return nil, err
	}

	conf := buildAwsConfig(secrets)
	awsSession, err := session.NewSession(conf)
	if err != nil {
		return nil, err
	}

	return s3.New(awsSession, &aws.Config{Endpoint: aws.String(s3endpoint)}), nil
}

func createBucket(svc s3iface.S3API, bucketName string) error {
	fmt.Printf("Creating bucket %q...\n", bucketName)
	in := &s3.CreateBucketInput{Bucket: aws.String(bucketName)}
	_, err := svc.CreateBucket(in)
	if err != nil {
		return err
	}
	fmt.Printf("Bucket %q created", bucketName)
	return nil
}

func listBucketContent(svc s3iface.S3API, bucketName string) error {
	in := &s3.ListObjectsV2Input{Bucket: aws.String(bucketName)}
	o, err := svc.ListObjectsV2(in)
	if err != nil {
		return err
	}
	fmt.Printf("%+v", o)
	return nil
}

func preparePut(fileName, content string) func(s3iface.S3API, string) error {
	return func(svc s3iface.S3API, bucketName string) error {
		fmt.Printf("Uploading %q to bucket %q...", fileName, bucketName)
		in := &s3.PutObjectInput{
			Body:   strings.NewReader(content),
			Key:    aws.String(fileName),
			Bucket: aws.String(bucketName),
			ACL:    aws.String(s3.BucketCannedACLPublicRead),
		}
		_, err := svc.PutObject(in)
		if err != nil {
			return err
		}
		fmt.Println("File uploaded successfully")
		return nil
	}
}

func buildAwsConfig(s s3SecretsJSON) *aws.Config {
	conf := aws.NewConfig()
	conf.WithRegion(s.Region)
	conf.WithCredentials(credentials.NewStaticCredentials(s.Id, s.Key, ""))
	return conf
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
