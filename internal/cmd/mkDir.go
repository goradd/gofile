// Copyright 2018 Shannon Pekary. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func mkDir(_ *cobra.Command, _ []string) error {
	for _, dir := range files {
		if err := os.MkdirAll(dir, os.FileMode(0777)); err != nil {
			return err
		} else if verbose {
			fmt.Printf("Created directory %s\n", dir)
		}
	}
	return nil
}
