package reporeader_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gitcha/gittest"
	"gitcha/reporeader"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCreatedDate(t *testing.T) {
	t.Run("given repository with commits should return time of oldest commit", func(t *testing.T) {
		ctx := context.Background()
		repo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		cIter, err := repo.Log(&git.LogOptions{Order: git.LogOrderCommitterTime})
		require.NoError(t, err)

		commits := make([]*object.Commit, 0)
		err = cIter.ForEach(func(c *object.Commit) error {
			commits = append(commits, c)
			return nil
		})
		require.NoError(t, err)

		expected := commits[len(commits)-1].Author.When
		actual, err := reporeader.GetCreatedDate(repo)

		assert.Equal(t, expected, actual)
		assert.NoError(t, err)
	})

	t.Run("given nil repository should return default time and error", func(t *testing.T) {
		expectedTime := time.Time{}
		expectedErr := fmt.Errorf("GetCreatedDate: received a nil repository")

		actualTime, err := reporeader.GetCreatedDate(nil)

		assert.Equal(t, expectedTime, actualTime)
		assert.ErrorContains(t, err, expectedErr.Error())
	})

	t.Run("given repository with nil Storer should return default time and error", func(t *testing.T) {
		expectedErr := fmt.Errorf("GetCreatedDate: invalid repository - Storer is nil")
		expectedTime := time.Time{}
		actualTime, err := reporeader.GetCreatedDate(&git.Repository{Storer: nil})

		assert.Equal(t, expectedTime, actualTime)
		assert.Error(t, err)
		assert.ErrorContains(t, err, expectedErr.Error())
	})

	t.Run("given empty repository should return default time and error", func(t *testing.T) {
		ctx := context.Background()
		repo, _ := gittest.CreateEmptyRepo(ctx, t)

		expectedErr := fmt.Errorf("GetCreatedDate: received a repository without a head")
		expectedTime := time.Time{}
		actualTime, err := reporeader.GetCreatedDate(repo)

		assert.Equal(t, expectedTime, actualTime)
		assert.ErrorContains(t, err, expectedErr.Error())
	})
}

func TestGetContributorsByCommits(t *testing.T) {
	t.Parallel()

	t.Run("given nil repository should return empty map and error", func(t *testing.T) {
		t.Parallel()
		expectedErr := fmt.Errorf("GetAuthorsByCommits: received a nil repository")

		actual, err := reporeader.GetAuthorsByCommits(nil)

		assert.Nil(t, actual)
		assert.ErrorContains(t, err, expectedErr.Error())
	})

	t.Run("given repository with nil Storer should empty map and error", func(t *testing.T) {
		t.Parallel()
		expectedErr := fmt.Errorf("GetAuthorsByCommits: invalid repository - Storer is nil")
		actual, err := reporeader.GetAuthorsByCommits(&git.Repository{Storer: nil})

		assert.Nil(t, actual)
		assert.Error(t, err)
		assert.ErrorContains(t, err, expectedErr.Error())
	})

	t.Run("given empty repository should return empty map and error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		repo, _ := gittest.CreateEmptyRepo(ctx, t)

		expectedErr := fmt.Errorf("GetAuthorsByCommits: received a repository without a head")
		actual, err := reporeader.GetAuthorsByCommits(repo)

		assert.Nil(t, actual)
		assert.ErrorContains(t, err, expectedErr.Error())
	})

	t.Run("given single commit author should return one contributor", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		repo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		actual, err := reporeader.GetAuthorsByCommits(repo)

		assert.NoError(t, err)
		assert.Len(t, actual, 1)
	})

	t.Run("given single commit author should return map with author as key", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		repo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		expectedAuthor := reporeader.Author{
			Name:  "gitcha-author-name",
			Email: "gitcha-author-email@gitcha.com",
		}

		actual, err := reporeader.GetAuthorsByCommits(repo)
		assert.NoError(t, err)
		assert.Contains(t, actual, expectedAuthor)
	})

	t.Run("given multiple commit authors should return map with each author as key", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		repo, err := gittest.CreateBasicMultiAuthorRepo(ctx, t)
		require.NoError(t, err)

		expectedAuthor1 := reporeader.Author{
			Name:  "Gitcha One",
			Email: "gitcha1@gitcha.com",
		}
		expectedAuthor2 := reporeader.Author{
			Name:  "Gitcha Two",
			Email: "gitcha2@gitcha.com",
		}
		expectedAuthor3 := reporeader.Author{
			Name:  "Gitcha Three",
			Email: "gitcha3@gitcha.com",
		}
		expectedAuthor4 := reporeader.Author{
			Name:  "Gitcha Four",
			Email: "gitcha4@gitcha.com",
		}

		actual, err := reporeader.GetAuthorsByCommits(repo)
		assert.NoError(t, err)

		assert.Contains(t, actual, expectedAuthor1)
		assert.Contains(t, actual, expectedAuthor2)
		assert.Contains(t, actual, expectedAuthor3)
		assert.Contains(t, actual, expectedAuthor4)
	})
}
