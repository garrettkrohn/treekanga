package cmd

import (
	"fmt"

	"github.com/garrettkrohn/treekanga/models"
	"github.com/garrettkrohn/treekanga/transformer"
)

type lister interface {
	list() ([]string, error)
}

func getLister(verbose, global, expand bool) lister {
	f := getFetcher(global)

	var t transformer.Transformer
	if verbose {
		t = &verboseTransformer{}
	} else {
		t = &simpleTransformer{}
	}

	return &listerImpl{fetcher: f, transformer: t}
}

type listerImpl struct {
	fetcher
	transformer transformer.Transformer
}

func (l *listerImpl) list() ([]string, error) {
	worktrees, err := l.fetcher.fetch()
	if err != nil {
		return nil, err
	}
	return l.transformer.Transform(worktrees)
}

type simpleTransformer struct{}

func (t *simpleTransformer) Transform(worktrees []models.Worktree) ([]string, error) {
	var worktreeStrings []string

	for _, worktree := range worktrees {
		if deps.AppConfig.ListDisplayMode == "directory" {
			worktreeStrings = append(worktreeStrings, worktree.Folder)
		} else {
			worktreeStrings = append(worktreeStrings, worktree.BranchName)
		}
	}

	return worktreeStrings, nil
}

type verboseTransformer struct{}

func (t *verboseTransformer) Transform(worktrees []models.Worktree) ([]string, error) {
	var worktreeBranches []string
	for _, worktree := range worktrees {
		branchDisplay := fmt.Sprintf("worktree: %s, branch: %s, fullPath: %s, commitHash: %s", worktree.Folder, worktree.BranchName, worktree.FullPath, worktree.CommitHash)
		worktreeBranches = append(worktreeBranches, branchDisplay)
	}
	return worktreeBranches, nil
}
