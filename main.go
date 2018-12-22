package main

import (
	"github.com/spekary/gofile/internal/cmd"
	"log"
)


func main() {
	rootCmd := cmd.MakeRootCommand()

	err := rootCmd.Execute()

	if err != nil {
		log.Fatal(err)
	}
}

