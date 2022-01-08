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
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/spf13/afero"
)

func TestNuvScan(t *testing.T) {
	t.Run("should have scan subcmd help", func(t *testing.T) {
		var cli CLI
		app := NewTestApp(t, &cli)
		require.PanicsWithValue(t, true, func() { // TODO: explain why needed
			_, err := app.Parse([]string{"scan", "--help"})
			require.NoError(t, err)
		})
	})

	t.Run("", func(t *testing.T) {
		var cli CLI
		app := NewTestApp(t, &cli)
		_, err := app.Parse([]string{"scan"})
		require.NoError(t, err)
	})
}

func Test_checkPackagesFolder(t *testing.T) {
	t.Run("should return true if packages folder is found", func(t *testing.T) {
		appFS := afero.NewMemMapFs()
		appFS.MkdirAll("packages", 0755)

		exists, err := checkPackagesFolder(appFS)

		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("should return false with no error when packages not found", func(t *testing.T) {
		appFS := afero.NewMemMapFs()

		exists, err := checkPackagesFolder(appFS)

		require.False(t, exists)
		require.NoError(t, err) // error in case file system operation failed
	})
}
