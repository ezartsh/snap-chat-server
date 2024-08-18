package cmd

import (
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Join(filepath.Dir(b), "..")

	rootCmd = &cobra.Command{
		Use:   "snap-cli",
		Short: "Cli for chat application",
		Long:  `This is a simple cli for chat.`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
