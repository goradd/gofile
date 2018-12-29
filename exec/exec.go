package exec

import (
	"github.com/goradd/gofile/internal/cmd"
	"log"
)

func Execute() {
	rootCmd := cmd.MakeRootCommand()

	err := rootCmd.Execute()

	if err != nil {
		log.Fatal(err)
	}
}
