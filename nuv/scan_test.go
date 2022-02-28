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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNuvScan(t *testing.T) {
	t.Run("should have scan subcmd help", func(t *testing.T) {
		var cli CLI
		app := NewTestApp(t, &cli)
		require.PanicsWithValue(t, true, func() {
			_, err := app.Parse([]string{"scan", "--help"})
			require.NoError(t, err)
		})
	})

	t.Run("should generate a Taskfile", func(t *testing.T) {
		var cli CLI
		app := NewTestApp(t, &cli)
		c, _ := app.Parse([]string{"scan", "test-embed/", "-o", "../"})
		err := c.Run()
		require.NoError(t, err)

		// TODO check task file
	})
}

const testFolder = "./test-embed/scan-tests/"
const emptyPkg = testFolder + "empty-pkg"
const foldersOnly = testFolder + "folders-only"

func Test_packagesFolderExists(t *testing.T) {
	t.Run("should return true if packages folder is found", func(t *testing.T) {
		exists, err := packagesFolderExists(testFolder)

		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should return false with error no such file or directory", func(t *testing.T) {
		exists, err := packagesFolderExists("non-existing-folder")
		assert.Errorf(t, err, "no such file or directory")
		assert.False(t, exists)
	})
}

//  *******************************************

func Test_visitScanFolder(t *testing.T) {
	// No tests if 'packages' does not exist cause checkPackagesFolder stops the pipeline in that case
	t.Run("should return a tree with just root node when packages folder is empty", func(t *testing.T) {
		root, err := visitScanFolder(emptyPkg)

		assert.Empty(t, root.folders)
		assert.Empty(t, root.mfActions)
		assert.Empty(t, root.sfActions)
		assert.Equal(t, ScanFolder, root.name)
		assert.NoError(t, err)
	})

	t.Run("should return a root with folders when packages has subfolders", func(t *testing.T) {
		expected1 := "subf1"
		expected2 := "subf2"

		root, err := visitScanFolder(foldersOnly)

		assert.NoError(t, err)
		assert.Empty(t, root.mfActions)
		assert.Empty(t, root.sfActions)
		assert.NotEmpty(t, root.folders)

		assert.Equal(t, expected1, root.folders[0].name)
		assert.Equal(t, expected2, root.folders[1].name)
	})

	// 	s.T().Run("should return a root with single file actions when packages has files", func(t *testing.T) {
	// 		testWithFs([]string{}, []string{"a.js", "b.py"}, func() {
	// 			root, err := scanPackagesFolder(testFolder)

	// 			s.Assert().NoError(err) // error in case file system operation failed
	// 			s.Assert().Empty(root.folders)
	// 			s.Assert().Empty(root.mfActions)
	// 			s.Assert().NotEmpty(root.sfActions)

	// 			s.Assert().Equal("a", root.sfActions[0].name)
	// 			s.Assert().Equal("b", root.sfActions[1].name)
	// 		})
	// 	})

	// 	s.T().Run("should return a root with sf actions and folders when 'packages' has both", func(t *testing.T) {
	// 		testWithFs([]string{"subf1", "subf2"}, []string{"a.js", "b.py"}, func() {
	// 			root, err := scanPackagesFolder(testFolder)

	// 			s.Assert().NoError(err) // error in case file system operation failed
	// 			s.Assert().NotEmpty(root.folders)
	// 			s.Assert().NotEmpty(root.sfActions)
	// 			s.Assert().Empty(root.mfActions)
	// 		})
	// 	})

	// 	s.T().Run("should return a tree with folders and mfActions when 'packages' has sub sub folders", func(t *testing.T) {
	// 		testWithFs(
	// 			[]string{"subf1/a1", "subf1/a2", "subf2/b1"},
	// 			[]string{"subf1/a1/package.json", "subf1/a1/a1.js", "subf1/a2/package.json", "subf1/a2/a2.js", "subf2/b1/assertments.txt", "subf2/b1/b1.py"},
	// 			func() {
	// 				root, err := scanPackagesFolder(testFolder)

	// 				s.Assert().NoError(err) // error in case file system operation failed
	// 				s.Assert().Empty(root.sfActions)
	// 				s.Assert().NotEmpty(root.folders)
	// 				s.Assert().Empty(root.mfActions)

	// 				s.Assert().Equal("a1", root.folders[0].mfActions[0].name)
	// 				s.Assert().Equal("a2", root.folders[0].mfActions[1].name)
	// 				s.Assert().Equal("b1", root.folders[1].mfActions[0].name)
	// 			})
	// 	})

	// 	s.T().Run("should return a complete tree representing the packages folder", func(t *testing.T) {
	// 		testWithFs(
	// 			[]string{"subf1", "subf1", "subf2/subsubf"},
	// 			[]string{"a.js", "b.go", "subf1/c.js", "subf2/subsubf/d.js"},
	// 			func() {

	// 				root, err := scanPackagesFolder(testFolder)

	// 				s.Assert().NoError(err) // error in case file system operation failed

	// 				s.Assert().Equal("a", root.sfActions[0].name)
	// 				s.Assert().Equal("b", root.sfActions[1].name)
	// 				s.Assert().Equal("c", root.folders[0].sfActions[0].name)
	// 				s.Assert().Equal("subsubf", root.folders[1].mfActions[0].name)

	// 				s.Assert().Len(root.sfActions, 2)
	// 				s.Assert().Len(root.folders, 2)
	// 				s.Assert().Empty(root.mfActions)
	// 				s.Assert().Len(root.folders[0].sfActions, 1)
	// 				s.Assert().Len(root.folders[1].mfActions, 1)
	// 			})
	// 	})
	// 	s.T().Run("folders should have the complete path to them", func(t *testing.T) {
	// 		testWithFs(
	// 			[]string{"subf1", "subf2"},
	// 			[]string{},
	// 			func() {

	// 				root, err := scanPackagesFolder(testFolder)

	// 				s.Assert().NoError(err) // error in case file system operation failed
	// 				s.Assert().Equal(filepath.Join(testFolder, "/packages/subf1"), root.folders[0].path)
	// 				s.Assert().Equal(filepath.Join(testFolder, "/packages/subf2"), root.folders[1].path)
	// 			})
	// 	})

	// 	s.T().Run("actions should have the complete path to the code", func(t *testing.T) {
	// 		testWithFs(
	// 			[]string{"subf", "subf/sub"},
	// 			[]string{"a.py", "subf/sub/b.js"},
	// 			func() {

	// 				root, err := scanPackagesFolder(testFolder)

	// 				s.Assert().NoError(err) // error in case file system operation failed
	// 				s.Assert().Equal(filepath.Join(testFolder, "/packages/a.py"), root.sfActions[0].path)
	// 				s.Assert().Equal(filepath.Join(testFolder, "/packages/subf/sub"), root.folders[0].mfActions[0].path)
	// 			})
	// 	})

	// 	s.T().Run("actions should hold runtime", func(t *testing.T) {
	// 		testWithFs(
	// 			[]string{"subf", "subf/sub"},
	// 			[]string{"a.py", "subf/sub/b.js"},
	// 			func() {
	// 				root, err := scanPackagesFolder(testFolder)

	// 				s.Assert().NoError(err) // error in case file system operation failed
	// 				s.Assert().Equal(".py", root.sfActions[0].runtime)
	// 				s.Assert().Equal(".js", root.folders[0].mfActions[0].runtime)
	// 			})
	// 	})
}

// //  *******************************************

// //  *** findMfaRuntime function tests ***
// func helpTestForRuntime(s *nuvScanTestSuite, searchFor, expectedRuntime string) {
// 	s.T().Helper()
// 	testWithFs([]string{}, []string{searchFor}, func() {
// 		runtime, err := findMfaRuntime(pkgPath)

// 		s.Assert().NoError(err)
// 		s.Assert().Equal(expectedRuntime, runtime)
// 	})
// }
// func (s *nuvScanTestSuite) Test_findMfaRuntime() {
// 	s.T().Run("should return error when no runtime found", func(t *testing.T) {
// 		testWithFs([]string{}, []string{}, func() {
// 			runtime, err := findMfaRuntime(testFolder)

// 			s.Assert().Empty(runtime)
// 			s.Assert().Errorf(err, "no supported runtime found")
// 		})
// 	})

// 	s.T().Run("should return correct runtime when present", func(t *testing.T) {
// 		helpTestForRuntime(s, "package.json", jsRuntime)
// 		helpTestForRuntime(s, "a.js", jsRuntime)
// 		helpTestForRuntime(s, "requirements.txt", pyRuntime)
// 		helpTestForRuntime(s, "a.py", pyRuntime)
// 		helpTestForRuntime(s, "pom.xml", javaRuntime)
// 		helpTestForRuntime(s, "a.java", javaRuntime)
// 		helpTestForRuntime(s, "go.mod", goRuntime)
// 		helpTestForRuntime(s, "a.go", goRuntime)
// 	})
// }

// //  *****************************************

// //  *** parseProjectTree function tests ***
// func (s *nuvScanTestSuite) Test_parseProjectTree() {
// 	s.T().Run("should return an empty slice of commands when given an empty project tree", func(t *testing.T) {
// 		root := ProjectTree{name: "packages"}

// 		res := parseProjectTree(&root)

// 		s.Assert().Empty(res)
// 	})

// 	s.T().Run("should return slice with 'wsk package update' when given root with folder", func(t *testing.T) {
// 		root := ProjectTree{name: "packages"}
// 		subf := ProjectTree{name: "subf"}
// 		root.folders = []*ProjectTree{&subf}

// 		res := parseProjectTree(&root)

// 		s.Assert().Equal("wsk package update subf", res[0])
// 	})

// 	s.T().Run("should return slice with 'wsk action update' when given root with file", func(t *testing.T) {
// 		testWithFs(
// 			[]string{},
// 			[]string{
// 				"helloGo.go",
// 				"helloJava.java",
// 				"helloJs.js",
// 				"helloPy.py",
// 			}, func() {
// 				expectedGo := fmt.Sprintf("wsk action update helloGo %s --kind go:default", filepath.Join(pkgPath, "helloGo.go"))
// 				expectedjava := fmt.Sprintf("wsk action update helloJava %s --kind java:default", filepath.Join(pkgPath, "helloJava.java"))
// 				expectedJs := fmt.Sprintf("wsk action update helloJs %s --kind nodejs:default", filepath.Join(pkgPath, "helloJs.js"))
// 				expectedPy := fmt.Sprintf("wsk action update helloPy %s --kind python:default", filepath.Join(pkgPath, "helloPy.py"))

// 				root, err := scanPackagesFolder(testFolder)

// 				s.Assert().NoError(err) // error in case file system operation failed

// 				res := parseProjectTree(&root)

// 				s.Assert().Equal(expectedGo, res[0])
// 				s.Assert().Equal(expectedjava, res[1])
// 				s.Assert().Equal(expectedJs, res[2])
// 				s.Assert().Equal(expectedPy, res[3])
// 			})
// 	})

// 	s.T().Run("should return slice with cmds for packages and actions when given tree with folders and files", func(t *testing.T) {
// 		testWithFs(
// 			[]string{"subf1"},
// 			[]string{
// 				"a.js",
// 			}, func() {
// 				expectedPkg := "wsk package update subf1"
// 				expectedJs := fmt.Sprintf("wsk action update a %s --kind nodejs:default", filepath.Join(pkgPath, "a.js"))

// 				root, err := scanPackagesFolder(testFolder)

// 				s.Assert().NoError(err) // error in case file system operation failed

// 				res := parseProjectTree(&root)

// 				s.Assert().Equal(expectedJs, res[0])
// 				s.Assert().Equal(expectedPkg, res[1])
// 			})
// 	})

// 	s.T().Run("should return slice with single file actions cmds in packages given tree with sub folders", func(t *testing.T) {
// 		testWithFs(
// 			[]string{"subf"},
// 			[]string{
// 				"subf/a.js",
// 			}, func() {

// 				expectedJs := fmt.Sprintf("wsk action update subf/a %s --kind nodejs:default", filepath.Join(pkgPath, "subf/a.js"))
// 				root, err := scanPackagesFolder(testFolder)

// 				s.Assert().NoError(err) // error in case file system operation failed

// 				res := parseProjectTree(&root)

// 				s.Assert().Equal(expectedJs, res[1])
// 			})
// 	})

// 	s.T().Run("should return slice with multi file action cmds given tree with mfActions", func(t *testing.T) {
// 		testWithFs(
// 			[]string{"subf/mf"},
// 			[]string{
// 				"subf/mf/a.js",
// 			}, func() {

// 				zipcmd := fmt.Sprintf("zip -r %s.zip %s/*", filepath.Join(pkgPath, "subf/mf/mf"), filepath.Join(pkgPath, "subf/mf"))
// 				mfacmd := fmt.Sprintf("wsk action update subf/mf %s --kind nodejs:default", filepath.Join(pkgPath, "subf/mf/mf.zip"))

// 				root, err := scanPackagesFolder(testFolder)
// 				s.Assert().NoError(err)

// 				res := parseProjectTree(&root)
// 				s.Assert().Equal("wsk package update subf", res[0])
// 				s.Assert().Equal(zipcmd, res[1])
// 				s.Assert().Equal(mfacmd, res[2])
// 			})
// 	})

// 	s.T().Run("should return slice with both sf and mf actions", func(t *testing.T) {
// 		testWithFs(
// 			[]string{"subf/mf"},
// 			[]string{"subf/a.js", "subf/mf/b.py"},
// 			func() {
// 				expectedSF := fmt.Sprintf("wsk action update subf/a %s --kind nodejs:default", filepath.Join(pkgPath, "subf/a.js"))

// 				expectedZip := fmt.Sprintf("zip -r %s.zip %s/*", filepath.Join(pkgPath, "subf/mf/mf"), filepath.Join(pkgPath, "subf/mf"))
// 				expectedMfa := fmt.Sprintf("wsk action update subf/mf %s --kind python:default", filepath.Join(pkgPath, "subf/mf/mf.zip"))

// 				root, err := scanPackagesFolder(testFolder)
// 				s.Assert().NoError(err)
// 				res := parseProjectTree(&root)

// 				expected := []string{"wsk package update subf", expectedSF, expectedZip, expectedMfa}
// 				s.Assert().ElementsMatch(res, expected)
// 			})
// 	})
// }
