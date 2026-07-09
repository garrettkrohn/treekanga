package git

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupStatusTestRepo(t *testing.T) (bareRepoPath, worktreePath string) {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "treekanga-status-test-*")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(tempDir) })

	bareRepoPath = filepath.Join(tempDir, "test.git")
	require.NoError(t, runCommand("git", "init", "--bare", bareRepoPath))
	require.NoError(t, ConfigureBare(bareRepoPath))

	worktreePath = filepath.Join(tempDir, "main")
	require.NoError(t, AddWorktree(bareRepoPath, tempDir, "main", []string{"-b", "main"}))
	require.NoError(t, runCommand("git", "-C", bareRepoPath, "symbolic-ref", "HEAD", "refs/heads/main"))

	require.NoError(t, runCommand("git", "-C", worktreePath, "config", "user.email", "test@example.com"))
	require.NoError(t, runCommand("git", "-C", worktreePath, "config", "user.name", "Test User"))
	require.NoError(t, runCommand("sh", "-c", fmt.Sprintf("cd %s && echo 'initial' > file.txt && git add file.txt && git commit -m 'initial commit'", worktreePath)))

	return bareRepoPath, worktreePath
}

func TestGetWorkingTreeStatus(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") != "" {
		t.Skip("Skipping integration test")
	}

	_, worktreePath := setupStatusTestRepo(t)

	staged, modified, untracked, err := GetWorkingTreeStatus(worktreePath)
	require.NoError(t, err)
	assert.False(t, staged)
	assert.False(t, modified)
	assert.False(t, untracked)

	// Untracked file
	require.NoError(t, runCommand("sh", "-c", fmt.Sprintf("cd %s && echo 'new' > untracked.txt", worktreePath)))
	staged, modified, untracked, err = GetWorkingTreeStatus(worktreePath)
	require.NoError(t, err)
	assert.False(t, staged)
	assert.False(t, modified)
	assert.True(t, untracked)

	// Staged file
	require.NoError(t, runCommand("git", "-C", worktreePath, "add", "untracked.txt"))
	staged, modified, untracked, err = GetWorkingTreeStatus(worktreePath)
	require.NoError(t, err)
	assert.True(t, staged)
	assert.False(t, modified)
	assert.False(t, untracked)

	// Modified tracked file (unstaged)
	require.NoError(t, runCommand("git", "-C", worktreePath, "commit", "-m", "add untracked file"))
	require.NoError(t, runCommand("sh", "-c", fmt.Sprintf("cd %s && echo 'changed' >> file.txt", worktreePath)))
	staged, modified, untracked, err = GetWorkingTreeStatus(worktreePath)
	require.NoError(t, err)
	assert.False(t, staged)
	assert.True(t, modified)
	assert.False(t, untracked)
}

func TestGetAheadBehind(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") != "" {
		t.Skip("Skipping integration test")
	}

	bareRepoPath, worktreePath := setupStatusTestRepo(t)

	ahead, behind, err := GetAheadBehind(worktreePath, "main")
	require.NoError(t, err)
	assert.Equal(t, 0, ahead)
	assert.Equal(t, 0, behind)

	// Create a feature branch ahead of main by one commit
	featurePath := filepath.Join(filepath.Dir(worktreePath), "feature")
	require.NoError(t, AddWorktree(bareRepoPath, filepath.Dir(worktreePath), "feature", []string{"-b", "feature"}))
	require.NoError(t, runCommand("git", "-C", featurePath, "config", "user.email", "test@example.com"))
	require.NoError(t, runCommand("git", "-C", featurePath, "config", "user.name", "Test User"))
	require.NoError(t, runCommand("sh", "-c", fmt.Sprintf("cd %s && echo 'feature' > feature.txt && git add feature.txt && git commit -m 'feature commit'", featurePath)))

	ahead, behind, err = GetAheadBehind(featurePath, "main")
	require.NoError(t, err)
	assert.Equal(t, 1, ahead)
	assert.Equal(t, 0, behind)

	// main is now behind feature by one commit
	ahead, behind, err = GetAheadBehind(worktreePath, "feature")
	require.NoError(t, err)
	assert.Equal(t, 0, ahead)
	assert.Equal(t, 1, behind)
}

func TestGetUpstreamBranch(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") != "" {
		t.Skip("Skipping integration test")
	}

	_, worktreePath := setupStatusTestRepo(t)

	upstream, err := GetUpstreamBranch(worktreePath)
	require.NoError(t, err)
	assert.Empty(t, upstream, "no upstream should be configured")

	// Fake a remote-tracking ref locally (no real network remote in this
	// test) and configure the branch to track it, mirroring what a real
	// fetch from origin would leave behind.
	require.NoError(t, runCommand("git", "-C", worktreePath, "update-ref", "refs/remotes/origin/main", "refs/heads/main"))
	require.NoError(t, SetUpstream(worktreePath, "main"))
	upstream, err = GetUpstreamBranch(worktreePath)
	require.NoError(t, err)
	assert.Equal(t, "origin/main", upstream)
}

func TestIsMergedAncestor(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") != "" {
		t.Skip("Skipping integration test")
	}

	bareRepoPath, worktreePath := setupStatusTestRepo(t)

	featurePath := filepath.Join(filepath.Dir(worktreePath), "feature")
	require.NoError(t, AddWorktree(bareRepoPath, filepath.Dir(worktreePath), "feature", []string{"-b", "feature"}))
	require.NoError(t, runCommand("git", "-C", featurePath, "config", "user.email", "test@example.com"))
	require.NoError(t, runCommand("git", "-C", featurePath, "config", "user.name", "Test User"))
	require.NoError(t, runCommand("sh", "-c", fmt.Sprintf("cd %s && echo 'feature' > feature.txt && git add feature.txt && git commit -m 'feature commit'", featurePath)))

	// Not merged yet: main hasn't seen feature's commit
	merged, err := IsMerged(worktreePath, "feature", "main")
	require.NoError(t, err)
	assert.False(t, merged)

	// Merge feature into main - now feature should be an ancestor of main
	require.NoError(t, runCommand("git", "-C", worktreePath, "merge", "feature", "--no-edit"))
	merged, err = IsMerged(worktreePath, "feature", "main")
	require.NoError(t, err)
	assert.True(t, merged)
}

func TestIsMergedSquash(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") != "" {
		t.Skip("Skipping integration test")
	}

	bareRepoPath, worktreePath := setupStatusTestRepo(t)

	featurePath := filepath.Join(filepath.Dir(worktreePath), "feature")
	require.NoError(t, AddWorktree(bareRepoPath, filepath.Dir(worktreePath), "feature", []string{"-b", "feature"}))
	require.NoError(t, runCommand("git", "-C", featurePath, "config", "user.email", "test@example.com"))
	require.NoError(t, runCommand("git", "-C", featurePath, "config", "user.name", "Test User"))
	require.NoError(t, runCommand("sh", "-c", fmt.Sprintf("cd %s && echo 'feature' > feature.txt && git add feature.txt && git commit -m 'feature commit'", featurePath)))

	merged, err := IsMerged(worktreePath, "feature", "main")
	require.NoError(t, err)
	assert.False(t, merged)

	// Squash-merge feature into main: same content diff, different commit history
	require.NoError(t, runCommand("git", "-C", worktreePath, "merge", "--squash", "feature"))
	require.NoError(t, runCommand("git", "-C", worktreePath, "commit", "-m", "squash merge feature"))

	merged, err = IsMerged(worktreePath, "feature", "main")
	require.NoError(t, err)
	assert.True(t, merged, "squash-merged branch should be detected as merged via patch-id content match")
}
