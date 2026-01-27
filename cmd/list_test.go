package cmd

import (
	"testing"

	"github.com/garrettkrohn/treekanga/adapters"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListCmd(t *testing.T) {
	// Mock the GetWorktrees method
	mockGit := adapters.NewMockGitAdapter(t)
	mockGit.On("GetWorktrees", mock.Anything).Return([]string{
		"/Users/gkrohn/code/development       abcdef12345 [branch1]",
		"/Users/gkrohn/code/featureBranch     abcdef12345 [branch2]",
	}, nil)

	// Replace the actual Git dependency with the mock
	deps.Git = mockGit

	// Execute the function
	worktrees, err := buildWorktreeStrings(false)
	assert.NoError(t, err)

	// Verify the output
	expectedOutput := []string{"branch1", "branch2"}
	assert.Equal(t, expectedOutput, worktrees)
}

func TestListCmdWithBranchDisplayMode(t *testing.T) {
	// Setup
	viper.Reset()
	mockGit := adapters.NewMockGitAdapter(t)
	mockGit.On("GetWorktrees", mock.Anything).Return([]string{
		"/Users/gkrohn/code/development       abcdef12345 [branch1]",
		"/Users/gkrohn/code/featureBranch     abcdef12345 [branch2]",
	}, nil)

	deps.Git = mockGit
	deps.AppConfig.RepoNameForConfig = "testRepo"

	// Set display mode to branch
	viper.Set("repos.testRepo.listDisplayMode", "branch")

	// Execute
	worktrees, err := buildWorktreeStrings(false)
	assert.NoError(t, err)

	// Verify - should display branch names
	expectedOutput := []string{"branch1", "branch2"}
	assert.Equal(t, expectedOutput, worktrees)
}

func TestListCmdWithDirectoryDisplayMode(t *testing.T) {
	// Setup
	viper.Reset()
	mockGit := adapters.NewMockGitAdapter(t)
	mockGit.On("GetWorktrees", mock.Anything).Return([]string{
		"/Users/gkrohn/code/development       abcdef12345 [branch1]",
		"/Users/gkrohn/code/featureBranch     abcdef12345 [branch2]",
	}, nil)

	deps.Git = mockGit
	deps.AppConfig.RepoNameForConfig = "testRepo"

	// Set display mode to directory
	viper.Set("repos.testRepo.listDisplayMode", "directory")

	// Execute
	worktrees, err := buildWorktreeStrings(false)
	assert.NoError(t, err)

	// Verify - should display directory names
	expectedOutput := []string{"development", "featureBranch"}
	assert.Equal(t, expectedOutput, worktrees)
}

func TestListCmdWithFolderDisplayMode(t *testing.T) {
	// Setup
	viper.Reset()
	mockGit := adapters.NewMockGitAdapter(t)
	mockGit.On("GetWorktrees", mock.Anything).Return([]string{
		"/Users/gkrohn/code/development       abcdef12345 [branch1]",
		"/Users/gkrohn/code/featureBranch     abcdef12345 [branch2]",
	}, nil)

	deps.Git = mockGit
	deps.AppConfig.RepoNameForConfig = "testRepo"

	// Set display mode to folder (alias for directory)
	viper.Set("repos.testRepo.listDisplayMode", "folder")

	// Execute
	worktrees, err := buildWorktreeStrings(false)
	assert.NoError(t, err)

	// Verify - should display directory names
	expectedOutput := []string{"development", "featureBranch"}
	assert.Equal(t, expectedOutput, worktrees)
}

func TestListCmdWithDefaultDisplayMode(t *testing.T) {
	// Setup
	viper.Reset()
	mockGit := adapters.NewMockGitAdapter(t)
	mockGit.On("GetWorktrees", mock.Anything).Return([]string{
		"/Users/gkrohn/code/development       abcdef12345 [branch1]",
		"/Users/gkrohn/code/featureBranch     abcdef12345 [branch2]",
	}, nil)

	deps.Git = mockGit
	deps.AppConfig.RepoNameForConfig = "testRepo"

	// Don't set any display mode - should default to branch

	// Execute
	worktrees, err := buildWorktreeStrings(false)
	assert.NoError(t, err)

	// Verify - should default to branch names
	expectedOutput := []string{"branch1", "branch2"}
	assert.Equal(t, expectedOutput, worktrees)
}

func TestListCmdVerboseOverridesDisplayMode(t *testing.T) {
	// Setup
	viper.Reset()
	mockGit := adapters.NewMockGitAdapter(t)
	mockGit.On("GetWorktrees", mock.Anything).Return([]string{
		"/Users/gkrohn/code/development       abcdef12345 [branch1]",
		"/Users/gkrohn/code/featureBranch     abcdef12345 [branch2]",
	}, nil)

	deps.Git = mockGit
	deps.AppConfig.RepoNameForConfig = "testRepo"

	// Set display mode to directory
	viper.Set("repos.testRepo.listDisplayMode", "directory")

	// Execute with verbose=true
	worktrees, err := buildWorktreeStrings(true)
	assert.NoError(t, err)

	// Verify - verbose should show all details regardless of display mode
	assert.Len(t, worktrees, 2)
	assert.Contains(t, worktrees[0], "worktree: development")
	assert.Contains(t, worktrees[0], "branch: branch1")
	assert.Contains(t, worktrees[0], "fullPath: /Users/gkrohn/code/development")
	assert.Contains(t, worktrees[0], "commitHash: abcdef12345")
}
