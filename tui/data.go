/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package tui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/garrettkrohn/treekanga/config"
	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/models"
	"github.com/garrettkrohn/treekanga/transformer"
	"github.com/garrettkrohn/treekanga/util"
)

// statusPlaceholder is shown in status columns before background loading
// resolves the real value (R9).
const statusPlaceholder = "…"

// FetchWorktrees fetches and sorts worktree data without computing status.
func FetchWorktrees(appConfig config.AppConfig) ([]models.Worktree, error) {
	rawWorktrees, err := git.ListWorktrees(appConfig.BareRepoPath)
	if err != nil {
		return nil, err
	}

	worktreeObjects := transformer.TransformWorktrees(rawWorktrees)
	util.SortWorktreesByModTime(worktreeObjects)

	return worktreeObjects, nil
}

// BuildWorktreeTableRows fetches worktree data into table rows. Status
// columns render placeholders until updated in the background (R9).
func BuildWorktreeTableRows(appConfig config.AppConfig) ([]table.Row, error) {
	worktreeObjects, err := FetchWorktrees(appConfig)
	if err != nil {
		return nil, err
	}

	return WorktreeTableRows(worktreeObjects), nil
}

// WorktreeTableRows renders table rows for already-fetched worktrees,
// showing computed status when available and a placeholder otherwise.
func WorktreeTableRows(worktrees []models.Worktree) []table.Row {
	rows := make([]table.Row, 0, len(worktrees))
	for _, worktree := range worktrees {
		rows = append(rows, table.Row{
			worktree.Folder,
			worktree.BranchName,
			worktree.FullPath,
			worktree.CommitHash,
			statusOrPlaceholder(worktree, transformer.DirtySymbols),
			statusOrPlaceholder(worktree, transformer.DefaultAheadBehindSymbols),
			statusOrPlaceholder(worktree, transformer.RemoteAheadBehindSymbols),
			statusOrPlaceholder(worktree, transformer.MergedSymbol),
		})
	}
	return rows
}

func statusOrPlaceholder(worktree models.Worktree, render func(models.Worktree) string) string {
	if !worktree.StatusLoaded {
		return statusPlaceholder
	}
	return render(worktree)
}
