package services

import (
	"testing"

	"github.com/garrettkrohn/treekanga/adapters"
	"github.com/garrettkrohn/treekanga/config"
	"github.com/garrettkrohn/treekanga/filter"
	"github.com/garrettkrohn/treekanga/form"
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

	cfg := config.AppConfig{
		FilterOnlyStaleBranches: false,
		DeleteBranch:            false,
		ForceDelete:             false,
	}

	// Act
	numOfWorktreesRemoved, err := DeleteWorktrees(mockGit, transformer, mockFilter, mockForm, mockZoxide, branches, cfg)

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

	branches := []string{"development"}

	cfg := config.AppConfig{
		FilterOnlyStaleBranches: false,
		DeleteBranch:            false,
		ForceDelete:             false,
	}

	// Act
	numOfWorktreesRemoved, err := DeleteWorktrees(mockGit, transformer, mockFilter, nil, mockZoxide, branches, cfg)

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

func getMockGit(t *testing.T) *adapters.MockGitAdapter {
	mockGit := adapters.NewMockGitAdapter(t)
	mockGit.On("GetWorktrees", mock.Anything).Return([]string{
		"/Users/gkrohn/code/development       abcdef12345 [branch1]",
		"/Users/gkrohn/code/featureBranch     abcdef12345 [branch2]",
	}, nil)
	mockGit.On("RemoveWorktree", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	return mockGit
}

func getMockSpinner(t *testing.T) *spinner.MockHuhSpinner {
	mockSpinner := spinner.NewMockHuhSpinner(t)
	// Spinner is no longer used in removeWorktrees, so no expectations are set
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
