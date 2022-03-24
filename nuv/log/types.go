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

// Logger defines the logging interface
// It is roughly a subset of github.com/kubernetes/klog
type Logger interface {
	// Debug is used for debug messages
	Debug(message string)
	// Debugf is used to write a Printf style debug message
	Debugf(format string, args ...interface{})

	// Info is used for normal user facing messages
	Info(message string)
	// Infof is used to write a Printf style user facing status message
	Infof(format string, args ...interface{})

	// StartSpinner starts the loading spinner with a message
	StartSpinner(status string)
	//EndSpinner stops the spinner and displays a checkmark for success or fail state
	EndSpinner(success bool)

	//EndSpinnerMsg works like EndSpinner but changes the msg displayed next to the checkmark
	EndSpinnerMsg(success bool, msg string)

	// Wraps a function with the start and end spinner around
	ActionWithSpinner(msg string, f func() bool)
}
