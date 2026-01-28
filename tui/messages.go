/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package tui

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
