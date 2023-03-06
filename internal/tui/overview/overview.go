package overview

import (
	"fmt"
	"sort"
	"strings"
	"time"

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
	view.WriteString(fmt.Sprintf("Repository Created Date - %s\n", o.RepoDetails.CreatedDate.Format(time.RFC822)))
	view.WriteString(fmt.Sprintf("Repository License - %s\n", o.RepoDetails.License))
	if len(o.orderedAuthorsByCommitCount) < topAuthorCount {
		for _, pair := range o.orderedAuthorsByCommitCount {
			view.WriteString(fmt.Sprintf("AuthorEmail - %s : Email - %s : Commit count - %d\n", pair.AuthorName, pair.AuthorEmail, len(pair.Commits)))
		}
	} else {
		for i := 0; i < topAuthorCount; i++ {
			name := o.orderedAuthorsByCommitCount[i].AuthorName
			email := o.orderedAuthorsByCommitCount[i].AuthorEmail
			count := len(o.orderedAuthorsByCommitCount[i].Commits)
			view.WriteString(fmt.Sprintf("AuthorEmail - %s : Email - %s : Commit count - %d\n", name, email, count))
		}
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
