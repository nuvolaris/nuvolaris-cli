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
}

func Test_checkPackagesFolder(t *testing.T) {
	t.Run("should return true if packages folder is found", func(t *testing.T) {
		appFS := afero.NewMemMapFs()
		appFS.MkdirAll("/packages", 0755)

		exists, err := checkPackagesFolder(appFS, "/")

		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("should return false with no error when packages not found", func(t *testing.T) {
		appFS := afero.NewMemMapFs()

		exists, err := checkPackagesFolder(appFS, "./")

		require.False(t, exists)
		require.NoError(t, err) // error in case file system operation failed
	})
}

func Test_scanPackagesFolder(t *testing.T) {
	// No tests if 'packages' does not exist cause checkPackagesFolder stops the pipeline in that case
	t.Run("should return a tree with just root node when packages folder is empty", func(t *testing.T) {
		appFS := afero.NewMemMapFs()
		appFS.Mkdir("/packages", 0755)

		root, err := scanPackagesFolder(appFS, "/")

		require.Empty(t, root.folders)
		require.Empty(t, root.files)
		require.Empty(t, root.parent)
		require.Equal(t, root.name, "packages")
		require.NoError(t, err) // error in case file system operation failed
	})

	t.Run("should return a root with folders when packages has subfolders", func(t *testing.T) {
		appFS := afero.NewMemMapFs()
		appFS.MkdirAll("/packages/subf1", 0755)
		appFS.MkdirAll("/packages/subf2", 0755)

		root, err := scanPackagesFolder(appFS, "/")

		require.NoError(t, err) // error in case file system operation failed
		require.Empty(t, root.files)
		require.NotEmpty(t, root.folders)

		require.Equal(t, root.folders[0].name, "subf1")
		require.Equal(t, root.folders[1].name, "subf2")
	})

	t.Run("folder children should have the root as parent", func(t *testing.T) {
		appFS := afero.NewMemMapFs()
		appFS.MkdirAll("/packages/subf1", 0755)
		appFS.MkdirAll("/packages/subf2", 0755)

		root, err := scanPackagesFolder(appFS, "/")

		require.NoError(t, err) // error in case file system operation failed
		for _, c := range root.folders {
			require.Equal(t, &root, c.parent)
		}
	})

	t.Run("should return a root with files children when packages has files", func(t *testing.T) {
		appFS := afero.NewMemMapFs()
		afero.WriteFile(appFS, "/packages/a", []byte("file a"), 0644)
		afero.WriteFile(appFS, "/packages/b", []byte("file b"), 0644)

		root, err := scanPackagesFolder(appFS, "/")

		require.NoError(t, err) // error in case file system operation failed
		require.Empty(t, root.folders)
		require.NotEmpty(t, root.files)

		require.Equal(t, root.files[0].name, "a")
		require.Equal(t, root.files[1].name, "b")
	})

	t.Run("should return a root with files and folders when packages has both", func(t *testing.T) {
		appFS := afero.NewMemMapFs()
		appFS.MkdirAll("/packages/subf1", 0755)
		appFS.MkdirAll("/packages/subf2", 0755)
		afero.WriteFile(appFS, "/packages/a", []byte("file a"), 0644)
		afero.WriteFile(appFS, "/packages/b", []byte("file b"), 0644)

		root, err := scanPackagesFolder(appFS, "/")

		require.NoError(t, err) // error in case file system operation failed
		require.NotEmpty(t, root.folders)
		require.NotEmpty(t, root.files)
	})

	t.Run("should return a complete tree representing the packages folder", func(t *testing.T) {
		appFS := afero.NewMemMapFs()

		appFS.MkdirAll("/packages/subf1", 0755)
		appFS.MkdirAll("/packages/subf2", 0755)
		appFS.MkdirAll("/packages/subf2/subsubf", 0755)
		afero.WriteFile(appFS, "/packages/a", []byte("file a"), 0644)
		afero.WriteFile(appFS, "/packages/b", []byte("file b"), 0644)
		afero.WriteFile(appFS, "/packages/subf1/c", []byte("file c"), 0644)
		afero.WriteFile(appFS, "/packages/subf2/subsubf/d", []byte("file d"), 0644)

		root, err := scanPackagesFolder(appFS, "/")

		require.NoError(t, err) // error in case file system operation failed

		require.Equal(t, root.folders[0].files[0].name, "c")
		require.Equal(t, root.folders[1].folders[0].name, "subsubf")
		require.Equal(t, root.folders[1].folders[0].files[0].name, "d")
	})
}
