package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	com "github.com/garrettkrohn/treekanga/common"
	"github.com/garrettkrohn/treekanga/execwrap"
	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/shell"
	spinnerhuh "github.com/garrettkrohn/treekanga/spinnerHuh"
	"github.com/garrettkrohn/treekanga/transformer"
	worktreeobj "github.com/garrettkrohn/treekanga/worktreeObj"
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
	realExec := execwrap.NewExec()
	realShell := shell.NewShell(realExec)
	realGit := git.NewGit(realShell)
	mockSpinner := &mockSpinner{}

	// Use a small public repository for testing
	testRepoURL := "https://github.com/octocat/Hello-World.git"
	expectedFolderName := "Hello-World.git_bare"

	t.Log("Step 1: Cloning bare repository...")

	// Run the clone command
	args := []string{testRepoURL}
	CloneBareRepo(realGit, mockSpinner, args)

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

	// Verify git config was set up correctly
	output, err := realShell.Cmd("git", "-C", bareRepoPath, "config", "--get", "remote.origin.fetch")
	assert.NoError(t, err, "Should be able to read git config")
	assert.Equal(t, "+refs/heads/*:refs/remotes/origin/*", output, "Git fetch config should be set correctly")

	t.Logf("✓ Successfully verified bare repository at: %s", bareRepoPath)

	t.Log("Step 2: Creating worktree from bare repository...")

	// First, fetch remote branches to determine what's available
	remoteBranches, err := realGit.GetRemoteBranches(&bareRepoPath)
	require.NoError(t, err, "Should be able to get remote branches")
	require.Greater(t, len(remoteBranches), 0, "Should have at least one remote branch")

	// Log all remote branches for debugging
	t.Logf("All remote branches: %v", remoteBranches)

	// Find a valid remote branch (must contain "origin/" prefix and have a branch name after it)
	var baseBranch string
	for _, branch := range remoteBranches {
		// Skip entries that are just "origin" or don't have the expected format
		if strings.HasPrefix(branch, "origin/") && len(branch) > 7 {
			// This is a valid remote branch like "origin/master" or "origin/main"
			baseBranch = branch[7:] // Remove "origin/" prefix
			break
		}
	}
	require.NotEmpty(t, baseBranch, "Should find at least one valid remote branch (origin/...)")
	t.Logf("Using base branch: %s", baseBranch)

	// Create a worktree using the add functionality
	testBranchName := "test_branch"
	worktreePath := filepath.Join(tempDir, testBranchName)

	// Create AddConfig for the worktree
	addConfig := &com.AddConfig{
		WorkingDir:        bareRepoPath,
		ParentDir:         tempDir,
		WorktreeTargetDir: worktreePath,
		GitInfo: com.GitInfo{
			NewBranchName:            testBranchName,
			BaseBranchName:           baseBranch,
			RepoName:                 "Hello-World",
			NewBranchExistsLocally:   false,
			NewBranchExistsRemotely:  false,
			BaseBranchExistsLocally:  false,
			BaseBranchExistsRemotely: true, // The base branch exists remotely
		},
		Flags: com.AddCmdFlags{
			Directory: &bareRepoPath,
		},
	}

	// Add the worktree
	err = realGit.AddWorktree(addConfig)
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
	worktrees, err := realGit.GetWorktrees(&bareRepoPath)
	assert.NoError(t, err, "Should be able to list worktrees")
	assert.Greater(t, len(worktrees), 0, "Should have at least one worktree")

	// Resolve symlinks in our worktree path (on macOS, /var is a symlink to /private/var)
	resolvedWorktreePath, err := filepath.EvalSymlinks(worktreePath)
	if err != nil {
		resolvedWorktreePath = worktreePath
	}

	// Check that our worktree is in the list
	// git worktree list returns lines like: "/path/to/worktree abcdef12345 [branch_name]"
	foundWorktree := false
	for _, wt := range worktrees {
		// Check if the worktree path is mentioned in this line
		if len(wt) > 0 {
			// Parse the worktree line to extract the path
			fields := splitWorktreeListLine(wt)
			if len(fields) > 0 {
				wtPath := fields[0]
				// Resolve symlinks in the worktree list path too
				resolvedWtPath, err := filepath.EvalSymlinks(wtPath)
				if err != nil {
					resolvedWtPath = wtPath
				}

				if resolvedWtPath == resolvedWorktreePath || wtPath == worktreePath {
					foundWorktree = true
					t.Logf("Found our worktree in git worktree list")
					break
				}
			}
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
	originalBareRepoPath := deps.BareRepoPath
	originalGit := deps.Git
	defer func() {
		deps.BareRepoPath = originalBareRepoPath
		deps.Git = originalGit
	}()

	// Set deps for the list command
	deps.BareRepoPath = bareRepoPath
	deps.Git = realGit

	// Get the list of worktrees
	worktreeList, err := buildWorktreeStrings(false)
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
	verboseList, err := buildWorktreeStrings(true)
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

	// Instead of using deleteWorktrees which calls util.CheckError (log.Fatal),
	// we'll call RemoveWorktree directly so we can handle errors properly in the test

	// Get the worktree objects to find our test_branch
	realTransformer := transformer.NewTransformer()
	worktreeStrings, err := realGit.GetWorktrees(&bareRepoPath)
	require.NoError(t, err, "Should be able to get worktrees")

	worktreeObjects := realTransformer.TransformWorktrees(worktreeStrings)

	// Find our test_branch worktree
	var testWorktree *worktreeobj.WorktreeObj
	for i, wt := range worktreeObjects {
		if wt.BranchName == testBranchName {
			testWorktree = &worktreeObjects[i]
			break
		}
	}
	require.NotNil(t, testWorktree, "Should find test_branch worktree object")

	// Resolve symlinks in bareRepoPath to match the FullPath format (macOS /var vs /private/var)
	resolvedBareRepoPath, err := filepath.EvalSymlinks(bareRepoPath)
	if err != nil {
		resolvedBareRepoPath = bareRepoPath
	}

	// Remove the worktree directly using the resolved path
	err = realGit.RemoveWorktree(testWorktree.FullPath, &resolvedBareRepoPath)
	assert.NoError(t, err, "Should successfully remove worktree")

	t.Logf("✓ Successfully deleted worktree")

	// Verify the worktree directory no longer exists
	_, err = os.Stat(worktreePath)
	assert.True(t, os.IsNotExist(err), "Worktree directory should no longer exist")

	// Verify the worktree is no longer in git worktree list
	worktreesAfterDelete, err := realGit.GetWorktrees(&bareRepoPath)
	assert.NoError(t, err, "Should be able to list worktrees after deletion")

	foundAfterDelete := false
	for _, wt := range worktreesAfterDelete {
		fields := splitWorktreeListLine(wt)
		if len(fields) > 0 {
			wtPath := fields[0]
			resolvedWtPath, _ := filepath.EvalSymlinks(wtPath)
			resolvedWorktreePath, _ := filepath.EvalSymlinks(worktreePath)

			if resolvedWtPath == resolvedWorktreePath || wtPath == worktreePath {
				foundAfterDelete = true
				break
			}
		}
	}
	assert.False(t, foundAfterDelete, "Worktree should not appear in git worktree list after deletion")

	// Verify the list command no longer shows it
	worktreeListAfterDelete, err := buildWorktreeStrings(false)
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

// assertPathExists is a helper function to check if a path exists
func assertPathExists(t *testing.T, path string) {
	t.Helper()
	_, err := os.Stat(path)
	assert.NoError(t, err, "Path should exist: %s", path)
}

// splitWorktreeListLine splits a git worktree list line into fields
// Format: "/path/to/worktree abcdef12345 [branch_name]"
func splitWorktreeListLine(line string) []string {
	var fields []string
	current := ""
	inBracket := false

	for i := 0; i < len(line); i++ {
		char := line[i]
		if char == ' ' && !inBracket {
			if current != "" {
				fields = append(fields, current)
				current = ""
			}
		} else if char == '[' {
			inBracket = true
			current += string(char)
		} else if char == ']' {
			inBracket = false
			current += string(char)
		} else {
			current += string(char)
		}
	}

	if current != "" {
		fields = append(fields, current)
	}

	return fields
}

// mockSpinner is a simple mock that does nothing - we don't need UI in tests
type mockSpinner struct{}

func (m *mockSpinner) Title(string) spinnerhuh.HuhSpinner { return m }
func (m *mockSpinner) Action(f func()) spinnerhuh.HuhSpinner {
	f() // Just execute the function immediately
	return m
}
func (m *mockSpinner) Run() error { return nil }
