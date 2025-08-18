/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/charmbracelet/log"

	"github.com/garrettkrohn/treekanga/transformer"
)

type Worktree struct {
	Path string
	Head string
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all worktrees",
	Long: `Display all worktrees in the current repository.

    Shows the branch name for each worktree in the repository.
    This is useful for getting an overview of all active worktrees
    before performing operations like deletion or cleanup.`,
	Run: func(cmd *cobra.Command, args []string) {
		worktrees, err := listWorktrees()
		if err != nil {
			log.Fatal(err)
		}
		for _, worktree := range worktrees {
			fmt.Println(worktree)
		}
	},
}

func listWorktrees() ([]string, error) {
	rawWorktrees, err := deps.Git.GetWorktrees()
	if err != nil {
		return nil, err
	}

	worktreetransformer := transformer.NewTransformer()
	worktreeObjects := worktreetransformer.TransformWorktrees(rawWorktrees)

	var worktreeBranches []string
	for _, worktree := range worktreeObjects {
		worktreeBranches = append(worktreeBranches, worktree.BranchName)
	}

	return worktreeBranches, nil
}

func init() {
}
