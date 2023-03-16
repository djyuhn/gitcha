package style

import (
	catppuccin "github.com/catppuccin/go"
	"github.com/charmbracelet/lipgloss"
)

type ThemeGeneral struct {
	BaseColor      lipgloss.AdaptiveColor
	PrimaryColor   lipgloss.AdaptiveColor
	SecondaryColor lipgloss.AdaptiveColor
}

type ThemeConfig struct {
	General ThemeGeneral
}

type Theme struct {
	General ThemeGeneral
}

func NewTheme(cfg *ThemeConfig) *Theme {
	if cfg == nil {
		return NewDefaultTheme()
	}

	configTheme := &Theme{
		General: cfg.General,
	}
	return configTheme
}

func NewDefaultTheme() *Theme {
	defaultTheme := &Theme{
		General: ThemeGeneral{
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
	return defaultTheme
}
