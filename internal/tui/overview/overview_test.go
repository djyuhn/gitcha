package overview_test

import (
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/djyuhn/gitcha/internal/reporeader"
	"github.com/djyuhn/gitcha/internal/tui/overview"

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

		orderedAuthors := getSortedAuthorsByCommitCount(authorCommits)
		expectedView := strings.Builder{}
		for i := 0; i < 3; i++ {
			name := orderedAuthors[i].AuthorName
			email := orderedAuthors[i].AuthorEmail
			count := len(orderedAuthors[i].Commits)
			expectedView.WriteString(fmt.Sprintf("AuthorEmail - %s : Email - %s : Commit count - %d\n", name, email, count))
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

		expectedView := fmt.Sprintf("AuthorEmail - %s : Email - %s : Commit count - %d\n", author.Name, author.Email, commitCount)
		actual := model.View()

		assert.Contains(t, actual, expectedView)
	})

	t.Run("given only 1 author with multiple names should return only 1 author and the name of their last commit in view", func(t *testing.T) {
		t.Parallel()

		authorCommits := make(map[string][]reporeader.Commit)

		authorName := "Author Name"
		authorEmail := "author@email.com"

		const commitCount = 10
		commits := make([]reporeader.Commit, 0, commitCount)
		for j := 0; j < commitCount; j++ {
			author := reporeader.Author{
				Name:  fmt.Sprintf("%s%d", authorName, j),
				Email: authorEmail,
			}
			commit := reporeader.Commit{
				Author:  author,
				Message: fmt.Sprintf("Message %d", j),
				Hash:    fmt.Sprintf("Hash %d", j),
			}
			commits = append(commits, commit)
		}
		authorCommits[authorEmail] = commits

		repoDetails := reporeader.RepoDetails{AuthorsCommits: authorCommits}
		model := overview.NewOverview(repoDetails)

		expectedName := fmt.Sprintf("%s%d", authorName, commitCount-1)

		expectedView := fmt.Sprintf("AuthorEmail - %s : Email - %s : Commit count - %d\n", expectedName, authorEmail, commitCount)
		actual := model.View()

		assert.Contains(t, actual, expectedView)
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

		expectedView := fmt.Sprintf("Repository Created Date - %s\n", repoDetails.CreatedDate.Format(time.RFC822))

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

		expectedView := fmt.Sprintf("Repository License - %s\n", repoDetails.License)

		actual := model.View()

		assert.Contains(t, actual, expectedView)
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
