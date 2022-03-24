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
package log

import (
	"io"
	"os"
	"runtime"

	"github.com/mattn/go-isatty"
)

// a fake TTY type for testing that can only be implemented within this package
type isTestFakeTTY interface {
	isTestFakeTTY()
}

// IsTerminal returns true if the writer w is a terminal
func IsTerminal(w io.Writer) bool {
	// check for internal fake type we can use for testing.
	if _, ok := (w).(isTestFakeTTY); ok {
		return true
	}
	// check for real terminals
	if v, ok := (w).(*os.File); ok {
		return isatty.IsTerminal(v.Fd())
	}
	return false
}

// IsSmartTerminal returns true if the writer w is a terminal AND
// we think that the terminal is smart enough to use VT escape codes etc.
func IsSmartTerminal(w io.Writer) bool {
	return isSmartTerminal(w, runtime.GOOS, os.LookupEnv)
}

func isSmartTerminal(w io.Writer, GOOS string, lookupEnv func(string) (string, bool)) bool {
	// Not smart if it's not a tty
	if !IsTerminal(w) {
		return false
	}

	// getenv helper for when we only care about the value
	getenv := func(e string) string {
		v, _ := lookupEnv(e)
		return v
	}

	// Explicit request for no ANSI escape codes
	// https://no-color.org/
	if _, set := lookupEnv("NO_COLOR"); set {
		return false
	}

	// Explicitly dumb terminals are not smart
	// https://en.wikipedia.org/wiki/Computer_terminal#Dumb_terminals
	term := getenv("TERM")
	if term == "dumb" {
		return false
	}
	// st has some bug 🤷‍♂️
	// https://github.com/kubernetes-sigs/kind/issues/1892
	if term == "st-256color" {
		return false
	}

	// On Windows WT_SESSION is set by the modern terminal component.
	// Older terminals have poor support for UTF-8, VT escape codes, etc.
	if GOOS == "windows" && getenv("WT_SESSION") == "" {
		return false
	}

	// OK, we'll assume it's smart now, given no evidence otherwise.
	return true
}
