package gitcha_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/djyuhn/gitcha/cmd/gitcha"
	"github.com/djyuhn/gitcha/gittest"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewApp(t *testing.T) {
	t.Parallel()

	t.Run("given directory with a valid repository should return non nil App and nil error and non-nil EntryModel RepoReader", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		basicRepo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		wt, err := basicRepo.Worktree()
		require.NoError(t, err)

		fs := wt.Filesystem

		dirPath := fs.Root()

		app, err := gitcha.NewApp(dirPath)

		assert.NoError(t, err)
		assert.NotNil(t, app)
		assert.NotNil(t, app.TuiModel.RepoReader)
	})

	t.Run("given directory with invalid repository should return nil App and error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		repoDir, _, err := gittest.CreateEmptyRepo(ctx, t)
		require.Error(t, err)

		expectedError := fmt.Errorf("NewApp: directory does not contain a repository")
		app, err := gitcha.NewApp(repoDir)

		assert.ErrorContains(t, err, expectedError.Error())
		assert.Nil(t, app)
	})
}

func TestApp_GitchaTui(t *testing.T) {
	t.Parallel()

	t.Run("given error when executing program should return error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		basicRepo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		wt, err := basicRepo.Worktree()
		require.NoError(t, err)

		fs := wt.Filesystem

		dirPath := fs.Root()

		var buf bytes.Buffer
		var in bytes.Buffer

		app, err := gitcha.NewApp(dirPath, tea.WithInput(&in), tea.WithOutput(&buf))
		require.NoError(t, err)

		go app.TuiProgram.Kill()

		expectedError := fmt.Errorf("GitchaTui: attempted to run program and received an error")
		err = app.GitchaTui()

		assert.ErrorContains(t, err, expectedError.Error())
	})

	t.Run("given program runs without error should return nil error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		basicRepo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		wt, err := basicRepo.Worktree()
		require.NoError(t, err)

		fs := wt.Filesystem

		dirPath := fs.Root()

		var buf bytes.Buffer
		var in bytes.Buffer

		app, err := gitcha.NewApp(dirPath, tea.WithInput(&in), tea.WithOutput(&buf))
		require.NoError(t, err)

		go app.TuiProgram.Send(tea.Quit())

		err = app.GitchaTui()

		assert.NoError(t, err)
	})
}

func TestGetDirectoryFromArgs(t *testing.T) {
	t.Parallel()

	t.Run("given args of length 0 should return working directory path and nil error", func(t *testing.T) {
		t.Parallel()

		expectedDirectory, err := os.Getwd()
		require.NoError(t, err)

		var args []string

		actual, err := gitcha.GetDirectoryFromArgs(args)

		assert.Equal(t, expectedDirectory, actual)
		assert.NoError(t, err)
	})

	t.Run("given an argument that is not a path should return empty string and error", func(t *testing.T) {
		t.Parallel()

		args := []string{"somePath1"}

		expectedDir := ""
		expectedError := fmt.Errorf("GetDirectoryFromArgs: argument %s is not a directory", args[0])
		actual, err := gitcha.GetDirectoryFromArgs(args)

		assert.Equal(t, expectedDir, actual)
		assert.ErrorContains(t, err, expectedError.Error())
	})

	t.Run("given an argument that is a path to a file should return empty string and error", func(t *testing.T) {
		t.Parallel()

		tempDir := t.TempDir()
		file, err := os.Create(filepath.Join(tempDir, "someFile.txt"))
		require.NoError(t, err)

		args := []string{file.Name()}

		expectedDir := ""
		expectedError := fmt.Errorf("GetDirectoryFromArgs: argument %s is not a directory", args[0])
		actual, err := gitcha.GetDirectoryFromArgs(args)

		assert.Equal(t, expectedDir, actual)
		assert.ErrorContains(t, err, expectedError.Error())
	})

	t.Run("given an argument that is a directory should return path and nil error", func(t *testing.T) {
		t.Parallel()

		expectedDir := t.TempDir()

		args := []string{expectedDir}

		actual, err := gitcha.GetDirectoryFromArgs(args)

		assert.Equal(t, expectedDir, actual)
		assert.NoError(t, err)
	})

	t.Run("given multiple arguments and first is not a directory and second is should return empty string and error", func(t *testing.T) {
		t.Parallel()

		tempDir := t.TempDir()

		args := []string{"somePath1", tempDir}

		expectedDir := ""
		expectedError := fmt.Errorf("GetDirectoryFromArgs: argument %s is not a directory", args[0])
		actual, err := gitcha.GetDirectoryFromArgs(args)

		assert.Equal(t, expectedDir, actual)
		assert.ErrorContains(t, err, expectedError.Error())
	})
}
