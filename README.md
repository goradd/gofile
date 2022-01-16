# gofile

Gofile is a simple file and directory manipulation tool primarily useful for building go applications and libraries. 

Go is module aware, and is cross-platform. Any directory can be represented as a module name,
followed by a subdirectory and the real path of the module will be substituted for 
the module name. Path names should use forward slashes to separate directories, and you can substitute 
environment variables into the path names using $VAR, or ${VAR} syntax.

Gofile is particularly useful for making build scripts for cross-platform projects. The cross-platform feature allows
you to create command-line scripts that will work on Windows and Unix based systems. The module aware feature
allows you to specify a path relative to a module's location, simply by starting the path with the module specifier.

If you put "replace" statements in the go.mod file, gofile will honor those too. 
For example, if your go.mod file has the following replace statement:

`
replace github.com/myproj/proj => /proj-src
`

and you execute the following command from within the source tree of the go.mod file:

`
gofile copy github.com/myproj/proj/README.md /a/b/
`

gofile will copy the README.md file from /proj-src to the /a/b/ directory.

Gofile exports the ModulePaths() function as a library so you can build your own module aware tools. 

## Installation

```shell
go get -u github.com/goradd/gofile/...
```

## Usage

```shell
gofile <command> <args...> [options] 
```

## Commands

With most of the following commands, -v will output status information while gofile is running.

-x specifies what to exclude when expanding a path that uses * or ? expansions. Note that the * pattern
does not match the path separator, so it only expands one level deep. The exclude pattern can include wild
cards like * and ?, but it will only compare against the last item in a path. For example, if you have
the following directories in the /tmp directory:

- test1
- test2

Then specifying `-x test2 /tmp/*` will result in only /tmp/test1 being used.

-o will force 

### Help
```shell
gofile -h
```

will print a complete description of the options and arguments of gofile.

### Copy
Copies a file or directory to another file or directory.

Usage:
```shell
gofile copy [-x excludes...] [-o|-n] <src> <dest> 
```
-x specifies names of files or directories you want to exclude from the source. This will match
file patterns as well. This is useful when copying whole directories but needing to exclude specific
types of files from the process.

Normally, a previously existing file will not be overwritten, and a previously existing directory will not be
deleted first but rather files will be added that are not present in the directory. The -o option will force
an overwrite of previously existing files, and the -n option will overwrite only if the new file is newer than
the old one. If you want to replace a previously existing directory, use the remove command described below first.

### Generate

Runs go generate on the given file.

Usage:
```shell
gofile generate [-x excludes...] <sources...>

```

-x specifies names of files or directories you want to exclude from the source. This is useful when
expanding a directory using '*'.

### Mkdir

Creates the named directory if it does not exist. Sets it to be writable.

Usage:
```shell
gofile mkdir <dest>
```

### Remove

Deletes the named directories or files.

Usage:
```shell
gofile remove [-x excludes...] <dest...> 
```

-x specifies names of files or directories you want to exclude from the destination. This is useful when
expanding a directory using '*'.

### GZip

Compresses the given files using the GZip method.

If a directory is specified, then the files inside that directory are individually
compressed. This will recursively do the same to directories within the specified directory.

Compressed files are placed in the same directory as the source file, and the
file name is appended with ".gz". If a file already has a ".gz" extension, gofile will
assume the file is already compressed and skip it.

Usage:
```shell
gofile gzip [-d] [-q level] [-x excludes...] <dest...> 
```

-x specifies names of files or directories you want to exclude from compression. For exmaple, 
"-x *.txt" will prevent all files ending in ".txt" from being compressed. Specifying the name
of a directory will exclude all the files within that directory.

-d will delete the source file after being compressed. Excluded files specified in the 
-x option are not removed.

-q specifies the compression level. Default is 9, which is the maximum.

### Brotli

Compresses the given files using the Brotli method.

If a directory is specified, then the files inside that directory are individually
compressed. This will recursively do the same to directories within the specified directory.

Compressed files are placed in the same directory as the source file, and the
file name is appended with ".br". If a file already has a ".br" extension, gofile will
assume the file is already compressed and skip it.

Usage:
```shell
gofile brotli [-d] [-q level] [-x excludes...] <dest...> 
```

-x specifies names of files or directories you want to exclude from compression. For exmaple,
"-x *.txt" will prevent all files ending in ".txt" from being compressed. Specifying the name
of a directory will exclude all the files within that directory.

-d will delete the source file after being compressed. Excluded files specified in the
-x option are not removed.

-q specifies the compression level. Default is 11, which is the maximum.
