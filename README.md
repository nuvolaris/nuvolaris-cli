<!--
  ~ Licensed to the Apache Software Foundation (ASF) under one
  ~ or more contributor license agreements.  See the NOTICE file
  ~ distributed with this work for additional information
  ~ regarding copyright ownership.  The ASF licenses this file
  ~ to you under the Apache License, Version 2.0 (the
  ~ "License"); you may not use this file except in compliance
  ~ with the License.  You may obtain a copy of the License at
  ~
  ~   http://www.apache.org/licenses/LICENSE-2.0
  ~
  ~ Unless required by applicable law or agreed to in writing,
  ~ software distributed under the License is distributed on an
  ~ "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
  ~ KIND, either express or implied.  See the License for the
  ~ specific language governing permissions and limitations
  ~ under the License.
  ~
-->
# Nuvolaris CLI

This repo discuss the Command Line Interface to Nuvolaris.

You can discuss it in the #[nuvolaris-cli](https://discord.gg/JWqFJJfvED) discord channel and in the forum under the category [CLI](https://github.com/nuvolaris/nuvolaris/discussions/categories/cli).

# Design

**NOTE** This design is work in progress and incomplete - feel free to propose improvements but **please document the feature and send a PR to this design BEFORE implementing the feature** to avoid rejections and time wasted.

## Goals

The tool `nuv` built in the repo [nuvolaris-cli](https://github.com/nuvolaris/nuvolaris-cli) it the Command Line interface for the Nuvolaris project.

It embeds the functionalities of the tool [wsk](https://github.com/apache/openwhisk-cli) for creating actions. As wsk is written in Go, we can directly include the code in `nuv`

It embeds the functionalities of the tool [task](https://taskfile.dev) for execution actions. As task is writtein in Go we can directly include the code in `nuv`

It also adds the some project conventions inspired by the [nim](https://github.com/nimbella/nimbella-cli) tool. But since nim is written in typescript we do not include that code, we will reimplement them. Most notably we want to reimplement the [project detection](#project-detection) euristic described below and nothing else.

The toll works scanning the current subtree, looking for actions and packages to deploy. It works generating a `Taskfile` (that can be inspected by the users) and then executing it.

It will be possible to add customizations of the task adding locally some `nuvolaris.yml` in the various subdirectories. This functionality will be described later.

There will be initially 4 commands:

- `nuv scan` scans the folder and generate a Taskfile 
- `nuv wsk` executes the wsk subcommand
- `nuv task` execute the task subcommand
- `nuv install` will also be able to execute a kubectl command that deploys the `nuvolaris-operator` that in turns inizialize openwhisk in any available kubernetes accessible with `kubectl` and initialize the `.wskprops` file used by `nuv wsk`

The expected workflow is that :
1. `nuv install` instal an openwhisk cluster using a configured `kubectl` in the path
2. `nuv scan` generates a `Taskfile`
3. `nuw task` execute the Taskfile that embeds a numer of `nuv wsk` commands
4. the various `nuv wsk` creates then a full project

An example of a a project to deploy can be [this](https://github.com/pagopa/io-sdk/tree/master/admin)

## Project Detection
`nuv` will scan the current directory looking for a folder named `packages` 

If it finds here a file it will create a package for each subfolder.

If it finds files in the file `packages` it will deploy them as [single file actions](#single-file-actions) in the package `default`. If it finds files in the subfolders of `packages` it will deploy them as  [single file actions](#single-file-actions) in packages named as the the subfolder. If it finds folders it will build [multi file actions](#multi-file-actions).

## Single File Actions

A single file actions is simply a file witn an extensions.

This extensions can be one of the supported ones: `.js`  `.py` `.go` `.java` 

This will cause the cration of an action with `--kind nodejs:default`, `--kind python:default`, `--kind go:default` and `--kind java:default` using the correct runtime.

The correct runtime is  described by `runtime.json` that can be downloaded from the api host configured.

If the extension is in format:  `.<version>.<extension>` it will deploy an action of  `--kind <language>:<version>`

## Multi File Actions

**initial draft**

A multi file action is stored in a subfolder of a subfolder of `packages`.

This is expected to be a file to build

`nuv` implements some heuristics to decide the correct type of the file to build.

Currently:

- if there is a `package.json`  or any `js` field in the folder then it is  `.js` and it builds with `npm install ; npm build`
- if there is a `requirements.txt` or any `.py` file then it is a pythn and it builds creating a virtual env as described in the python runtime documentation
- if there is `pom.xml` the it builds using `mvn install`
- if there is a `go.mod` then it builds using `go build`

then it will zip the folder and send as an action of the current type to the runtime.





