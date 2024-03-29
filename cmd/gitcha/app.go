package gitcha

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/djyuhn/gitcha/internal/reporeader"
	"github.com/djyuhn/gitcha/internal/tui"
)

type App struct {
	TuiModel   tui.EntryModel
	TuiProgram tea.Program
}

func NewApp(repoDirPath string, opts ...tea.ProgramOption) (*App, error) {
	repoReader, err := reporeader.NewRepoReader(repoDirPath)
	if err != nil {
		return nil, fmt.Errorf("NewApp: directory does not contain a repository: %w", err)
	}

	entryModel, err := tui.NewEntryModel(repoReader)
	if err != nil {
		return nil, fmt.Errorf("NewApp: error during creation of tui model: %w", err)
	}

	program := tea.NewProgram(entryModel, opts...)

	return &App{TuiModel: entryModel, TuiProgram: *program}, nil
}

// GitchaTui will start up the TUI program for Gitcha.
func (a *App) GitchaTui() error {
	if _, err := a.TuiProgram.Run(); err != nil {
		return fmt.Errorf("GitchaTui: attempted to run program and received an error: %w", err)
	}

	return nil
}

// GetDirectoryFromArgs will attempt to get a directory from the given args.
//   - If no args are provided the working directory is returned with a nil error.
//   - If multiple args are provided the first argument alone will be evaluated.
//   - If the first argument is not a directory an empty string will be returned with a non-nil error.
func GetDirectoryFromArgs(args []string) (string, error) {
	if len(args) == 0 {
		return os.Getwd()
	}

	path := args[0]

	isDirectory, err := isDirectory(path)
	if err != nil {
		return "", fmt.Errorf("GetDirectoryFromArgs: argument %s is not a directory: %w", path, err)
	}
	if !isDirectory {
		return "", fmt.Errorf("GetDirectoryFromArgs: argument %s is not a directory", path)
	}

	return path, nil
}

// isDirectory determines if the provided path is a directory or not.
func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}
