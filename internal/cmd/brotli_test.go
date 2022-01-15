// Copyright 2018 Shannon Pekary. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package cmd

import (
	brotlilib "github.com/andybalholm/brotli"
	"io"
	"os"
	"path/filepath"
	"testing"
)


func TestBrotli(t *testing.T) {
	f := filepath.Join(os.TempDir(), "brotliTest")
	fbr := f + ".br"
	const testString = "This is a test"

	err := os.WriteFile(f, []byte(testString), 0o666)
	if err != nil {
		t.Error("Test file not created :" + err.Error())
	}

	cmd := MakeRootCommand()
	cmd.SetArgs([]string{"brotli", f})
	err = cmd.Execute()
	if err != nil {
		t.Error(err)
	}
	_,err = os.Stat(fbr)
	if err != nil  {
		t.Error("Brotli file not created")
	}
	var r *os.File
	r, err = os.Open(fbr)
	if err != nil {
		t.Error(err)
	}

	zr := brotlilib.NewReader(r)
	var out []byte
	out, err = io.ReadAll(zr)

	if string(out) != testString {
		t.Error("unBrotli comparison failed: " + string(out))
	}
	_ = os.Remove(f)
	_ = os.Remove(fbr)
}

func TestBrotliDir(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "brotliTestDir")
	_ = os.Mkdir(dir, 0o777)

	f := filepath.Join(dir, "brotliTest")
	f2 := filepath.Join(dir, "brotliTest2.txt")
	f3 := filepath.Join(dir, "brotliTest3.abc")
	fgz := f + ".br"
	f2gz := f2 + ".br"
	f3gz := f3 + ".br"
	const testString = "This is another test"

	if err := os.WriteFile(f, []byte(testString), 0o666); err != nil {
		t.Error("Test file not created :" + err.Error())
	}
	if err := os.WriteFile(f2, []byte(testString), 0o666); err != nil {
		t.Error("Test file #2 not created :" + err.Error())
	}
	if err := os.WriteFile(f3, []byte(testString), 0o666); err != nil {
		t.Error("Test file #3 not created :" + err.Error())
	}

	cmd := MakeRootCommand()
	cmd.SetArgs([]string{"brotli", "-v", "-d", "-q", "4", "-x", "*.abc", dir})
	if err := cmd.Execute(); err != nil {
		t.Error(err)
	}

	if _,err := os.Stat(f); err == nil  {
		t.Error("Source file not removed")
	}

	if _,err := os.Stat(fgz); err != nil  {
		t.Error("Brotli file not created")
	}
	if _,err := os.Stat(f2gz); err != nil  {
		t.Error("Brotli file #2 not created")
	}

	if _,err := os.Stat(f3gz); err == nil  {
		t.Error("Brotli file #3 was created, but should have been skipped")
	}
	if _,err := os.Stat(f3); err != nil  {
		t.Error("Source file #3 was removed, but should have been left intact")
	}


	if r, err := os.Open(fgz); err != nil {
		t.Error(err)
	} else {
		zr := brotlilib.NewReader(r)
		out, _ := io.ReadAll(zr)
		if string(out) != testString {
			t.Error("Brotli decompress comparison failed: " + string(out))
		}
	}

	if r, err := os.Open(f2gz); err != nil {
		t.Error(err)
	} else {
		zr := brotlilib.NewReader(r)
		out, _ := io.ReadAll(zr)
		if string(out) != testString {
			t.Error("Brotli decompress #2 comparison failed: " + string(out))
		}
	}

	if err := os.RemoveAll(dir) ; err != nil {
		t.Error(err)
	}
}

