// Copyright 2018 Shannon Pekary. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package sys

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// SplitList will split a file and/or directory list into individual items in a cross-platform way. In
// other words, the list can be specified in a Unix(colons) or Windows(semicolon) friendly way, and
// it will split the list regardless of the platform the list was designed for, or run on. This means on
// Unix, you can't use a semicolon in a file name, and on windows, you can't use a colon, but that should
// not be an issue for most people.
func SplitList(s string) (list []string) {
	s = strings.Replace(s, ";", ":", -1)
	for _, item := range strings.Split(s, ":") {
		if item != "" {
			list = append(list, item)
		}
	}
	return
}

// ModuleExpandFileList will do the following given a list of arguments that represent command line arguments
// that would be a list of
// files, directories, or glob patterns:
//   replace any environment variables with their values
//   replace any items that start with a module with the actual location on disk
//   expand any glob patterns
//   remove duplicates
// It will return the final list.
//
// modules is the list of modules returned from ModulePaths.
//
// Glob patterns will be matched, and if nothing is found, no file will be generated.
// However, if a path does not have a glob pattern, but does not exist, it will be left in the list,
// since it might refer to a file or directory you want to add.
// The list out is not necessarily in the same order as the list in.
func ModuleExpandFileList(args []string, modules map[string]string) (list []string) {
	files2 := make(map[string]bool)

	for _, arg := range args {
		arg = os.ExpandEnv(arg)
		arg, _ = GetModulePath(arg, modules)
		var files []string
		if hasMeta(arg) {
			files, _ = filepath.Glob(arg)
		} else {
			files = append(files, arg)
		}

		for _, f := range files {
			files2[f] = true
		}
	}
	for k := range files2 {
		list = append(list, k)
	}
	return
}

// hasMeta reports whether path contains any of the magic characters
// recognized by Match.
// This is an copied from the unexported function from filepath.
func hasMeta(path string) bool {
	magicChars := `*?[`
	if runtime.GOOS != "windows" {
		magicChars = `*?[\`
	}
	return strings.ContainsAny(path, magicChars)
}

// CopyOverwriteType is used by the CopyFiles function to determine how to
// treat file collisions when copying over a file that already exists.
type CopyOverwriteType int

const (
	// CopyDoNotOverwrite will prevent copying if the destination file already exists
	CopyDoNotOverwrite CopyOverwriteType = 0
	// CopyOverwrite will cause the destination to be replaced if it exists
	CopyOverwrite = 1
	// CopyOverwriteOnlyIfNewer will cause the destination to be overwritten only if the source is newer than the destination.
	CopyOverwriteOnlyIfNewer = 2
)

// CopyFiles copies the src files or directories to the destination
//
// If there is more than one source, the destination must be a directory that exists. The items listed
// will be copied inside the destination directory.
//
// If there is only one source, the destination must be:
//   - A directory that exists, in which case the source will be placed in the destination directory
//   - A file that exists, in which case the source will overwrite the destination. The source must also be a single file.
//	 - A file that does not exist, but whose parent directory does exist, in which case the file will be copied
//     and renamed to the destination.
// If overwrite is true, files that already exist will be overwritten. If overwrite is false, only new files
// will be created. If a directory is over-writing another directory, this will determine what happens when
// file names are duplicates. Note that old files in a directory will not be deleted when a directory
// overwrites another directory. If you want old files to be deleted, empty the destination directory first.
func CopyFiles(dst string, overwrite CopyOverwriteType, src ...string) (err error) {
	// Sanity checks
	if dst == "" {
		return fmt.Errorf("no destination specified")
	}

	if len(src) == 0 {
		return fmt.Errorf("no source files specified")
	}

	dstInfo, destErr := os.Stat(dst)
	srcInfo, srcErr := os.Stat(src[0])

	if srcErr != nil {
		return fmt.Errorf("error with source: %s", srcErr.Error())
	}

	if len(src) > 1 || srcInfo.IsDir() {
		if destErr != nil {
			return destErr // path doesn't exist?
		}
		if !dstInfo.IsDir() {
			return fmt.Errorf("when copying multiple files, the destination must be a directory: %s", dst)
		}

		for _, f := range src {
			err = copyTo(f, dst, "", overwrite)
			if err != nil {
				return
			}
		}
	} else {
		if os.IsPathSeparator(dst[len(dst)-1]) {
			// Definitely trying to point to a directory
			if os.IsNotExist(destErr) {
				return fmt.Errorf("the destination directory does not exist: %s", dst)
			}
			err = copyTo(src[0], dst, "", overwrite)
			if err != nil {
				return
			}
		} else {
			// might be a destination directory, or a file
			if os.IsNotExist(destErr) {
				// Since it doesn't exist, we are going to assume we are trying to specify a file, since
				// we have already checked above to see if we are trying to specify a directory with a slash at end.

				// Check on parent directory
				parentDir, fileName := filepath.Split(dst)
				_, parentErr := os.Stat(parentDir)
				if parentErr != nil {
					return fmt.Errorf("the parent directory of a new file must exist: %s", dst)
				}
				// We are writing to a new file
				err = copyFileTo(src[0], parentDir, fileName, overwrite)
				if err != nil {
					return
				}
			} else {
				// destination is a file or a directory that already exists
				if dstInfo.IsDir() {
					err = copyTo(src[0], dst, "", overwrite)
					if err != nil {
						return
					}
				} else {
					parentDir, fileName := filepath.Split(dst)
					err = copyTo(src[0], parentDir, fileName, overwrite)
					if err != nil {
						return
					}
				}
			}
		}
	}
	return
}

// copyFileTo copies the given file to the destination directory.
//If a name is given, it will rename the file.
// It does no checks to see if the destination directory exists.
// overwrite would prevent the file from being copied, then the copy does not happen and
// error is nil.
func copyFileTo(src string, destDir string, name string, overwrite CopyOverwriteType) error {
	var count int64

	srcInfo, srcErr := os.Stat(src)
	if srcErr != nil {
		return srcErr
	}
	if srcInfo.IsDir() {
		return fmt.Errorf(`source "%s"" is not a file`, src)
	}
	var perm os.FileMode

	perm = srcInfo.Mode() & os.ModePerm

	if name == "" {
		name = filepath.Base(src)
	}
	destName := filepath.Join(destDir, name)
	destInfo, destErr := os.Stat(destName)
	if destErr == nil {
		// destination exists
		if overwrite == CopyDoNotOverwrite {
			return nil
		} else if overwrite == CopyOverwriteOnlyIfNewer {
			modSrc := srcInfo.ModTime()
			modDest := destInfo.ModTime()

			if modSrc.Before(modDest) || modSrc.Equal(modDest) {
				return nil
			}
		}
		// prepare for copy by deleting in case permissions are different
		if err := os.Remove(destName); err != nil {
			return err
		}
	}

	from, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		_ = from.Close()
	}()

	to, err := os.OpenFile(destName, os.O_RDWR|os.O_CREATE, perm)
	if err != nil {
		return err
	}

	defer func() {
		_ = to.Close()
	}()

	count, err = io.Copy(to, from)
	if err != nil {
		//to.Close()
		return err
	}
	err = to.Truncate(count) // chop end of file in case file gets smaller
	if err != nil {
		return err
	}

	return to.Close()
}

// copyTo copies the src to the destination directory. The source can be a file or directory.
// if a name is specified, src must be a file. The name will be the name of the file in the new directory.
func copyTo(src string, destDir string, name string, overwrite CopyOverwriteType) error {
	srcInfo, srcErr := os.Stat(src)

	if srcErr != nil {
		return srcErr
	}

	if srcInfo.IsDir() && name != "" {
		if name != filepath.Base(src) {
			return fmt.Errorf("cannot copy a directory to a file")
		}
		p := filepath.Join(destDir, name)
		dstInfo, dstErr := os.Stat(p)
		if (dstErr == nil || !os.IsNotExist(dstErr)) &&
			!dstInfo.IsDir() {
			return fmt.Errorf("cannot copy a directory onto a file that already exists: %s", p)
		}
	}

	if !srcInfo.IsDir() {
		return copyFileTo(src, destDir, name, overwrite)
	}

	// src is a directory, and destination is a directory
	return CopyDirectory(src, destDir, overwrite)
}

// CopyDirectory copies the src directory to the destination directory. The destination directory will be the parent of
// the resulting directory, and the result will have the same name as the source. If the destination already exists,
// it will perform a kind of merge, where existing files will not be touched, and only new files will be copied.
// If you want to replace the destination, delete it first. dst must exist.
func CopyDirectory(src, dst string, overwrite CopyOverwriteType) (err error) {
	dstInfo, dstErr := os.Stat(dst)
	_, srcErr := os.Stat(src)

	if srcErr != nil {
		return fmt.Errorf("source directory error: %s", srcErr.Error())
	}

	if dstErr != nil {
		return fmt.Errorf("destination directory error: %s", dstErr.Error())
	}

	dirDest := filepath.Dir(dst)
	if len(src) <= len(dirDest) && dirDest[:len(src)] == src { // does dst start with src?
		return fmt.Errorf("destination directory is not allowed to be in the src directory")
	}

	if !dstInfo.Mode().IsDir() {
		return fmt.Errorf("source %s is a file, not a directory", dst)
	}

	// create destination if needed
	newPath := filepath.Join(dst, filepath.Base(src))

	newInfo, err := os.Stat(newPath)
	if err == nil || !os.IsNotExist(err) {
		// path exists
		if !newInfo.IsDir() {
			return fmt.Errorf("path %s is a directory in the source, but %s is a file in the destination", src, newPath)
		}
	} else {
		//perm := srcInfo.Mode().Perm()	// copy the permission
		err = os.Mkdir(newPath, 0755)
		if err != nil {
			return fmt.Errorf("error creating directory %s: %s", newPath, err.Error())
		}
	}

	f, err := os.Open(src)
	if err != nil {
		return err
	}
	list, err := f.Readdir(-1)
	_ = f.Close()

	for _, item := range list {
		itemName := item.Name()
		itemPath := filepath.Join(src, itemName)
		err = copyTo(itemPath, newPath, itemName, overwrite)
		if err != nil {
			return err
		}
	}

	return
}

// CopyDirectoryEx copies the src directory to the destination directory excluding files that match a specified pattern.
//
// The destination directory will be the parent of
// the resulting directory, and the result will have the same name as the source. If the destination already exists,
// it will perform a kind of merge, where existing files will not be touched, and only new files will be copied.
// If you want to replace the destination, delete it first. dst must exist.
func CopyDirectoryEx(src, dst string, overwrite CopyOverwriteType, excludes []string) (err error) {
	dstInfo, dstErr := os.Stat(dst)
	_, srcErr := os.Stat(src)

	if srcErr != nil {
		return fmt.Errorf("source directory error: %s", srcErr.Error())
	}

	if dstErr != nil {
		return fmt.Errorf("destination directory error: %s", dstErr.Error())
	}

	dirDest := filepath.Dir(dst)
	if len(src) <= len(dirDest) && dirDest[:len(src)] == src { // does dst start with src?
		return fmt.Errorf("destination directory is not allowed to be in the src directory")
	}

	if !dstInfo.Mode().IsDir() {
		return fmt.Errorf("source %s is a file, not a directory", dst)
	}

	// create destination if needed
	newPath := filepath.Join(dst, filepath.Base(src))

	newInfo, err := os.Stat(newPath)
	if err == nil || !os.IsNotExist(err) {
		// path exists
		if !newInfo.IsDir() {
			return fmt.Errorf("path %s is a directory in the source, but %s is a file in the destination", src, newPath)
		}
	} else {
		//perm := srcInfo.Mode().Perm()	// copy the permission
		err = os.Mkdir(newPath, 0755)
		if err != nil {
			return fmt.Errorf("error creating directory %s: %s", newPath, err.Error())
		}
	}

	f, err := os.Open(src)
	if err != nil {
		return err
	}
	list, err := f.Readdir(-1)
	_ = f.Close()

	for _, item := range list {
		itemName := item.Name()
		if isExcluded(itemName, excludes) {
			continue // skip
		}

		itemPath := filepath.Join(src, itemName)
		err = copyToEx(itemPath, newPath, itemName, overwrite, excludes)
		if err != nil {
			return err
		}
	}

	return
}

// copyTo copies the src to the destination directory. The source can be a file or directory.
// if a name is specified, src must be a file. The name will be the name of the file in the new directory.
func copyToEx(src string, destDir string, name string, overwrite CopyOverwriteType, excludes []string) error {
	srcInfo, srcErr := os.Stat(src)

	if srcErr != nil {
		return srcErr
	}

	if srcInfo.IsDir() && name != "" {
		if name != filepath.Base(src) {
			return fmt.Errorf("cannot copy a directory to a file")
		}
		p := filepath.Join(destDir, name)
		dstInfo, dstErr := os.Stat(p)
		if (dstErr == nil || !os.IsNotExist(dstErr)) &&
			!dstInfo.IsDir() {
			return fmt.Errorf("cannot copy a directory onto a file that already exists: %s", p)
		}
	}

	if !srcInfo.IsDir() {
		if isExcluded(src, excludes) {
			return nil // skip file
		}
		return copyFileTo(src, destDir, name, overwrite)
	}

	// src is a directory, and destination is a directory
	return CopyDirectoryEx(src, destDir, overwrite, excludes)
}

// CopyFilesEx copies the src files or directories to the destination excluding files matching the exclusions slice.
//
// If there is more than one source, the destination must be a directory that exists. The items listed
// will be copied inside the destination directory.
//
// If there is only one source, the destination must be:
//   - A directory that exists, in which case the source will be placed in the destination directory
//   - A file that exists, in which case the source will overwrite the destination. The source must also be a single file.
//	 - A file that does not exist, but whose parent directory does exist, in which case the file will be copied
//     and renamed to the destination.
// If overwrite is true, files that already exist will be overwritten. If overwrite is false, only new files
// will be created. If a directory is over-writing another directory, this will determine what happens when
// file names are duplicates. Note that old files in a directory will not be deleted when a directory
// overwrites another directory. If you want old files to be deleted, empty the destination directory first.
func CopyFilesEx(dst string, overwrite CopyOverwriteType, exclusions []string, src ...string) (err error) {
	// Sanity checks
	if dst == "" {
		return fmt.Errorf("no destination specified")
	}

	if len(src) == 0 {
		return fmt.Errorf("no source files specified")
	}

	dstInfo, destErr := os.Stat(dst)
	srcInfo, srcErr := os.Stat(src[0])

	if srcErr != nil {
		return fmt.Errorf("error with source: %s", srcErr.Error())
	}

	if len(src) > 1 || srcInfo.IsDir() {
		if destErr != nil {
			return destErr // path doesn't exist?
		}
		if !dstInfo.IsDir() {
			return fmt.Errorf("when copying multiple files, the destination must be a directory: %s", dst)
		}

		for _, f := range src {
			err = copyToEx(f, dst, "", overwrite, exclusions)
			if err != nil {
				return
			}
		}
	} else {
		if os.IsPathSeparator(dst[len(dst)-1]) {
			// Definitely trying to point to a directory
			if os.IsNotExist(destErr) {
				return fmt.Errorf("the destination directory does not exist: %s", dst)
			}
			err = copyToEx(src[0], dst, "", overwrite, exclusions)
			if err != nil {
				return
			}
		} else {
			// might be a destination directory, or a file
			if os.IsNotExist(destErr) {
				// Since it doesn't exist, we are going to assume we are trying to specify a file, since
				// we have already checked above to see if we are trying to specify a directory with a slash at end.

				// Check on parent directory
				parentDir, fileName := filepath.Split(dst)
				_, parentErr := os.Stat(parentDir)
				if parentErr != nil {
					return fmt.Errorf("the parent directory of a new file must exist: %s", dst)
				}
				// We are writing to a new file
				if isExcluded(src[0], exclusions) {
					return nil
				}
				err = copyFileTo(src[0], parentDir, fileName, overwrite)
				if err != nil {
					return
				}
			} else {
				// destination is a file or a directory that already exists
				if dstInfo.IsDir() {
					err = copyToEx(src[0], dst, "", overwrite, exclusions)
					if err != nil {
						return
					}
				} else {
					parentDir, fileName := filepath.Split(dst)
					if isExcluded(src[0], exclusions) {
						return nil
					}
					err = copyTo(src[0], parentDir, fileName, overwrite)
					if err != nil {
						return
					}
				}
			}
		}
	}
	return
}

// IsDir returns true if the given path exists and is a directory
func IsDir(path string) bool {
	dstInfo, err := os.Stat(path)
	if err == nil { // file exists
		if dstInfo.Mode().IsDir() {
			return true
		}
	}
	return false
}

// isExcluded returns true if the given file matches one of the exclusion strings
func isExcluded(file string, excludes []string) bool {
	for _, e := range excludes {
		m, _ := filepath.Match(e, filepath.Base(file))
		if m {
			return true
		}
	}
	return false
}
