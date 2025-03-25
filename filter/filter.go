package filter

import (
	"github.com/garrettkrohn/treekanga/worktreeObj"
	"slices"
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
		if !slices.Contains(remoteBranches, worktree.BranchName) {
			nonMatchingWorktrees = append(nonMatchingWorktrees, worktree)
		}
	}

	return nonMatchingWorktrees
}

func (f *RealFilter) BranchExistsInSlice(branches []string, newBranch string) bool {

	return slices.Contains(branches, newBranch)
}

func (f *RealFilter) GetBranchMatchList(selectedBranchNames []string, allWorktrees []worktreeobj.WorktreeObj) []worktreeobj.WorktreeObj {
	var selectedWorktreeObj []worktreeobj.WorktreeObj
	for _, worktreeobj := range allWorktrees {
		if slices.Contains(selectedBranchNames, worktreeobj.BranchName) {
			selectedWorktreeObj = append(selectedWorktreeObj, worktreeobj)
		}
	}
	return selectedWorktreeObj
}
