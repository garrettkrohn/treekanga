package transformer

import (
	"strings"

	"github.com/garrettkrohn/treekanga/models"
)

type Transformer interface {
	Transform(worktrees []models.Worktree) ([]string, error)
}

func TransformWorktrees(worktreeStrings []string) []models.Worktree {
	var worktrees []models.Worktree

	for _, worktreeString := range worktreeStrings {
		parts := strings.Fields(worktreeString)

		if len(parts) < 3 {
			continue
		}

		fullPath := parts[0]
		commitHash := parts[1]

		folder := strings.Split(fullPath, "/")[len(strings.Split(fullPath, "/"))-1]

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

func RemoveOriginPrefix(branchStrings []string) []string {
	for i, branch := range branchStrings {
		branchStrings[i] = strings.TrimSpace(strings.Replace(branch, "origin/", "", -1))
	}
	return branchStrings
}

func TransformWorktreesToBranchNames(worktreeObjs []models.Worktree) []string {
	var stringWorktrees []string
	for _, worktreeObj := range worktreeObjs {
		stringWorktrees = append(stringWorktrees, worktreeObj.Folder)
	}
	return stringWorktrees
}

func RemoveQuotes(branchStrings []string) []string {
	for i, branch := range branchStrings {
		branchStrings[i] = strings.TrimSpace(strings.Replace(branch, "'", "", -1))
	}
	return branchStrings
}
