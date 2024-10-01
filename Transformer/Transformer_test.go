package transformer

import (
	"testing"

	worktreeobj "github.com/garrettkrohn/treekanga/worktreeObj"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	worktreeStrings := []string{
		"/Users/gkrohn/code/platform_work/platform_bare                                     (bare)",
		"/Users/gkrohn/code/platform_work/add_asset_regression                              94dbf65923 [add_asset_regression_fix]",
	}

	expected := []worktreeobj.WorktreeObj{
		{
			FullPath:   "/Users/gkrohn/code/platform_work/add_asset_regression",
			Folder:     "add_asset_regression",
			BranchName: "add_asset_regression_fix",
			CommitHash: "94dbf65923",
		},
	}

	t.Run("test worktree transformer", func(t *testing.T) {
		wt := &RealTransformer{}
		result := wt.TransformWorktrees(worktreeStrings)
		assert.Equal(t, result, expected)
	})
}
