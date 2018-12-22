# gofile

Gofile is a simple file and directory manipulation tool primarily useful for building go applications. 

Go is module and GOPATH aware, and is cross-platform. Any directory can be represented as a module name,
followed by a subdirectory and the real path of the module will be substituted for 
the module name. Path names should use forward slashes to separate directories, and you can substitute 
environment variables into the path names using $VAR, or ${VAR} syntax.

Gofile is particularly useful for making build scripts for open-source projects. The cross-platform feature allows
you to create command-line scripts that will work on Windows and Unix based systems. The module aware feature
allows you to specify a path relative to a modules location, simply by starting the path with the module specifier.
If you are not working in Module-aware mode, it will use the GOPATH to locate packages instead.

If you are using modules, there may be times when you would like to work on packages and modules that are in your
go.mod file's require statements. However, in module-aware mode, these files are located in your GOPATH and the 
GOPATH is write-protected, so you cannot edit those files. The solution is to use a replace statement in your go.mod
file to temporarily point the module to different location on your disk. For example, the following replace statement
will change where the go build system will look for the gofile source:

`
replace github.com/spekary/gofile => ../gofile-src // to work with a local version of goradd
`

However, this can pose other problems if your build system is copying files out of one of these modules. If you use
the GOPATH environment variable in your build scripts to locate a module, it will not know about the change you
made in the go.mod file. Gofile solves this problem by substituting a module name for its disk location. For example,

`
gofile copy github.com/spekary/gofile/README.md /a/b/
`
will copy this README file to the /a/b/ directory on your disk, even if you have changed the location of the gofile
module using a replace statement. Gofile looks for the go.mod file by searching for it in the current working directory,
and parent directories until it finds one. Note that this is standard behavior of the go build tools and not under
gofile's control. 

Gofile exports the ModulePaths() function as a library so you can build your own module aware tools. 

Usage:
```bash
gofile <command> <args...> [options] 
```

## Commands
### Help
```bash
gofile -h
```

will print a complete description of the options and arguments of gofile.

### Copy
Copies a file or directory to another file or directory.

Usage:
```bash
gofile copy <src> <dest> [-x excludes...]
```
-x specifies names of files or directories you want to exclude from the source. This is useful when
expanding a directory using '*'.

### Generate

Runs go generate on the given file.

Usage:
```bash
gofile generate <sources...> [-x excludes...]

```

-x specifies names of files or directories you want to exclude from the source. This is useful when
expanding a directory using '*'.

### Mkdir

Creates the named directory if it does not exist. Sets it to be writable.

Usage:
```bash
gofile mkdir <dest>
```

### Remove

Deletes the named directories or files.

Usage:
```bash
gofile remove <dest...> [-x excludes...]
```

-x specifies names of files or directories you want to exclude from the destination. This is useful when
expanding a directory using '*'.

-v will output status information while gofile is running.

