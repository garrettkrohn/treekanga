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

func TestDeleteWorktreesWithoutArgs(t *testing.T) {
	// Arrange

	mockGit := getMockGit(t)

	transformer := transformer.NewTransformer()
	mockFilter := filter.NewMockFilter(t)

	mockSpinner := getMockSpinner(t)

	mockForm := getMockForm(t)

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

	var branches []string

	// Act
	numOfWorktreesRemoved, err := deleteWorktrees(mockGit, transformer, mockFilter, mockSpinner, mockForm, mockZoxide, branches, false, false)

	// Assert
	assert.NoError(t, err)

	// Verify the result
	expectedNumOfWorktreesRemoved := 1 // Adjust based on your test case
	assert.Equal(t, expectedNumOfWorktreesRemoved, numOfWorktreesRemoved)

	// Ensure all expectations are met
	mockGit.AssertExpectations(t)
	mockFilter.AssertExpectations(t)
	mockForm.AssertExpectations(t)
	mockSpinner.AssertExpectations(t) // Added to ensure spinner expectations are also checked
	mockZoxide.AssertExpectations(t)
}

func TestDeleteWorktreesWithArgs(t *testing.T) {

	// Arrange
	mockGit := getMockGit(t)

	transformer := transformer.NewTransformer()
	mockFilter := filter.NewMockFilter(t)
	mockFilter.On("GetBranchMatchList", mock.Anything, mock.Anything).Return([]worktreeobj.WorktreeObj{
		{
			FullPath:   "/Users/gkrohn/code/development",
			Folder:     "development",
			BranchName: "branch1",
			CommitHash: "abcdef12345",
		},
	})

	mockSpinner := getMockSpinner(t)

	mockZoxide := zoxide.NewMockZoxide(t)
	mockZoxide.On("RemovePath", mock.Anything).Return(nil)

	branches := []string{"branch1"}

	// Act
	numOfWorktreesRemoved, err := deleteWorktrees(mockGit, transformer, mockFilter, mockSpinner, nil, mockZoxide, branches, false, false)

	// Assert
	assert.NoError(t, err)

	// Verify the result
	expectedNumOfWorktreesRemoved := 1 // Adjust based on your test case
	assert.Equal(t, expectedNumOfWorktreesRemoved, numOfWorktreesRemoved)

	// Ensure all expectations are met
	mockGit.AssertExpectations(t)
	mockFilter.AssertExpectations(t)
	mockSpinner.AssertExpectations(t) // Added to ensure spinner expectations are also checked
	mockZoxide.AssertExpectations(t)
}

func getMockGit(t *testing.T) *git.MockGit {
	mockGit := git.NewMockGit(t)
	mockGit.On("GetWorktrees").Return([]string{
		"/Users/gkrohn/code/development       abcdef12345 [branch1]",
		"/Users/gkrohn/code/featureBranch     abcdef12345 [branch2]",
	}, nil)
	mockGit.On("RemoveWorktree", mock.Anything).Return("", nil)
	return mockGit
}

func getMockSpinner(t *testing.T) *spinner.MockHuhSpinner {
	mockSpinner := spinner.NewMockHuhSpinner(t)
	mockSpinner.On("Title", mock.Anything).Return(mockSpinner)
	mockSpinner.On("Action", mock.Anything).Run(func(args mock.Arguments) {
		// Call the action function
		args.Get(0).(func())()
	}).Return(mockSpinner)
	mockSpinner.On("Run").Run(func(args mock.Arguments) {
	}).Return(nil).Once()
	return mockSpinner
}

func getMockForm(t *testing.T) *form.MockHuhForm {
	mockForm := form.NewMockHuhForm(t)
	mockForm.On("SetSelections", mock.Anything).Run(func(args mock.Arguments) {
		// Modify the selections variable
		*args.Get(0).(*[]string) = append(*args.Get(0).(*[]string), "branch1")
	}).Return()
	mockForm.On("SetOptions", mock.Anything).Once()
	mockForm.On("Run").Return(nil)
	return mockForm
}
