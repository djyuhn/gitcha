package reporeader_test

import (
	"context"
	"fmt"

	"gitcha/gittest"
	"gitcha/reporeader"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing"
	"time"
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
