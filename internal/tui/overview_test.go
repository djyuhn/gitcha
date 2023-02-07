package tui_test

import (
	"testing"

	"github.com/djyuhn/gitcha/internal/tui"

	"github.com/stretchr/testify/assert"
)

func TestNewOverview(t *testing.T) {
	t.Parallel()

	t.Run("should return default Overview model", func(t *testing.T) {
		t.Parallel()

		actual := tui.NewOverview()

		assert.Equal(t, tui.Overview{}, actual)
	})
}

func TestOverview_Init(t *testing.T) {
	t.Parallel()

	t.Run("should return nil", func(t *testing.T) {
		t.Parallel()

		overview := tui.NewOverview()

		cmd := overview.Init()

		assert.Nil(t, cmd)
	})
}

func TestOverview_Update(t *testing.T) {
	t.Parallel()

	t.Run("given nil msg should return model and nil cmd", func(t *testing.T) {
		t.Parallel()

		model := tui.NewOverview()

		actual, cmd := model.Update(nil)

		assert.Equal(t, model, actual)
		assert.Nil(t, cmd)
	})
}

func TestOverview_View(t *testing.T) {
	t.Parallel()

	t.Run("should return view as OVERVIEW", func(t *testing.T) {
		t.Parallel()

		model := tui.NewOverview()

		expectedView := "OVERVIEW"
		actual := model.View()

		assert.Equal(t, actual, expectedView)
	})
}
