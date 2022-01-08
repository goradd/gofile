// Copyright 2018 Shannon Pekary. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package sys

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCopyFile(t *testing.T) {
	dir, err := os.UserCacheDir()
	if err != nil {
		dir = os.TempDir()
	}
	dir = filepath.Join(dir, "gofileTest")
	err = os.Mkdir(dir, 0777)
	if err != nil {
		os.RemoveAll(dir)
		t.Fatal(err)
		return
	}
	defer os.RemoveAll(dir)

	src := filepath.Join(dir, "test.txt")

	testContent := "I am a test \n and the 2nd line"

	err = ioutil.WriteFile(src, []byte(testContent), 0755)
	if err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(dir, "test2.txt")

	time.Sleep(2 * time.Second) // delay to make sure file is created later

	err = CopyFiles(dst, CopyDoNotOverwrite, src)
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := ioutil.ReadFile(dst)
	if string(bytes) != testContent {
		t.Error("File content does not match")
	}

	time.Sleep(2 * time.Second) // delay to make sure file is created later

	// test grow
	testContent2 := "I am a test \n and the 2nd line\nand a 3rd"
	src2 := filepath.Join(dir, "test2-2.txt")
	err = ioutil.WriteFile(src2, []byte(testContent2), 0755)
	if err != nil {
		t.Fatal(err)
	}

	// We should not overwrite a file that already exists
	err = CopyFiles(dst, CopyDoNotOverwrite, src2)
	if err != nil {
		t.Fatal(err)
	}
	bytes, err = ioutil.ReadFile(dst)
	if string(bytes) != testContent {
		t.Error("File content should not have changed")
	}

	// We should overwrite a file that already exists here
	err = CopyFiles(dst, CopyOverwriteOnlyIfNewer, src2)
	if err != nil {
		t.Fatal(err)
	}
	bytes, err = ioutil.ReadFile(dst)
	if string(bytes) != testContent2 {
		t.Error("File should overwrite")
	}

	// We should not overwrite a newer file with an older file here
	err = CopyFiles(dst, CopyOverwriteOnlyIfNewer, src)
	if err != nil {
		t.Fatal(err)
	}
	bytes, err = ioutil.ReadFile(dst)
	if string(bytes) != testContent2 {
		t.Error("File should not overwrite")
	}

	// We should overwrite the file here
	err = CopyFiles(dst, CopyOverwrite, src)
	if err != nil {
		t.Fatal(err)
	}
	bytes, err = ioutil.ReadFile(dst)
	if string(bytes) != testContent {
		t.Error("File should overwrite")
	}

	// test shrink of file after copying over
	testContent3 := "I am a smaller test"
	src3 := filepath.Join(dir, "test3.txt")
	err = ioutil.WriteFile(src3, []byte(testContent3), 0755)

	err = CopyFiles(dst, CopyOverwrite, src3)
	if err != nil {
		t.Fatal(err)
	}
	bytes, err = ioutil.ReadFile(dst)
	if string(bytes) != testContent3 {
		t.Error("File content does not match")
	}

	// Test the copy directory capability of CopyFiles
	dst2 := filepath.Join(os.TempDir(), "gofileTest2")
	defer os.RemoveAll(dst2)

	err = CopyFiles(dst2, CopyDoNotOverwrite, dir)
	if err == nil {
		t.Error("Copying a directory to a non-existant location should fail")
	}

	err = CopyFiles(dst2+"/", CopyDoNotOverwrite, dir)
	if err == nil {
		t.Error("Copying a directory to a non-existant location should fail")
	}

	err = os.Mkdir(dst2, 0777)
	if err != nil {
		t.Fatal(err)
	}

	err = CopyFiles(dst2, CopyDoNotOverwrite, dir)
	if err != nil {
		t.Fatal(err)
	}

	items, _ := ioutil.ReadDir(dst2)
	if items[0].Name() != "gofileTest" {
		t.Fatal("First item in directory is not gofileTest")
	}

	// Now copy individual items
	var items2 []string
	items, _ = ioutil.ReadDir(dir)
	for _, item := range items {
		items2 = append(items2, filepath.Join(dir, item.Name()))
	}
	err = CopyFiles(dst2, CopyDoNotOverwrite, items2...)
	if err != nil {
		t.Fatal(err)
	}

	items, _ = ioutil.ReadDir(dst2)
	if len(items) != 5 {
		t.Error("Items were not copied")
	}
}

func TestDirectoryCopy(t *testing.T) {
	dir1 := filepath.Join(os.TempDir(), "dir1")
	dir2 := filepath.Join(os.TempDir(), "dir2")

	err := os.Mkdir(dir1, 0777)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer os.RemoveAll(dir1)

	err = os.Mkdir(dir2, 0777)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer os.RemoveAll(dir2)

	subdir := filepath.Join(dir1, "subdir1")
	if err := os.Mkdir(subdir, 0777); err != nil {
		t.Fatal(err)
	}

	// set up the test directory
	testContent := "I am a test"
	if err := ioutil.WriteFile(filepath.Join(dir1, "test1"), []byte(testContent), 0755); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(filepath.Join(dir1, "test2"), []byte(testContent), 0755); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(filepath.Join(subdir, "test3"), []byte(testContent), 0755); err != nil {
		t.Fatal(err)
	}

	if err := CopyDirectory(dir1, dir2, CopyOverwrite); err != nil {
		t.Fatal(err)
	}

	items, err := ioutil.ReadDir(dir2)
	if err != nil {
		t.Fatal(err)
	}
	if items[0].Name() != "dir1" {
		t.Fatal("First item in directory is not dir1")
	}
	items, err = ioutil.ReadDir(filepath.Join(dir2, "dir1"))
	if err != nil {
		t.Fatal(err)
	}

	if items[0].Name() != "subdir1" {
		t.Fatal("First item in directory is not subdir1, but rather: " + items[0].Name())
	}
	if items[1].Name() != "test1" {
		t.Fatal("Second item in directory is not test1, but rather: " + items[1].Name())
	}
}

func TestSplitList(t *testing.T) {
	s := "a;b"
	l := SplitList(s)
	if len(l) != 2 {
		t.Error("List is not split correctly")
	}

	s = "a:b:c"
	l = SplitList(s)
	if len(l) != 3 {
		t.Error("List is not split correctly")
	}

}

func TestModuleExpandFileList(t *testing.T) {
	modules, err := ModulePaths()
	if err != nil {
		t.Fatal(err)
	}

	fileList := []string{
		"github.com/goradd/gofile/*.md", // readme file
		"github.com/goradd/gofile/LICENSE",
		"/a/b/c", // non-existent items should be left alone
	}
	l := ModuleExpandFileList(fileList, modules)

	if len(l) != 3 {
		t.Error("Not the correct list size.")
	}

	for _, i := range l {
		if i[:6] == "github" {
			t.Error("Item 1 or 2 was not changed.")
		}
	}

	if !listContains(l, filepath.FromSlash("/a/b/c")) {
		t.Error("Item 3 was changed.")
	}
}

func listContains(list []string, val string) bool {
	for _, l := range list {
		if l == val {
			return true
		}
	}
	return false
}
