package services

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/garrettkrohn/treekanga/config"
	"github.com/garrettkrohn/treekanga/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenameWorktree(t *testing.T) {
	// Skip if running in CI without git
	if os.Getenv("SKIP_INTEGRATION_TESTS") != "" {
		t.Skip("Skipping integration test")
	}

	t.Run("successful rename with simple branch name", func(t *testing.T) {
		// Create a temporary directory for test repo
		tempDir, err := os.MkdirTemp("", "treekanga-rename-simple-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Initialize a bare repo
		bareRepoPath := filepath.Join(tempDir, "test.git")
		err = git.CloneBare("https://github.com/octocat/Hello-World.git", bareRepoPath)
		require.NoError(t, err)

		err = git.ConfigureBare(bareRepoPath)
		require.NoError(t, err)

		// Create a worktree
		worktreePath := filepath.Join(tempDir, "old-branch")
		err = git.AddWorktree(bareRepoPath, tempDir, "old-branch", []string{"-b", "old-branch", "origin/master", "--no-track"})
		require.NoError(t, err)

		// Create config
		cfg := config.AppConfig{
			BareRepoPath:     bareRepoPath,
			WorktreeTargetDir: tempDir,
		}

		// Rename the worktree
		err = RenameWorktree(cfg, "new-branch", worktreePath)
		assert.NoError(t, err, "Should successfully rename worktree")

		// Verify branch was renamed
		branches, err := git.GetLocalBranches(bareRepoPath)
		require.NoError(t, err)
		assert.Contains(t, branches, "new-branch", "New branch should exist")
		assert.NotContains(t, branches, "old-branch", "Old branch should not exist")

		// Verify folder was moved
		newPath := filepath.Join(tempDir, "new-branch")
		_, err = os.Stat(newPath)
		assert.NoError(t, err, "New folder should exist")

		_, err = os.Stat(worktreePath)
		assert.True(t, os.IsNotExist(err), "Old folder should not exist")
	})

	t.Run("successful rename with branch containing slashes", func(t *testing.T) {
		// Create a temporary directory for test repo
		tempDir, err := os.MkdirTemp("", "treekanga-rename-slash-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Initialize a bare repo
		bareRepoPath := filepath.Join(tempDir, "test.git")
		err = git.CloneBare("https://github.com/octocat/Hello-World.git", bareRepoPath)
		require.NoError(t, err)

		err = git.ConfigureBare(bareRepoPath)
		require.NoError(t, err)

		// Create a worktree with sanitized folder name
		worktreePath := filepath.Join(tempDir, "old-branch")
		err = git.AddWorktree(bareRepoPath, tempDir, "old-branch", []string{"-b", "old/branch", "origin/master", "--no-track"})
		require.NoError(t, err)

		// Create config
		cfg := config.AppConfig{
			BareRepoPath:     bareRepoPath,
			WorktreeTargetDir: tempDir,
		}

		// Rename to a branch with slashes
		err = RenameWorktree(cfg, "feature/api/users", worktreePath)
		assert.NoError(t, err, "Should successfully rename worktree with slashes")

		// Verify branch was renamed (with slashes preserved)
		branches, err := git.GetLocalBranches(bareRepoPath)
		require.NoError(t, err)
		assert.Contains(t, branches, "feature/api/users", "New branch should exist with slashes")

		// Verify folder was moved (with slashes converted to dashes)
		newPath := filepath.Join(tempDir, "feature-api-users")
		_, err = os.Stat(newPath)
		assert.NoError(t, err, "New folder should exist with sanitized name")
	})

	t.Run("error when new branch already exists locally", func(t *testing.T) {
		// Create a temporary directory for test repo
		tempDir, err := os.MkdirTemp("", "treekanga-rename-exists-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Initialize a bare repo
		bareRepoPath := filepath.Join(tempDir, "test.git")
		err = git.CloneBare("https://github.com/octocat/Hello-World.git", bareRepoPath)
		require.NoError(t, err)

		err = git.ConfigureBare(bareRepoPath)
		require.NoError(t, err)

		// Create two worktrees
		worktreePath := filepath.Join(tempDir, "old-branch")
		err = git.AddWorktree(bareRepoPath, tempDir, "old-branch", []string{"-b", "old-branch", "origin/master", "--no-track"})
		require.NoError(t, err)

		err = git.AddWorktree(bareRepoPath, tempDir, "existing-branch", []string{"-b", "existing-branch", "origin/master", "--no-track"})
		require.NoError(t, err)

		// Create config
		cfg := config.AppConfig{
			BareRepoPath:     bareRepoPath,
			WorktreeTargetDir: tempDir,
		}

		// Try to rename to existing branch
		err = RenameWorktree(cfg, "existing-branch", worktreePath)
		assert.Error(t, err, "Should error when new branch already exists")
		assert.Contains(t, err.Error(), "already exists", "Error should mention branch exists")
	})

	t.Run("error when target folder already exists", func(t *testing.T) {
		// Create a temporary directory for test repo
		tempDir, err := os.MkdirTemp("", "treekanga-rename-folder-exists-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Initialize a bare repo
		bareRepoPath := filepath.Join(tempDir, "test.git")
		err = git.CloneBare("https://github.com/octocat/Hello-World.git", bareRepoPath)
		require.NoError(t, err)

		err = git.ConfigureBare(bareRepoPath)
		require.NoError(t, err)

		// Create a worktree
		worktreePath := filepath.Join(tempDir, "old-branch")
		err = git.AddWorktree(bareRepoPath, tempDir, "old-branch", []string{"-b", "old-branch", "origin/master", "--no-track"})
		require.NoError(t, err)

		// Create a conflicting folder
		conflictingFolder := filepath.Join(tempDir, "new-folder")
		err = os.Mkdir(conflictingFolder, 0755)
		require.NoError(t, err)

		// Create config
		cfg := config.AppConfig{
			BareRepoPath:     bareRepoPath,
			WorktreeTargetDir: tempDir,
		}

		// Try to rename to existing folder
		err = RenameWorktree(cfg, "new/folder", worktreePath)
		assert.Error(t, err, "Should error when target folder already exists")
		assert.Contains(t, err.Error(), "already exists", "Error should mention folder exists")
	})
}
