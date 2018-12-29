/*
Gofile is a simple file and directory manipulation tool primarily useful for building go applications.

Go is module and GOPATH aware, and is cross-platform. It has subcommands to make directories, copy files, delete files,
and run go generate on files. Files and directories can be specified relative to the location of modules.
Any directory can be represented as a module name, followed by a subdirectory and the real path of the module will be
substituted for the module name.

Gofile exports the ModulePaths() function as a library so you can build your own module aware tools.

For complete documentation of the command-line tool, see the README file.

 */
package main
