package transformer

import (
	"github.com/garrettkrohn/treekanga/worktreeObj"
	"strings"
)

type transformer interface {
	TransformWorktrees([]string) []worktreeobj.WorktreeObj
	RemoveOriginPrefix([]string) []string
}

type RealTransformer struct {
}

func NewWorktreeTransformer() *RealTransformer {
	return &RealTransformer{}
}

func (r *RealTransformer) TransformWorktrees(worktreeStrings []string) []worktreeobj.WorktreeObj {

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

func (r *RealTransformer) RemoveOriginPrefix(branchStrings []string) []string {
	for i, branch := range branchStrings {
		branchStrings[i] = strings.TrimSpace(strings.Replace(branch, "origin/", "", -1))
	}
	return branchStrings
}
