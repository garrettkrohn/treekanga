package cmd

import (
	"testing"

	"github.com/garrettkrohn/treekanga/git"
	"github.com/stretchr/testify/assert"
)

func TestListCmd(t *testing.T) {
	// Mock the GetWorktrees method
	mockGit := git.NewMockGit(t)
	mockGit.On("GetWorktrees").Return([]string{
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
