// Copyright 2018 Shannon Pekary. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func removeFiles(cmd *cobra.Command, args []string) error {
	if len(files) == 0 {
		if verbose {
			fmt.Printf("No source files were specified in a remove operation.")
		}
		return nil
	}

	for _, f := range files {
		err := os.RemoveAll(f)
		if err != nil {
			return err
		} else if verbose {
			fmt.Printf("Removed %s\n", f)
		}
	}
	return nil
}
