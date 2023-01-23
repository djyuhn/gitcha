package cmd_test

import (
	"testing"

	"gitcha/cmd"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestNewRootCmd(t *testing.T) {
	t.Parallel()

	t.Run("should return cobra command with gitcha descriptions", func(t *testing.T) {
		t.Parallel()

		expected := cmd.RootCmd{
			Command: cobra.Command{
				Use:   "gitcha",
				Short: "A command-line tool to get Git information.",
				Long:  "Gitcha is a Git CLI tool to get Git information for repositories.",
			},
		}

		actual := cmd.NewRootCmd()

		assert.Equal(t, expected, actual)
	})
}
