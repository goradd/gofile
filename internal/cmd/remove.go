// Copyright 2018 Shannon Pekary. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func removeFiles(_ *cobra.Command, _ []string) error {
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
