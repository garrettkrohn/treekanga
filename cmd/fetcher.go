package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/models"
	"github.com/garrettkrohn/treekanga/transformer"
)

type fetcher interface {
	fetch() ([]models.Worktree, error)
}

type simpleFetcher struct{}

func getFetcher(global bool) fetcher {
	if global {
		return &globalFetcher{}
	}
	return &simpleFetcher{}
}

func (f *simpleFetcher) fetch() ([]models.Worktree, error) {
	rawWorktrees, err := git.ListWorktrees(deps.AppConfig.BareRepoPath)
	if err != nil {
		return nil, err
	}

	worktreeObjects := transformer.TransformWorktrees(rawWorktrees)
	sortWorktreesByModTime(worktreeObjects)

	return worktreeObjects, nil
}

type globalFetcher struct{}

func (f *globalFetcher) fetch() ([]models.Worktree, error) {
	var allWorktrees []models.Worktree

	for _, worktreeDir := range deps.AppConfig.AllBareRepoPaths {
		bareRepoPath, err := findBareRepoFromWorktreeDir(worktreeDir)
		if err != nil {
			continue
		}

		rawWorktrees, err := git.ListWorktrees(bareRepoPath)
		if err != nil {
			continue
		}

		worktreeObjects := transformer.TransformWorktrees(rawWorktrees)
		allWorktrees = append(allWorktrees, worktreeObjects...)
	}

	sortWorktreesByModTime(allWorktrees)

	return allWorktrees, nil
}

func findBareRepoFromWorktreeDir(targetDir string) (string, error) {
	// First check if there's a .bare directory directly in targetDir
	bareRepoPath := filepath.Join(targetDir, ".bare")
	if info, err := os.Stat(bareRepoPath); err == nil && info.IsDir() {
		return bareRepoPath, nil
	}

	// Otherwise, search for a worktree and extract the bare repo path from it
	entries, err := os.ReadDir(targetDir)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		gitFilePath := filepath.Join(targetDir, entry.Name(), ".git")
		data, err := os.ReadFile(gitFilePath)
		if err != nil {
			continue
		}

		content := string(data)
		if strings.HasPrefix(content, "gitdir:") {
			gitDir := strings.TrimSpace(strings.TrimPrefix(content, "gitdir:"))
			// gitDir is like "/path/.bare/worktrees/branch-name"
			// We need to extract just "/path/.bare"
			// Find the "/worktrees/" part and strip everything after .bare
			if idx := strings.Index(gitDir, "/worktrees/"); idx != -1 {
				bareRepoPath := gitDir[:idx]
				return bareRepoPath, nil
			}
			return gitDir, nil
		}
	}

	return "", fmt.Errorf("no bare repo found for %s", targetDir)
}
