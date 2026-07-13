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

func TestSetConfigForAddServiceFetchesRemoteBranchForCheckoutRemote(t *testing.T) {
	// Skip if running in CI without git
	if os.Getenv("SKIP_INTEGRATION_TESTS") != "" {
		t.Skip("Skipping integration test")
	}

	// setupOriginAndBareClone creates a local "origin" repo with an initial
	// commit on main, then bare-clones it the way treekanga does. Returns
	// paths to both so a test can push further commits/branches to origin
	// after the bare clone's initial fetch has already happened.
	setupOriginAndBareClone := func(t *testing.T, tempDir string) (originPath, bareRepoPath string) {
		originPath = filepath.Join(tempDir, "origin")
		_, err := runCommandOutput("git", "init", "-b", "main", originPath)
		require.NoError(t, err)

		readmePath := filepath.Join(originPath, "README.md")
		require.NoError(t, os.WriteFile(readmePath, []byte("hello"), 0644))
		_, err = runCommandOutput("git", "-C", originPath, "add", "README.md")
		require.NoError(t, err)
		_, err = runCommandOutput("git", "-C", originPath, "-c", "user.email=test@test.com", "-c", "user.name=test", "commit", "-m", "initial commit")
		require.NoError(t, err)

		bareRepoPath = filepath.Join(tempDir, "test.git")
		err = git.CloneBare(originPath, bareRepoPath)
		require.NoError(t, err)
		err = git.ConfigureBare(bareRepoPath)
		require.NoError(t, err)

		return originPath, bareRepoPath
	}

	t.Run("finds branch pushed to remote after last fetch when --remote flag is used", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "treekanga-add-fetch-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		originPath, bareRepoPath := setupOriginAndBareClone(t, tempDir)

		// Push a new branch to origin *after* the bare repo's initial fetch.
		_, err = runCommandOutput("git", "-C", originPath, "checkout", "-b", "feature/late-push")
		require.NoError(t, err)
		readmePath := filepath.Join(originPath, "README.md")
		require.NoError(t, os.WriteFile(readmePath, []byte("hello again"), 0644))
		_, err = runCommandOutput("git", "-C", originPath, "add", "README.md")
		require.NoError(t, err)
		_, err = runCommandOutput("git", "-C", originPath, "-c", "user.email=test@test.com", "-c", "user.name=test", "commit", "-m", "late push")
		require.NoError(t, err)

		// Sanity check: the bare repo's cached remote-tracking refs don't know about it yet.
		cachedBranches, err := git.GetRemoteBranches(bareRepoPath)
		require.NoError(t, err)
		assert.NotContains(t, cachedBranches, "feature/late-push", "branch should not be visible without a fetch")

		cfg := config.AppConfig{
			BareRepoPath:   bareRepoPath,
			CheckoutRemote: true,
		}

		cfg = SetConfigForAddService(cfg, []string{"feature/late-push"})

		assert.True(t, cfg.NewBranchExistsRemotely, "branch pushed after last fetch should be found once --remote triggers a fetch")
	})

	t.Run("does not fatal and leaves branch missing when it truly doesn't exist remotely", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "treekanga-add-fetch-missing-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		_, bareRepoPath := setupOriginAndBareClone(t, tempDir)

		cfg := config.AppConfig{
			BareRepoPath:   bareRepoPath,
			CheckoutRemote: true,
		}

		cfg = SetConfigForAddService(cfg, []string{"does-not-exist-anywhere"})

		assert.False(t, cfg.NewBranchExistsRemotely, "nonexistent branch should not be found, and the failed targeted fetch should not crash config setup")
	})

	t.Run("does not fetch when --remote flag is not used", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "treekanga-add-no-fetch-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		originPath, bareRepoPath := setupOriginAndBareClone(t, tempDir)

		// Push a new branch to origin after the bare repo's initial fetch.
		_, err = runCommandOutput("git", "-C", originPath, "checkout", "-b", "feature/late-push")
		require.NoError(t, err)
		readmePath := filepath.Join(originPath, "README.md")
		require.NoError(t, os.WriteFile(readmePath, []byte("hello again"), 0644))
		_, err = runCommandOutput("git", "-C", originPath, "add", "README.md")
		require.NoError(t, err)
		_, err = runCommandOutput("git", "-C", originPath, "-c", "user.email=test@test.com", "-c", "user.name=test", "commit", "-m", "late push")
		require.NoError(t, err)

		cfg := config.AppConfig{
			BareRepoPath:   bareRepoPath,
			CheckoutRemote: false,
			BaseBranch:     "main",
		}

		cfg = SetConfigForAddService(cfg, []string{"feature/late-push"})

		assert.False(t, cfg.NewBranchExistsRemotely, "branch pushed after last fetch should stay hidden when --remote isn't used, since no fetch should be triggered")
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
			BaseBranchExistsRemotely:   true,
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
			BaseBranchExistsRemotely:   true,
			NewBranchExistsLocally:     false,
			NewBranchExistsRemotely:    false,
			PullBeforeCuttingNewBranch: true,
		}

		args := GetAddWorktreeArguements(params)
		expected := []string{"-b", "feature/new", "--no-track", "origin/main"}
		assert.Equal(t, expected, args, "Should create from remote base branch when pull flag is true")
	})

	t.Run("remote mode checks out existing remote branch", func(t *testing.T) {
		params := AddWorktreeConfig{
			BareRepoPath:             "/test/path",
			WorktreeTargetDirectory:  "/test/worktrees",
			CheckoutRemote:           true,
			CheckoutLocal:            false,
			NewBranchName:            "existing-remote-branch",
			BaseBranch:               "main",
			NewBranchExistsLocally:   false,
			NewBranchExistsRemotely:  true,
			BaseBranchExistsLocally:  true,
			BaseBranchExistsRemotely: true,
		}

		args := GetAddWorktreeArguements(params)
		expected := []string{"existing-remote-branch"}
		assert.Equal(t, expected, args, "Should checkout existing remote branch")
	})

	t.Run("local mode checks out existing local branch", func(t *testing.T) {
		params := AddWorktreeConfig{
			BareRepoPath:             "/test/path",
			WorktreeTargetDirectory:  "/test/worktrees",
			CheckoutRemote:           false,
			CheckoutLocal:            true,
			NewBranchName:            "existing-local-branch",
			BaseBranch:               "main",
			NewBranchExistsLocally:   true,
			NewBranchExistsRemotely:  false,
			BaseBranchExistsLocally:  true,
			BaseBranchExistsRemotely: true,
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
			BaseBranchExistsRemotely:   true,
			NewBranchExistsLocally:     false,
			NewBranchExistsRemotely:    false,
			PullBeforeCuttingNewBranch: false,
		}

		args := GetAddWorktreeArguements(params)
		expected := []string{"-b", "feature/new", "--no-track", "origin/main"}
		assert.Equal(t, expected, args, "Should create from remote when base doesn't exist locally")
	})

	t.Run("default mode creates from local when pull flag is true but base only exists locally", func(t *testing.T) {
		params := AddWorktreeConfig{
			BareRepoPath:               "/test/path",
			WorktreeTargetDirectory:    "/test/worktrees",
			CheckoutRemote:             false,
			CheckoutLocal:              false,
			NewBranchName:              "feature/new",
			BaseBranch:                 "local-only-branch",
			BaseBranchExistsLocally:    true,
			BaseBranchExistsRemotely:   false,
			NewBranchExistsLocally:     false,
			NewBranchExistsRemotely:    false,
			PullBeforeCuttingNewBranch: true,
		}

		args := GetAddWorktreeArguements(params)
		expected := []string{"-b", "feature/new", "--no-track", "local-only-branch"}
		assert.Equal(t, expected, args, "Should create from local branch when pull flag is true but base doesn't exist remotely")
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
