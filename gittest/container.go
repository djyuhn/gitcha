package gittest

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
)

func CreateEmptyRepo(ctx context.Context, t *testing.T) (*git.Repository, error) {
	t.Helper()

	const exposedPort = "9418/tcp"

	pool, err := dockertest.NewPool("")
	require.NoError(t, err)
	require.NoError(t, pool.Client.Ping())

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository:   "gitcha/empty_repo",
		ExposedPorts: []string{exposedPort},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
	})
	require.NoError(t, err)

	const containerLifeTimeSeconds = 60
	require.NoError(t, resource.Expire(containerLifeTimeSeconds))

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

func CreateBasicRepo(ctx context.Context, t *testing.T) (*git.Repository, error) {
	t.Helper()

	const exposedPort = "9418/tcp"

	pool, err := dockertest.NewPool("")
	require.NoError(t, err)
	require.NoError(t, pool.Client.Ping())

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository:   "gitcha/basic_repo_single_author",
		ExposedPorts: []string{exposedPort},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
	})
	require.NoError(t, err)

	const containerLifeTimeSeconds = 60
	require.NoError(t, resource.Expire(containerLifeTimeSeconds))

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

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository:   "gitcha/basic_repo_multiple_authors",
		ExposedPorts: []string{exposedPort},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
	})
	require.NoError(t, err)

	const containerLifeTimeSeconds = 30
	require.NoError(t, resource.Expire(containerLifeTimeSeconds))

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

func CreateMultiNamedAuthorRepo(ctx context.Context, t *testing.T) (*git.Repository, error) {
	t.Helper()

	const exposedPort = "9418/tcp"

	pool, err := dockertest.NewPool("")
	require.NoError(t, err)
	require.NoError(t, pool.Client.Ping())

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository:   "gitcha/repo_multi_named_authors",
		ExposedPorts: []string{exposedPort},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
	})
	require.NoError(t, err)

	const containerLifeTimeSeconds = 30
	require.NoError(t, resource.Expire(containerLifeTimeSeconds))

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
