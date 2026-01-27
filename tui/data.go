/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package tui

import (
	"os"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/log"
	"github.com/garrettkrohn/treekanga/adapters"
	"github.com/garrettkrohn/treekanga/config"
	"github.com/garrettkrohn/treekanga/models"
	"github.com/garrettkrohn/treekanga/transformer"
)

// BuildWorktreeTableRows fetches and transforms worktree data into table rows
func BuildWorktreeTableRows(git adapters.GitAdapter, appConfig config.AppConfig) ([]table.Row, error) {
	var rawWorktrees []string
	var err error

	if appConfig.BareRepoPath != "" {
		log.Debug("Using bare repo path for worktree list", "path", appConfig.BareRepoPath)
		rawWorktrees, err = git.GetWorktrees(&appConfig.BareRepoPath)
	} else {
		log.Debug("No bare repo path set, using current directory")
		rawWorktrees, err = git.GetWorktrees(nil)
	}

	if err != nil {
		return nil, err
	}

	worktreetransformer := transformer.NewTransformer()
	worktreeObjects := worktreetransformer.TransformWorktrees(rawWorktrees)

	// Sort worktrees by most recently modified
	sortWorktreesByModTime(worktreeObjects)

	var worktreeBranches []table.Row
	for _, worktree := range worktreeObjects {
		worktreeBranches = append(worktreeBranches, table.Row{worktree.Folder, worktree.BranchName, worktree.FullPath, worktree.CommitHash})
	}

	return worktreeBranches, nil
}

// getPopupItems returns the list of items to display in the popup
func getPopupItems(zoxideEntries []string) []list.Item {
	var returnItems []list.Item
	for _, item := range zoxideEntries {
		returnItems = append(returnItems, popupItem{
			title: item,
			desc:  "", // add description if needed
		})
	}
	return returnItems
}

// sortWorktreesByModTime sorts worktrees by modification time (most recent first)
func sortWorktreesByModTime(worktrees []models.Worktree) {
	sort.Slice(worktrees, func(i, j int) bool {
		statI, errI := os.Stat(worktrees[i].FullPath)
		statJ, errJ := os.Stat(worktrees[j].FullPath)

		// If there's an error accessing either path, push it to the end
		if errI != nil {
			log.Debug("Error stat'ing worktree", "path", worktrees[i].FullPath, "error", errI)
			return false
		}
		if errJ != nil {
			log.Debug("Error stat'ing worktree", "path", worktrees[j].FullPath, "error", errJ)
			return true
		}

		// Sort by modification time, most recent first
		return statI.ModTime().After(statJ.ModTime())
	})
}
