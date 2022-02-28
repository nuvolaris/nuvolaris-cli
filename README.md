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

This repo discusses the Command Line Interface to Nuvolaris.

You can discuss it in the #[nuvolaris-cli](https://discord.gg/JWqFJJfvED) discord channel and in the forum under the category [CLI](https://github.com/nuvolaris/nuvolaris/discussions/categories/cli).

# Design

**NOTE** This design is work in progress and incomplete - feel free to propose improvements but **please document the feature and send a PR to this design BEFORE implementing the feature** to avoid rejections and time wasted.

## Goals

The tool `nuv` built in the repo [nuvolaris-cli](https://github.com/nuvolaris/nuvolaris-cli) it the Command Line Interface for the Nuvolaris project.

It embeds the functionalities of the tool [wsk](https://github.com/apache/openwhisk-cli) for creating actions. Since wsk is written in Go, we can directly include the code in `nuv`

It embeds the functionalities of the tool [task](https://taskfile.dev) for execution actions. Since task is written in Go, we can directly include the code in `nuv`

It also adds some project conventions inspired by the [nim](https://github.com/nimbella/nimbella-cli) tool. But since nim is written in typescript we do not include it, but will reimplement it. Most notably we want to reimplement the [project detection](#project-detection) heuristic described below and nothing else.

The tool works by scanning the current subtree, looking for actions and packages to deploy. It generates a `Taskfile` (that can be inspected by the users) and then executes it.

It will be possible to add customizations of the task adding locally `nuvolaris.yml` in the various subdirectories. This functionality will be described later.

Initially, there will be 4 commands:

- `nuv scan` scans the folder and generates a Taskfile 
- `nuv wsk` executes the wsk subcommand
- `nuv task` executes the task subcommand
- `nuv setup` will also be able to execute a kubectl command that deploys the `nuvolaris-operator` that in turns inizializes openwhisk in any available kubernetes accessible with `kubectl` and initialize the `.wskprops` file used by `nuv wsk`

The expected workflow is that:
1. `nuv setup` installs an openwhisk cluster using `kubectl` configured in the path
2. `nuv scan` generates a `Taskfile` 
3. `nuv task` executes the `Taskfile` that embeds many `nuv wsk` commands
4. the various `nuv wsk` create then a full project

An example of a project to deploy can be [this](https://github.com/pagopa/io-sdk/tree/master/admin)

## Project Detection
`nuv` will scan the current directory looking for a folder named `packages` 

If it finds here a file, it will create a package for each subfolder.

If it finds files in the file `packages`, it will deploy them as [single file actions](#single-file-actions) in the package `default`. If it finds files in the subfolders of `packages` it will deploy them as [single file actions](#single-file-actions) in packages named as the the subfolder. If it finds folders it will build [multi file actions](#multi-file-actions).

## Single File Actions

A single file actions is simply a file with an extension.

This extension can be one of the supported ones: `.js`  `.py` `.go` `.java` 

This will cause the creation of an action with `--kind nodejs:default`, `--kind python:default`, `--kind go:default` and `--kind java:default` using the correct runtime.

The correct runtime is described by `runtime.json` that can be downloaded from the configured api host.

If the extension is in format:  `.<version>.<extension>`, it will deploy an action of  `--kind <language>:<version>`

## Static frontend

Nuv is also able to deploy static frontends. A static front-end is a collection of static asset under a given folder that will be published in a web server under a path. 

A folder containing static (web) assets is always named `web` and can be placed in different parts in the folder hierarchy. The path in the website where is published depends on the location in the hierarchy, as described below.

Before publishing, `nuv` executes some build commands.

### Hostname

In general, for each namespace there will be a `https://<namespace>.<domain>` website where to publish the resources. For the local deployment there will be a website `http://127.0.0.1:8080` where the resources are published, with the namespace and the domain ignored.

### Path detection

The path where the assets are published depends on the path in the action hierarchy.

The sub-folder `web` is published as "/".

Any subfolder `web` under `packages/<package>/web` is published unser `/<packages>/`.

Any subfolder `web` under `packages/default/<action>\web` is published as `/<action>`.

Any subfolder `web` under `packages/<package>/<action>/web` is published as `/<package>/<action>`

What is published (files collected) and how it is built is defined by the next paragraph.

### Building and Collecting

In every folder `web` it will check if there is a `nuvolaris.json`

If there is not a `nuvolaris.json` and not a `package.json` it will assume this base `nuvolaris.json`:

```
{
  "collect": ".",
  "install": "echo nothing to install",
  "build": "echo nothing to build"
}
```

If instead there is `packages.json`, it will assume this base `nuvolaris.json`:

```
{
  "collect": "public",
  "install": "npm install",
  "build": "npm run build"
}
```

The it will read the `nuvolaris.json` replacing the keys in it with the default ones.

The generated taskfile will execute at deployment step:

- the command defined by `install` only if there is not a `node_modules`
- the command defined by `build` always
- then it will collect for publishing (creating a crd instance) the files in the folder defined by `collect`

It is recommended that `nuv scan` does not execute directy the command but instead it delegates to another command like `nuv build` and in turn the creation of `crd` to another `nuv crd` subcommand, after changing to the corresponding suddirectory. All those commands should work by default in current directory. 




## Multi File Actions

**initial draft**

A multi-file action is stored in a subfolder of a subfolder of `packages`.

This is expected to be a file to build.

`nuv` implements some heuristics to decide the correct type of the file to build.

Currently:

- if there is a `package.json`  or any `js` field in the folder then it is  `.js` and it builds with `npm install ; npm build`
- if there is a `requirements.txt` or any `.py` file then it is python and it builds creating a virtual env as described in the python runtime documentation
- if there is `pom.xml` then it builds using `mvn install`
- if there is a `go.mod` then it builds using `go build`

then it will zip the folder and send as an action of the current type to the runtime.







