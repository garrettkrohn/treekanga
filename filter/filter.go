package filter

import (
	"github.com/garrettkrohn/treekanga/worktreeObj"
)

type Filter interface {
	GetBranchNoMatchList([]string, []worktreeobj.WorktreeObj) []worktreeobj.WorktreeObj
	BranchExistsInSlice([]string, string) bool
	GetBranchMatchList([]string, []worktreeobj.WorktreeObj) []worktreeobj.WorktreeObj
}

type RealFilter struct{}

func NewFilter() *RealFilter {
	return &RealFilter{}
}

func (f *RealFilter) GetBranchNoMatchList(remoteBranches []string, worktreeBranches []worktreeobj.WorktreeObj) []worktreeobj.WorktreeObj {
	var nonMatchingWorktrees []worktreeobj.WorktreeObj

	for _, worktree := range worktreeBranches {
		if !contains(remoteBranches, worktree.BranchName) {
			nonMatchingWorktrees = append(nonMatchingWorktrees, worktree)
		}
	}

	return nonMatchingWorktrees
}

func (f *RealFilter) BranchExistsInSlice(branches []string, newBranch string) bool {

	return contains(branches, newBranch)
}

func (f *RealFilter) GetBranchMatchList(selectedBranchNames []string, allWorktrees []worktreeobj.WorktreeObj) []worktreeobj.WorktreeObj {
	var selectedWorktreeObj []worktreeobj.WorktreeObj
	for _, worktreeobj := range allWorktrees {
		for _, str := range selectedBranchNames {
			if worktreeobj.BranchName == str {
				selectedWorktreeObj = append(selectedWorktreeObj, worktreeobj)
				break
			}
		}
	}
	return selectedWorktreeObj
}

// contains checks if a slice contains a specific string.
func contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}
