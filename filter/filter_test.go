package filter

import (
	// "fmt"
	"testing"

	"github.com/garrettkrohn/treekanga/worktreeObj"
	"github.com/stretchr/testify/assert"
)

func TestFilterWorktreesAndBranches(t *testing.T) {
	remoteBranches := []string{"branch1", "branch2", "branch4"}
	worktreeObjs := []worktreeobj.WorktreeObj{
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

	expected := []worktreeobj.WorktreeObj{
		{
			FullPath:   "/path/to/repo3",
			Folder:     "repo3",
			BranchName: "branch3",
			CommitHash: "hash3",
		},
	}

	t.Run("TestFilterWorktreesAndBranches", func(t *testing.T) {
		f := &RealFilter{}
		result := f.GetBranchNoMatchList(remoteBranches, worktreeObjs)

		assert.Equal(t, result, expected)
	})
	t.Run("Test branch exists on remote", func(t *testing.T) {
		f := &RealFilter{}
		result := f.BranchExistsInSlice(remoteBranches, "branch1")

		assert.Equal(t, result, true)
	})

	t.Run("Test branch does not exist on remote", func(t *testing.T) {
		f := &RealFilter{}
		result := f.BranchExistsInSlice(remoteBranches, "branch5")

		assert.Equal(t, result, false)
	})

}
