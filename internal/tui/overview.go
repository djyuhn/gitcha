package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Overview struct {
}

var _ tea.Model = Overview{}

func NewOverview() Overview {
	return Overview{}
}

func (o Overview) Init() tea.Cmd {
	return nil
}

func (o Overview) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return o, nil
}

func (o Overview) View() string {
	return "OVERVIEW"
}
