package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func run(t *testing.T, args ...string) {
	t.Helper()
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "command %v failed: %s", args, output)
}

func setupComputeStatusRepo(t *testing.T) (bareRepoPath, worktreePath string) {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "treekanga-compute-status-test-*")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(tempDir) })

	bareRepoPath = filepath.Join(tempDir, "test.git")
	run(t, "git", "init", "--bare", bareRepoPath)
	require.NoError(t, git.ConfigureBare(bareRepoPath))

	worktreePath = filepath.Join(tempDir, "main")
	require.NoError(t, git.AddWorktree(bareRepoPath, tempDir, "main", []string{"-b", "main"}))
	run(t, "git", "-C", bareRepoPath, "symbolic-ref", "HEAD", "refs/heads/main")

	run(t, "git", "-C", worktreePath, "config", "user.email", "test@example.com")
	run(t, "git", "-C", worktreePath, "config", "user.name", "Test User")
	run(t, "sh", "-c", fmt.Sprintf("cd %s && echo 'initial' > file.txt && git add file.txt && git commit -m 'initial commit'", worktreePath))
	run(t, "git", "-C", worktreePath, "update-ref", "refs/remotes/origin/main", "refs/heads/main")

	return bareRepoPath, worktreePath
}

func TestComputeWorktreeStatus(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") != "" {
		t.Skip("Skipping integration test")
	}

	bareRepoPath, worktreePath := setupComputeStatusRepo(t)

	featurePath := filepath.Join(filepath.Dir(worktreePath), "feature")
	require.NoError(t, git.AddWorktree(bareRepoPath, filepath.Dir(worktreePath), "feature", []string{"-b", "feature"}))
	run(t, "git", "-C", featurePath, "config", "user.email", "test@example.com")
	run(t, "git", "-C", featurePath, "config", "user.name", "Test User")
	run(t, "sh", "-c", fmt.Sprintf("cd %s && echo 'feature' > feature.txt && git add feature.txt && git commit -m 'feature commit'", featurePath))
	run(t, "sh", "-c", fmt.Sprintf("cd %s && echo 'dirty' > dirty.txt && git add dirty.txt", featurePath))
	run(t, "sh", "-c", fmt.Sprintf("cd %s && echo 'untracked' > untracked.txt", featurePath))

	worktree := models.Worktree{
		FullPath:   featurePath,
		Folder:     "feature",
		BranchName: "feature",
	}

	result := ComputeWorktreeStatus(worktree, "main")

	assert.True(t, result.StatusLoaded)
	assert.True(t, result.HasStaged)
	assert.False(t, result.HasModified)
	assert.True(t, result.HasUntracked)
	assert.Equal(t, 1, result.AheadDefault)
	assert.Equal(t, 0, result.BehindDefault)
	assert.False(t, result.HasUpstream, "feature worktree has no upstream configured")
	assert.Equal(t, models.MergeStatusNotMerged, result.Merged)
}

func TestComputeWorktreeStatusMerged(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") != "" {
		t.Skip("Skipping integration test")
	}

	bareRepoPath, worktreePath := setupComputeStatusRepo(t)

	featurePath := filepath.Join(filepath.Dir(worktreePath), "feature")
	require.NoError(t, git.AddWorktree(bareRepoPath, filepath.Dir(worktreePath), "feature", []string{"-b", "feature"}))
	run(t, "git", "-C", featurePath, "config", "user.email", "test@example.com")
	run(t, "git", "-C", featurePath, "config", "user.name", "Test User")
	run(t, "sh", "-c", fmt.Sprintf("cd %s && echo 'feature' > feature.txt && git add feature.txt && git commit -m 'feature commit'", featurePath))

	run(t, "git", "-C", worktreePath, "merge", "feature", "--no-edit")
	run(t, "git", "-C", worktreePath, "update-ref", "refs/remotes/origin/main", "refs/heads/main")

	worktree := models.Worktree{
		FullPath:   featurePath,
		Folder:     "feature",
		BranchName: "feature",
	}

	result := ComputeWorktreeStatus(worktree, "main")

	assert.True(t, result.StatusLoaded)
	assert.False(t, result.HasStaged)
	assert.False(t, result.HasModified)
	assert.False(t, result.HasUntracked)
	assert.Equal(t, models.MergeStatusMerged, result.Merged)
}
