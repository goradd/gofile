package exec

import (
	"fmt"
	"github.com/goradd/gofile/internal/cmd"
	"os"
)

func Execute() {
	rootCmd,err := cmd.MakeRootCommand()

	if err != nil {
		_,_ = fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	err = rootCmd.Execute()

	if err != nil {
		_,_ = fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)	}
}
