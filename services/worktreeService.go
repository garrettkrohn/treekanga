package services

import (
	"github.com/charmbracelet/log"
	"github.com/garrettkrohn/treekanga/adapters"
	"github.com/garrettkrohn/treekanga/models"
	"github.com/garrettkrohn/treekanga/transformer"
)

func getWorktrees(git adapters.GitAdapter, transformer *transformer.RealTransformer, bareRepoPath string) []models.Worktree {
	var worktreeStrings []string
	var wError error

	if bareRepoPath != "" {
		worktreeStrings, wError = git.GetWorktrees(&bareRepoPath)
	} else {
		worktreeStrings, wError = git.GetWorktrees(nil)
	}

	if wError != nil {
		log.Fatal(wError)
	}

	worktrees := transformer.TransformWorktrees(worktreeStrings)

	// Sort worktrees by most recently modified
	sortWorktreesByModTime(worktrees)

	return worktrees
}
