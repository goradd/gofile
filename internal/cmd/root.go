// Copyright 2018 Shannon Pekary. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/goradd/gofile/pkg/sys"
	"github.com/spf13/cobra"
)

var excludes []string
var exclude string
var modules map[string]string
var files []string
var copyOverwrite bool
var copyOverwriteIfNewer bool
var verbose bool
var deleteAfterZip bool
var gzipCompressionLevel int
var brotliCompressionLevel int

// MakeRootCommand creates the command tree for cobra.
func MakeRootCommand() (*cobra.Command, error) {
	var err error

	modules, err = sys.ModulePaths()
	if err != nil {
		return nil, err
	}

	var rootCmd = &cobra.Command{
		Use:   "gofile",
		Short: "gofile is a module-aware, cross-platform, go file manipulation tool",
		Long: `gofile is a module-aware, cross-platform, go file manipulation tool.
After each command, specify a file or group of files to process. In each file description, you can
use standard GLOB identifiers (like * to match any string). If an identifier starts with a module
identifier (e.g. github.com/repo/proj), gofile will look for that file or directory in the module
specified. Environment variables can be specified with $NAME or ${NAME}. Separate paths with
forward slash to be cross-platform compatible.`,
		PersistentPreRun: processExclude,
	}

	rootCmd.PersistentFlags().StringVarP(&exclude, "exclude", "x", "", "list of pattern match expressions, separated by semicolons or colons, that when matched, will be excluded from the list of files to process. The pattern match is the same as file GLOB pattern matching.")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	var cmdRemove = &cobra.Command{
		Use:    "remove [files to remove]",
		Short:  "Deletes the given files.",
		Long:   `Deletes the listed files and directories permanently. Use with care.`,
		Args:   cobra.MinimumNArgs(1),
		PreRun: processFileListArgs,
		RunE:   removeFiles,
	}

	var cmdGenerate = &cobra.Command{
		Use:    "generate [files to hand off to go generate]",
		Short:  "go generate the given files.",
		Long:   `Passes the given files to go generate.`,
		Args:   cobra.MinimumNArgs(1),
		PreRun: processFileListArgs,
		RunE:   generateFiles,
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
		Use:    "mkdir [directory to create]",
		Short:  "Create the given directory.",
		Long:   `Create the given directory.`,
		Args:   cobra.MinimumNArgs(1),
		PreRun: processFileListArgs,
		RunE:   mkDir,
	}

	var cmdGZip = &cobra.Command{
		Use:    "gzip [files or directories to zip]",
		Short:  "GZip the given files or directories.",
		Long:   `GZips the given files, or all the files in the specified directories, placing zipped files alongside the given files, with .gz suffixes. Uses the maximum compression algorithm.`,
		Args:   cobra.MinimumNArgs(1),
		PreRun: processExpandedFileListArgs,
		RunE:   gzip,
	}
	cmdGZip.Flags().BoolVarP(&deleteAfterZip, "delete", "d", false, "Compressed source files will be deleted, leaving only the compressed version.")
	cmdGZip.Flags().IntVarP(&gzipCompressionLevel, "quality", "q", 9, "The compression level to use. Higher numbers offer higher compression and slower compression speed, but have negligible effect on decompression speed.")

	var cmdBrotli = &cobra.Command{
		Use:    "brotli [files or directories to compress]",
		Short:  "Brotli compress the given files or directories.",
		Long:   `Compresses the given files with the Brotli method, or all the files in the specified directories, placing compressed files alongside the given files, with .br suffixes.`,
		Args:   cobra.MinimumNArgs(1),
		PreRun: processExpandedFileListArgs,
		RunE:   brotli,
	}
	cmdBrotli.Flags().BoolVarP(&deleteAfterZip, "delete", "d", false, "Compressed source files will be deleted, leaving only the compressed version.")
	cmdBrotli.Flags().IntVarP(&brotliCompressionLevel, "quality", "q", 11, "The compression level to use. Higher numbers offer higher compression and slower compression speed, and have negligible effect on decompression speed.")

	var cmdPath = &cobra.Command{
		Use:   "path [path to convert]",
		Short: "Converts a module relative path to its absolute path.",
		Long:  `Converts a module relative path to its absolute path and sends it to stdout.`,
		Args:  cobra.ExactArgs(1),
		RunE:  outPath,
	}

	rootCmd.AddCommand(cmdRemove, cmdGenerate, cmdCopy, cmdMkDir, cmdGZip, cmdBrotli, cmdPath)

	return rootCmd, nil
}

func processExclude(_ *cobra.Command, _ []string) {
	exclude = os.ExpandEnv(exclude)
	excludes = sys.SplitList(exclude)
}

// processFileListArgs accepts the group of arguments that would represent files, directories
// etc., processes them, removes excluded files, and sets the files global to this list
// non-existent names are left intact so that we can create them.
func processFileListArgs(_ *cobra.Command, args []string) {
	files2 := sys.ModuleExpandFileList(args, modules)

	if excludes == nil {
		files = files2
		return
	}

	files = nil
Files2:
	for _, f := range files2 {
		for _, e := range excludes {
			m, _ := filepath.Match(e, filepath.Base(f))
			if m {
				continue Files2
			}
		}
		files = append(files, f)
	}
}

// processExistingFileListArgs accepts the group of arguments that would represent files, directories
// etc., expands the list based on the current modules, expands directories to the list of files in those
// directories, removes excluded files, and sets the files global to this list.
// non-existent names are removed
func processExpandedFileListArgs(_ *cobra.Command, args []string) {
	files2 := sys.ModuleExpandFileList(args, modules)

	files = nil

	for _, f := range files2 {
		if sys.IsDir(f) {
			_ = filepath.WalkDir(f, func(path string, d fs.DirEntry, err error) error {
				if err == nil {
					if d.IsDir() {
						if isExcluded(path) {
							return filepath.SkipDir
						}
					} else {
						if !isExcluded(path) {
							files = append(files, path)
						}
					}
				}
				return nil
			})
		} else {
			if isExcluded(f) {
				continue
			}
			files = append(files, f)
		}
	}
}

// isExcluded returns true if the given file matches one of the exclusion strings
func isExcluded(file string) bool {
	for _, e := range excludes {
		m, _ := filepath.Match(e, filepath.Base(file))
		if m {
			return true
		}
	}
	return false
}

func processFileArg(arg string) string {
	arg = os.ExpandEnv(arg)
	arg2, _ := sys.GetModulePath(arg, modules)
	arg2 = filepath.FromSlash(arg2)

	if verbose && arg2 != arg {
		fmt.Printf("module path " + arg + " found at " + arg2)
	}
	return arg2
}
