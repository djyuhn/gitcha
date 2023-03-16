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

		_, repo, err := gittest.CreateEmptyRepo(ctx, t)
		assert.Error(t, err)
		assert.ErrorContains(t, err, transport.ErrEmptyRemoteRepository.Error())
		assert.NotNil(t, repo)
	})

	t.Run("should return directory of empty repository", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		actual, repo, err := gittest.CreateEmptyRepo(ctx, t)
		require.Error(t, err)

		wt, err := repo.Worktree()
		require.NoError(t, err)

		fs := wt.Filesystem

		dirPath := fs.Root()

		assert.Equal(t, dirPath, actual)
	})
}

func TestCreateBasicRepo(t *testing.T) {
	t.Parallel()

	t.Run("should return basic repository with three commits with dates", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		_, repo, err := gittest.CreateBasicRepo(ctx, t)
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

		_, repo, err := gittest.CreateBasicRepo(ctx, t)
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

	t.Run("should return directory of the basic repository", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		actual, repo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		wt, err := repo.Worktree()
		require.NoError(t, err)

		fs := wt.Filesystem

		dirPath := fs.Root()

		assert.Equal(t, dirPath, actual)
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

func TestCreateMultiNamedAuthorRepo(t *testing.T) {
	t.Parallel()

	t.Run("should return repository with 10 commits from 4 emails with author names differing", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		const expectedAuthorEmail1 = "gitcha1@gitcha.com"
		const expectedAuthorEmail2 = "gitcha2@gitcha.com"
		const expectedAuthorEmail3 = "gitcha3@gitcha.com"
		const expectedAuthorEmail4 = "gitcha4@gitcha.com"

		repo, err := gittest.CreateMultiNamedAuthorRepo(ctx, t)
		require.NoError(t, err)

		cIter, err := repo.Log(&git.LogOptions{All: true})
		require.NoError(t, err)

		commits := make([]*object.Commit, 0)
		err = cIter.ForEach(func(c *object.Commit) error {
			commits = append(commits, c)
			return nil
		})
		require.NoError(t, err)

		authorEmailToCommits := make(map[string][]*object.Commit)
		for _, c := range commits {
			if val, ok := authorEmailToCommits[c.Author.Email]; ok {
				authorEmailToCommits[c.Author.Email] = append(val, c)
			} else {
				authorEmailToCommits[c.Author.Email] = []*object.Commit{c}
			}
		}

		assert.Equal(t, 4, len(authorEmailToCommits))

		assert.Len(t, authorEmailToCommits[expectedAuthorEmail1], 1)
		assert.Len(t, authorEmailToCommits[expectedAuthorEmail2], 2)
		assert.Len(t, authorEmailToCommits[expectedAuthorEmail3], 3)
		assert.Len(t, authorEmailToCommits[expectedAuthorEmail4], 4)
	})
}
