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
		require.Equal(t, "packages", root.name)
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

		require.Equal(t, "subf1", root.folders[0].name)
		require.Equal(t, "subf2", root.folders[1].name)
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
		afero.WriteFile(appFS, "/packages/a.js", []byte("file a"), 0644)
		afero.WriteFile(appFS, "/packages/b.py", []byte("file b"), 0644)

		root, err := scanPackagesFolder(appFS, "/")

		require.NoError(t, err) // error in case file system operation failed
		require.Empty(t, root.folders)
		require.NotEmpty(t, root.files)

		require.Equal(t, "a", root.files[0].name)
		require.Equal(t, "b", root.files[1].name)
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
		afero.WriteFile(appFS, "/packages/subf2/subsubf/d.js", []byte("file d"), 0644)

		root, err := scanPackagesFolder(appFS, "/")

		require.NoError(t, err) // error in case file system operation failed

		require.Equal(t, "c", root.folders[0].files[0].name)
		require.Equal(t, "subsubf", root.folders[1].folders[0].name)
		require.Equal(t, "d", root.folders[1].folders[0].files[0].name)
	})

	t.Run("folders should have the complete path to them", func(t *testing.T) {
		appFS := afero.NewMemMapFs()
		appFS.MkdirAll("/packages/subf1", 0755)
		appFS.MkdirAll("/packages/subf2", 0755)

		root, err := scanPackagesFolder(appFS, "/")

		require.NoError(t, err) // error in case file system operation failed
		require.Equal(t, "/packages/subf1", root.folders[0].path)
		require.Equal(t, "/packages/subf2", root.folders[1].path)
	})

	t.Run("files should have the complete path to the code", func(t *testing.T) {
		appFS := afero.NewMemMapFs()
		afero.WriteFile(appFS, "/packages/a.py", []byte("file a"), 0644)
		afero.WriteFile(appFS, "/packages/b.js", []byte("file b"), 0644)

		root, err := scanPackagesFolder(appFS, "/")

		require.NoError(t, err) // error in case file system operation failed
		require.Equal(t, "/packages/a.py", root.files[0].path)
		require.Equal(t, "/packages/b.js", root.files[1].path)
	})
}

func Test_parseProjectTree(t *testing.T) {
	t.Run("should return an empty tree (just root node) of task commands when given an empty project tree", func(t *testing.T) {
		root := ProjectTree{name: "packages"}

		res := parseProjectTree(&root)

		require.Empty(t, res.parent)
		require.Empty(t, res.tasks)
		require.Empty(t, res.command)
	})

	// t.Run("should return tree with child 'wsk package update' when given root with folder", func(t *testing.T) {
	// 	root := ProjectTree{name: "packages"}
	// 	subf := ProjectTree{name: "subf"}
	// 	root.folders = []*ProjectTree{&subf}

	// 	res := parseProjectTree(&root)

	// 	require.Equal(t, "wsk package update subf", res.tasks[0].command)
	// })

	t.Run("should return tree with child 'wsk action update' when given root with file", func(t *testing.T) {
		appFS := afero.NewMemMapFs()
		afero.WriteFile(appFS, "/packages/helloGo.go", []byte("hey there"), 0644)
		afero.WriteFile(appFS, "/packages/helloJava.java", []byte("hey there"), 0644)
		afero.WriteFile(appFS, "/packages/helloJs.js", []byte("hey there"), 0644)
		afero.WriteFile(appFS, "/packages/helloPy.py", []byte("hey there"), 0644)

		root, err := scanPackagesFolder(appFS, "/")

		require.NoError(t, err) // error in case file system operation failed

		res := parseProjectTree(&root)

		require.Equal(t, "wsk action update helloGo /packages/helloGo.go --kind go:default", res.tasks[0].command)
		require.Equal(t, "wsk action update helloJava /packages/helloJava.java --kind java:default", res.tasks[1].command)
		require.Equal(t, "wsk action update helloJs /packages/helloJs.js --kind nodejs:default", res.tasks[2].command)
		require.Equal(t, "wsk action update helloPy /packages/helloPy.py --kind python:default", res.tasks[3].command)
	})

	t.Run("should return tree with cmds for packages and actions when given tree with folders and files", func(t *testing.T) {
		appFS := afero.NewMemMapFs()

		appFS.MkdirAll("/packages/subf1", 0755)
		appFS.MkdirAll("/packages/subf2", 0755)
		afero.WriteFile(appFS, "/packages/a.js", []byte("file a"), 0644)
		afero.WriteFile(appFS, "/packages/b.py", []byte("file b"), 0644)

		root, err := scanPackagesFolder(appFS, "/")

		require.NoError(t, err) // error in case file system operation failed

		res := parseProjectTree(&root)

		require.Equal(t, "wsk action update a /packages/a.js --kind nodejs:default", res.tasks[0].command)
		require.Equal(t, "wsk action update b /packages/b.py --kind python:default", res.tasks[1].command)
		require.Equal(t, "wsk package update subf1", res.tasks[2].command)
		require.Equal(t, "wsk package update subf2", res.tasks[3].command)
	})

	t.Run("should return tree with cmds to create single file actions in packages given tree with sub folders", func(t *testing.T) {
		appFS := afero.NewMemMapFs()

		appFS.MkdirAll("/packages/subf", 0755)
		afero.WriteFile(appFS, "/packages/subf/a.js", []byte("file a"), 0644)
		afero.WriteFile(appFS, "/packages/subf/b.py", []byte("file b"), 0644)

		root, err := scanPackagesFolder(appFS, "/")

		require.NoError(t, err) // error in case file system operation failed

		res := parseProjectTree(&root)

		require.Equal(t, "wsk action update subf/a /packages/subf/a.js --kind nodejs:default", res.tasks[0].tasks[0].command)
		require.Equal(t, "wsk action update subf/b /packages/subf/b.py --kind python:default", res.tasks[0].tasks[1].command)
	})
}
