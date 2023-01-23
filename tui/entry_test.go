package tui_test

import (
	"testing"

	"github.com/djyuhn/gitcha/tui"

	"github.com/stretchr/testify/assert"
)

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

	t.Run("should return same model and nil msg", func(t *testing.T) {
		t.Parallel()
		expectedModel := tui.EntryModel{}

		model, cmd := expectedModel.Update(nil)

		assert.Equal(t, expectedModel, model)
		assert.Nil(t, cmd)
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
