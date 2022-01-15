// Copyright 2018 Shannon Pekary. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"os"
	"path/filepath"
	"testing"
)


func TestCopy(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "gofileTest")

	cmd := MakeRootCommand()
	cmd.SetArgs([]string{"mkdir", "-v", dir})
	err := cmd.Execute()
	if err != nil {
		t.Error(err)
	}
	info,err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		t.Error("Directory not created")
	}

	cmd.SetArgs([]string{"copy", "-v", "-x", "b", "github.com/goradd/gofile/internal/cmd/testdata/copytest/*", dir})
	err = cmd.Execute()
	if err != nil {
		t.Error(err)
	}

	if _, err = os.Stat(filepath.Join(dir, "a", "t1.txt")); err != nil {
		t.Error(err)
	}

	if _, err = os.Stat(filepath.Join(dir, "b")); err == nil {
		t.Error("Directory b was not supposed to be copied")
	}

	if _, err = os.Stat(filepath.Join(dir, "c")); err != nil {
		t.Error(err)
	}

	cmd.SetArgs([]string{"remove", "-v", dir})
	err = cmd.Execute()

	_,err = os.Stat(dir)
	if err == nil {
		t.Error("Directory not removed")
	}
}

func TestSubCopy(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "gofileTest2")

	cmd := MakeRootCommand()
	cmd.SetArgs([]string{"mkdir", dir})
	err := cmd.Execute()
	if err != nil {
		t.Error(err)
	}
	info,err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		t.Error("Directory not created")
	}

	t.Logf("Copying testdata/emptyTest1/* to %s", dir)

	cmd.SetArgs([]string{"copy", "-v", "testdata/emptyTest1/*", dir})
	err = cmd.Execute()
	if err != nil {
		t.Error(err)
	}

	if _, err = os.Stat(filepath.Join(dir, "d", "t4.txt")); err != nil {
		t.Error(err)
	}

	cmd.SetArgs([]string{"remove", dir})
	err = cmd.Execute()

	_,err = os.Stat(dir)
	if err == nil {
		t.Error("Directory not removed")
	}
}

/*
func TestModuleDirectoryCopy(t *testing.T) {
	modules, err := ModulePaths()

	dir1 := filepath.Join(os.TempDir(), "dir1")
	//dir1 := "github.com/goradd/gofile/test/dir1"
	//dir1,_ = GetModulePath(dir1, modules)

	os.RemoveAll(dir1)

	err = os.Mkdir(dir1, 0777)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer os.RemoveAll(dir1)

	if !IsDir(dir1) {
		t.Fatal("directory not detected")
	}

	dir2 := "github.com/goradd/gofile/internal/testdir"
	dir2Actual, _ := GetModulePath(dir2, modules)

	if err = CopyDirectory(dir2Actual, dir1, CopyOverwrite); err != nil {
		t.Fatal(err)
	}

	if !IsDir(filepath.Join(dir1, "testdir")) {
		t.Fatal("test directory not detected")
	}
}*/

