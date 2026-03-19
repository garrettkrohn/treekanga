/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package tui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/garrettkrohn/treekanga/config"
	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/transformer"
	"github.com/garrettkrohn/treekanga/util"
)

// BuildWorktreeTableRows fetches and transforms worktree data into table rows
func BuildWorktreeTableRows(appConfig config.AppConfig) ([]table.Row, error) {
	rawWorktrees, err := git.ListWorktrees(appConfig.BareRepoPath)
	if err != nil {
		return nil, err
	}

	worktreeObjects := transformer.TransformWorktrees(rawWorktrees)

	// Sort worktrees by most recently modified
	util.SortWorktreesByModTime(worktreeObjects)

	var worktreeBranches []table.Row
	for _, worktree := range worktreeObjects {
		worktreeBranches = append(worktreeBranches, table.Row{worktree.Folder, worktree.BranchName, worktree.FullPath, worktree.CommitHash})
	}

	return worktreeBranches, nil
}
