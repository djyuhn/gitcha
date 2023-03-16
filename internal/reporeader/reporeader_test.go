package reporeader_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/djyuhn/gitcha/gittest"
	"github.com/djyuhn/gitcha/internal/reporeader"

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

		repoDir, _, err := gittest.CreateEmptyRepo(ctx, t)
		require.Error(t, err)

		expectedError := fmt.Errorf("NewRepoReader: error detected in attempting to open repository")
		reader, err := reporeader.NewRepoReader(repoDir)

		assert.Nil(t, reader)
		assert.Errorf(t, err, expectedError.Error())
	})

	t.Run("given a directory with a repository should return RepoReader and nil error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		dirPath, _, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		reader, err := reporeader.NewRepoReader(dirPath)

		assert.NotNil(t, reader)
		assert.NoError(t, err)
	})
}

func TestNewRepoReaderRepository(t *testing.T) {
	t.Parallel()

	t.Run("given an invalid repository should return nil RepoReader and error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		_, emptyRepo, err := gittest.CreateEmptyRepo(ctx, t)
		require.Error(t, err)

		expectedError := fmt.Errorf("NewRepoReaderRepository: received an invalid repository")
		reader, err := reporeader.NewRepoReaderRepository(emptyRepo)

		assert.Nil(t, reader)
		assert.Errorf(t, err, expectedError.Error())
	})

	t.Run("given a valid repository should return RepoReader and nil error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		_, basicRepo, err := gittest.CreateBasicRepo(ctx, t)
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
		_, repo, err := gittest.CreateBasicRepo(ctx, t)
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

	t.Run("given single commit author should return map with author commits", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		_, repo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		head, err := repo.Head()
		require.NoError(t, err)

		commits := make([]*object.Commit, 0)
		cIter, _ := repo.Log(&git.LogOptions{From: head.Hash(), Order: git.LogOrderCommitterTime})
		_ = cIter.ForEach(func(c *object.Commit) error {
			commits = append(commits, c)
			return nil
		})

		expectedCommits := make([]reporeader.Commit, 0, len(commits))
		for _, commit := range commits {
			author := reporeader.Author{
				Name:  commit.Author.Name,
				Email: commit.Author.Email,
			}
			commit := reporeader.Commit{
				Author:  author,
				Message: commit.Message,
				Hash:    commit.Hash.String(),
			}

			expectedCommits = append(expectedCommits, commit)
		}

		expectedAuthor := reporeader.Author{
			Name:  "gitcha-author-name",
			Email: "gitcha-author-email@gitcha.com",
		}

		repoReader, err := reporeader.NewRepoReaderRepository(repo)
		require.NoError(t, err)

		actual, err := repoReader.GetRepoDetails()

		assert.NoError(t, err)
		assert.Contains(t, actual.AuthorsCommits, expectedAuthor.Email)
		assert.ElementsMatch(t, actual.AuthorsCommits[expectedAuthor.Email], expectedCommits)
	})

	t.Run("given multiple commit authors should return map with each author email as key", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		_, repo, err := gittest.CreateBasicMultiAuthorRepo(ctx, t)
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

		assert.Contains(t, actual.AuthorsCommits, expectedAuthor1.Email)
		assert.Contains(t, actual.AuthorsCommits, expectedAuthor2.Email)
		assert.Contains(t, actual.AuthorsCommits, expectedAuthor3.Email)
		assert.Contains(t, actual.AuthorsCommits, expectedAuthor4.Email)
	})

	t.Run("given basic repository with LICENSE file at root should return MIT license and nil error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		_, repo, err := gittest.CreateBasicRepo(ctx, t)
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
		_, repo, err := gittest.CreateBasicRepo(ctx, t)
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
		_, repo, err := gittest.CreateBasicRepo(ctx, t)
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

	t.Run("given single commit author should return map with author commits", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		_, repo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		head, err := repo.Head()
		require.NoError(t, err)

		commits := make([]*object.Commit, 0)
		cIter, _ := repo.Log(&git.LogOptions{From: head.Hash(), Order: git.LogOrderCommitterTime})
		_ = cIter.ForEach(func(c *object.Commit) error {
			commits = append(commits, c)
			return nil
		})

		expectedCommits := make([]reporeader.Commit, 0, len(commits))
		for _, commit := range commits {
			author := reporeader.Author{
				Name:  commit.Author.Name,
				Email: commit.Author.Email,
			}
			commit := reporeader.Commit{
				Author:  author,
				Message: commit.Message,
				Hash:    commit.Hash.String(),
			}

			expectedCommits = append(expectedCommits, commit)
		}

		expectedAuthor := reporeader.Author{
			Name:  "gitcha-author-name",
			Email: "gitcha-author-email@gitcha.com",
		}

		repoReader, err := reporeader.NewRepoReaderRepository(repo)
		require.NoError(t, err)

		actual, err := repoReader.GetRepoDetails()

		assert.NoError(t, err)
		assert.Contains(t, actual.AuthorsCommits, expectedAuthor.Email)
		assert.ElementsMatch(t, actual.AuthorsCommits[expectedAuthor.Email], expectedCommits)
	})

	t.Run("given multiple commit authors should return map with each author email as key", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		_, repo, err := gittest.CreateBasicMultiAuthorRepo(ctx, t)
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

		assert.Contains(t, actual, expectedAuthor1.Email)
		assert.Contains(t, actual, expectedAuthor2.Email)
		assert.Contains(t, actual, expectedAuthor3.Email)
		assert.Contains(t, actual, expectedAuthor4.Email)
	})

	t.Run("given multiple commit authors with pseudonyms should return map with each author email as key and the total number of their commits", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		repo, err := gittest.CreateMultiNamedAuthorRepo(ctx, t)
		require.NoError(t, err)

		repoReader, err := reporeader.NewRepoReaderRepository(repo)
		require.NoError(t, err)

		const expectedAuthorEmail1 = "gitcha1@gitcha.com"
		const expectedAuthorCommitCount1 = 1

		const expectedAuthorEmail2 = "gitcha2@gitcha.com"
		const expectedAuthorCommitCount2 = 2

		const expectedAuthorEmail3 = "gitcha3@gitcha.com"
		const expectedAuthorCommitCount3 = 3

		const expectedAuthorEmail4 = "gitcha4@gitcha.com"
		const expectedAuthorCommitCount4 = 4

		actual, err := repoReader.GetAuthorsByCommits()
		assert.NoError(t, err)

		assert.Contains(t, actual, expectedAuthorEmail1)
		assert.Len(t, actual[expectedAuthorEmail1], expectedAuthorCommitCount1)

		assert.Contains(t, actual, expectedAuthorEmail2)
		assert.Len(t, actual[expectedAuthorEmail2], expectedAuthorCommitCount2)

		assert.Contains(t, actual, expectedAuthorEmail3)
		assert.Len(t, actual[expectedAuthorEmail3], expectedAuthorCommitCount3)

		assert.Contains(t, actual, expectedAuthorEmail4)
		assert.Len(t, actual[expectedAuthorEmail4], expectedAuthorCommitCount4)
	})
}

func TestRepoReader_GetLicense(t *testing.T) {
	t.Parallel()

	t.Run("given basic repository with LICENSE file at root should return MIT license and nil error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		_, repo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		repoReader, err := reporeader.NewRepoReaderRepository(repo)
		require.NoError(t, err)

		actual, err := repoReader.GetLicense()

		assert.Equal(t, "MIT", actual)
		assert.NoError(t, err)
	})

	t.Run("given basic repository with LICENSE.md file at root should return MIT license and nil error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		_, repo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		wt, err := repo.Worktree()
		require.NoError(t, err)
		fs := wt.Filesystem

		require.NoError(t, fs.Rename("LICENSE", "LICENSE.md"))

		repoReader, err := reporeader.NewRepoReaderRepository(repo)
		require.NoError(t, err)

		actual, err := repoReader.GetLicense()

		assert.Equal(t, "MIT", actual)
		assert.NoError(t, err)
	})

	t.Run("given repository with no LICENSE file should return NO LICENSE string and nil error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		_, repo, err := gittest.CreateBasicRepo(ctx, t)
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
		_, repo, err := gittest.CreateEmptyRepo(ctx, t)
		require.Error(t, err)

		expectedErr := fmt.Errorf("ValidateRepository: received a repository without a head")
		head, err := reporeader.ValidateRepository(repo)

		assert.ErrorContains(t, err, expectedErr.Error())
		assert.Nil(t, head)
	})

	t.Run("given valid repository should return head reference and nil error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		_, repo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		expected, err := repo.Head()
		require.NoError(t, err)

		head, err := reporeader.ValidateRepository(repo)

		assert.NoError(t, err)
		assert.Equal(t, expected, head)
	})
}
