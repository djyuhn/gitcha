package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/djyuhn/gitcha/reporeader"

	tea "github.com/charmbracelet/bubbletea"
)

type EntryModel struct {
	RepoReader reporeader.RepoReader

	RepoDetails reporeader.RepoDetails
}

func NewEntryModel(repoReader *reporeader.RepoReader) (EntryModel, error) {
	if repoReader == nil {
		return EntryModel{}, fmt.Errorf("NewEntryModel: received a nil RepoReader")
	}

	return EntryModel{RepoReader: *repoReader}, nil
}

var _ tea.Model = EntryModel{}

type RepoDetailsMsg struct {
	Err         error
	RepoDetails reporeader.RepoDetails
}

func (m EntryModel) Init() tea.Cmd {
	return m.processRepo()
}

func (m EntryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	case RepoDetailsMsg:
		if msg.Err == nil {
			m.RepoDetails = msg.RepoDetails
			return m, nil
		}
	default:
		return m, nil
	}

	return m, nil
}

func (m EntryModel) View() string {
	view := strings.Builder{}
	view.WriteString(fmt.Sprintf("Repository Created Date - %s\n", m.RepoDetails.CreatedDate.Format(time.RFC822)))
	view.WriteString(fmt.Sprintf("Repository License - %s\n", m.RepoDetails.License))
	for author, commits := range m.RepoDetails.AuthorsCommits {
		view.WriteString(fmt.Sprintf("Author - %s : Email - %s : Commit count - %d \n", author.Name, author.Email, len(commits)))
	}
	return view.String()
}

func (m EntryModel) processRepo() tea.Cmd {
	return func() tea.Msg {
		details, err := m.RepoReader.GetRepoDetails()

		detailsMsg := RepoDetailsMsg{
			Err:         err,
			RepoDetails: details,
		}

		return detailsMsg
	}
}
