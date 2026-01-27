package filter

import (
	"slices"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/garrettkrohn/treekanga/models"
)

type Filter interface {
	GetBranchNoMatchList([]string, []models.Worktree) []models.Worktree
	GetBranchMatchList([]string, []models.Worktree) []models.Worktree
}

type RealFilter struct{}

func NewFilter() *RealFilter {
	return &RealFilter{}
}

func (f *RealFilter) GetBranchNoMatchList(remoteBranches []string, worktreeBranches []models.Worktree) []models.Worktree {
	log.Debug(remoteBranches, worktreeBranches)
	var nonMatchingWorktrees []models.Worktree

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

func (f *RealFilter) GetBranchMatchList(selectedBranchNames []string, allWorktrees []models.Worktree) []models.Worktree {
	var selectedWorktreeObj []models.Worktree
	for _, worktree := range allWorktrees {
		if slices.Contains(selectedBranchNames, worktree.Folder) {
			selectedWorktreeObj = append(selectedWorktreeObj, worktree)
		}
	}
	return selectedWorktreeObj
}
