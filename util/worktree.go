package util

import (
	"os"
	"strings"
	"time"

	"github.com/garrettkrohn/treekanga/models"
)

// ParseWorktrees parses git worktree list output into Worktree structs
func ParseWorktrees(worktreeStrings []string) []models.Worktree {
	var worktrees []models.Worktree

	for _, worktreeString := range worktreeStrings {
		parts := strings.Fields(worktreeString)

		// Skip bare repo and empty lines
		if len(parts) < 3 {
			continue
		}

		fullPath := parts[0]
		commitHash := parts[1]
		
		// Get folder name from path
		folder := strings.Split(fullPath, "/")[len(strings.Split(fullPath, "/"))-1]
		
		// Remove brackets from branch name
		branchName := strings.Trim(parts[2], "[]")

		worktrees = append(worktrees, models.Worktree{
			FullPath:   fullPath,
			Folder:     folder,
			BranchName: branchName,
			CommitHash: commitHash,
		})
	}

	return worktrees
}

// SortWorktreesByModTime sorts worktrees by modification time (most recent first)
func SortWorktreesByModTime(worktrees []models.Worktree) {
	// Get mod times for all worktrees
	modTimes := make(map[string]time.Time)
	for _, wt := range worktrees {
		if info, err := os.Stat(wt.FullPath); err == nil {
			modTimes[wt.FullPath] = info.ModTime()
		}
	}
	
	// Sort by mod time
	for i := 0; i < len(worktrees); i++ {
		for j := i + 1; j < len(worktrees); j++ {
			if modTimes[worktrees[i].FullPath].Before(modTimes[worktrees[j].FullPath]) {
				worktrees[i], worktrees[j] = worktrees[j], worktrees[i]
			}
		}
	}
}

// FilterStaleBranches filters worktrees to only those without remote branches
func FilterStaleBranches(worktrees []models.Worktree, remoteBranches []string) []models.Worktree {
	var filtered []models.Worktree

	for _, wt := range worktrees {
		isStale := true
		for _, remoteBranch := range remoteBranches {
			if wt.BranchName == remoteBranch {
				isStale = false
				break
			}
		}
		if isStale {
			filtered = append(filtered, wt)
		}
	}

	return filtered
}

// SanitizeForSessionName sanitizes a string to be safe for tmux session names
// Replaces characters that are problematic in tmux session names
func SanitizeForSessionName(name string) string {
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, ":", "-")
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ReplaceAll(name, " ", "_")
	return name
}
