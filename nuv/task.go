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

//go:build !fast
// +build !fast

package main

import (
	taskmain "github.com/go-task/task/cmd/task"
)

type TaskCmd struct {
	Args []string `arg:"" optional:"" name:"args" help:"task subcommand args"`
}

// Task runs embedded task
func Task(args ...string) error {
	taskmain.TaskMain(append([]string{"task"}, args...))
	return nil
}

func (task *TaskCmd) Run() error {
	return Task(task.Args...)
}
