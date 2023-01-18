package reporeader_test

import (
	"context"
	"fmt"
	"testing"

	"gitcha/gittest"
	"gitcha/reporeader"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLicense(t *testing.T) {
	t.Parallel()

	t.Run("given an invalid repository should return empty string and error", func(t *testing.T) {
		t.Parallel()

		expectedErr := fmt.Errorf("GetLicense: received an invalid repository")
		actual, err := reporeader.GetLicense(nil)

		assert.Equal(t, "", actual)
		assert.ErrorContains(t, err, expectedErr.Error())
	})

	t.Run("given basic repository with LICENSE file at root should return MIT license and nil error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		repo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		actual, err := reporeader.GetLicense(repo)

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

		actual, err := reporeader.GetLicense(repo)

		assert.Equal(t, "NO LICENSE", actual)
		assert.NoError(t, err)
	})
}
