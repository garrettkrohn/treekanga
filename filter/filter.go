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
		trimmedBranchName := strings.TrimSpace(worktree.BranchName)
		log.Debug(fmt.Sprintf("Checking worktree branch: '%s'", trimmedBranchName))

		if slices.Contains(remoteBranches, trimmedBranchName) {
			log.Debug(fmt.Sprintf("'%s' exists in remote branches", trimmedBranchName))
		} else {
			for _, remoteBranch := range remoteBranches {
				log.Debug(fmt.Sprintf("Comparing '%s' with remote branch '%s'", trimmedBranchName, remoteBranch))
			}
			log.Debug(fmt.Sprintf("'%s' not in remote branches: %s", trimmedBranchName, remoteBranches))
			nonMatchingWorktrees = append(nonMatchingWorktrees, worktree)
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
