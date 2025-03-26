package filter

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/garrettkrohn/treekanga/worktreeObj"
)

type Filter interface {
	GetBranchNoMatchList([]string, []worktreeobj.WorktreeObj) []worktreeobj.WorktreeObj
	GetBranchMatchList([]string, []worktreeobj.WorktreeObj) []worktreeobj.WorktreeObj
}

type RealFilter struct{}

func NewFilter() *RealFilter {
	return &RealFilter{}
}

func (f *RealFilter) GetBranchNoMatchList(remoteBranches []string, worktreeBranches []worktreeobj.WorktreeObj) []worktreeobj.WorktreeObj {
	log.Debug(remoteBranches, worktreeBranches)
	var nonMatchingWorktrees []worktreeobj.WorktreeObj

	for _, worktree := range worktreeBranches {
		if !slices.Contains(remoteBranches, strings.TrimSpace(worktree.BranchName)) {
			log.Debug(fmt.Sprintf("%s not in remote branches", worktree.BranchName))
			nonMatchingWorktrees = append(nonMatchingWorktrees, worktree)
		} else {
			log.Debug(fmt.Sprintf("%s exists in remote branches", worktree.BranchName))
		}
	}

	return nonMatchingWorktrees
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
