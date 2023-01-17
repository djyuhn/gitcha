package reporeader

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// GetCreatedDate returns the time that the repository was first created.
func GetCreatedDate(repo *git.Repository) (time.Time, error) {
	if repo == nil {
		return time.Time{}, fmt.Errorf("GetCreatedDate: received a nil repository")
	}
	if repo.Storer == nil {
		return time.Time{}, fmt.Errorf("GetCreatedDate: invalid repository - Storer is nil")
	}

	head, err := repo.Head()
	if err != nil {
		return time.Time{}, fmt.Errorf("GetCreatedDate: received a repository without a head: %w", err)
	}
	if head == nil {
		return time.Time{}, fmt.Errorf("GetCreatedDate: received a repository without a head")
	}

	commits := make([]*object.Commit, 0)
	cIter, _ := repo.Log(&git.LogOptions{From: head.Hash(), Order: git.LogOrderCommitterTime})
	_ = cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		return nil
	})

	return commits[len(commits)-1].Author.When, nil
}

type Author struct {
	Name  string
	Email string
}

// GetAuthorsByCommits returns the contributors and their commits they made.
func GetAuthorsByCommits(repo *git.Repository) (map[Author][]object.Commit, error) {
	if repo == nil {
		return nil, fmt.Errorf("GetAuthorsByCommits: received a nil repository")
	}
	if repo.Storer == nil {
		return nil, fmt.Errorf("GetAuthorsByCommits: invalid repository - Storer is nil")
	}

	head, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("GetAuthorsByCommits: received a repository without a head: %w", err)
	}
	if head == nil {
		return nil, fmt.Errorf("GetAuthorsByCommits: received a repository without a head")
	}

	contributorCommits := make(map[Author][]object.Commit)

	cIter, _ := repo.Log(&git.LogOptions{From: head.Hash(), Order: git.LogOrderCommitterTime})
	_ = cIter.ForEach(func(c *object.Commit) error {
		contributor := Author{
			Name:  c.Author.Name,
			Email: c.Author.Email,
		}

		contributorCommits[contributor] = append(contributorCommits[contributor], *c)
		return nil
	})

	return contributorCommits, nil
}

// ValidateRepository validates the given repository and returns the head of the repository if valid and a nil error.
// If the repository is invalid a nil head reference and a non-nil error are returned.
func ValidateRepository(repo *git.Repository) (*plumbing.Reference, error) {
	if repo == nil {
		return nil, fmt.Errorf("ValidateRepository: received a nil repository")
	}
	if repo.Storer == nil {
		return nil, fmt.Errorf("ValidateRepository: invalid repository - Storer is nil")
	}

	head, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("ValidateRepository: received a repository without a head: %w", err)
	}
	if head == nil {
		return nil, fmt.Errorf("ValidateRepository: received a repository without a head")
	}

	return head, nil
}
