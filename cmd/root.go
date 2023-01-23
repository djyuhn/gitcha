package cmd

import (
	"os"

	"gitcha/cmd/gitcha"
	"gitcha/tui"

	"github.com/spf13/cobra"
)

type RootCmd struct {
	cobra.Command
}

func NewRootCmd() RootCmd {
	return RootCmd{
		Command: cobra.Command{
			Use:     "gitcha",
			Short:   "A command-line tool to get Git information.",
			Long:    "Gitcha is a Git CLI tool to get Git information for repositories.",
			Example: "gitcha",
			RunE: func(cmd *cobra.Command, args []string) error {
				app, err := gitcha.NewApp(tui.EntryModel{})
				if err != nil {
					return err
				}

				if err := app.GitchaTui(); err != nil {
					return err
				}

				return nil
			},
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
