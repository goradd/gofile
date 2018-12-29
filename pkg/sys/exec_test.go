// Copyright 2018 Shannon Pekary. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package sys

import (
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


