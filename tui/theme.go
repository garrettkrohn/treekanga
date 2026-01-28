/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package tui

import "github.com/charmbracelet/lipgloss"

// Theme defines all colors used in the application UI.
type Theme struct {
	Accent    lipgloss.Color
	AccentFg  lipgloss.Color // Foreground color for text on Accent background
	AccentDim lipgloss.Color
	Border    lipgloss.Color
	BorderDim lipgloss.Color
	MutedFg   lipgloss.Color
	TextFg    lipgloss.Color
	SuccessFg lipgloss.Color
	WarnFg    lipgloss.Color
	ErrorFg   lipgloss.Color
	Cyan      lipgloss.Color
}

// Catppuccin returns the Catppuccin theme (Latte)
func Catppuccin() *Theme {
	return &Theme{
		Accent:    lipgloss.Color("#40a02b"), // Mauve
		AccentFg:  lipgloss.Color("#eff1f5"), // Base (light text on accent)
		AccentDim: lipgloss.Color("#ccd0da"), // Surface 0
		Border:    lipgloss.Color("#7c7f93"), // Overlay 2
		BorderDim: lipgloss.Color("#bcc0cc"), // Surface 1
		MutedFg:   lipgloss.Color("#6c6f85"), // Subtext 0
		TextFg:    lipgloss.Color("#4c4f69"), // Text
		SuccessFg: lipgloss.Color("#40a02b"), // Green
		WarnFg:    lipgloss.Color("#df8e1d"), // Yellow
		ErrorFg:   lipgloss.Color("#d20f39"), // Red
		Cyan:      lipgloss.Color("#209fb5"), // Sapphire
	}
}

// RosePine returns the Rosé Pine theme (Dark)
func RosePine() *Theme {
	return &Theme{
		Accent:    lipgloss.Color("#C4A7E7"), // Iris
		AccentFg:  lipgloss.Color("#191724"), // Dark text on accent
		AccentDim: lipgloss.Color("#26233A"), // Selection
		Border:    lipgloss.Color("#403D52"), // Border
		BorderDim: lipgloss.Color("#26233A"), // Selection
		MutedFg:   lipgloss.Color("#6E6A86"), // Muted
		TextFg:    lipgloss.Color("#E0DEF4"), // Foreground
		SuccessFg: lipgloss.Color("#9CCFD8"), // Foam
		WarnFg:    lipgloss.Color("#F6C177"), // Gold
		ErrorFg:   lipgloss.Color("#EB6F92"), // Love
		Cyan:      lipgloss.Color("#31748F"), // Pine
	}
}

// DefaultTheme returns the default theme
func DefaultTheme() *Theme {
	return Catppuccin()
}
