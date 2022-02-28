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
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func Test_packagesFolderExists(t *testing.T) {
	t.Run("should return true if packages folder is found", func(t *testing.T) {
		fakeFS := fstest.MapFS{ScanFolder: {Mode: fs.ModeDir}}
		exists, err := packagesFolderExists(fakeFS)

		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should return false with error no such file or directory", func(t *testing.T) {
		fakeFS := fstest.MapFS{} // empty
		exists, err := packagesFolderExists(fakeFS)
		assert.Errorf(t, err, "no such file or directory")
		assert.False(t, exists)
	})
}

//  *******************************************

func Test_visitScanFolder(t *testing.T) {

	// No tests if 'packages' does not exist cause checkPackagesFolder stops the pipeline in that case
	t.Run("should return empty tree when ScanFolder is empty", func(t *testing.T) {
		emptyScan := fstest.MapFS{ScanFolder: {Mode: fs.ModeDir}}
		root, _ := visitScanFolder(emptyScan)

		assert.Empty(t, root.packages)
		assert.Empty(t, root.mfActions)
		assert.Empty(t, root.sfActions)
		assert.Equal(t, ScanFolder, root.name)
	})

	t.Run("should return a tree with packages when subfolders are present", func(t *testing.T) {
		packagesExample := fstest.MapFS{
			ScanFolder + "/subf1": {Mode: fs.ModeDir},
			ScanFolder + "/subf2": {Mode: fs.ModeDir},
		}
		expected1 := "subf1"
		expected2 := "subf2"

		root, _ := visitScanFolder(packagesExample)

		assert.Empty(t, root.mfActions)
		assert.Empty(t, root.sfActions)
		assert.NotEmpty(t, root.packages)

		assert.Equal(t, expected1, root.packages[0].name)
		assert.Equal(t, expected2, root.packages[1].name)
	})

	t.Run("should return a tree with single file actions when files are present", func(t *testing.T) {
		sfaExample := fstest.MapFS{
			ScanFolder + "/a.js": {Data: []byte{}},
			ScanFolder + "/b.py": {Data: []byte{}},
		}
		root, _ := visitScanFolder(sfaExample)

		assert.Empty(t, root.packages)
		assert.Empty(t, root.mfActions)
		assert.NotEmpty(t, root.sfActions)

		assert.Equal(t, "a", root.sfActions[0].name)
		assert.Equal(t, "b", root.sfActions[1].name)
	})

	t.Run("should return a tree with sf actions and packages when present", func(t *testing.T) {
		packagesAndSfaExample := fstest.MapFS{
			ScanFolder + "/subf1": {Mode: fs.ModeDir},
			ScanFolder + "/a.js":  {Data: []byte{}},
		}
		root, _ := visitScanFolder(packagesAndSfaExample)

		assert.NotEmpty(t, root.packages)
		assert.NotEmpty(t, root.sfActions)
		assert.Empty(t, root.mfActions)
	})

	t.Run("should return a tree with packages and mfActions with sub sub folders", func(t *testing.T) {
		mfaExample := fstest.MapFS{
			ScanFolder + "/subf1/mfa/package.json": {Data: []byte{}},
			ScanFolder + "/subf1/mfa/a.js":         {Data: []byte{}},
		}
		root, _ := visitScanFolder(mfaExample)

		assert.Empty(t, root.sfActions)
		assert.NotEmpty(t, root.packages)
		assert.Empty(t, root.mfActions)

		assert.Equal(t, "mfa", root.packages[0].mfActions[0].name)
	})

	t.Run("packages should have the complete path to them", func(t *testing.T) {
		sub1Path := ScanFolder + "/subf1"
		sub2Path := ScanFolder + "/subf2"
		pathExample := fstest.MapFS{
			sub1Path: {Mode: fs.ModeDir},
			sub2Path: {Mode: fs.ModeDir},
		}
		root, _ := visitScanFolder(pathExample)

		assert.Equal(t, sub1Path, root.packages[0].path)
		assert.Equal(t, sub2Path, root.packages[1].path)
	})

	t.Run("actions should have the complete path to the code", func(t *testing.T) {
		subSFAPath := ScanFolder + "/a.py"
		subMFAPath := ScanFolder + "/subf1/mfa"
		pathExample := fstest.MapFS{
			subSFAPath:           {Data: []byte{}},
			subMFAPath + "/b.js": {Data: []byte{}},
		}
		root, _ := visitScanFolder(pathExample)

		assert.Equal(t, subSFAPath, root.sfActions[0].path)
		assert.Equal(t, subMFAPath, root.packages[0].mfActions[0].path)
	})

	t.Run("actions should hold runtime", func(t *testing.T) {
		subSFAPath := ScanFolder + "/a.py"
		subMFAPath := ScanFolder + "/subf1/mfa"
		runtimeExample := fstest.MapFS{
			subSFAPath:           {Data: []byte{}},
			subMFAPath + "/b.js": {Data: []byte{}},
		}
		root, _ := visitScanFolder(runtimeExample)

		assert.Equal(t, ".py", root.sfActions[0].runtime)
		assert.Equal(t, ".js", root.packages[0].mfActions[0].runtime)
	})
}

func Test_findMfaRuntime(t *testing.T) {
	t.Run("should return error when no runtime found", func(t *testing.T) {
		emptyScan := fstest.MapFS{ScanFolder: {Mode: fs.ModeDir}}
		runtime, err := findMfaRuntime(emptyScan, ScanFolder)

		assert.Empty(t, runtime)
		assert.Errorf(t, err, "no supported runtime found")
	})

	t.Run("should return js runtime when present", func(t *testing.T) {
		rtExample := fstest.MapFS{"a.js": {Data: []byte{}}}
		checkIfRuntimePresent(t, rtExample, jsRuntime)

		rtExample = fstest.MapFS{"package.json": {Data: []byte{}}}
		checkIfRuntimePresent(t, rtExample, jsRuntime)
	})

	t.Run("should return python runtime when present", func(t *testing.T) {
		rtExample := fstest.MapFS{"a.py": {Data: []byte{}}}
		checkIfRuntimePresent(t, rtExample, pyRuntime)

		rtExample = fstest.MapFS{"requirements.txt": {Data: []byte{}}}
		checkIfRuntimePresent(t, rtExample, pyRuntime)
	})

	t.Run("should return java runtime when present", func(t *testing.T) {
		rtExample := fstest.MapFS{"a.java": {Data: []byte{}}}
		checkIfRuntimePresent(t, rtExample, javaRuntime)

		rtExample = fstest.MapFS{"pom.xml": {Data: []byte{}}}
		checkIfRuntimePresent(t, rtExample, javaRuntime)
	})

	t.Run("should return go runtime when present", func(t *testing.T) {
		rtExample := fstest.MapFS{"a.go": {Data: []byte{}}}
		checkIfRuntimePresent(t, rtExample, goRuntime)

		rtExample = fstest.MapFS{"go.mod": {Data: []byte{}}}
		checkIfRuntimePresent(t, rtExample, goRuntime)
	})

}
func checkIfRuntimePresent(t *testing.T, rtExample fs.FS, expectedRuntime string) {
	t.Helper()
	runtime, err := findMfaRuntime(rtExample, "")
	assert.Equal(t, expectedRuntime, runtime)
	assert.NoError(t, err)
}

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
