package transformer

import (
	"testing"

	worktreeobj "github.com/garrettkrohn/treekanga/worktreeObj"
	"github.com/stretchr/testify/assert"
)

func TestTransformer(t *testing.T) {
	worktreeStrings := []string{
		"/Users/gkrohn/code/platform_work/platform_bare                                     (bare)",
		"/Users/gkrohn/code/platform_work/add_asset_regression                              94dbf65923 [add_asset_regression_fix]",
	}

	expectedWt := []worktreeobj.WorktreeObj{
		{
			FullPath:   "/Users/gkrohn/code/platform_work/add_asset_regression",
			Folder:     "add_asset_regression",
			BranchName: "add_asset_regression_fix",
			CommitHash: "94dbf65923",
		},
	}

	t.Run("test worktree transformer", func(t *testing.T) {
		transformer := &RealTransformer{}
		result := transformer.TransformWorktrees(worktreeStrings)
		assert.Equal(t, result, expectedWt)
	})

	branchStrings := []string{
		"  origin/main",
		"origin/develop",
	}

	expectedB := []string{
		"main",
		"develop",
	}

	t.Run("test clean branch names", func(t *testing.T) {
		transformer := &RealTransformer{}
		result := transformer.RemoveOriginPrefix(branchStrings)
		assert.Equal(t, result, expectedB)
	})
}
