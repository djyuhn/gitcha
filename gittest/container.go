package gittest

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
)

func CreateBasicRepo(ctx context.Context, t *testing.T) (*git.Repository, error) {
	t.Helper()

	const exposedPort = "9418/tcp"

	pool, err := dockertest.NewPool("")
	require.NoError(t, err)
	require.NoError(t, pool.Client.Ping())

	buildOpts := docker.BuildImageOptions{
		Name:         "gitcha/basic_repo_single_author",
		Dockerfile:   "basic_repo_single_author.Dockerfile",
		OutputStream: io.Discard,
		ContextDir:   "testdata/image",
	}
	require.NoError(t, pool.Client.BuildImage(buildOpts))

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository:   "gitcha/basic_repo_single_author",
		ExposedPorts: []string{exposedPort},
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		err := pool.Purge(resource)
		require.NoError(t, err)
	})

	port := resource.GetHostPort(exposedPort)
	url := fmt.Sprintf("git://%s/testdata", port)
	testDir := t.TempDir()
	repo, err := git.PlainCloneContext(ctx, testDir, false, &git.CloneOptions{
		URL: url,
	})

	return repo, err
}

func CreateBasicMultiAuthorRepo(ctx context.Context, t *testing.T) (*git.Repository, error) {
	t.Helper()

	const exposedPort = "9418/tcp"

	pool, err := dockertest.NewPool("")
	require.NoError(t, err)
	require.NoError(t, pool.Client.Ping())

	buildOpts := docker.BuildImageOptions{
		Name:         "gitcha/basic_repo_multiple_authors",
		Dockerfile:   "basic_repo_multiple_authors.Dockerfile",
		OutputStream: io.Discard,
		ContextDir:   "testdata/image",
	}
	require.NoError(t, pool.Client.BuildImage(buildOpts))

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository:   "gitcha/basic_repo_multiple_authors",
		ExposedPorts: []string{exposedPort},
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, pool.Purge(resource))
	})

	port := resource.GetHostPort(exposedPort)
	url := fmt.Sprintf("git://%s/testdata", port)
	testDir := t.TempDir()
	repo, err := git.PlainCloneContext(ctx, testDir, false, &git.CloneOptions{
		URL: url,
	})

	return repo, err
}
