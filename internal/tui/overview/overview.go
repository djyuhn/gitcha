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
	orderedLanguagesByFileCount []reporeader.LanguageDetails
}

var _ tea.Model = Overview{}

func NewOverview(repoDetails reporeader.RepoDetails) Overview {
	topAuthorsByCommits := getSortedAuthorsByCommitCount(repoDetails.AuthorsCommits)
	orderedLanguagesByFileCount := getSortedLanguagesByFileCount(repoDetails.LanguageDetails)

	defaultTheme := style.NewDefaultTheme()

	return Overview{
		RepoDetails:                 repoDetails,
		theme:                       *defaultTheme,
		orderedAuthorsByCommitCount: topAuthorsByCommits,
		orderedLanguagesByFileCount: orderedLanguagesByFileCount,
	}
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
	view.WriteString(o.buildLanguageView() + "\n")

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

	primaryColorStyle := lipgloss.NewStyle().Foreground(o.theme.General.PrimaryColor)
	secondaryColorStyle := lipgloss.NewStyle().Foreground(o.theme.General.SecondaryColor)

	view.WriteString(primaryColorStyle.Render("Author:"))
	for i := 0; i < authorCount; i++ {
		name := secondaryColorStyle.Render(o.orderedAuthorsByCommitCount[i].AuthorName)
		email := secondaryColorStyle.Render(o.orderedAuthorsByCommitCount[i].AuthorEmail)
		count := secondaryColorStyle.Render(fmt.Sprintf("%d", len(o.orderedAuthorsByCommitCount[i].Commits)))

		view.WriteString(fmt.Sprintf(" %s %s %s", name, email, count))
	}

	return view.String()
}

func (o Overview) buildLanguageView() string {
	view := strings.Builder{}

	primaryColorStyle := lipgloss.NewStyle().Foreground(o.theme.General.PrimaryColor)
	secondaryColorStyle := lipgloss.NewStyle().Foreground(o.theme.General.SecondaryColor)

	view.WriteString(primaryColorStyle.Render("Language:"))
	for _, details := range o.orderedLanguagesByFileCount {
		language := secondaryColorStyle.Render(string(details.Language))
		count := secondaryColorStyle.Render(fmt.Sprintf("%d", details.FileCount))

		view.WriteString(fmt.Sprintf(" %s %s", language, count))
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

// getSortedLanguagesByFileCount iterates through languageToDetails and returns an ordered slice of reporeader.LanguageDetails.
//
// The slice is ordered by the highest to the lowest file count.
func getSortedLanguagesByFileCount(languageToDetails map[reporeader.Language]reporeader.LanguageDetails) []reporeader.LanguageDetails {
	languageDetails := make([]reporeader.LanguageDetails, 0, len(languageToDetails))
	for _, details := range languageToDetails {
		languageDetails = append(languageDetails, details)
	}

	// Want to order language details by the highest file count to the lowest
	sort.SliceStable(languageDetails, func(i, j int) bool {
		return languageDetails[i].FileCount > languageDetails[j].FileCount
	})

	return languageDetails
}
