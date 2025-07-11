// Copyright 2018 Shannon Pekary. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package sys

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestModules(t *testing.T) {

	// since modules are very dependent on the environment, this is mostly just
	// a sanity check that it works without errors.
	modules, err := ModulePaths()
	if err != nil {
		t.Fatal(err)
		return
	}

	newPath, err := GetModulePath("github.com/goradd/gofile/modules", modules)
	if err != nil {
		t.Fatal(err)
		return
	}

	if newPath == "github.com/goradd/gofile/modules" {
		t.Error("Module path was not changed to an absolute path")
	}

	newPath, err = GetModulePath("github.com/goradd/gofile/*/modules", modules)
	if err != nil {
		t.Fatal(err)
		return
	}

	if newPath == "github.com/goradd/gofile/*/modules" {
		t.Error("Module path was not changed to an absolute path")
	}

	var newPath2 string
	newPath2, err = GetModulePath("github.com/goradd/gofile", modules)
	if err != nil {
		t.Fatal(err)
		return
	}
	if filepath.Join(newPath2, "/*/modules") != newPath {
		t.Error("Module path is not rooted")
	}
}

func TestSplitCommandParts(t *testing.T) {
	parts, err := splitCommandParts("a b c")
	if fmt.Sprint(parts) != "[a b c]" || err != nil {
		t.Error(fmt.Sprint(parts))
	}

	parts, err = splitCommandParts("ab bc cd")
	if fmt.Sprint(parts) != "[ab bc cd]" || err != nil {
		t.Error(fmt.Sprint(parts))
	}

	parts, err = splitCommandParts(`"ab bc cd"`)
	if len(parts) != 1 || fmt.Sprint(parts) != "[ab bc cd]" || err != nil {
		t.Error(fmt.Sprint(parts))
	}

	parts, err = splitCommandParts(`'ab bc cd"'`)
	if len(parts) != 1 || fmt.Sprint(parts) != `[ab bc cd"]` || err != nil {
		t.Error(fmt.Sprint(parts))
	}

	parts, err = splitCommandParts(`ab 'b c' cd`)
	if len(parts) != 3 || fmt.Sprint(parts) != `[ab b c cd]` || err != nil {
		t.Error(fmt.Sprint(parts))
	}

	parts, err = splitCommandParts(`ab'b c'cd`)
	if len(parts) != 1 || fmt.Sprint(parts) != `[abb ccd]` || err != nil {
		t.Error(fmt.Sprint(parts))
	}
}

func TestExecuteShellCommand(t *testing.T) {
	_, err := ExecuteShellCommand(`abc "`)
	if err == nil {
		t.Error("error expected")
	}
}

func TestImportPath(t *testing.T) {
	modules, err := ModulePaths()
	if err != nil {
		t.Fatal(err)
		return
	}

	var newPath string
	newPath, err = GetModulePath("github.com/goradd/gofile/pkg/sys/testdata", modules)
	if err != nil {
		t.Fatal(err)
		return
	}

	var s string
	s, err = ImportPath(filepath.Join(newPath, "t1.txt"))
	if err != nil {
		t.Fatal(err)
		return
	}
	if s != "github.com/goradd/gofile/pkg/sys/testdata" {
		t.Error("ImportPath is not correct.")
		return
	}

	s, err = ImportPath(newPath)
	if err != nil {
		t.Fatal(err)
		return
	}
	if s != "github.com/goradd/gofile/pkg/sys/testdata" {
		t.Error("ImportPath is not correct.")
		return
	}

}
