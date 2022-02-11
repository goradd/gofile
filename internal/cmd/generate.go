// Copyright 2018 Shannon Pekary. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"github.com/goradd/gofile/pkg/sys"
	"github.com/spf13/cobra"
)

var generateResult []byte

func generateFiles(_ *cobra.Command, _ []string) error {
	for _, f := range files {
		var err error
		generateResult, err = sys.ExecuteShellCommand("go generate " + f)
		if err != nil {
			return err
		} else if verbose {
			fmt.Printf("Generated %s: %s\n", f, string(generateResult))
		}
	}
	return nil
}
