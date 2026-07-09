package models

import "github.com/charmbracelet/lipgloss"

// MergeStatus represents whether a worktree's branch content is already
// present in the default branch.
type MergeStatus int

const (
	MergeStatusUnknown MergeStatus = iota
	MergeStatusMerged
	MergeStatusNotMerged
)

type Worktree struct {
	FullPath   string
	Folder     string
	BranchName string
	CommitHash string

	// Working tree state (R1)
	HasStaged    bool
	HasModified  bool
	HasUntracked bool

	// Ahead/behind the default branch (R2)
	AheadDefault int
	BehindDefault int

	// Ahead/behind the remote tracking branch (R3)
	HasUpstream  bool
	AheadRemote  int
	BehindRemote int

	// Merge status against origin/<default-branch> (R4)
	Merged MergeStatus

	// StatusLoaded is true once the R1-R4 fields above have been computed.
	// Used by the TUI to distinguish "not yet loaded" from "loaded, all clear".
	StatusLoaded bool
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
