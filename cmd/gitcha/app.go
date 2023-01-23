package gitcha

import (
	"fmt"

	"github.com/djyuhn/gitcha/tui"

	tea "github.com/charmbracelet/bubbletea"
)

type Program interface {
	Run() (tea.Model, error)
}

var _ Program = &tea.Program{}

type App struct {
	TuiModel   tui.EntryModel
	TuiProgram Program
}

func NewApp(tuiModel tui.EntryModel) (*App, error) {
	program := tea.NewProgram(tuiModel)

	return &App{TuiModel: tuiModel, TuiProgram: program}, nil
}

func NewAppProgram(program Program) (*App, error) {
	if program == nil {
		return nil, fmt.Errorf("NewAppProgram: received nil program")
	}

	return &App{TuiProgram: program}, nil
}

// GitchaTui will start up the TUI for Gitcha.
func (a *App) GitchaTui() error {
	if _, err := a.TuiProgram.Run(); err != nil {
		return fmt.Errorf("GitchaTui: attempted to run program and received an error: %w", err)
	}

	return nil
}
