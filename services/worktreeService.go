package services

import (
	"github.com/charmbracelet/log"
	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/models"
	"github.com/garrettkrohn/treekanga/util"
)

func getWorktrees(bareRepoPath string) []models.Worktree {
	worktreeStrings, err := git.ListWorktrees(bareRepoPath)
	if err != nil {
		log.Fatal(err)
	}

	worktrees := util.ParseWorktrees(worktreeStrings)

	// Sort worktrees by most recently modified
	util.SortWorktreesByModTime(worktrees)

	return worktrees
}
