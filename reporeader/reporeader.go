package reporeader

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
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

	cIter, _ := repo.Log(&git.LogOptions{From: head.Hash(), Order: git.LogOrderCommitterTime})
	commits := make([]*object.Commit, 0)
	_ = cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		return nil
	})

	return commits[len(commits)-1].Author.When, nil
}
