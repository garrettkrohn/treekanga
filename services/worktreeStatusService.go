package services

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/models"
)

// FetchDefaultBranch fetches origin/<defaultBranch> so subsequent merge and
// ahead/behind comparisons reflect the remote's current state (R5).
func FetchDefaultBranch(bareRepoPath, defaultBranch string) error {
	return git.Fetch(bareRepoPath, defaultBranch)
}

// ComputeWorktreeStatus fills in the R1-R4 status fields on a worktree by
// shelling out to git. defaultBranch is AppConfig.BaseBranch; callers should
// fetch it first via FetchDefaultBranch for an up-to-date merge comparison.
func ComputeWorktreeStatus(worktree models.Worktree, defaultBranch string) models.Worktree {
	staged, modified, untracked, err := git.GetWorkingTreeStatus(worktree.FullPath)
	if err != nil {
		log.Debug("Failed to get working tree status", "worktree", worktree.Folder, "error", err)
	}
	worktree.HasStaged = staged
	worktree.HasModified = modified
	worktree.HasUntracked = untracked

	aheadDefault, behindDefault, err := git.GetAheadBehind(worktree.FullPath, defaultBranch)
	if err != nil {
		log.Debug("Failed to get ahead/behind default branch", "worktree", worktree.Folder, "error", err)
	}
	worktree.AheadDefault = aheadDefault
	worktree.BehindDefault = behindDefault

	upstream, err := git.GetUpstreamBranch(worktree.FullPath)
	if err != nil {
		log.Debug("Failed to get upstream branch", "worktree", worktree.Folder, "error", err)
	}
	if upstream != "" {
		worktree.HasUpstream = true
		aheadRemote, behindRemote, err := git.GetAheadBehind(worktree.FullPath, upstream)
		if err != nil {
			log.Debug("Failed to get ahead/behind remote", "worktree", worktree.Folder, "error", err)
		}
		worktree.AheadRemote = aheadRemote
		worktree.BehindRemote = behindRemote
	}

	targetRef := fmt.Sprintf("origin/%s", defaultBranch)
	merged, err := git.IsMerged(worktree.FullPath, worktree.BranchName, targetRef)
	if err != nil {
		log.Debug("Failed to compute merge status", "worktree", worktree.Folder, "error", err)
		worktree.Merged = models.MergeStatusUnknown
	} else if merged {
		worktree.Merged = models.MergeStatusMerged
	} else {
		worktree.Merged = models.MergeStatusNotMerged
	}

	worktree.StatusLoaded = true
	return worktree
}

// ComputeAllWorktreeStatuses fetches the default branch once, then computes
// status for every worktree. Intended for the CLI's synchronous -v path.
func ComputeAllWorktreeStatuses(bareRepoPath, defaultBranch string, worktrees []models.Worktree) []models.Worktree {
	if err := FetchDefaultBranch(bareRepoPath, defaultBranch); err != nil {
		log.Debug("Failed to fetch default branch before computing status", "branch", defaultBranch, "error", err)
	}

	result := make([]models.Worktree, len(worktrees))
	for i, wt := range worktrees {
		result[i] = ComputeWorktreeStatus(wt, defaultBranch)
	}
	return result
}
