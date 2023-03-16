package overview

import (
	"fmt"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/djyuhn/gitcha/internal/reporeader"
	"github.com/djyuhn/gitcha/internal/tui/style"
)

const topAuthorCount = 3

type Overview struct {
	RepoDetails reporeader.RepoDetails
	theme       style.Theme

	orderedAuthorsByCommitCount []AuthorCommitsPair
}

var _ tea.Model = Overview{}

func NewOverview(repoDetails reporeader.RepoDetails) Overview {
	topAuthorsByCommits := getSortedAuthorsByCommitCount(repoDetails.AuthorsCommits)

	defaultTheme := style.NewDefaultTheme()

	return Overview{RepoDetails: repoDetails, orderedAuthorsByCommitCount: topAuthorsByCommits, theme: *defaultTheme}
}

func (o Overview) Init() tea.Cmd {
	return nil
}

func (o Overview) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return o, nil
}

func (o Overview) View() string {
	view := strings.Builder{}

	view.WriteString(o.buildRepoCreatedDateView() + "\n")
	view.WriteString(o.buildLicenseView() + "\n")
	view.WriteString(o.buildAuthorView() + "\n")

	return view.String()
}

func (o Overview) buildRepoCreatedDateView() string {
	view := strings.Builder{}

	primaryColorStyle := lipgloss.NewStyle().Foreground(o.theme.General.PrimaryColor)
	secondaryColorStyle := lipgloss.NewStyle().Foreground(o.theme.General.SecondaryColor)

	labelView := primaryColorStyle.Render("Created:")
	createdDateView := secondaryColorStyle.Render(o.RepoDetails.CreatedDate.Format(time.RFC822))

	view.WriteString(fmt.Sprintf("%s %s", labelView, createdDateView))

	return view.String()
}

func (o Overview) buildLicenseView() string {
	view := strings.Builder{}

	primaryColorStyle := lipgloss.NewStyle().Foreground(o.theme.General.PrimaryColor)
	secondaryColorStyle := lipgloss.NewStyle().Foreground(o.theme.General.SecondaryColor)

	labelView := primaryColorStyle.Render("License:")
	licenseView := secondaryColorStyle.Render(o.RepoDetails.License)

	view.WriteString(fmt.Sprintf("%s %s", labelView, licenseView))

	return view.String()
}

func (o Overview) buildAuthorView() string {
	view := strings.Builder{}

	authorCount := topAuthorCount

	if len(o.orderedAuthorsByCommitCount) < authorCount {
		authorCount = len(o.orderedAuthorsByCommitCount)
	}

	for i := 0; i < authorCount; i++ {
		primaryColorStyle := lipgloss.NewStyle().Foreground(o.theme.General.PrimaryColor)
		secondaryColorStyle := lipgloss.NewStyle().Foreground(o.theme.General.SecondaryColor)

		label := primaryColorStyle.Render("Author:")
		name := secondaryColorStyle.Render(o.orderedAuthorsByCommitCount[i].AuthorName)
		email := secondaryColorStyle.Render(o.orderedAuthorsByCommitCount[i].AuthorEmail)
		count := secondaryColorStyle.Render(fmt.Sprintf("%d", len(o.orderedAuthorsByCommitCount[i].Commits)))

		view.WriteString(fmt.Sprintf("%s %s %s %s\n", label, name, email, count))
	}

	return view.String()
}

type AuthorCommitsPair struct {
	AuthorName  string
	AuthorEmail string
	Commits     []reporeader.Commit
}

// getSortedAuthorsByCommitCount iterates through authorCommits and returns an ordered slice of AuthorCommitsPair.
//
// The slice is ordered by the highest to the lowest commit count.
func getSortedAuthorsByCommitCount(authorCommits map[string][]reporeader.Commit) []AuthorCommitsPair {
	authorCommitPairs := make([]AuthorCommitsPair, 0, len(authorCommits))
	for email, commits := range authorCommits {
		if len(commits) == 0 {
			continue
		}
		pair := AuthorCommitsPair{
			AuthorName:  commits[len(commits)-1].Author.Name,
			AuthorEmail: email,
			Commits:     commits,
		}
		authorCommitPairs = append(authorCommitPairs, pair)
	}

	// Want to order authors by the highest to the lowest commit count
	sort.SliceStable(authorCommitPairs, func(i, j int) bool {
		return len(authorCommitPairs[i].Commits) > len(authorCommitPairs[j].Commits)
	})

	return authorCommitPairs
}
