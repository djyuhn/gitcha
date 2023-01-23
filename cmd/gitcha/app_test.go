package gitcha_test

import (
	"fmt"
	"testing"

	"github.com/djyuhn/gitcha/cmd/gitcha"
	"github.com/djyuhn/gitcha/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockProgram struct {
	mock.Mock
}

func (m *MockProgram) Run() (tea.Model, error) {
	args := m.Called()
	return args.Get(0).(tea.Model), args.Error(1)
}

var _ gitcha.Program = &MockProgram{}

func TestNewApp(t *testing.T) {
	t.Parallel()

	t.Run("given model should return non nil App and nil error", func(t *testing.T) {
		t.Parallel()

		app, err := gitcha.NewApp(tui.EntryModel{})

		assert.NotNil(t, app)
		assert.NoError(t, err)
	})

	t.Run("given model should create new program and return nil error", func(t *testing.T) {
		t.Parallel()

		app, err := gitcha.NewApp(tui.EntryModel{})

		assert.NotNil(t, app.TuiProgram)
		assert.NoError(t, err)
	})
}

func TestNewAppProgram(t *testing.T) {
	t.Parallel()

	t.Run("given nil program should return nil App and error", func(t *testing.T) {
		t.Parallel()

		expectedErr := fmt.Errorf("NewAppProgram: received nil program")
		app, err := gitcha.NewAppProgram(nil)

		assert.Nil(t, app)
		assert.ErrorContains(t, err, expectedErr.Error())
	})

	t.Run("given program should return App and nil error", func(t *testing.T) {
		t.Parallel()

		program := tea.Program{}
		app, err := gitcha.NewAppProgram(&program)

		assert.NotNil(t, app)
		assert.NoError(t, err)
	})
}

func TestApp_GitchaTui(t *testing.T) {
	t.Parallel()

	t.Run("given error when executing program should return error", func(t *testing.T) {
		t.Parallel()

		mockProgram := new(MockProgram)
		mockProgram.On("Run").Return(tui.EntryModel{}, fmt.Errorf("error in program"))

		app, err := gitcha.NewAppProgram(mockProgram)
		require.NoError(t, err)

		expectedError := fmt.Errorf("GitchaTui: attempted to run program and received an error")
		err = app.GitchaTui()

		assert.ErrorContains(t, err, expectedError.Error())
	})

	t.Run("given program runs without error should return nil error", func(t *testing.T) {
		t.Parallel()

		mockProgram := new(MockProgram)
		mockProgram.On("Run").Return(tui.EntryModel{}, nil)

		app, err := gitcha.NewAppProgram(mockProgram)
		require.NoError(t, err)

		err = app.GitchaTui()

		assert.NoError(t, err)
	})
}
