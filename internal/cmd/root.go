// Copyright 2018 Shannon Pekary. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/goradd/gofile/pkg/sys"
	"github.com/spf13/cobra"
	"os"
	"path"
)

var excludes []string
var exclude string
var modules map[string]string
var files []string
var copyOverwrite bool
var copyOverwriteIfNewer bool
var verbose bool

func MakeRootCommand() *cobra.Command {
	var err error

	modules, err = sys.ModulePaths()
	if err != nil {
		panic(err)
	}

	var rootCmd = &cobra.Command{
		Use:   "gofile",
		Short: "gofile is a module-aware, cross-platform, go file manipulation tool",
		Long: `gofile is a module-aware, cross-platform, go file manipulation tool.
After each command, list a file, or group of files to process. In each file description, you can
use standard GLOB identifiers (like * to match any string). If an identifier starts with a module
identifier (e.g. github.com/repo/proj), gofile will look for that file or directory in the module
specified. Environment variables can be specified with $NAME or ${NAME}. Separate paths with
forward slash to be cross-platform compatible.`,
		PersistentPreRun: processExclude,

	}

	rootCmd.PersistentFlags().StringVarP(&exclude, "exclude", "x", "", "list of pattern match expressions, separated by semicolons or colons, that when matched, will be excluded from the list of files to process. The pattern match is the same as file GLOB pattern matching.")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	var cmdRemove = &cobra.Command{
		Use:   "remove [files to remove]",
		Short: "Deletes the given files.",
		Long: `Deletes the listed files and directories permanently. Use with care.`,
		Args: cobra.MinimumNArgs(1),
		PreRun: processFileListArgs,
		RunE: removeFiles,
	}

	var cmdGenerate = &cobra.Command{
		Use:   "generate [files to hand off to go generate]",
		Short: "go generate the given files.",
		Long: `Passes the given files to go generate.`,
		Args: cobra.MinimumNArgs(1),
		PreRun: processFileListArgs,
		RunE: generateFiles,
	}

	var cmdCopy = &cobra.Command{
		Use:   "copy [files or directories to copy] [destination file or directory]",
		Short: "Copy the files to the given location.",
		Long: `Copies the files or directories to the given location. If you are copying one
file, the destination can be a file name that does not exist, but whose parent exists. If 
copying more than one file, the destination must be a directory that exists.`,
		Args: cobra.MinimumNArgs(2),
		RunE: copyFiles,
	}
	cmdCopy.Flags().BoolVarP(&copyOverwrite, "overwrite", "o", false, "Files will overwrite previous files when copying.")
	cmdCopy.Flags().BoolVarP(&copyOverwriteIfNewer, "newer", "n", false, "Files will overwrite previous files when copying only if the new file is newer than the old.")

	var cmdMkDir = &cobra.Command{
		Use:   "mkdir [directory to create]",
		Short: "Create the given directory.",
		Long: `Create the given directory.`,
		Args: cobra.MinimumNArgs(1),
		PreRun: processFileListArgs,
		RunE: mkDir,
	}

	rootCmd.AddCommand(cmdRemove, cmdGenerate, cmdCopy, cmdMkDir)

	return rootCmd
}

func processExclude(cmd *cobra.Command, args []string) {
	exclude = os.ExpandEnv(exclude)
	excludes = sys.SplitList(exclude)
}

// processFileListArgs accepts the group of arguments that would represent files, directories
// etc., processes them, removes excluded files, and sets the files global to this list
func processFileListArgs(cmd *cobra.Command, args []string) {
	files2 := sys.ModuleExpandFileList(args, modules)

	if excludes == nil || files == nil {
		files = files2
		return
	}

	files = nil
	for _, f := range files2 {
		for _,e := range excludes {
			m,_ := path.Match(e, f)
			if !m {
				files = append(files, f)
			}
		}
	}
}

func processFileArg(arg string) string {
	arg = os.ExpandEnv(arg)
	arg, _ = sys.GetModulePath(arg, modules)
	return arg
}


