package tui

import (
	"fmt"

	"github.com/djyuhn/gitcha/reporeader"

	tea "github.com/charmbracelet/bubbletea"
)

type EntryModel struct {
	RepoReader *reporeader.RepoReader
}

func NewEntryModel(repoReader *reporeader.RepoReader) (EntryModel, error) {
	if repoReader == nil {
		return EntryModel{}, fmt.Errorf("NewEntryModel: received a nil RepoReader")
	}

	return EntryModel{RepoReader: repoReader}, nil
}

var _ tea.Model = EntryModel{}

func (m EntryModel) Init() tea.Cmd {
	return nil
}

func (m EntryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	default:
		return m, nil
	}

	return m, nil
}

func (m EntryModel) View() string {
	return "Entry View"
}
