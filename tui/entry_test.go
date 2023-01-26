package tui_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/djyuhn/gitcha/gittest"
	"github.com/djyuhn/gitcha/reporeader"
	"github.com/djyuhn/gitcha/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEntryModel(t *testing.T) {
	t.Parallel()

	t.Run("given nil RepoReader should return default EntryModel and error", func(t *testing.T) {
		t.Parallel()

		expectedError := fmt.Errorf("NewEntryModel: received a nil RepoReader")
		actual, err := tui.NewEntryModel(nil)

		assert.Equal(t, tui.EntryModel{}, actual)
		assert.ErrorContains(t, err, expectedError.Error())
	})

	t.Run("given RepoReader should return EntryModel with RepoReader and nil error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		repoReader, err := reporeader.NewRepoReaderRepository(repo)
		require.NoError(t, err)

		expected := tui.EntryModel{RepoReader: repoReader}
		actual, err := tui.NewEntryModel(repoReader)

		assert.Equal(t, expected, actual)
		assert.NoError(t, err)
	})
}

func TestEntryModel_Init(t *testing.T) {
	t.Parallel()

	t.Run("should return nil command", func(t *testing.T) {
		t.Parallel()
		model := tui.EntryModel{}

		actual := model.Init()

		assert.Nil(t, actual)
	})
}

func TestEntryModel_Update(t *testing.T) {
	t.Parallel()

	t.Run("given nil msg should return same model and nil msg", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		repoReader, err := reporeader.NewRepoReaderRepository(repo)
		require.NoError(t, err)

		model, err := tui.NewEntryModel(repoReader)
		require.NoError(t, err)

		actual, cmd := model.Update(nil)

		assert.Equal(t, model, actual)
		assert.Nil(t, cmd)
	})

	t.Run("given Ctrl+C message should emit quit message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		repoReader, err := reporeader.NewRepoReaderRepository(repo)
		require.NoError(t, err)

		model, err := tui.NewEntryModel(repoReader)
		require.NoError(t, err)

		quitCmd := tea.KeyMsg{
			Type: tea.KeyCtrlC,
		}

		actual, cmd := model.Update(quitCmd)

		assert.Equal(t, model, actual)
		assert.Equal(t, tea.Quit(), cmd())
	})
}

func TestEntryModel_View(t *testing.T) {
	t.Parallel()

	t.Run("should return Entry View string", func(t *testing.T) {
		t.Parallel()
		model := tui.EntryModel{}

		actual := model.View()

		assert.Equal(t, "Entry View", actual)
	})
}
