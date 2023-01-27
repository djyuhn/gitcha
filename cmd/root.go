package cmd

import (
	"os"

	"github.com/djyuhn/gitcha/cmd/gitcha"

	"github.com/spf13/cobra"
)

type RootCmd struct {
	cobra.Command
}

func NewRootCmd() RootCmd {
	return RootCmd{
		Command: cobra.Command{
			Use:     "gitcha [-D dir]",
			Short:   "A command-line tool to get Git information.",
			Long:    "Gitcha is a Git CLI tool to get Git information for repositories.",
			Example: "gitcha",
			Args:    cobra.MaximumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				path, err := gitcha.GetDirectoryFromArgs(args)
				if err != nil {
					return err
				}

				app, err := gitcha.NewApp(path)
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
