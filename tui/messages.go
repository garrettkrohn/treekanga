/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package tui

import "github.com/garrettkrohn/treekanga/models"

// statusFetchDoneMsg is sent once the default branch has been fetched from
// origin (R5), signalling it's safe to start computing per-worktree status.
type statusFetchDoneMsg struct{}

// worktreeStatusMsg is sent when a single worktree's status (R1-R4) has
// finished computing in the background, so its table row can be updated
// without blocking the rest of the table (R9).
type worktreeStatusMsg struct {
	fullPath string
	worktree models.Worktree
}

// deleteCompleteMsg is sent when deletion is complete
type deleteCompleteMsg struct {
	err          error
	worktreeName string
	output       string
}

// deleteErrorMsg is sent when deletion fails
type deleteErrorMsg struct {
	err          error
	worktreePath string
	worktreeName string
	branchName   string
	output       string
}

// addCompleteMsg is sent when add worktree is complete
type addCompleteMsg struct {
	err        error
	branchName string
	output     string
}

// addErrorMsg is sent when add worktree fails
type addErrorMsg struct {
	err        error
	branchName string
	output     string
}

// popupItem represents an item in the popup list
type popupItem struct {
	title string
	desc  string
}

func (i popupItem) Title() string       { return i.title }
func (i popupItem) Description() string { return i.desc }
func (i popupItem) FilterValue() string { return i.title }

// branchSelectionReadyMsg is sent when branch list is ready for selection
type branchSelectionReadyMsg struct {
	branches []string
}

// folderSelectionReadyMsg is sent when folder list is ready for selection
type folderSelectionReadyMsg struct {
	folders []string
}
