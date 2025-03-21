package cmd

import (
	"testing"

	"github.com/garrettkrohn/treekanga/filter"
	"github.com/garrettkrohn/treekanga/form"
	"github.com/garrettkrohn/treekanga/git"
	spinner "github.com/garrettkrohn/treekanga/spinnerHuh"
	"github.com/garrettkrohn/treekanga/transformer"
	worktreeobj "github.com/garrettkrohn/treekanga/worktreeObj"
	"github.com/garrettkrohn/treekanga/zoxide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteWorktreesWithClean(t *testing.T) {
	// Setup dependencies
	mockGit := git.NewMockGit(t)
	mockGit.On("GetWorktrees").Return([]string{
		"/Users/gkrohn/code/development       abcdef12345 [branch1]",
		"/Users/gkrohn/code/featureBranch     abcdef12345 [branch2]",
	}, nil)
	mockGit.On("GetRemoteBranches", mock.Anything).Return([]string{
		"/Users/gkrohn/code/development       abcdef12345 [branch1]",
	}, nil)
	mockGit.On("RemoveWorktree", "development").Return("", nil)

	transformer := transformer.NewTransformer()
	mockFilter := filter.NewMockFilter(t)
	mockFilter.On("GetBranchNoMatchList", mock.Anything, mock.Anything).Return([]worktreeobj.WorktreeObj{
		{
			FullPath:   "/Users/gkrohn/code/development",
			Folder:     "development",
			BranchName: "branch1",
			CommitHash: "abcdef12345",
		},
	}, nil)

	mockSpinner := spinner.NewMockHuhSpinner(t)
	mockSpinner.On("Title", mock.Anything).Return(mockSpinner)
	mockSpinner.On("Action", mock.Anything).Run(func(args mock.Arguments) {
		// Call the action function
		args.Get(0).(func())()
	}).Return(mockSpinner)
	mockSpinner.On("Run").Run(func(args mock.Arguments) {}).Return(nil)

	mockForm := form.NewMockHuhForm(t)
	mockForm.On("SetSelections", mock.Anything).Run(func(args mock.Arguments) {
		// Modify the selections variable
		*args.Get(0).(*[]string) = append(*args.Get(0).(*[]string), "branch1")
	}).Return()
	mockForm.On("SetOptions", mock.Anything).Once()
	mockForm.On("Run", mock.Anything).Return(nil)

	mockFilter.On("GetBranchMatchList", mock.Anything, mock.Anything).Return([]worktreeobj.WorktreeObj{
		{
			FullPath:   "/Users/gkrohn/code/development",
			Folder:     "development",
			BranchName: "branch1",
			CommitHash: "abcdef12345",
		},
	})
	mockZoxide := zoxide.NewMockZoxide(t)
	mockZoxide.On("RemovePath", mock.Anything).Return(nil)

	// Execute the function
	numOfWorktreesRemoved, err := cleanWorktrees(mockGit, transformer, mockFilter, mockSpinner, mockForm, mockZoxide)
	assert.NoError(t, err)

	// Verify the result
	expectedNumOfWorktreesRemoved := 1 // Adjust based on your test case
	assert.Equal(t, expectedNumOfWorktreesRemoved, numOfWorktreesRemoved)

	// Ensure all expectations are met
	// mockGit.AssertExpectations(t)
	// mockFilter.AssertExpectations(t)
	// mockForm.AssertExpectations(t)
	// mockSpinner.AssertExpectations(t) // Added to ensure spinner expectations are also checked
}
