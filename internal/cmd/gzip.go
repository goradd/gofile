// Copyright 2018 Shannon Pekary. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package cmd

import (
	ziplib "compress/gzip"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
)

func gzip(_ *cobra.Command, _ []string) error {
	if len(files) == 0 {
		if verbose {
			fmt.Printf("No source files were specified in a gzip operation.")
		}
		return nil
	}

	for _,f := range files {
		if err := zipFile(f); err != nil {
			if filepath.Ext(f) == ".gz" {
				continue // do not compress a file that is already compressed
			}
			return fmt.Errorf("error zipping file %s: %s", f, err.Error())
		}
		if deleteAfterZip {
			if err := os.Remove(f); err != nil {
				return fmt.Errorf("error deleting file %s: %s", f, err.Error())
			}
		}
		if verbose {
			fmt.Printf("Zipped %s\n", f)
		}
	}
	return nil
}

func zipFile(fileName string) error {
	f, err := os.Create(fileName + ".gz")
	if err != nil {
		return err
	}
	defer func () {
		_ = f.Close()
	}()

	var r *os.File
	r, err = os.Open(fileName)
	if err != nil {
		return err
	}
	defer func () {
		_ = r.Close()
	}()

	var buf []byte
	buf, err = io.ReadAll(r)
	var w *ziplib.Writer
	w, err = ziplib.NewWriterLevel(f, gzipCompressionLevel)
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	if err != nil {
		return err
	}
	err = w.Close()
	return err
}