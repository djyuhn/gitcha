package tui

import tea "github.com/charmbracelet/bubbletea"

type EntryModel struct {
}

var _ tea.Model = EntryModel{}

func (m EntryModel) Init() tea.Cmd {
	return nil
}

func (m EntryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m EntryModel) View() string {
	return "Entry View"
}
