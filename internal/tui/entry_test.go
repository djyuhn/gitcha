package tui_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/djyuhn/gitcha/gittest"
	"github.com/djyuhn/gitcha/internal/reporeader"
	"github.com/djyuhn/gitcha/internal/tui"
	"github.com/djyuhn/gitcha/internal/tui/overview"

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
		_, repo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		repoReader, err := reporeader.NewRepoReaderRepository(repo)
		require.NoError(t, err)

		actual, err := tui.NewEntryModel(repoReader)

		assert.Equal(t, repoReader, &actual.RepoReader)
		assert.NoError(t, err)
	})

	t.Run("given RepoReader should return EntryModel with non default spinner and nil error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		_, repo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		repoReader, err := reporeader.NewRepoReaderRepository(repo)
		require.NoError(t, err)

		actual, err := tui.NewEntryModel(repoReader)

		var basicSpinner spinner.Model

		assert.NotEqual(t, basicSpinner, actual.Spinner)
		assert.NoError(t, err)
	})

	t.Run("given RepoReader should return EntryModel with IsLoading true", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		_, repo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		repoReader, err := reporeader.NewRepoReaderRepository(repo)
		require.NoError(t, err)

		actual, err := tui.NewEntryModel(repoReader)

		assert.True(t, actual.IsLoading)
		assert.NoError(t, err)
	})
}

func TestEntryModel_Init(t *testing.T) {
	t.Parallel()

	t.Run("should return RepoDetailsMsg as part of batched cmds", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		_, repo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		repoReader, err := reporeader.NewRepoReaderRepository(repo)
		require.NoError(t, err)

		entryModel, err := tui.NewEntryModel(repoReader)
		require.NoError(t, err)

		expectedDetails, expectedErr := repoReader.GetRepoDetails()

		expectedMsg := tui.RepoDetailsMsg{
			Err:         expectedErr,
			RepoDetails: expectedDetails,
		}
		cmd := entryModel.Init()
		require.NotNil(t, cmd)

		batchedMsg := cmd()

		assert.IsType(t, tea.BatchMsg{}, batchedMsg)

		var msgs []tea.Msg
		for _, batchedCmd := range batchedMsg.(tea.BatchMsg) {
			msgs = append(msgs, batchedCmd())
		}
		assert.Contains(t, msgs, expectedMsg)
	})

	t.Run("should return spinner tick msg as part of batched cmds", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		_, repo, err := gittest.CreateBasicRepo(ctx, t)
		require.NoError(t, err)

		repoReader, err := reporeader.NewRepoReaderRepository(repo)
		require.NoError(t, err)

		entryModel, err := tui.NewEntryModel(repoReader)
		require.NoError(t, err)

		cmd := entryModel.Init()
		require.NotNil(t, cmd)

		batchedMsg := cmd()

		assert.IsType(t, tea.BatchMsg{}, batchedMsg)

		var spinnerMsg spinner.TickMsg
		for _, batchedCmd := range batchedMsg.(tea.BatchMsg) {
			if msg, ok := batchedCmd().(spinner.TickMsg); ok {
				spinnerMsg = msg
				break
			}
		}
		assert.Equal(t, entryModel.Spinner.ID(), spinnerMsg.ID)
	})
}

func TestEntryModel_Update(t *testing.T) {
	t.Parallel()

	t.Run("given nil msg should return same model and nil msg", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		_, repo, err := gittest.CreateBasicRepo(ctx, t)
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
		_, repo, err := gittest.CreateBasicRepo(ctx, t)
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

	t.Run("given RepoDetailsMsg and Err is nil", func(t *testing.T) {
		t.Parallel()

		t.Run("should update entry model repoDetails and return LoadingRepoMsg as false", func(t *testing.T) {
			t.Parallel()

			repoDetails := reporeader.RepoDetails{
				CreatedDate:    time.Date(2023, time.January, 26, 3, 2, 1, 0, time.UTC),
				AuthorsCommits: nil,
				License:        "SOME LICENSE",
			}

			msg := tui.RepoDetailsMsg{
				Err:         nil,
				RepoDetails: repoDetails,
			}
			LoadingRepoMsg := tui.LoadingRepoMsg{IsLoading: false}

			model := tui.EntryModel{}

			updatedModel, cmd := model.Update(msg)

			actual, ok := updatedModel.(tui.EntryModel)
			require.True(t, ok)

			assert.Equal(t, msg.RepoDetails, actual.RepoDetails)

			require.NotNil(t, cmd)
			require.IsType(t, tui.LoadingRepoMsg{}, cmd())
			assert.Equal(t, LoadingRepoMsg, cmd())
		})

		t.Run("should update Overview model and return LoadingRepoMsg as false", func(t *testing.T) {
			t.Parallel()

			repoDetails := reporeader.RepoDetails{
				CreatedDate:    time.Date(2023, time.January, 26, 3, 2, 1, 0, time.UTC),
				AuthorsCommits: nil,
				License:        "SOME LICENSE",
			}

			msg := tui.RepoDetailsMsg{
				Err:         nil,
				RepoDetails: repoDetails,
			}
			LoadingRepoMsg := tui.LoadingRepoMsg{IsLoading: false}

			model := tui.EntryModel{}

			expectedOverview := overview.NewOverview(repoDetails)
			updatedModel, cmd := model.Update(msg)

			actual, ok := updatedModel.(tui.EntryModel)
			require.True(t, ok)

			assert.Equal(t, expectedOverview, actual.Overview)

			require.NotNil(t, cmd)
			require.IsType(t, tui.LoadingRepoMsg{}, cmd())
			assert.Equal(t, LoadingRepoMsg, cmd())
		})
	})

	t.Run("given RepoDetailsMsg and Err is not nil should update entry model RepoErr and return LoadingRepoMsg as false", func(t *testing.T) {
		t.Parallel()

		msg := tui.RepoDetailsMsg{
			Err:         fmt.Errorf("some error reading repository"),
			RepoDetails: reporeader.RepoDetails{},
		}

		model := tui.EntryModel{}

		updatedModel, cmd := model.Update(msg)

		actual, ok := updatedModel.(tui.EntryModel)
		require.True(t, ok)

		require.NotNil(t, cmd)
		require.IsType(t, tui.LoadingRepoMsg{}, cmd())
		assert.Equal(t, msg.Err, actual.RepoError)
	})

	t.Run("given spinner tick msg and model IsLoading is true should return spinner tick msg", func(t *testing.T) {
		t.Parallel()

		model := tui.EntryModel{IsLoading: true}

		tickMsg := model.Spinner.Tick()

		updatedModel, cmd := model.Update(tickMsg)

		actualModel, ok := updatedModel.(tui.EntryModel)
		require.True(t, ok)

		require.NotNil(t, cmd)
		msg := cmd()

		actualMsg, ok := msg.(spinner.TickMsg)
		require.True(t, ok)
		assert.Equal(t, actualModel.Spinner.ID(), actualMsg.ID)
	})

	t.Run("given spinner tick msg and model IsLoading is false should return nil msg", func(t *testing.T) {
		t.Parallel()

		model := tui.EntryModel{IsLoading: false}

		tickMsg := model.Spinner.Tick()

		updatedModel, cmd := model.Update(tickMsg)

		_, ok := updatedModel.(tui.EntryModel)
		require.True(t, ok)

		assert.Nil(t, cmd)
	})

	t.Run("given LoadingRepoMsg should set IsLoading on model and return nil cmd", func(t *testing.T) {
		t.Parallel()

		model := tui.EntryModel{IsLoading: false}

		loadingRepoMsg := tui.LoadingRepoMsg{IsLoading: true}

		updatedModel, cmd := model.Update(loadingRepoMsg)

		actual, ok := updatedModel.(tui.EntryModel)
		require.True(t, ok)

		assert.Equal(t, loadingRepoMsg.IsLoading, actual.IsLoading)
		assert.Nil(t, cmd)
	})
}

func TestEntryModel_View(t *testing.T) {
	t.Parallel()

	t.Run("given IsLoading is true should show Processing... and spinner view", func(t *testing.T) {
		t.Parallel()

		model := tui.EntryModel{IsLoading: true, Spinner: spinner.New()}

		expectedView := strings.Builder{}
		expectedView.WriteString(model.Spinner.View())
		expectedView.WriteString(" Processing...")
		actual := model.View()

		assert.Contains(t, actual, expectedView.String())
	})

	t.Run("given RepoError is not nil should show message saying error occurred", func(t *testing.T) {
		t.Parallel()

		model := tui.EntryModel{RepoError: fmt.Errorf("some error")}

		expectedView := "An error occurred while processing the repository."
		actual := model.View()

		assert.Contains(t, actual, expectedView)
	})

	t.Run("given not loading should return Overview view", func(t *testing.T) {
		t.Parallel()

		authorCommits := make(map[string][]reporeader.Commit)
		author := reporeader.Author{
			Name:  "FirstName LastName",
			Email: "authorname@gitcha.com",
		}
		commits := []reporeader.Commit{
			{
				Author:  author,
				Message: "commit message",
				Hash:    "someHash",
			},
		}
		authorCommits[author.Email] = commits

		repoDetails := reporeader.RepoDetails{
			CreatedDate:    time.Date(2023, time.January, 26, 3, 2, 1, 0, time.UTC),
			AuthorsCommits: authorCommits,
			License:        "SOME LICENSE",
		}
		model := tui.EntryModel{
			IsLoading: false,
			Overview:  overview.NewOverview(repoDetails),
		}

		actual := model.View()

		assert.Contains(t, actual, model.Overview.View())
	})
}
