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
package commands

import (
	"context"
	"fmt"

	"github.com/go-task/task/v3"
	"github.com/go-task/task/v3/taskfile"
)

type TaskCmd struct {
	Color       bool   `help:"colored output. Enabled by default. Set flag to false or use NO_COLOR=1 to disable (default true)" default:"true" short:"c" env:"NO_COLOR"`
	Concurrency int    `help:"limit number tasks to run concurrently" short:"C"`
	Dir         string `help:"sets directory of execution" short:"d"`
	Dry         bool   `help:"compiles and prints tasks in the order that they would be run, without executing them"`
	Force       bool   `help:"forces execution even when the task is up-to-date" short:"f"`
	Init        bool   `help:"creates a new Taskfile.yml in the current folder" short:"i"`
	List        bool   `help:"lists tasks with description of current Taskfile" short:"l"`
	Output      string `help:"sets output style: [interleaved|group|prefixed]" short:"o"`
	Parallel    bool   `help:"executes tasks provided on command line in parallel" short:"p"`
	Silent      bool   `help:"disables echoing" short:"s"`
	Status      bool   `help:"exits with non-zero exit code if any of the given tasks is not up-to-date"`
	Summary     bool   `help:"show summary about a task"`
	Taskfile    string `help:"choose which Taskfile to run. Defaults to Taskfile.yml" type:"existingfile"`
	Verbose     bool   `help:"enables verbose mode" short:"v"`
	Version     bool   `help:"show Task version"`
	Watch       bool   `help:"enables watch of the given task" short:"w"`
}

func (t *TaskCmd) Run() error {
	runTaskInteractionSample()
	return nil
}

func runTaskInteractionSample() {
	fmt.Println("Doing something with task...")

	te := task.Executor{}
	err := te.Setup()
	if err != nil {
		fmt.Println(err)
		return
	}

	te.RunTask(context.Background(), taskfile.Call{Task: "setup"})
}
