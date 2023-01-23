package gittest_test

import (
	"context"
	"testing"
	"time"

	"github.com/djyuhn/gitcha/gittest"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateEmptyRepo(t *testing.T) {
	t.Parallel()

	t.Run("should return empty repository", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		repo, err := gittest.CreateEmptyRepo(ctx, t)
		assert.Error(t, err)
		assert.ErrorContains(t, err, transport.ErrEmptyRemoteRepository.Error())
		assert.NotNil(t, repo)
	})
}

func TestCreateBasicRepo(t *testing.T) {
	t.Parallel()

	t.Run("should return basic repository with three commits with dates", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		repo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		cIter, err := repo.Log(&git.LogOptions{All: true})
		require.NoError(t, err)

		commits := make([]*object.Commit, 0)
		err = cIter.ForEach(func(c *object.Commit) error {
			commits = append(commits, c)
			return nil
		})
		assert.NoError(t, err)
		assert.Equal(t, 3, len(commits))

		expectedThirdCommitTime, err := time.Parse(time.RFC3339, "2022-11-10T08:20:00-06:00")
		require.NoError(t, err)
		assert.Equal(t, expectedThirdCommitTime.Format(time.RFC3339), commits[0].Author.When.Format(time.RFC3339))

		expectedSecondCommitTime, err := time.Parse(time.RFC3339, "2022-11-10T08:10:00-06:00")
		require.NoError(t, err)
		assert.Equal(t, expectedSecondCommitTime.Format(time.RFC3339), commits[1].Author.When.Format(time.RFC3339))

		expectedFirstCommitTime, err := time.Parse(time.RFC3339, "2022-11-10T08:00:00-06:00")
		require.NoError(t, err)
		assert.Equal(t, expectedFirstCommitTime.Format(time.RFC3339), commits[2].Author.When.Format(time.RFC3339))
	})

	t.Run("should return basic repository with one author", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		repo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		cIter, err := repo.Log(&git.LogOptions{All: true})
		require.NoError(t, err)

		commits := make([]*object.Commit, 0)
		err = cIter.ForEach(func(c *object.Commit) error {
			commits = append(commits, c)
			return nil
		})

		uniqueAuthors := make(map[string]struct{})
		for _, c := range commits {
			uniqueAuthors[c.Author.Name] = struct{}{}
		}

		assert.NoError(t, err)
		assert.Equal(t, 1, len(uniqueAuthors))
	})
}

func TestCreateBasicMultiAuthorRepo(t *testing.T) {
	t.Parallel()

	t.Run("should return basic repository with four commits with dates", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		repo, err := gittest.CreateBasicMultiAuthorRepo(ctx, t)
		require.NoError(t, err)

		cIter, err := repo.Log(&git.LogOptions{All: true})
		require.NoError(t, err)

		commits := make([]*object.Commit, 0)
		err = cIter.ForEach(func(c *object.Commit) error {
			commits = append(commits, c)
			return nil
		})
		assert.NoError(t, err)
		assert.Equal(t, 4, len(commits))

		expectedFourthCommitTime, err := time.Parse(time.RFC3339, "2022-11-10T08:30:00-06:00")
		require.NoError(t, err)
		assert.Equal(t, expectedFourthCommitTime.Format(time.RFC3339), commits[0].Author.When.Format(time.RFC3339))

		expectedThirdCommitTime, err := time.Parse(time.RFC3339, "2022-11-10T08:20:00-06:00")
		require.NoError(t, err)
		assert.Equal(t, expectedThirdCommitTime.Format(time.RFC3339), commits[1].Author.When.Format(time.RFC3339))

		expectedSecondCommitTime, err := time.Parse(time.RFC3339, "2022-11-10T08:10:00-06:00")
		require.NoError(t, err)
		assert.Equal(t, expectedSecondCommitTime.Format(time.RFC3339), commits[2].Author.When.Format(time.RFC3339))

		expectedFirstCommitTime, err := time.Parse(time.RFC3339, "2022-11-10T08:00:00-06:00")
		require.NoError(t, err)
		assert.Equal(t, expectedFirstCommitTime.Format(time.RFC3339), commits[3].Author.When.Format(time.RFC3339))
	})

	t.Run("should return basic repository with four authors", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		repo, err := gittest.CreateBasicMultiAuthorRepo(ctx, t)
		require.NoError(t, err)

		cIter, err := repo.Log(&git.LogOptions{All: true})
		require.NoError(t, err)

		commits := make([]*object.Commit, 0)
		err = cIter.ForEach(func(c *object.Commit) error {
			commits = append(commits, c)
			return nil
		})

		uniqueAuthors := make(map[string]struct{})
		for _, c := range commits {
			uniqueAuthors[c.Author.Name] = struct{}{}
		}

		assert.NoError(t, err)
		assert.Equal(t, 4, len(uniqueAuthors))
	})
}
