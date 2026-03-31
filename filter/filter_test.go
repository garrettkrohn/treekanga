package filter

import (
	"testing"

	"github.com/garrettkrohn/treekanga/models"
	"github.com/stretchr/testify/assert"
)

func TestFilterWorktreesAndBranches(t *testing.T) {
	remoteBranches := []string{"repo1", "repo2", "repo4"}
	worktreeObjs := []models.Worktree{
		{
			FullPath:   "/path/to/repo1",
			Folder:     "repo1",
			BranchName: "branch1",
			CommitHash: "hash1",
		},
		{
			FullPath:   "/path/to/repo2",
			Folder:     "repo2",
			BranchName: "branch2",
			CommitHash: "hash2",
		},
		{
			FullPath:   "/path/to/repo3",
			Folder:     "repo3",
			BranchName: "branch3",
			CommitHash: "hash3",
		},
	}

	t.Run("TestGetNoBranchMatch", func(t *testing.T) {
		expected := []models.Worktree{
			{
				FullPath:   "/path/to/repo3",
				Folder:     "repo3",
				BranchName: "branch3",
				CommitHash: "hash3",
			},
		}

		f := &RealFilter{}
		result := f.GetBranchNoMatchList(remoteBranches, worktreeObjs)

		assert.Equal(t, result, expected)
	})

	t.Run("TestGetBranchMatch", func(t *testing.T) {
		selectedBranches := []string{"branch1", "branch2"}
		expected := []models.Worktree{
			{
				FullPath:   "/path/to/repo1",
				Folder:     "repo1",
				BranchName: "branch1",
				CommitHash: "hash1",
			},
			{
				FullPath:   "/path/to/repo2",
				Folder:     "repo2",
				BranchName: "branch2",
				CommitHash: "hash2",
			},
		}

		f := &RealFilter{}
		result := f.GetBranchMatchList(selectedBranches, worktreeObjs)

		assert.Equal(t, result, expected)
	})

	t.Run("TestGetBranchMatchWithSlashes", func(t *testing.T) {
		worktreesWithSlashes := []models.Worktree{
			{
				FullPath:   "/path/to/feature-abc",
				Folder:     "feature-abc",
				BranchName: "feature/abc",
				CommitHash: "hash1",
			},
			{
				FullPath:   "/path/to/bugfix-xyz",
				Folder:     "bugfix-xyz",
				BranchName: "bugfix/xyz",
				CommitHash: "hash2",
			},
			{
				FullPath:   "/path/to/main",
				Folder:     "main",
				BranchName: "main",
				CommitHash: "hash3",
			},
		}

		selectedBranches := []string{"feature/abc", "bugfix/xyz"}
		expected := []models.Worktree{
			{
				FullPath:   "/path/to/feature-abc",
				Folder:     "feature-abc",
				BranchName: "feature/abc",
				CommitHash: "hash1",
			},
			{
				FullPath:   "/path/to/bugfix-xyz",
				Folder:     "bugfix-xyz",
				BranchName: "bugfix/xyz",
				CommitHash: "hash2",
			},
		}

		f := &RealFilter{}
		result := f.GetBranchMatchList(selectedBranches, worktreesWithSlashes)

		assert.Equal(t, expected, result)
	})

}
