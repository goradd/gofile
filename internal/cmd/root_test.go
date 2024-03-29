// Copyright 2018 Shannon Pekary. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestModules(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "gofileTest")

	cmd, _ := MakeRootCommand()
	cmd.SetArgs([]string{"mkdir", dir})
	err := cmd.Execute()
	if err != nil {
		t.Error(err)
	}
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		t.Error("Directory not created")
	}

	cmd.SetArgs([]string{"generate", "testdata/generateTestFile.go"})
	err = cmd.Execute()
	if err != nil {
		t.Error(err)
	}

	if string(generateResult)[:10] != "go version" {
		t.Error("Go generate failed.")
	}

	cmd.SetArgs([]string{"copy", "testdata/generateTestFile.go", dir})
	err = cmd.Execute()
	if err != nil {
		t.Error(err)
	}

	_, err = os.Stat(filepath.Join(dir, "generateTestFile.go"))
	if err != nil {
		t.Error(err)
	}

	cmd.SetArgs([]string{"remove", dir})
	err = cmd.Execute()

	_, err = os.Stat(dir)
	if err == nil {
		t.Error("Directory not removed")
	}

}
