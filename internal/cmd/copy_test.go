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
	cmd.SetArgs([]string{"mkdir", dir})
	err := cmd.Execute()
	if err != nil {
		t.Error(err)
	}
	info,err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		t.Error("Directory not created")
	}

	cmd.SetArgs([]string{"copy", "-x", "b", "testdata/copytest/*", dir})
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

	cmd.SetArgs([]string{"remove", dir})
	err = cmd.Execute()

	_,err = os.Stat(dir)
	if err == nil {
		t.Error("Directory not removed")
	}

}


