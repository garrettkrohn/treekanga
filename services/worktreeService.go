package services

import (
	"github.com/charmbracelet/log"
	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/transformer"
	worktreeobj "github.com/garrettkrohn/treekanga/worktreeObj"
)

func getWorktrees(git git.GitAdapter, transformer *transformer.RealTransformer, bareRepoPath string) []worktreeobj.WorktreeObj {
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
