package overview

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/djyuhn/gitcha/internal/reporeader"
)

const topAuthorCount = 3

type Overview struct {
	RepoDetails reporeader.RepoDetails

	orderedAuthorsByCommitCount []AuthorCommitsPair
}

var _ tea.Model = Overview{}

func NewOverview(repoDetails reporeader.RepoDetails) Overview {
	topAuthorsByCommits := getSortedAuthorsByCommitCount(repoDetails.AuthorsCommits)

	return Overview{RepoDetails: repoDetails, orderedAuthorsByCommitCount: topAuthorsByCommits}
}

func (o Overview) Init() tea.Cmd {
	return nil
}

func (o Overview) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return o, nil
}

func (o Overview) View() string {
	view := strings.Builder{}
	if len(o.orderedAuthorsByCommitCount) < topAuthorCount {
		for _, pair := range o.orderedAuthorsByCommitCount {
			view.WriteString(fmt.Sprintf("Author - %s : Email - %s : Commit count - %d\n", pair.Author.Name, pair.Author.Email, len(pair.Commits)))
		}
	} else {
		for i := 0; i < topAuthorCount; i++ {
			name := o.orderedAuthorsByCommitCount[i].Author.Name
			email := o.orderedAuthorsByCommitCount[i].Author.Email
			count := len(o.orderedAuthorsByCommitCount[i].Commits)
			view.WriteString(fmt.Sprintf("Author - %s : Email - %s : Commit count - %d\n", name, email, count))
		}
	}

	return view.String()
}

type AuthorCommitsPair struct {
	Author  reporeader.Author
	Commits []reporeader.Commit
}

// getSortedAuthorsByCommitCount iterates through authorCommits and returns an ordered slice of AuthorCommitsPair.
//
// The slice is ordered by the highest to the lowest commit count.
func getSortedAuthorsByCommitCount(authorCommits map[reporeader.Author][]reporeader.Commit) []AuthorCommitsPair {
	authorCommitPairs := make([]AuthorCommitsPair, 0, len(authorCommits))
	for author, commits := range authorCommits {
		pair := AuthorCommitsPair{
			Author:  author,
			Commits: commits,
		}
		authorCommitPairs = append(authorCommitPairs, pair)
	}

	// Want to order authors by the highest to the lowest commit count
	sort.SliceStable(authorCommitPairs, func(i, j int) bool {
		return len(authorCommitPairs[i].Commits) > len(authorCommitPairs[j].Commits)
	})

	return authorCommitPairs
}
