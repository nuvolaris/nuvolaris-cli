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

// Name is the name of the CLI
const Name = "nuv"

// Description is the description of the nuv CLI
const Description = "nuv is the command line tool to manage Nuvolaris"

// MinDockerVersion required
const MinDockerVersion = "18.06.3-ce"

// MinDockerMem is the minimum amount of memory required by docker
const MinDockerMem = (4 * 1000 * 1000 * 1000) - 1

// TimeoutInSec is global timeout of 5 mins when polling
const TimeoutInSec = 600

const ScanFolder = "packages"

const WskPropsFilename = ".wskprops"

// NuvolarisNamespace is Kubernetes namespace where nuvolaris components are deployed
const NuvolarisNamespace = "nuvolaris"

// NuvolarisConfigmapName is the config map from where to read annotations
const NuvolarisConfigmapName = "config"
