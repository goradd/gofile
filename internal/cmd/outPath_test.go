package cmd

import (
	"bytes"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/spf13/cobra"
)

func Test_outPath(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	file = filepath.FromSlash(file)

	var cmd *cobra.Command
	cmd, _ = MakeRootCommand()
	cmd.SetArgs([]string{"path", "github.com/goradd/gofile/internal/cmd/outPath_test.go"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	_ = cmd.Execute()
	if buf.String() != file {
		t.Error("Files do not match")
	}
}
