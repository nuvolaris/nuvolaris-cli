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
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/nuvolaris/nuvolaris-cli/nuv/log"
)

// Logger is the kind cli's log.Logger implementation
type Logger struct {
	writer     io.Writer
	writerMu   sync.Mutex
	bufferPool *log.BufferPool
	// kind special additions
	isSmartWriter bool

	spinner       *log.Spinner
	spinnerStatus string
	// for controlling coloring etc
	successFormat string
	failureFormat string
}

var _ log.Logger = &Logger{}

// Debug is part of the log.NuvLogger interface
func (l *Logger) Debug(message string) {
	l.debug(message)
}

// Debugf is part of the log.NuvLogger interface
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.debugf(format, args...)
}

// Info is part of the log.NuvLogger interface
func (l *Logger) Info(message string) {
	l.print(message)
}

// Infof is part of the log.NuvLogger interface
func (l *Logger) Infof(format string, args ...interface{}) {
	l.printf(format, args...)
}

func (l *Logger) ActionWithSpinner(msg string, f func() bool) {
	l.StartSpinner(msg)
	ok := f()
	l.EndSpinner(ok)
}

// StartSpinner starts a new phase of the status, if attached to a terminal
// there will be a loading spinner with this status
func (l *Logger) StartSpinner(msg string) {
	l.EndSpinner(true)
	// set new status
	l.spinnerStatus = msg
	if l.spinner != nil {
		l.spinner.SetSuffix(fmt.Sprintf(" %s ", l.spinnerStatus))
		l.spinner.Start()
	} else {
		l.Infof(" • %s  ...\n", l.spinnerStatus)
	}
}

// EndSpinner completes the current status, ending any previous spinning and
// marking the status as success or failure
func (l *Logger) EndSpinner(success bool) {
	if l.spinnerStatus == "" {
		return
	}

	if l.spinner != nil {
		l.spinner.Stop()
		l.spinner.Write([]byte("\r"))
		// fmt.Fprint(l.spinner.writer, "\r")
	}
	if success {
		l.Infof(l.successFormat, l.spinnerStatus)
	} else {
		l.Infof(l.failureFormat, l.spinnerStatus)
	}

	l.spinnerStatus = ""
}

func (l *Logger) EndSpinnerMsg(success bool, msg string) {
	if l.spinnerStatus == "" {
		return
	}
	l.spinnerStatus = msg
	l.EndSpinner(success)
}

// NewLogger returns the standard logger used by the CLI
func NewLogger() *Logger {
	var writer io.Writer = os.Stdout
	if log.IsSmartTerminal(writer) {
		writer = log.NewSpinner(writer)
	}

	l := &Logger{
		bufferPool: log.NewBufferPool(),
	}
	l.setWriter(writer)

	return l
}

// setWriter sets the output writer
func (l *Logger) setWriter(w io.Writer) {
	l.writerMu.Lock()
	defer l.writerMu.Unlock()

	l.writer = w
	if v2, ok := w.(*log.Spinner); ok {
		l.spinner = v2
		l.isSmartWriter = ok
		// use colored success / failure messages
		l.successFormat = " \x1b[32m✓\x1b[0m %s\n"
		l.failureFormat = " \x1b[31m✗\x1b[0m %s\n"
	} else {
		l.isSmartWriter = log.IsSmartTerminal(w)
	}
}

// synchronized write to the inner writer
func (l *Logger) write(p []byte) (n int, err error) {
	l.writerMu.Lock()
	defer l.writerMu.Unlock()
	return l.writer.Write(p)
}

// writeBuffer writes buf with write, ensuring there is a trailing newline
func (l *Logger) writeBuffer(buf *bytes.Buffer) {
	// ensure trailing newline
	if buf.Len() == 0 || buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	// TODO: should we handle this somehow??
	// Who logs for the logger? 🤔
	_, _ = l.write(buf.Bytes())
}

// print writes a simple string to the log writer
func (l *Logger) print(message string) {
	buf := bytes.NewBufferString(message)
	l.writeBuffer(buf)
}

// printf is roughly fmt.Fprintf against the log writer
func (l *Logger) printf(format string, args ...interface{}) {
	buf := l.bufferPool.Get()
	fmt.Fprintf(buf, format, args...)
	l.writeBuffer(buf)
	l.bufferPool.Put(buf)
}

// addDebugHeader inserts the debug line header to buf
func addDebugHeader(buf *bytes.Buffer) {
	_, file, line, ok := runtime.Caller(3)
	// lifted from klog
	if !ok {
		file = "???"
		line = 1
	} else {
		if slash := strings.LastIndex(file, "/"); slash >= 0 {
			path := file
			file = path[slash+1:]
			if dirsep := strings.LastIndex(path[:slash], "/"); dirsep >= 0 {
				file = path[dirsep+1:]
			}
		}
	}
	buf.Grow(len(file) + 12) // we know at least this many bytes are needed
	buf.WriteString("DEBUG: ")
	buf.WriteString(file)
	buf.WriteByte(':')
	fmt.Fprintf(buf, "%d", line)
	buf.WriteString(" - ")
}

// debug is like print but with a debug log header
func (l *Logger) debug(message string) {
	buf := l.bufferPool.Get()
	addDebugHeader(buf)
	buf.WriteString(message)
	l.writeBuffer(buf)
	l.bufferPool.Put(buf)
}

// debugf is like printf but with a debug log header
func (l *Logger) debugf(format string, args ...interface{}) {
	buf := l.bufferPool.Get()
	addDebugHeader(buf)
	fmt.Fprintf(buf, format, args...)
	l.writeBuffer(buf)
	l.bufferPool.Put(buf)
}
