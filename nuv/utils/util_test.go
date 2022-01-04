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
package utils

import "fmt"

func ExampleSysErr() {
	SysErr("/bin/echo 1 2 3")
	SysErr("/bin/echo 3", "4", "5")
	SysErr("@sh -c", "echo foo >/tmp/foo")
	out_foo, _ := SysErr("cat /tmp/foo")
	fmt.Print(out_foo)
	SysErr("@sh -c", "echo bar >/tmp/bar")
	out_bar, _ := SysErr("@cat /tmp/bar")
	fmt.Print(out_bar)
	_, err := SysErr("false")
	fmt.Println("ERR", err)
	_, err = SysErr("donotexist")
	fmt.Println("ERR", err)
	// Output:
	// 1 2 3
	// 3 4 5
	// foo
	// foo
	// bar
	// ERR exit status 1
	// ERR exec: "donotexist": executable file not found in $PATH
}

func ExampleDryRunSysErr() {
	DryRunPush("first", "second", "!third")
	out, err := DryRunSysErr("dummy")
	fmt.Println(1, out, err)
	out, err = DryRunSysErr("dummy", "alpha", "beta")
	fmt.Println(2, out, err)
	out, err = DryRunSysErr("dummy")
	fmt.Println(3, "out", out, "err", err)
	out, err = DryRunSysErr("dummy")
	fmt.Println(4, "out", out, "err", err)
	// Output:
	// dummy
	// 1 first <nil>
	// dummy alpha beta
	// 2 second <nil>
	// dummy
	// 3 out  err third
	// dummy
	// 4 out  err <nil>
}
