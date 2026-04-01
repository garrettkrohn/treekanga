package git

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenameBranch(t *testing.T) {
	// Skip if running in CI without git
	if os.Getenv("SKIP_INTEGRATION_TESTS") != "" {
		t.Skip("Skipping integration test")
	}

	// Create a temporary directory for test repo
	tempDir, err := os.MkdirTemp("", "treekanga-rename-branch-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Initialize a bare repo
	bareRepoPath := filepath.Join(tempDir, "test.git")
	err = runCommand("git", "init", "--bare", bareRepoPath)
	require.NoError(t, err)

	// Configure the bare repo
	err = ConfigureBare(bareRepoPath)
	require.NoError(t, err)

	// Create a worktree with a branch
	worktreePath := filepath.Join(tempDir, "test-branch")
	err = AddWorktree(bareRepoPath, tempDir, "test-branch", []string{"-b", "old-branch"})
	require.NoError(t, err)

	// Create an initial commit so the branch actually exists
	err = runCommand("sh", "-c", fmt.Sprintf("cd %s && echo 'test' > test.txt && git add test.txt && git commit -m 'initial commit'", worktreePath))
	require.NoError(t, err)

	// Rename the branch
	err = RenameBranch(bareRepoPath, "old-branch", "new-branch")
	assert.NoError(t, err, "Should successfully rename branch")

	// Verify the new branch exists
	branches, err := GetLocalBranches(bareRepoPath)
	require.NoError(t, err)
	assert.Contains(t, branches, "new-branch", "New branch should exist")
	assert.NotContains(t, branches, "old-branch", "Old branch should not exist")
}

func TestMoveWorktree(t *testing.T) {
	// Skip if running in CI without git
	if os.Getenv("SKIP_INTEGRATION_TESTS") != "" {
		t.Skip("Skipping integration test")
	}

	// Create a temporary directory for test repo
	tempDir, err := os.MkdirTemp("", "treekanga-move-worktree-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Initialize a bare repo
	bareRepoPath := filepath.Join(tempDir, "test.git")
	err = runCommand("git", "init", "--bare", bareRepoPath)
	require.NoError(t, err)

	// Configure the bare repo
	err = ConfigureBare(bareRepoPath)
	require.NoError(t, err)

	// Create a worktree
	oldPath := filepath.Join(tempDir, "old-folder")
	err = AddWorktree(bareRepoPath, tempDir, "old-folder", []string{"-b", "test-branch"})
	require.NoError(t, err)

	// Verify old path exists
	_, err = os.Stat(oldPath)
	require.NoError(t, err, "Old worktree path should exist")

	// Move the worktree
	newPath := filepath.Join(tempDir, "new-folder")
	err = MoveWorktree(bareRepoPath, oldPath, newPath)
	assert.NoError(t, err, "Should successfully move worktree")

	// Verify new path exists and old path doesn't
	_, err = os.Stat(newPath)
	assert.NoError(t, err, "New worktree path should exist")

	_, err = os.Stat(oldPath)
	assert.True(t, os.IsNotExist(err), "Old worktree path should not exist")
}

func TestGetCurrentBranch(t *testing.T) {
	// Skip if running in CI without git
	if os.Getenv("SKIP_INTEGRATION_TESTS") != "" {
		t.Skip("Skipping integration test")
	}

	// Create a temporary directory for test repo
	tempDir, err := os.MkdirTemp("", "treekanga-current-branch-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Initialize a bare repo
	bareRepoPath := filepath.Join(tempDir, "test.git")
	err = runCommand("git", "init", "--bare", bareRepoPath)
	require.NoError(t, err)

	// Configure the bare repo
	err = ConfigureBare(bareRepoPath)
	require.NoError(t, err)

	// Create a worktree with a specific branch
	worktreePath := filepath.Join(tempDir, "test-branch")
	err = AddWorktree(bareRepoPath, tempDir, "test-branch", []string{"-b", "feature/test-branch"})
	require.NoError(t, err)

	// Get current branch from the worktree
	currentBranch, err := GetCurrentBranch(worktreePath)
	assert.NoError(t, err, "Should successfully get current branch")
	assert.Equal(t, "feature/test-branch", currentBranch, "Should return correct branch name")

	// Bare repos technically can report a branch (the default branch like master/main)
	// but that's different from a worktree's current branch. We just verify it doesn't error.
	_, err = GetCurrentBranch(bareRepoPath)
	// This may or may not error depending on git version, so we just ensure the function completes
	t.Logf("GetCurrentBranch from bare repo returned: %v", err)
}
