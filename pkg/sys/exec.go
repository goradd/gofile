// Copyright 2018 Shannon Pekary. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package sys

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ExecuteShellCommand executes a shell command in the current working directory and returns its output, if any.
// The result is stdOut. If you get an error, you can cast err to (*exec.ExitError) and read the Stderr member to see
// the error message that was generated.
// The command string is the shell command, complete with all arguments.
// To add string that has a space in it, enclose it in single or double quotes. Linux permits backslash escaped spaces, but that
// will not work here, since a backslash represents something else in windows. However, quotes work in all OS's.
// To include a quote character, use the other kind of quote. For example, to include a single quote, surround with double quotes.
// This does not support recursive quotes. For that, you just will need to revert to the exec.Command function.
func ExecuteShellCommand(command string) (result []byte, err error) {
	var parts []string
	parts, err = splitCommandParts(command)
	if len(parts) == 0 || err != nil {
		return
	}

	cmd := exec.Command(parts[0], parts[1:]...)

	result, err = cmd.Output()
	return
}

func splitCommandParts(command string) (parts []string, err error) {
	cur := command
	for cur != "" {
		i := strings.IndexAny(cur, ` '"`)
		if i == -1 {
			parts = append(parts, cur)
			break
		} else if cur[i] == ' ' {
			parts = append(parts, cur[:i])
			cur = cur[i+1:]
		} else {
			lookFor := cur[i : i+1]
			i2 := strings.Index(cur[i+1:], lookFor)
			if i2 == -1 {
				// An error, an unterminated quote
				err = fmt.Errorf("unterminated quote at: %s", cur[i+1:])
				return
			}
			var parts2 []string
			parts2, err = splitCommandParts(cur[i+i2+2:])
			next := cur[:i] + cur[i+1:i+i2+1]
			if len(parts2) == 0 {
				parts = append(parts, next)
				return
			} else {
				next += parts2[0]
				parts = append(parts, next)
			}
			if len(parts2) == 1 {
				return
			}
			parts = append(parts, parts2[1:]...)
			return
		}
	}
	return
}

type pathType11 struct {
	Path string
	Dir  string
}

// ModulePaths returns a listing of the paths of all the modules included in the go.mod file,
// keyed by module name, from the perspective of the
// current working directory.
//
// If we are running without module support, it will return an error.
func ModulePaths() (ret map[string]string, err error) {
	var outText []byte

	outText, err = ExecuteShellCommand("go list -m -json all")

	if err == nil {
		if outText != nil && len(outText) > 0 {
			ret = make(map[string]string)
			dec := json.NewDecoder(bytes.NewReader(outText))
			for {
				var v pathType11
				if err2 := dec.Decode(&v); err2 != nil {
					if err2 == io.EOF {
						return
					}
					return nil, fmt.Errorf("error unpacking json from go list command.\n%s\n%s", string(outText), err.Error())
				}
				ret[v.Path] = v.Dir
			}
		}
		return
	} else {
		// unpack standard error
		stdErr := string(err.(*exec.ExitError).Stderr)
		return nil, fmt.Errorf("error getting module list %s", stdErr)
	}
}

// GetModulePath compares the given path with the list of modules and if the path begins with a module name, it will
// substitute the absolute path for the module name.
// modules is the output from ModulePaths. Module paths always use forward slashes. The resulting
// path uses the native path separator.
//
// If the path is not a module path, the path will be returned unchanged.
func GetModulePath(path string, modules map[string]string) (newPath string, err error) {
	newPath = path
	for modPath, dir := range modules {
		if len(modPath) <= len(path) && path[:len(modPath)] == modPath { // if the path starts with a module path, replace it with the actual directory
			if dir == "" {
				err = fmt.Errorf("module %s is in the cache, but is not installed. Possibly you only installed its application? "+
					"Install the module again using go get -u %[1]s", modPath)
			}
			path = filepath.Join(dir, path[len(modPath):])
			newPath = filepath.FromSlash(path)
			break
		}
	}

	return
}

// ImportPath will convert a path to a go source file to the equivalent module oriented import path.
// This is useful for doing code generation to convert relative paths to absolute module paths, since
// go does not support relative import paths.
// fp can be a directory, in which case it will be treated as a path to a package that you want the import path for.
// This can be a lengthy process, so minimize the number of times you call this.
func ImportPath(fp string) (importPath string, err error) {
	fp, err = filepath.Abs(fp)
	if err != nil {
		return
	}
	var parentDir string
	if IsDir(fp) {
		parentDir = fp
	} else {
		parentDir = filepath.Dir(fp)
	}
	var oldWd string
	oldWd, err = os.Getwd()
	if err != nil {
		return
	}
	err = os.Chdir(parentDir)
	if err != nil {
		return
	}
	defer func() { err = os.Chdir(oldWd) }()

	// Run `go list` to get the module root and path
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Path}}")
	cmd.Dir = parentDir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err = cmd.Run(); err != nil {
		return
	}

	modulePaths := strings.Fields(out.String())

	// Run `go list` again to find the module's directory
	cmd = exec.Command("go", "list", "-m", "-f", "{{.Dir}}")
	cmd.Dir = parentDir
	out.Reset()
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err = cmd.Run(); err != nil {
		return
	}

	moduleDirs := strings.Fields(out.String())

	// Find the dir that matches our own directory
	for i, moduleDir := range moduleDirs {
		if strings.HasPrefix(fp, moduleDir) {
			var relativePath string
			relativePath, err = filepath.Rel(moduleDir, parentDir)
			if err != nil {
				return
			}
			importPath = filepath.ToSlash(filepath.Join(modulePaths[i], relativePath))
			return
		}
	}
	return
}
