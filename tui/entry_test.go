package tui_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

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

		expected := tui.EntryModel{RepoReader: *repoReader}
		actual, err := tui.NewEntryModel(repoReader)

		assert.Equal(t, expected, actual)
		assert.NoError(t, err)
	})
}

func TestEntryModel_Init(t *testing.T) {
	t.Parallel()

	t.Run("should return RepoDetailsMsg", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repo, err := gittest.CreateBasicRepo(ctx, t)
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

		assert.Equal(t, expectedMsg, cmd())
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

	t.Run("given RepoDetailsMsg with no error should update entry model RepoDetails", func(t *testing.T) {
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

		model := tui.EntryModel{}

		updatedModel, cmd := model.Update(msg)

		actual, ok := updatedModel.(tui.EntryModel)
		require.True(t, ok)

		assert.Equal(t, msg.RepoDetails, actual.RepoDetails)
		assert.Nil(t, cmd)
	})
}

func TestEntryModel_View(t *testing.T) {
	t.Parallel()

	t.Run("should return repo details in separate lines", func(t *testing.T) {
		t.Parallel()

		authorCommits := make(map[reporeader.Author][]reporeader.Commit)
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
		authorCommits[author] = commits

		repoDetails := reporeader.RepoDetails{
			CreatedDate:    time.Date(2023, time.January, 26, 3, 2, 1, 0, time.UTC),
			AuthorsCommits: authorCommits,
			License:        "SOME LICENSE",
		}
		model := tui.EntryModel{
			RepoDetails: repoDetails,
		}

		expectedView := strings.Builder{}
		expectedView.WriteString(fmt.Sprintf("Repository Created Date - %s\n", repoDetails.CreatedDate.Format(time.RFC822)))
		expectedView.WriteString(fmt.Sprintf("Repository License - %s\n", repoDetails.License))
		expectedView.WriteString(fmt.Sprintf("Author - %s : Email - %s : Commit count - %d \n", author.Name, author.Email, len(commits)))

		actual := model.View()

		assert.Equal(t, expectedView.String(), actual)
	})
}
