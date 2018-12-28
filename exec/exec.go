package exec

import (
	"github.com/spekary/gofile/internal/cmd"
	"log"
)

func Execute() {
	rootCmd := cmd.MakeRootCommand()

	err := rootCmd.Execute()

	if err != nil {
		log.Fatal(err)
	}
}
