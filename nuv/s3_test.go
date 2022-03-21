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
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_readS3Secrets(t *testing.T) {

	t.Run("should return error when unable to read secrets.json", func(t *testing.T) {
		emptyConfig := fstest.MapFS{"/": {Mode: fs.ModeDir}}
		_, err := readS3Secrets(emptyConfig)
		assert.Error(t, err)
	})

	t.Run("should return s3Secrets when secrets.json is valid", func(t *testing.T) {
		expected := s3SecretsJSON{
			Id:     "some-id",
			Key:    "some-key",
			Region: "some-region",
		}
		secrets := `
{
	"id": "some-id",
	"key": "some-key",
	"region": "some-region"
}`
		fakeFS := fstest.MapFS{"secrets.json": {Data: []byte(secrets)}}

		config, err := readS3Secrets(fakeFS)
		assert.NoError(t, err)
		assert.Equal(t, expected, config)
	})
}

func Test_buildAwsConfig(t *testing.T) {
	secrets := s3SecretsJSON{
		Id:     "some-id",
		Key:    "some-key",
		Region: "some-region",
	}

	expected := &aws.Config{
		Region: aws.String(secrets.Region),
		Credentials: credentials.NewStaticCredentials(
			secrets.Id,
			secrets.Key,
			"",
		),
	}

	config := buildAwsConfig(secrets)
	assert.Equal(t, expected, config)
}

type mockS3Client struct {
	s3iface.S3API
	mock.Mock
}

func (m *mockS3Client) CreateBucket(in *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
	args := m.Called(in)
	return args.Get(0).(*s3.CreateBucketOutput), args.Error(1)
}

func (m *mockS3Client) ListObjectsV2(in *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	args := m.Called(in)
	return args.Get(0).(*s3.ListObjectsV2Output), args.Error(1)
}
func (m *mockS3Client) PutObject(in *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	args := m.Called(in)
	return args.Get(0).(*s3.PutObjectOutput), args.Error(1)
}

func Test_createBucket(t *testing.T) {
	t.Run("should use CreateBucket and return an error if unable to create bucket", func(t *testing.T) {
		mockSvc := new(mockS3Client)
		mockSvc.On("CreateBucket", mock.Anything).Return(&s3.CreateBucketOutput{}, errors.New("failed"))
		bucketName := "some-bucket"
		err := createBucket(mockSvc, bucketName)
		mockSvc.AssertCalled(t, "CreateBucket", &s3.CreateBucketInput{Bucket: aws.String(bucketName)})
		assert.Error(t, err)
	})

	t.Run("should use CreateBucket to create a bucket when no error occurs", func(t *testing.T) {
		mockSvc := new(mockS3Client)
		mockSvc.On("CreateBucket", mock.Anything).Return(&s3.CreateBucketOutput{}, nil)
		bucketName := "some-bucket"
		err := createBucket(mockSvc, bucketName)
		mockSvc.AssertCalled(t, "CreateBucket", &s3.CreateBucketInput{Bucket: aws.String(bucketName)})
		assert.NoError(t, err)
	})
}

func Test_listBucketContent(t *testing.T) {
	t.Run("should use ListObjectsV2 and return an error if error occurred", func(t *testing.T) {
		mockSvc := new(mockS3Client)
		mockSvc.On("ListObjectsV2", mock.Anything).Return(&s3.ListObjectsV2Output{}, errors.New("failed"))
		bucketName := "some-bucket"
		err := listBucketContent(mockSvc, bucketName)
		mockSvc.AssertCalled(t, "ListObjectsV2", &s3.ListObjectsV2Input{Bucket: aws.String(bucketName)})
		assert.Error(t, err)
	})
}

func Test_putFile(t *testing.T) {
	t.Run("should use PutObjet and return an error if error occurred", func(t *testing.T) {
		mockSvc := new(mockS3Client)
		mockSvc.On("PutObject", mock.Anything).Return(&s3.PutObjectOutput{}, errors.New("failed"))
		bucketName := "some-bucket"
		file := "some-file"
		err := preparePut(file, "hello")(mockSvc, bucketName)
		in := &s3.PutObjectInput{
			Body:   strings.NewReader("hello"),
			Key:    aws.String(file),
			Bucket: aws.String(bucketName),
			ACL:    aws.String(s3.BucketCannedACLPublicRead),
		}
		mockSvc.AssertCalled(t, "PutObject", in)
		assert.Error(t, err)
	})
}

func Test_saveSecrets(t *testing.T) {
	id := "some-id"
	key := "some-key"
	region := "some-region"
	s := secrets{Id: id, Key: key, Region: region}
	err := s.Run()
	assert.NoError(t, err)

	path, _ := GetOrCreateNuvolarisConfigDir()
	path = filepath.Join(path, secretFile)
	content, err := os.ReadFile(path)
	assert.NoError(t, err)

	var sj s3SecretsJSON
	json.Unmarshal(content, &sj)
	assert.Equal(t, s3SecretsJSON{Id: id, Key: key, Region: region}, sj)
}
