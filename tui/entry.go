package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/djyuhn/gitcha/reporeader"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type EntryModel struct {
	RepoReader  reporeader.RepoReader
	RepoDetails reporeader.RepoDetails
	RepoError   error

	Spinner spinner.Model

	IsLoading bool
}

func NewEntryModel(repoReader *reporeader.RepoReader) (EntryModel, error) {
	if repoReader == nil {
		return EntryModel{}, fmt.Errorf("NewEntryModel: received a nil RepoReader")
	}

	sp := spinner.New()

	return EntryModel{RepoReader: *repoReader, Spinner: sp, IsLoading: true}, nil
}

var _ tea.Model = EntryModel{}

type RepoDetailsMsg struct {
	Err         error
	RepoDetails reporeader.RepoDetails
}

type LoadingRepoMsg struct {
	IsLoading bool
}

func (m EntryModel) Init() tea.Cmd {
	return tea.Batch(
		m.Spinner.Tick,
		m.processRepo,
	)
}

func (m EntryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	case spinner.TickMsg:
		if m.IsLoading {
			var cmd tea.Cmd
			m.Spinner, cmd = m.Spinner.Update(msg)
			return m, cmd
		}
		return m, nil
	case RepoDetailsMsg:
		m.RepoDetails = msg.RepoDetails
		m.RepoError = msg.Err
		return m, createLoadingRepoCmd(false)
	case LoadingRepoMsg:
		m.IsLoading = msg.IsLoading
		return m, nil
	default:
		return m, nil
	}

	return m, nil
}

func (m EntryModel) View() string {
	if m.IsLoading {
		return m.Spinner.View() + " Processing..."
	}
	if m.RepoError != nil {
		return "An error occurred while processing the repository."
	}
	view := strings.Builder{}
	view.WriteString(fmt.Sprintf("Repository Created Date - %s\n", m.RepoDetails.CreatedDate.Format(time.RFC822)))
	view.WriteString(fmt.Sprintf("Repository License - %s\n", m.RepoDetails.License))
	for author, commits := range m.RepoDetails.AuthorsCommits {
		view.WriteString(fmt.Sprintf("Author - %s : Email - %s : Commit count - %d \n", author.Name, author.Email, len(commits)))
	}
	return view.String()
}

func (m EntryModel) processRepo() tea.Msg {
	details, err := m.RepoReader.GetRepoDetails()

	detailsMsg := RepoDetailsMsg{
		Err:         err,
		RepoDetails: details,
	}

	return detailsMsg
}

func createLoadingRepoCmd(isLoading bool) tea.Cmd {
	return func() tea.Msg {
		return LoadingRepoMsg{IsLoading: isLoading}
	}
}
