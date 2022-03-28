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
	"bytes"
	"sync"
)

// BufferPool is a type safe sync.Pool of *byte.Buffer, guaranteed to be Reset
type BufferPool struct {
	sync.Pool
}

// NewBufferPool returns a new bufferPool
func NewBufferPool() *BufferPool {
	return &BufferPool{
		sync.Pool{
			New: func() interface{} {
				// The Pool's New function should generally only return pointer
				// types, since a pointer can be put into the return interface
				// value without an allocation:
				return new(bytes.Buffer)
			},
		},
	}
}

// Get obtains a buffer from the pool
func (b *BufferPool) Get() *bytes.Buffer {
	return b.Pool.Get().(*bytes.Buffer)
}

// Put returns a buffer to the pool, resetting it first
func (b *BufferPool) Put(x *bytes.Buffer) {
	// only store small buffers to avoid pointless allocation
	// avoid keeping arbitrarily large buffers
	if x.Len() > 256 {
		return
	}
	x.Reset()
	b.Pool.Put(x)
}
