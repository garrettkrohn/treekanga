package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/garrettkrohn/treekanga/config"
	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/services"
	spinnerhuh "github.com/garrettkrohn/treekanga/spinnerHuh"
	"github.com/garrettkrohn/treekanga/transformer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCloneAndAddIntegration is an integration test that verifies
// the clone command successfully creates a bare repository and
// the add command successfully creates a worktree from it
func TestCloneAndAddIntegration(t *testing.T) {
	// Skip if running in CI without git
	if os.Getenv("SKIP_INTEGRATION_TESTS") != "" {
		t.Skip("Skipping integration test")
	}

	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "treekanga-integration-test-*")
	require.NoError(t, err, "Failed to create temp directory")
	defer os.RemoveAll(tempDir)

	// Save current directory and change to temp directory
	originalDir, err := os.Getwd()
	require.NoError(t, err, "Failed to get working directory")
	defer os.Chdir(originalDir)

	err = os.Chdir(tempDir)
	require.NoError(t, err, "Failed to change to temp directory")

	// Set up real dependencies
	mockSpinner := &mockSpinner{}

	// Use a small public repository for testing
	testRepoURL := "https://github.com/octocat/Hello-World.git"
	expectedFolderName := "Hello-World.git_bare"

	t.Log("Step 1: Cloning bare repository...")

	// Run the clone command
	args := []string{testRepoURL}
	CloneBareRepo(mockSpinner, args)

	// Verify the bare repository was created
	bareRepoPath := filepath.Join(tempDir, expectedFolderName)

	// Check that the directory exists
	_, err = os.Stat(bareRepoPath)
	assert.NoError(t, err, "Bare repository directory should exist")

	// Verify it's a bare repository by checking for key files/directories
	assertPathExists(t, filepath.Join(bareRepoPath, "HEAD"))
	assertPathExists(t, filepath.Join(bareRepoPath, "refs"))
	assertPathExists(t, filepath.Join(bareRepoPath, "objects"))
	assertPathExists(t, filepath.Join(bareRepoPath, "config"))

	// Verify there's no working directory (no .git subdirectory)
	gitSubdir := filepath.Join(bareRepoPath, ".git")
	_, err = os.Stat(gitSubdir)
	assert.True(t, os.IsNotExist(err), "Bare repo should not have .git subdirectory")

	t.Logf("✓ Successfully verified bare repository at: %s", bareRepoPath)

	t.Log("Step 2: Creating worktree from bare repository...")

	// First, fetch remote branches to determine what's available
	remoteBranches, err := git.GetRemoteBranches(bareRepoPath)
	require.NoError(t, err, "Should be able to get remote branches")
	require.Greater(t, len(remoteBranches), 0, "Should have at least one remote branch")

	// Log all remote branches for debugging
	t.Logf("All remote branches: %v", remoteBranches)

	// Find a valid remote branch
	var baseBranch string
	for _, branch := range remoteBranches {
		if branch != "" && !strings.Contains(branch, "HEAD") {
			baseBranch = branch
			break
		}
	}
	require.NotEmpty(t, baseBranch, "Should find at least one valid remote branch")
	t.Logf("Using base branch: %s", baseBranch)

	// Create a worktree using the add functionality
	testBranchName := "test_branch"
	worktreePath := filepath.Join(tempDir, testBranchName)

	// Create AddWorktreeConfig for the worktree
	worktreeConfig := services.AddWorktreeConfig{
		BareRepoPath:               bareRepoPath,
		WorktreeTargetDirectory:    tempDir,
		NewBranchExistsLocally:     false,
		NewBranchExistsRemotely:    false,
		BaseBranchExistsLocally:    false,
		NewBranchName:              testBranchName,
		PullBeforeCuttingNewBranch: false,
		BaseBranch:                 baseBranch,
		NewWorktreeName:            testBranchName,
	}

	// Get the branch arguments from the service function
	branchArgs := services.GetAddWorktreeArguements(worktreeConfig)

	// Add the worktree
	err = git.AddWorktree(bareRepoPath, tempDir, testBranchName, branchArgs)
	require.NoError(t, err, "Should successfully add worktree")

	t.Logf("✓ Successfully created worktree at: %s", worktreePath)

	t.Log("Step 3: Verifying worktree was created correctly...")

	// Verify the worktree directory exists
	_, err = os.Stat(worktreePath)
	assert.NoError(t, err, "Worktree directory should exist")

	// Verify the worktree has a .git file (not a directory)
	gitFile := filepath.Join(worktreePath, ".git")
	gitFileInfo, err := os.Stat(gitFile)
	assert.NoError(t, err, ".git file should exist in worktree")
	assert.False(t, gitFileInfo.IsDir(), ".git should be a file, not a directory in a worktree")

	// Verify the .git file points to the correct location
	gitFileContent, err := os.ReadFile(gitFile)
	assert.NoError(t, err, "Should be able to read .git file")
	assert.Contains(t, string(gitFileContent), "gitdir:", ".git file should contain gitdir reference")
	assert.Contains(t, string(gitFileContent), expectedFolderName, ".git file should reference the bare repo")

	// Verify we can see the worktree in git worktree list
	rawWorktrees, err := git.ListWorktrees(bareRepoPath)
	assert.NoError(t, err, "Should be able to list worktrees")
	assert.Greater(t, len(rawWorktrees), 0, "Should have at least one worktree")

	worktrees := transformer.TransformWorktrees(rawWorktrees)

	// Resolve symlinks in our worktree path (on macOS, /var is a symlink to /private/var)
	resolvedWorktreePath, err := filepath.EvalSymlinks(worktreePath)
	if err != nil {
		resolvedWorktreePath = worktreePath
	}

	// Check that our worktree is in the list
	foundWorktree := false
	for _, wt := range worktrees {
		resolvedWtPath, _ := filepath.EvalSymlinks(wt.FullPath)

		if resolvedWtPath == resolvedWorktreePath || wt.FullPath == worktreePath {
			foundWorktree = true
			t.Logf("Found our worktree in git worktree list")
			break
		}
	}
	assert.True(t, foundWorktree, "Should find our worktree in git worktree list")

	// Verify the worktree has actual git repository contents
	assertPathExists(t, filepath.Join(worktreePath, "README"))

	t.Logf("✓ Successfully verified worktree at: %s", worktreePath)

	t.Log("Step 4: Verifying list command shows the new worktree...")

	// Use the buildWorktreeStrings function (or directly get worktrees and transform)
	// We need to set up deps for the list command to work
	// Save the original deps values
	originalBareRepoPath := deps.AppConfig.BareRepoPath
	defer func() {
		deps.AppConfig.BareRepoPath = originalBareRepoPath
	}()

	// Set deps for the list command
	deps.AppConfig.BareRepoPath = bareRepoPath

	// Get the list of worktrees
	worktreeList, err := buildWorktreeStrings(false, false, false, false)
	assert.NoError(t, err, "Should be able to build worktree strings")
	assert.Greater(t, len(worktreeList), 0, "Should have at least one worktree in the list")

	// Check that our test_branch appears in the list
	foundInList := false
	for _, wt := range worktreeList {
		t.Logf("Found in list: %s", wt)
		if wt == testBranchName {
			foundInList = true
			break
		}
	}
	assert.True(t, foundInList, "Should find test_branch in the worktree list output")

	// Also test verbose mode
	verboseList, err := buildWorktreeStrings(true, false, false, false)
	assert.NoError(t, err, "Should be able to build verbose worktree strings")
	assert.Greater(t, len(verboseList), 0, "Should have at least one worktree in verbose list")

	// Check that our test_branch appears in the verbose list
	foundInVerboseList := false
	for _, wt := range verboseList {
		t.Logf("Found in verbose list: %s", wt)
		// Verbose output contains the branch name within the string
		if len(wt) > 0 {
			// Check if the line contains our test branch name
			if len(wt) >= len(testBranchName) {
				for i := 0; i <= len(wt)-len(testBranchName); i++ {
					if wt[i:i+len(testBranchName)] == testBranchName {
						foundInVerboseList = true
						break
					}
				}
			}
		}
		if foundInVerboseList {
			break
		}
	}
	assert.True(t, foundInVerboseList, "Should find test_branch in the verbose worktree list output")

	t.Logf("✓ Successfully verified list command shows the new worktree")

	t.Log("Step 5: Deleting the worktree and verifying it's gone...")

	// Remove the worktree directly
	err = git.RemoveWorktree(bareRepoPath, worktreePath, true)
	assert.NoError(t, err, "Should successfully remove worktree")

	t.Logf("✓ Successfully deleted worktree")

	// Verify the worktree directory no longer exists
	_, err = os.Stat(worktreePath)
	assert.True(t, os.IsNotExist(err), "Worktree directory should no longer exist")

	// Verify the worktree is no longer in git worktree list
	rawWorktreesAfterDelete, err := git.ListWorktrees(bareRepoPath)
	assert.NoError(t, err, "Should be able to list worktrees after deletion")

	worktreesAfterDelete := transformer.TransformWorktrees(rawWorktreesAfterDelete)

	foundAfterDelete := false
	for _, wt := range worktreesAfterDelete {
		resolvedWtPath, _ := filepath.EvalSymlinks(wt.FullPath)
		resolvedWorktreePath, _ := filepath.EvalSymlinks(worktreePath)

		if resolvedWtPath == resolvedWorktreePath || wt.FullPath == worktreePath {
			foundAfterDelete = true
			break
		}
	}
	assert.False(t, foundAfterDelete, "Worktree should not appear in git worktree list after deletion")

	// Verify the list command no longer shows it
	worktreeListAfterDelete, err := buildWorktreeStrings(false, false, false, false)
	assert.NoError(t, err, "Should be able to build worktree strings after deletion")

	foundInListAfterDelete := false
	for _, wt := range worktreeListAfterDelete {
		if wt == testBranchName {
			foundInListAfterDelete = true
			break
		}
	}
	assert.False(t, foundInListAfterDelete, "test_branch should not appear in list command after deletion")

	t.Logf("✓ Successfully verified worktree is deleted and no longer appears anywhere")
	t.Log("✅ Integration test completed successfully!")
}

// TestRenameWorktreeIntegration is an integration test that verifies
// the rename command successfully renames both the branch and worktree folder
func TestRenameWorktreeIntegration(t *testing.T) {
	// Skip if running in CI without git
	if os.Getenv("SKIP_INTEGRATION_TESTS") != "" {
		t.Skip("Skipping integration test")
	}

	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "treekanga-rename-integration-*")
	require.NoError(t, err, "Failed to create temp directory")
	defer os.RemoveAll(tempDir)

	// Save current directory and change to temp directory
	originalDir, err := os.Getwd()
	require.NoError(t, err, "Failed to get working directory")
	defer os.Chdir(originalDir)

	err = os.Chdir(tempDir)
	require.NoError(t, err, "Failed to change to temp directory")

	// Use a small public repository for testing
	testRepoURL := "https://github.com/octocat/Hello-World.git"
	expectedFolderName := "Hello-World.git_bare"

	t.Log("Step 1: Cloning bare repository...")

	// Clone the bare repo
	mockSpinner := &mockSpinner{}
	args := []string{testRepoURL}
	CloneBareRepo(mockSpinner, args)

	bareRepoPath := filepath.Join(tempDir, expectedFolderName)
	_, err = os.Stat(bareRepoPath)
	require.NoError(t, err, "Bare repository should exist")

	t.Log("Step 2: Creating initial worktree...")

	// Set up deps for commands
	deps.AppConfig.BareRepoPath = bareRepoPath

	// Get remote branches
	remoteBranches, err := git.GetRemoteBranches(bareRepoPath)
	require.NoError(t, err)
	require.Greater(t, len(remoteBranches), 0)

	var baseBranch string
	for _, branch := range remoteBranches {
		if branch != "" && !strings.Contains(branch, "HEAD") {
			baseBranch = branch
			break
		}
	}
	require.NotEmpty(t, baseBranch)

	// Create initial worktree
	initialBranch := "feature/old-name"
	initialFolder := "feature-old-name"
	worktreePath := filepath.Join(tempDir, initialFolder)

	worktreeConfig := services.AddWorktreeConfig{
		BareRepoPath:               bareRepoPath,
		WorktreeTargetDirectory:    tempDir,
		NewBranchExistsLocally:     false,
		NewBranchExistsRemotely:    false,
		BaseBranchExistsLocally:    false,
		NewBranchName:              initialBranch,
		PullBeforeCuttingNewBranch: false,
		BaseBranch:                 baseBranch,
		NewWorktreeName:            initialFolder,
	}

	branchArgs := services.GetAddWorktreeArguements(worktreeConfig)
	err = git.AddWorktree(bareRepoPath, tempDir, initialFolder, branchArgs)
	require.NoError(t, err)

	t.Logf("✓ Created worktree at: %s", worktreePath)

	// Verify initial state
	branches, err := git.GetLocalBranches(bareRepoPath)
	require.NoError(t, err)
	assert.Contains(t, branches, initialBranch, "Initial branch should exist")

	_, err = os.Stat(worktreePath)
	assert.NoError(t, err, "Initial worktree folder should exist")

	t.Log("Step 3: Renaming the worktree...")

	// Change into the worktree directory
	err = os.Chdir(worktreePath)
	require.NoError(t, err, "Should be able to cd into worktree")

	// Set up config for rename
	newBranchName := "feature/new-name"
	newFolderName := "feature-new-name"
	newWorktreePath := filepath.Join(tempDir, newFolderName)

	cfg := config.AppConfig{
		BareRepoPath:     bareRepoPath,
		WorktreeTargetDir: tempDir,
	}

	// Execute rename (pass nil for connector and confirmer since we don't need tmux handling in tests)
	err = services.RenameWorktree(cfg, newBranchName, worktreePath, nil, nil)
	assert.NoError(t, err, "Should successfully rename worktree")

	t.Log("Step 4: Verifying the rename...")

	// Verify branch was renamed
	branches, err = git.GetLocalBranches(bareRepoPath)
	require.NoError(t, err)
	assert.Contains(t, branches, newBranchName, "New branch should exist")
	assert.NotContains(t, branches, initialBranch, "Old branch should not exist")

	// Verify folder was moved
	_, err = os.Stat(newWorktreePath)
	assert.NoError(t, err, "New worktree folder should exist")

	_, err = os.Stat(worktreePath)
	assert.True(t, os.IsNotExist(err), "Old worktree folder should not exist")

	// Verify worktree list shows the new branch
	rawWorktrees, err := git.ListWorktrees(bareRepoPath)
	require.NoError(t, err)

	worktrees := transformer.TransformWorktrees(rawWorktrees)
	foundNewBranch := false
	for _, wt := range worktrees {
		if wt.BranchName == newBranchName {
			foundNewBranch = true
			t.Logf("Found renamed branch in worktree list: %s", wt.BranchName)
			break
		}
	}
	assert.True(t, foundNewBranch, "New branch should appear in git worktree list")

	t.Logf("✓ Successfully verified rename: %s → %s", initialBranch, newBranchName)
	t.Logf("✓ Successfully verified folder move: %s → %s", initialFolder, newFolderName)

	// Change back to temp dir before cleanup
	err = os.Chdir(tempDir)
	require.NoError(t, err)

	t.Log("✅ Rename integration test completed successfully!")
}

// assertPathExists is a helper function to check if a path exists
func assertPathExists(t *testing.T, path string) {
	t.Helper()
	_, err := os.Stat(path)
	assert.NoError(t, err, "Path should exist: %s", path)
}

// mockSpinner is a simple mock that does nothing - we don't need UI in tests
type mockSpinner struct{}

func (m *mockSpinner) Title(string) spinnerhuh.HuhSpinner { return m }
func (m *mockSpinner) Action(f func()) spinnerhuh.HuhSpinner {
	f() // Just execute the function immediately
	return m
}
func (m *mockSpinner) Run() error { return nil }
