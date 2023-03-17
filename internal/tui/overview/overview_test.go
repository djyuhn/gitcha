package overview_test

import (
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/djyuhn/gitcha/internal/reporeader"
	"github.com/djyuhn/gitcha/internal/tui/overview"
	"github.com/djyuhn/gitcha/internal/tui/style"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

func TestNewOverview(t *testing.T) {
	t.Parallel()

	t.Run("should return overview model with RepoDetails", func(t *testing.T) {
		t.Parallel()

		repoDetails := reporeader.RepoDetails{
			CreatedDate:    time.Date(2023, time.January, 26, 3, 2, 1, 0, time.UTC),
			AuthorsCommits: nil,
			License:        "SOME LICENSE",
		}
		actual := overview.NewOverview(repoDetails)

		assert.Equal(t, repoDetails, actual.RepoDetails)
	})
}

func TestOverview_Init(t *testing.T) {
	t.Parallel()

	t.Run("should return nil", func(t *testing.T) {
		t.Parallel()

		repoDetails := reporeader.RepoDetails{}
		model := overview.NewOverview(repoDetails)

		cmd := model.Init()

		assert.Nil(t, cmd)
	})
}

func TestOverview_Update(t *testing.T) {
	t.Parallel()

	t.Run("given nil msg should return model and nil cmd", func(t *testing.T) {
		t.Parallel()

		repoDetails := reporeader.RepoDetails{}
		model := overview.NewOverview(repoDetails)

		actual, cmd := model.Update(nil)

		assert.Equal(t, model, actual)
		assert.Nil(t, cmd)
	})
}

func TestOverview_View(t *testing.T) {
	t.Parallel()

	t.Run("given multiple authors should return view with top 3 authors by commit count", func(t *testing.T) {
		t.Parallel()

		authorCommits := make(map[string][]reporeader.Commit)
		// Create 10 authors with increasing commit numbers
		for i := 1; i < 11; i++ {
			authorName := fmt.Sprintf("AuthorEmail %d", i)
			authorEmail := fmt.Sprintf("author%d@email.com", i)
			author := reporeader.Author{
				Name:  authorName,
				Email: authorEmail,
			}
			commits := make([]reporeader.Commit, 0, i)
			for j := 0; j < i; j++ {
				commit := reporeader.Commit{
					Author:  author,
					Message: fmt.Sprintf("Message %d", j),
					Hash:    fmt.Sprintf("Hash %d", j),
				}
				commits = append(commits, commit)
			}
			authorCommits[author.Email] = commits
		}

		repoDetails := reporeader.RepoDetails{AuthorsCommits: authorCommits}
		model := overview.NewOverview(repoDetails)

		defaultTheme := style.NewDefaultTheme()

		orderedAuthors := getSortedAuthorsByCommitCount(authorCommits)
		primaryColorStyle := lipgloss.NewStyle().Foreground(defaultTheme.General.PrimaryColor)
		secondaryColorStyle := lipgloss.NewStyle().Foreground(defaultTheme.General.SecondaryColor)

		expectedView := strings.Builder{}
		expectedView.WriteString(primaryColorStyle.Render("Author:"))
		for i := 0; i < 3; i++ {
			name := secondaryColorStyle.Render(orderedAuthors[i].AuthorName)
			email := secondaryColorStyle.Render(orderedAuthors[i].AuthorEmail)
			count := secondaryColorStyle.Render(fmt.Sprintf("%d", len(orderedAuthors[i].Commits)))

			expectedView.WriteString(fmt.Sprintf(" %s %s %s", name, email, count))
		}

		actual := model.View()

		assert.Contains(t, actual, expectedView.String())
	})

	t.Run("given only 1 author should return only 1 author and their commit count in view", func(t *testing.T) {
		t.Parallel()

		authorCommits := make(map[string][]reporeader.Commit)

		authorName := "AuthorEmail"
		authorEmail := "author@email.com"
		author := reporeader.Author{
			Name:  authorName,
			Email: authorEmail,
		}

		const commitCount = 10
		commits := make([]reporeader.Commit, 0, commitCount)
		for j := 0; j < commitCount; j++ {
			commit := reporeader.Commit{
				Author:  author,
				Message: fmt.Sprintf("Message %d", j),
				Hash:    fmt.Sprintf("Hash %d", j),
			}
			commits = append(commits, commit)
		}
		authorCommits[author.Email] = commits

		repoDetails := reporeader.RepoDetails{AuthorsCommits: authorCommits}
		model := overview.NewOverview(repoDetails)

		defaultTheme := style.NewDefaultTheme()

		orderedAuthors := getSortedAuthorsByCommitCount(authorCommits)
		primaryColorStyle := lipgloss.NewStyle().Foreground(defaultTheme.General.PrimaryColor)
		secondaryColorStyle := lipgloss.NewStyle().Foreground(defaultTheme.General.SecondaryColor)

		expectedView := strings.Builder{}
		expectedView.WriteString(primaryColorStyle.Render("Author:"))
		for i := 0; i < len(orderedAuthors); i++ {
			name := secondaryColorStyle.Render(orderedAuthors[i].AuthorName)
			email := secondaryColorStyle.Render(orderedAuthors[i].AuthorEmail)
			count := secondaryColorStyle.Render(fmt.Sprintf("%d", len(orderedAuthors[i].Commits)))

			expectedView.WriteString(fmt.Sprintf(" %s %s %s", name, email, count))
		}

		actual := model.View()

		assert.Contains(t, actual, expectedView.String())
	})

	t.Run("given repository created date should return created date formatted as RFC822 in view", func(t *testing.T) {
		t.Parallel()

		authorCommits := make(map[string][]reporeader.Commit)
		repoDetails := reporeader.RepoDetails{
			CreatedDate:    time.Date(2023, time.January, 26, 3, 2, 1, 0, time.UTC),
			AuthorsCommits: authorCommits,
			License:        "SOME LICENSE",
		}
		model := overview.NewOverview(repoDetails)

		defaultTheme := style.NewDefaultTheme()

		primaryColorStyle := lipgloss.NewStyle().Foreground(defaultTheme.General.PrimaryColor)
		secondaryColorStyle := lipgloss.NewStyle().Foreground(defaultTheme.General.SecondaryColor)

		labelView := primaryColorStyle.Render("Created:")
		createdDate := secondaryColorStyle.Render(repoDetails.CreatedDate.Format(time.RFC822))

		expectedView := fmt.Sprintf("%s %s", labelView, createdDate)

		actual := model.View()

		assert.Contains(t, actual, expectedView)
	})

	t.Run("given license should return license in view", func(t *testing.T) {
		t.Parallel()

		authorCommits := make(map[string][]reporeader.Commit)
		repoDetails := reporeader.RepoDetails{
			CreatedDate:    time.Date(2023, time.January, 26, 3, 2, 1, 0, time.UTC),
			AuthorsCommits: authorCommits,
			License:        "SOME LICENSE",
		}
		model := overview.NewOverview(repoDetails)

		defaultTheme := style.NewDefaultTheme()

		primaryColorStyle := lipgloss.NewStyle().Foreground(defaultTheme.General.PrimaryColor)
		secondaryColorStyle := lipgloss.NewStyle().Foreground(defaultTheme.General.SecondaryColor)

		labelView := primaryColorStyle.Render("License:")
		licenseView := secondaryColorStyle.Render(repoDetails.License)

		expectedView := fmt.Sprintf("%s %s", labelView, licenseView)

		actual := model.View()

		assert.Contains(t, actual, expectedView)
	})

	t.Run("given language details should return view with languages and their file count in decreasing order", func(t *testing.T) {
		t.Parallel()

		authorCommits := make(map[string][]reporeader.Commit)
		repoDetails := reporeader.RepoDetails{
			CreatedDate:    time.Date(2023, time.January, 26, 3, 2, 1, 0, time.UTC),
			AuthorsCommits: authorCommits,
			License:        "SOME LICENSE",
			LanguageDetails: map[reporeader.Language]reporeader.LanguageDetails{
				"Go":     {Language: "Go", FileCount: 10},
				"Elixir": {Language: "Elixir", FileCount: 5},
				"Rust":   {Language: "Rust", FileCount: 2},
			},
		}
		model := overview.NewOverview(repoDetails)

		defaultTheme := style.NewDefaultTheme()

		orderedLanguages := getSortedLanguagesByFileCount(repoDetails.LanguageDetails)
		primaryColorStyle := lipgloss.NewStyle().Foreground(defaultTheme.General.PrimaryColor)
		secondaryColorStyle := lipgloss.NewStyle().Foreground(defaultTheme.General.SecondaryColor)

		expectedView := strings.Builder{}
		expectedView.WriteString(primaryColorStyle.Render("Language:"))
		for _, details := range orderedLanguages {
			language := secondaryColorStyle.Render(string(details.Language))
			count := secondaryColorStyle.Render(fmt.Sprintf("%d", details.FileCount))

			expectedView.WriteString(fmt.Sprintf(" %s %s", language, count))
		}

		actual := model.View()

		assert.Contains(t, actual, expectedView.String())
	})
}

func getSortedAuthorsByCommitCount(authorCommits map[string][]reporeader.Commit) []overview.AuthorCommitsPair {
	authorCommitPairs := make([]overview.AuthorCommitsPair, 0, len(authorCommits))
	for email, commits := range authorCommits {
		if len(commits) == 0 {
			continue
		}
		pair := overview.AuthorCommitsPair{
			AuthorName:  commits[0].Author.Name,
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
