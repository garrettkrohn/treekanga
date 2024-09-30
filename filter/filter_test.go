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

	// mockPathwrap := new(pathwrap.MockPath)
	// mockGit := new(git.MockGit)
	// n := NewNamer(mockPathwrap, mockGit)
	//
	// t.Run("name for git repo", func(t *testing.T) {
	// 	mockGit.On("ShowTopLevel", "/Users/josh/c/dotfiles/.config/neovim").Return(true, "/Users/josh/c/dotfiles", nil)
	// 	mockGit.On("GitCommonDir", "/Users/josh/c/dotfiles/.config/neovim").Return(true, "", nil)
	// 	mockPathwrap.On("Base", "/Users/josh/c/dotfiles").Return("dotfiles")
	// 	name, _ := n.Name("/Users/josh/c/dotfiles/.config/neovim")
	// 	assert.Equal(t, "dotfiles/_config/neovim", name)
	// })
	//
	// t.Run("name for git worktree", func(t *testing.T) {
	// 	mockGit.On("ShowTopLevel", "/Users/josh/c/sesh/main").Return(true, "/Users/josh/c/sesh/main", nil)
	// 	mockGit.On("GitCommonDir", "/Users/josh/c/sesh/main").Return(true, "/Users/josh/c/sesh/.bare", nil)
	// 	mockPathwrap.On("Base", "/Users/josh/c/sesh").Return("sesh")
	// 	name, _ := n.Name("/Users/josh/c/sesh/main")
	// 	assert.Equal(t, "sesh/main", name)
	// })
	//
	// t.Run("returns base on non-git dir", func(t *testing.T) {
	// 	mockGit.On("ShowTopLevel", "/Users/josh/.config/neovim").Return(false, "", fmt.Errorf("not a git repository (or any of the parent"))
	// 	mockGit.On("GitCommonDir", "/Users/josh/.config/neovim").Return(false, "", fmt.Errorf("not a git repository (or any of the parent"))
	// 	mockPathwrap.On("Base", "/Users/josh/.config/neovim").Return("neovim")
	// 	name, _ := n.Name("/Users/josh/.config/neovim")
	// 	assert.Equal(t, "neovim", name)
	// })
}
