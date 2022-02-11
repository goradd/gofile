// Copyright 2018 Shannon Pekary. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func outPath(cmd *cobra.Command, args []string) error {
	f := processFileArg(args[0])
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), f)
	return nil
}
