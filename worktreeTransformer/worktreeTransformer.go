package worktreetransformer

import (
	"github.com/garrettkrohn/treekanga/worktreeObj"
	"strings"
)

type worktreeTransformer interface {
	TransformWorktrees([]string) []worktreeobj.WorktreeObj
}

type RealWorktreeTransformer struct {
}

func NewWorktreeTransformer() *RealWorktreeTransformer {
	return &RealWorktreeTransformer{}
}

func (r *RealWorktreeTransformer) TransformWorktrees(worktreeStrings []string) []worktreeobj.WorktreeObj {

	var worktrees []worktreeobj.WorktreeObj

	for _, worktreeString := range worktreeStrings {
		parts := strings.Fields(worktreeString)

		// takes care of bare repo and mysterious empty last worktree
		if len(parts) < 3 {
			continue
		}

		FullPath := parts[0]
		CommitHash := parts[1]

		// Split the FullPath by "/" and get the last part
		Folder := strings.Split(FullPath, "/")[len(strings.Split(FullPath, "/"))-1]

		// Remove the brackets from the branch name
		BranchName := strings.Trim(parts[2], "[]")

		worktrees = append(worktrees, worktreeobj.WorktreeObj{
			FullPath:   FullPath,
			Folder:     Folder,
			BranchName: BranchName,
			CommitHash: CommitHash,
		})
	}

	return worktrees
}
