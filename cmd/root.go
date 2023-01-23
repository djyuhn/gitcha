package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

type RootCmd struct {
	cobra.Command
}

func NewRootCmd() RootCmd {
	return RootCmd{
		Command: cobra.Command{
			Use:   "gitcha",
			Short: "A command-line tool to get Git information.",
			Long:  "Gitcha is a Git CLI tool to get Git information for repositories.",
		},
	}
}

func Execute() {
	rootCmd := NewRootCmd()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
