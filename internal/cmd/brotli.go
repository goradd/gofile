// Copyright 2022 Shannon Pekary. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	brotlilib "github.com/andybalholm/brotli"
	"github.com/spf13/cobra"
)

func brotli(_ *cobra.Command, _ []string) error {
	if len(files) == 0 {
		if verbose {
			fmt.Printf("No source files were specified in a gzip operation.")
		}
		return nil
	}

	for _, f := range files {
		if err := brotliFile(f); err != nil {
			if filepath.Ext(f) == ".br" {
				continue // do not compress a file that is already compressed
			}
			return fmt.Errorf("error compressing file %s: %s", f, err.Error())
		}
		if deleteAfterZip {
			if err := os.Remove(f); err != nil {
				return fmt.Errorf("error deleting file %s: %s", f, err.Error())
			}
		}
		if verbose {
			fmt.Printf("Brotli compressed %s\n", f)
		}
	}
	return nil
}

func brotliFile(fileName string) error {

	if brotliCompressionLevel < 0 || brotliCompressionLevel > 11 {
		return fmt.Errorf("compression level must be between 0 and 11")
	}

	f, err := os.Create(fileName + ".br")
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	var r *os.File
	r, err = os.Open(fileName)
	if err != nil {
		return err
	}
	defer func() {
		_ = r.Close()
	}()

	var buf []byte
	buf, err = io.ReadAll(r)
	w := brotlilib.NewWriterLevel(f, brotliCompressionLevel)
	_, err = w.Write(buf)
	if err != nil {
		return err
	}
	err = w.Close()
	return err
}
