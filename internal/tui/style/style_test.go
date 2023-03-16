package style_test

import (
	"testing"

	catppuccin "github.com/catppuccin/go"
	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"

	"github.com/djyuhn/gitcha/internal/tui/style"
)

func TestNewTheme(t *testing.T) {
	t.Run("given nil theme config should return default general colors", func(t *testing.T) {
		expected := style.NewDefaultTheme()

		actual := style.NewTheme(nil)

		assert.Equal(t, expected.General, actual.General)
	})

	t.Run("given theme config with general values should return theme with theme config general values", func(t *testing.T) {
		cfg := &style.ThemeConfig{General: style.ThemeGeneral{
			BaseColor:      lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#000000"},
			PrimaryColor:   lipgloss.AdaptiveColor{Light: "#123456", Dark: "#abcdef"},
			SecondaryColor: lipgloss.AdaptiveColor{Light: "#987654", Dark: "#fedcba"},
		}}

		actual := style.NewTheme(cfg)

		assert.Equal(t, cfg.General, actual.General)
	})
}

func TestNewDefaultTheme(t *testing.T) {
	t.Run("should return default general colors", func(t *testing.T) {
		expected := &style.Theme{
			General: style.ThemeGeneral{
				BaseColor: lipgloss.AdaptiveColor{
					Light: catppuccin.Latte.Base().Hex,
					Dark:  catppuccin.Mocha.Base().Hex,
				},
				PrimaryColor: lipgloss.AdaptiveColor{
					Light: catppuccin.Latte.Sapphire().Hex,
					Dark:  catppuccin.Mocha.Sapphire().Hex,
				},
				SecondaryColor: lipgloss.AdaptiveColor{
					Light: catppuccin.Latte.Rosewater().Hex,
					Dark:  catppuccin.Mocha.Rosewater().Hex,
				},
			},
		}

		actual := style.NewDefaultTheme()

		assert.Equal(t, expected.General, actual.General)
	})
}
