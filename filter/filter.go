package filter

import (
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
		trimmedBranchName := strings.TrimSpace(worktree.Folder)
		log.Debug("Checking worktree branch", "branch", trimmedBranchName)

		if slices.Contains(remoteBranches, trimmedBranchName) {
			log.Debug("branch exists in remote branches", "branch", trimmedBranchName)
		} else {
			for _, remoteBranch := range remoteBranches {
				log.Debug("Comparing with remote branch", "local", trimmedBranchName, "remote", remoteBranch)
			}
			log.Debug("branch not in remote branches", "branch", trimmedBranchName, "remoteBranches", remoteBranches)
			nonMatchingWorktrees = append(nonMatchingWorktrees, worktree)
		}
	}

	return nonMatchingWorktrees
}

func (f *RealFilter) GetBranchMatchList(selectedBranchNames []string, allWorktrees []worktreeobj.WorktreeObj) []worktreeobj.WorktreeObj {
	var selectedWorktreeObj []worktreeobj.WorktreeObj
	for _, worktreeobj := range allWorktrees {
		if slices.Contains(selectedBranchNames, worktreeobj.Folder) {
			selectedWorktreeObj = append(selectedWorktreeObj, worktreeobj)
		}
	}
	return selectedWorktreeObj
}
