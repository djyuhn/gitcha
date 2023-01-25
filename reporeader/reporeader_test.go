package reporeader_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/djyuhn/gitcha/gittest"
	"github.com/djyuhn/gitcha/reporeader"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRepoReader(t *testing.T) {
	t.Parallel()

	t.Run("given a directory with an invalid repository should return nil RepoReader and error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		emptyRepo, err := gittest.CreateEmptyRepo(ctx, t)
		require.Error(t, err)
		wt, err := emptyRepo.Worktree()
		require.NoError(t, err)

		fs := wt.Filesystem

		expectedError := fmt.Errorf("NewRepoReader: error detected in attempting to open repository")
		reader, err := reporeader.NewRepoReader(fs.Root())

		assert.Nil(t, reader)
		assert.Errorf(t, err, expectedError.Error())
	})

	t.Run("given a directory with a repository should return RepoReader and nil error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		basicRepo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)
		wt, err := basicRepo.Worktree()
		require.NoError(t, err)

		fs := wt.Filesystem

		reader, err := reporeader.NewRepoReader(fs.Root())

		assert.NotNil(t, reader)
		assert.NoError(t, err)
	})
}

func TestNewRepoReaderRepository(t *testing.T) {
	t.Parallel()

	t.Run("given an invalid repository should return nil RepoReader and error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		emptyRepo, err := gittest.CreateEmptyRepo(ctx, t)
		require.Error(t, err)

		expectedError := fmt.Errorf("NewRepoReaderRepository: received an invalid repository")
		reader, err := reporeader.NewRepoReaderRepository(emptyRepo)

		assert.Nil(t, reader)
		assert.Errorf(t, err, expectedError.Error())
	})

	t.Run("given a valid repository should return RepoReader and nil error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		basicRepo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		reader, err := reporeader.NewRepoReaderRepository(basicRepo)

		assert.NotNil(t, reader)
		assert.NoError(t, err)
	})
}

func TestRepoReader_GetRepoDetails(t *testing.T) {
	t.Parallel()

	t.Run("given repository with commits should return time of oldest commit and nil error", func(t *testing.T) {
		t.Parallel()
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

		repoReader, err := reporeader.NewRepoReaderRepository(repo)
		require.NoError(t, err)

		expected := commits[len(commits)-1].Author.When
		actual, err := repoReader.GetRepoDetails()

		assert.Equal(t, expected, actual.CreatedDate)
		assert.NoError(t, err)
	})

	t.Run("given single commit author should return map with author as key", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		repo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		repoReader, err := reporeader.NewRepoReaderRepository(repo)
		require.NoError(t, err)

		expectedAuthor := reporeader.Author{
			Name:  "gitcha-author-name",
			Email: "gitcha-author-email@gitcha.com",
		}
		actual, err := repoReader.GetRepoDetails()

		assert.NoError(t, err)
		assert.Contains(t, actual.AuthorsCommits, expectedAuthor)
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

		repoReader, err := reporeader.NewRepoReaderRepository(repo)
		require.NoError(t, err)

		actual, err := repoReader.GetRepoDetails()
		assert.NoError(t, err)

		assert.Contains(t, actual.AuthorsCommits, expectedAuthor1)
		assert.Contains(t, actual.AuthorsCommits, expectedAuthor2)
		assert.Contains(t, actual.AuthorsCommits, expectedAuthor3)
		assert.Contains(t, actual.AuthorsCommits, expectedAuthor4)
	})

	t.Run("given basic repository with LICENSE file at root should return MIT license and nil error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		repo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		repoReader, err := reporeader.NewRepoReaderRepository(repo)
		require.NoError(t, err)

		actual, err := repoReader.GetRepoDetails()

		assert.Equal(t, "MIT", actual.License)
		assert.NoError(t, err)
	})

	t.Run("given repository with no LICENSE file should return NO LICENSE string and nil error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		repo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		wt, err := repo.Worktree()
		require.NoError(t, err)
		fs := wt.Filesystem

		require.NoError(t, fs.Remove("LICENSE"))

		repoReader, err := reporeader.NewRepoReaderRepository(repo)
		require.NoError(t, err)

		actual, err := repoReader.GetRepoDetails()

		assert.Equal(t, "NO LICENSE", actual.License)
		assert.NoError(t, err)
	})
}

func TestRepoReader_GetCreatedDate(t *testing.T) {
	t.Parallel()

	t.Run("given repository with commits should return time of oldest commit", func(t *testing.T) {
		t.Parallel()
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

		repoReader, err := reporeader.NewRepoReaderRepository(repo)
		require.NoError(t, err)

		expected := commits[len(commits)-1].Author.When
		actual, err := repoReader.GetCreatedDate()

		assert.Equal(t, expected, actual)
		assert.NoError(t, err)
	})
}

func TestRepoReader_GetAuthorsByCommits(t *testing.T) {
	t.Parallel()

	t.Run("given single commit author should return map with author as key", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		repo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		repoReader, err := reporeader.NewRepoReaderRepository(repo)
		require.NoError(t, err)

		expectedAuthor := reporeader.Author{
			Name:  "gitcha-author-name",
			Email: "gitcha-author-email@gitcha.com",
		}
		actual, err := repoReader.GetAuthorsByCommits()

		assert.NoError(t, err)
		assert.Contains(t, actual, expectedAuthor)
	})

	t.Run("given multiple commit authors should return map with each author as key", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		repo, err := gittest.CreateBasicMultiAuthorRepo(ctx, t)
		require.NoError(t, err)

		repoReader, err := reporeader.NewRepoReaderRepository(repo)
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

		actual, err := repoReader.GetAuthorsByCommits()
		assert.NoError(t, err)

		assert.Contains(t, actual, expectedAuthor1)
		assert.Contains(t, actual, expectedAuthor2)
		assert.Contains(t, actual, expectedAuthor3)
		assert.Contains(t, actual, expectedAuthor4)
	})
}

func TestRepoReader_GetLicense(t *testing.T) {
	t.Parallel()

	t.Run("given basic repository with LICENSE file at root should return MIT license and nil error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		repo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		repoReader, err := reporeader.NewRepoReaderRepository(repo)
		require.NoError(t, err)

		actual, err := repoReader.GetLicense()

		assert.Equal(t, "MIT", actual)
		assert.NoError(t, err)
	})

	t.Run("given repository with no LICENSE file should return NO LICENSE string and nil error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		repo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		wt, err := repo.Worktree()
		require.NoError(t, err)
		fs := wt.Filesystem

		require.NoError(t, fs.Remove("LICENSE"))

		repoReader, err := reporeader.NewRepoReaderRepository(repo)
		require.NoError(t, err)

		actual, err := repoReader.GetLicense()

		assert.Equal(t, "NO LICENSE", actual)
		assert.NoError(t, err)
	})
}

func TestValidateRepository(t *testing.T) {
	t.Parallel()

	t.Run("given a nil repository should nil head reference and error", func(t *testing.T) {
		t.Parallel()

		expectedErr := fmt.Errorf("ValidateRepository: received a nil repository")
		head, err := reporeader.ValidateRepository(nil)

		assert.ErrorContains(t, err, expectedErr.Error())
		assert.Nil(t, head)
	})

	t.Run("given repository with nil Storer should return error", func(t *testing.T) {
		t.Parallel()

		expectedErr := fmt.Errorf("ValidateRepository: invalid repository - Storer is nil")
		head, err := reporeader.ValidateRepository(&git.Repository{Storer: nil})

		assert.ErrorContains(t, err, expectedErr.Error())
		assert.Nil(t, head)
	})

	t.Run("given empty repository should return error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		repo, err := gittest.CreateEmptyRepo(ctx, t)
		require.Error(t, err)

		expectedErr := fmt.Errorf("ValidateRepository: received a repository without a head")
		head, err := reporeader.ValidateRepository(repo)

		assert.ErrorContains(t, err, expectedErr.Error())
		assert.Nil(t, head)
	})

	t.Run("given valid repository should return head reference and nil error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		repo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		expected, err := repo.Head()
		require.NoError(t, err)

		head, err := reporeader.ValidateRepository(repo)

		assert.NoError(t, err)
		assert.Equal(t, expected, head)
	})
}
