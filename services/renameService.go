package services

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/garrettkrohn/treekanga/config"
	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/utility"
)

// RenameWorktree renames the current worktree's branch and folder
func RenameWorktree(cfg config.AppConfig, newBranchName, currentWorktreePath string) error {
	log.Debug("Starting worktree rename", "newBranchName", newBranchName)

	// Get current branch
	currentBranch, err := git.GetCurrentBranch(currentWorktreePath)
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	if currentBranch == "" {
		return fmt.Errorf("not on a branch (detached HEAD state) - cannot rename")
	}

	log.Debug("Current branch", "branch", currentBranch)

	// Validate new branch doesn't already exist
	localBranches, err := git.GetLocalBranches(cfg.BareRepoPath)
	if err != nil {
		return fmt.Errorf("failed to get local branches: %w", err)
	}

	remoteBranches, err := git.GetRemoteBranches(cfg.BareRepoPath)
	if err != nil {
		return fmt.Errorf("failed to get remote branches: %w", err)
	}

	if slices.Contains(localBranches, newBranchName) {
		return fmt.Errorf("branch '%s' already exists locally", newBranchName)
	}

	if slices.Contains(remoteBranches, newBranchName) {
		return fmt.Errorf("branch '%s' already exists on remote", newBranchName)
	}

	// Sanitize new branch name for folder (replace / with -)
	newFolderName := strings.ReplaceAll(newBranchName, "/", "-")
	newWorktreePath := filepath.Join(cfg.WorktreeTargetDir, newFolderName)

	// Check if target folder already exists
	if _, err := os.Stat(newWorktreePath); err == nil {
		return fmt.Errorf("target folder '%s' already exists", newWorktreePath)
	}

	log.Debug("Renaming branch", "from", currentBranch, "to", newBranchName)

	// Rename the branch
	err = git.RenameBranch(cfg.BareRepoPath, currentBranch, newBranchName)
	if err != nil {
		return fmt.Errorf("failed to rename branch: %w", err)
	}

	log.Debug("Moving worktree", "from", currentWorktreePath, "to", newWorktreePath)

	// Move the worktree folder
	err = git.MoveWorktree(cfg.BareRepoPath, currentWorktreePath, newWorktreePath)
	if err != nil {
		// Try to rollback the branch rename
		rollbackErr := git.RenameBranch(cfg.BareRepoPath, newBranchName, currentBranch)
		if rollbackErr != nil {
			log.Error("Failed to rollback branch rename after worktree move failure", "error", rollbackErr)
		}
		return fmt.Errorf("failed to move worktree: %w", err)
	}

	log.Info("Worktree renamed successfully",
		"oldBranch", currentBranch,
		"newBranch", newBranchName,
		"newPath", newWorktreePath)

	// Inform user about path change
	fmt.Printf("\n✓ Worktree renamed successfully!\n")
	fmt.Printf("  Branch: %s → %s\n", currentBranch, newBranchName)
	fmt.Printf("  Folder: %s → %s\n", filepath.Base(currentWorktreePath), newFolderName)
	fmt.Printf("\nNote: Your current directory is now invalid. Navigate to the new location:\n")
	fmt.Printf("  cd %s\n\n", newWorktreePath)

	return nil
}

// GetCurrentWorktreePath returns the current worktree path by calling GetBareRepoPath and resolving
func GetCurrentWorktreePath() (string, error) {
	// Get the git common dir (which points to the bare repo or main .git)
	gitCommonDir, err := git.GetBareRepoPath("")
	if err != nil {
		return "", fmt.Errorf("failed to get git directory: %w", err)
	}

	// Check if we're in a worktree or bare repo
	// In a worktree, gitCommonDir will be something like /path/to/bare.git/worktrees/branch-name
	// We want the actual worktree directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	// If gitCommonDir contains "/worktrees/", we're in a worktree
	if strings.Contains(gitCommonDir, "/worktrees/") {
		// We're in a worktree, return the current directory
		return currentDir, nil
	}

	// Check if we're in a bare repo by seeing if we have a working tree
	_, err = git.GetCurrentBranch(currentDir)
	if err != nil {
		return "", fmt.Errorf("not in a worktree - cannot rename from bare repository")
	}

	return currentDir, nil
}

// SetConfigForRenameService sets up configuration for rename command
func SetConfigForRenameService(cfg config.AppConfig) (config.AppConfig, error) {
	log.Debug("Running configuration for rename command")

	// Get current worktree path
	currentWorktreePath, err := GetCurrentWorktreePath()
	if err != nil {
		return cfg, err
	}

	// Determine WorktreeTargetDir from current path
	cfg.WorktreeTargetDir = filepath.Dir(currentWorktreePath)
	log.Debug("Set WorktreeTargetDir", "dir", cfg.WorktreeTargetDir)

	return cfg, nil
}

// ValidateRenameArgs validates the arguments for rename command
func ValidateRenameArgs(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("please provide new branch name as an argument")
	}

	if len(args) > 1 {
		return "", fmt.Errorf("too many arguments - expected 1, got %d", len(args))
	}

	newBranchName := strings.TrimSpace(args[0])
	if newBranchName == "" {
		return "", fmt.Errorf("branch name cannot be empty")
	}

	return newBranchName, nil
}

// ExecuteRename executes the full rename workflow
func ExecuteRename(cfg config.AppConfig, args []string) error {
	// Validate arguments
	newBranchName, err := ValidateRenameArgs(args)
	if err != nil {
		return err
	}

	// Set up configuration
	cfg, err = SetConfigForRenameService(cfg)
	if err != nil {
		return err
	}

	// Get current worktree path
	currentWorktreePath, err := GetCurrentWorktreePath()
	if err != nil {
		return err
	}

	// Execute rename
	err = RenameWorktree(cfg, newBranchName, currentWorktreePath)
	utility.CheckError(err)

	return nil
}
