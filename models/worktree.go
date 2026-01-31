package models

import "github.com/charmbracelet/lipgloss"

type Worktree struct {
	FullPath   string
	Folder     string
	BranchName string
	CommitHash string
}

type CustomThemeData struct {
	Base      string
	Accent    string
	AccentFg  string
	AccentDim string
	Border    string
	BorderDim string
	MutedFg   string
	TextFg    string
	SuccessFg string
	WarnFg    string
	ErrorFg   string
	Cyan      string
}

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
