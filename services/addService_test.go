package services

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/garrettkrohn/treekanga/config"
	"github.com/garrettkrohn/treekanga/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddWorktreeWithPullFlag(t *testing.T) {
	// Skip if running in CI without git
	if os.Getenv("SKIP_INTEGRATION_TESTS") != "" {
		t.Skip("Skipping integration test")
	}

	t.Run("adds worktree from remote base branch with pull flag", func(t *testing.T) {
		// Create a temporary directory for test repo
		tempDir, err := os.MkdirTemp("", "treekanga-add-pull-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Clone a real repo to have remote branches
		bareRepoPath := filepath.Join(tempDir, "test.git")
		err = git.CloneBare("https://github.com/octocat/Hello-World.git", bareRepoPath)
		require.NoError(t, err)

		err = git.ConfigureBare(bareRepoPath)
		require.NoError(t, err)

		// Get the commit hash of origin/master before fetch
		beforeFetchCommit, err := getCommitHash(bareRepoPath, "origin/master")
		require.NoError(t, err)

		// Simulate the config with pull flag enabled
		cfg := config.AppConfig{
			BareRepoPath:               bareRepoPath,
			WorktreeTargetDir:          tempDir,
			NewBranchName:              "feature/new-branch",
			NewWorktreeName:            "feature-new-branch",
			BaseBranch:                 "master",
			BaseBranchExistsLocally:    false,
			BaseBranchExistsRemotely:   true,
			NewBranchExistsLocally:     false,
			NewBranchExistsRemotely:    false,
			PullBeforeCuttingNewBranch: true,
		}

		// Call AddWorktree with connector and shell as nil (not needed for this test)
		AddWorktree(nil, nil, cfg)

		// Verify the worktree was created
		worktreePath := filepath.Join(tempDir, "feature-new-branch")
		_, err = os.Stat(worktreePath)
		assert.NoError(t, err, "Worktree directory should exist")

		// Verify the new branch was created
		branches, err := git.GetLocalBranches(bareRepoPath)
		require.NoError(t, err)
		assert.Contains(t, branches, "feature/new-branch", "New branch should exist")

		// Verify the branch was created from origin/master
		currentBranch, err := git.GetCurrentBranch(worktreePath)
		require.NoError(t, err)
		assert.Equal(t, "feature/new-branch", currentBranch, "Should be on the new branch")

		// Verify commit hash matches what we fetched from origin/master
		newBranchCommit, err := getCommitHash(bareRepoPath, "feature/new-branch")
		require.NoError(t, err)
		assert.Equal(t, beforeFetchCommit, newBranchCommit, "New branch should be cut from origin/master")
	})

	t.Run("does not fetch when pull flag is false", func(t *testing.T) {
		// Create a temporary directory for test repo
		tempDir, err := os.MkdirTemp("", "treekanga-add-no-pull-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Clone a real repo
		bareRepoPath := filepath.Join(tempDir, "test.git")
		err = git.CloneBare("https://github.com/octocat/Hello-World.git", bareRepoPath)
		require.NoError(t, err)

		err = git.ConfigureBare(bareRepoPath)
		require.NoError(t, err)

		// Create a local base branch to work from
		err = git.AddWorktree(bareRepoPath, tempDir, "base-worktree", []string{"-b", "local-base", "origin/master", "--no-track"})
		require.NoError(t, err)

		// Simulate the config WITHOUT pull flag
		cfg := config.AppConfig{
			BareRepoPath:               bareRepoPath,
			WorktreeTargetDir:          tempDir,
			NewBranchName:              "feature/no-pull",
			NewWorktreeName:            "feature-no-pull",
			BaseBranch:                 "local-base",
			BaseBranchExistsLocally:    true,
			BaseBranchExistsRemotely:   false,
			NewBranchExistsLocally:     false,
			NewBranchExistsRemotely:    false,
			PullBeforeCuttingNewBranch: false,
		}

		// Call AddWorktree
		AddWorktree(nil, nil, cfg)

		// Verify the worktree was created
		worktreePath := filepath.Join(tempDir, "feature-no-pull")
		_, err = os.Stat(worktreePath)
		assert.NoError(t, err, "Worktree directory should exist")

		// Verify the new branch exists
		branches, err := git.GetLocalBranches(bareRepoPath)
		require.NoError(t, err)
		assert.Contains(t, branches, "feature/no-pull", "New branch should exist")
	})
}

func TestGetAddWorktreeArguments(t *testing.T) {
	t.Run("default mode creates from local branch when pull flag is false", func(t *testing.T) {
		params := AddWorktreeConfig{
			BareRepoPath:               "/test/path",
			WorktreeTargetDirectory:    "/test/worktrees",
			CheckoutRemote:             false,
			CheckoutLocal:              false,
			NewBranchName:              "feature/new",
			BaseBranch:                 "main",
			BaseBranchExistsLocally:    true,
			NewBranchExistsLocally:     false,
			NewBranchExistsRemotely:    false,
			PullBeforeCuttingNewBranch: false,
		}

		args := GetAddWorktreeArguements(params)
		expected := []string{"-b", "feature/new", "--no-track", "main"}
		assert.Equal(t, expected, args, "Should create from local base branch")
	})

	t.Run("default mode creates from remote branch when pull flag is true", func(t *testing.T) {
		params := AddWorktreeConfig{
			BareRepoPath:               "/test/path",
			WorktreeTargetDirectory:    "/test/worktrees",
			CheckoutRemote:             false,
			CheckoutLocal:              false,
			NewBranchName:              "feature/new",
			BaseBranch:                 "main",
			BaseBranchExistsLocally:    true,
			NewBranchExistsRemotely:    false,
			NewBranchExistsLocally:     false,
			PullBeforeCuttingNewBranch: true,
		}

		args := GetAddWorktreeArguements(params)
		expected := []string{"-b", "feature/new", "--no-track", "origin/main"}
		assert.Equal(t, expected, args, "Should create from remote base branch when pull flag is true")
	})

	t.Run("remote mode checks out existing remote branch", func(t *testing.T) {
		params := AddWorktreeConfig{
			BareRepoPath:            "/test/path",
			WorktreeTargetDirectory: "/test/worktrees",
			CheckoutRemote:          true,
			CheckoutLocal:           false,
			NewBranchName:           "existing-remote-branch",
			BaseBranch:              "main",
			NewBranchExistsLocally:  false,
			NewBranchExistsRemotely: true,
			BaseBranchExistsLocally: true,
		}

		args := GetAddWorktreeArguements(params)
		expected := []string{"existing-remote-branch"}
		assert.Equal(t, expected, args, "Should checkout existing remote branch")
	})

	t.Run("local mode checks out existing local branch", func(t *testing.T) {
		params := AddWorktreeConfig{
			BareRepoPath:            "/test/path",
			WorktreeTargetDirectory: "/test/worktrees",
			CheckoutRemote:          false,
			CheckoutLocal:           true,
			NewBranchName:           "existing-local-branch",
			BaseBranch:              "main",
			NewBranchExistsLocally:  true,
			NewBranchExistsRemotely: false,
			BaseBranchExistsLocally: true,
		}

		args := GetAddWorktreeArguements(params)
		expected := []string{"existing-local-branch"}
		assert.Equal(t, expected, args, "Should checkout existing local branch")
	})

	t.Run("default mode creates from remote when base only exists remotely", func(t *testing.T) {
		params := AddWorktreeConfig{
			BareRepoPath:               "/test/path",
			WorktreeTargetDirectory:    "/test/worktrees",
			CheckoutRemote:             false,
			CheckoutLocal:              false,
			NewBranchName:              "feature/new",
			BaseBranch:                 "main",
			BaseBranchExistsLocally:    false,
			NewBranchExistsLocally:     false,
			NewBranchExistsRemotely:    false,
			PullBeforeCuttingNewBranch: false,
		}

		args := GetAddWorktreeArguements(params)
		expected := []string{"-b", "feature/new", "--no-track", "origin/main"}
		assert.Equal(t, expected, args, "Should create from remote when base doesn't exist locally")
	})
}

// Helper function to get commit hash for a branch
func getCommitHash(bareRepoPath, ref string) (string, error) {
	output, err := runCommandOutput("git", "-C", bareRepoPath, "rev-parse", ref)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

// Helper function for running git commands in tests
func runCommandOutput(cmd string, args ...string) (string, error) {
	// Use exec.Command directly since git package's runCommandOutput is unexported
	c := exec.Command(cmd, args...)
	output, err := c.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
