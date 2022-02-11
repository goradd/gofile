// Copyright 2018 Shannon Pekary. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"github.com/goradd/gofile/pkg/sys"
	"github.com/spf13/cobra"
)

func copyFiles(cmd *cobra.Command, args []string) error {
	// Cobra will guarantee we have at least 2 arguments
	dest := args[len(args) - 1]
	dest = processFileArg(dest)

	args = args[:len(args) - 1]
	processFileListArgs(cmd, args) // puts the list of files in the files global

	var overwrite = sys.CopyDoNotOverwrite

	if copyOverwrite {
		overwrite = sys.CopyOverwrite
	} else if copyOverwriteIfNewer {
		overwrite = sys.CopyOverwriteOnlyIfNewer
	}

	if len(files) == 0 {
		if verbose {
			fmt.Printf("No source files were specified in a copy operation.")
		}
		return nil
	}

	err := sys.CopyFilesEx(dest, overwrite, excludes, files...)

	if err != nil {
		return err
	} else if verbose {
		fmt.Printf("Copied %v to %s\n", args, dest)
	}
	return nil
}