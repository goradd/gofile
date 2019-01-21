// Copyright 2018 Shannon Pekary. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package sys

import (
	"fmt"
	"testing"
)


func TestModules(t *testing.T) {

	// since modules are very dependent on the environment, this is mostly just
	// a sanity check that it works without errors.

	modules,err := ModulePaths()
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
}

func TestSplitCommandParts(t *testing.T) {
	parts,err := splitCommandParts("a b c")
	if fmt.Sprint(parts) != "[a b c]" || err != nil {
		t.Error(fmt.Sprint(parts))
	}

	parts,err = splitCommandParts("ab bc cd")
	if fmt.Sprint(parts) != "[ab bc cd]" || err != nil {
		t.Error(fmt.Sprint(parts))
	}

	parts,err = splitCommandParts(`"ab bc cd"`)
	if len(parts) != 1 || fmt.Sprint(parts) != "[ab bc cd]" || err != nil {
		t.Error(fmt.Sprint(parts))
	}

	parts,err = splitCommandParts(`'ab bc cd"'`)
	if len(parts) != 1 || fmt.Sprint(parts) != `[ab bc cd"]` || err != nil {
		t.Error(fmt.Sprint(parts))
	}

	parts,err = splitCommandParts(`ab 'b c' cd`)
	if len(parts) != 3 || fmt.Sprint(parts) != `[ab b c cd]` || err != nil {
		t.Error(fmt.Sprint(parts))
	}

	parts,err = splitCommandParts(`ab'b c'cd`)
	if len(parts) != 1 || fmt.Sprint(parts) != `[abb ccd]` || err != nil {
		t.Error(fmt.Sprint(parts))
	}



}

