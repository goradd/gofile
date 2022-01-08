// Copyright 2018 Shannon Pekary. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package cmd

import (
	ziplib "compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestGZip(t *testing.T) {
	f := filepath.Join(os.TempDir(), "gzipTest")
	fgz := f + ".gz"
	const testString = "This is a test"

	err := os.WriteFile(f, []byte(testString), 0o666)
	if err != nil {
		t.Error("Test file not created :" + err.Error())
	}

	cmd := MakeRootCommand()
	cmd.SetArgs([]string{"gzip", f})
	err = cmd.Execute()
	if err != nil {
		t.Error(err)
	}
	_, err = os.Stat(fgz)
	if err != nil {
		t.Error("Zip file not created")
	}
	var r *os.File
	r, err = os.Open(fgz)
	if err != nil {
		t.Error(err)
	}

	zr, _ := ziplib.NewReader(r)
	var out []byte
	out, err = io.ReadAll(zr)

	if string(out) != testString {
		t.Error("unzip comparison failed: " + string(out))
	}
	_ = os.Remove(f)
	_ = os.Remove(fgz)
}

func TestGZipDir(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "gzipTestDir")
	if err := os.Mkdir(dir, 0o777); err != nil {
		// if it already exists from prior failed test, remove it and try again
		os.Remove(dir)
		os.Mkdir(dir, 0o777)
	}

	f := filepath.Join(dir, "gzipTest")
	f2 := filepath.Join(dir, "gzipTest2.txt")
	f3 := filepath.Join(dir, "gzipTest3.abc")
	fgz := f + ".gz"
	f2gz := f2 + ".gz"
	f3gz := f3 + ".gz"
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
	cmd.SetArgs([]string{"gzip", "-v", "-d", "-x", "*.abc", dir})
	if err := cmd.Execute(); err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(f); err == nil {
		t.Error("Source file not removed")
	}

	if _, err := os.Stat(fgz); err != nil {
		t.Error("Zip file not created")
	}
	if _, err := os.Stat(f2gz); err != nil {
		t.Error("Zip file #2 not created")
	}

	if _, err := os.Stat(f3gz); err == nil {
		t.Error("Zip file #3 was created, but should have been skipped")
	}
	if _, err := os.Stat(f3); err != nil {
		t.Error("Source file #3 was removed, but should have been left intact")
	}

	if r, err := os.Open(fgz); err != nil {
		t.Error(err)
	} else {
		zr, _ := ziplib.NewReader(r)
		out, _ := io.ReadAll(zr)
		r.Close()
		if string(out) != testString {
			t.Error("unzip comparison failed: " + string(out))
		}
	}

	if r, err := os.Open(f2gz); err != nil {
		t.Error(err)
	} else {
		zr, _ := ziplib.NewReader(r)
		out, _ := io.ReadAll(zr)
		r.Close()
		if string(out) != testString {
			t.Error("unzip #2 comparison failed: " + string(out))
		}
	}

	if err := os.RemoveAll(dir); err != nil {
		t.Error(err)
	}
}
