/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/charmbracelet/log"

	"github.com/garrettkrohn/treekanga/transformer"
	util "github.com/garrettkrohn/treekanga/utility"
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
		verbose, err := cmd.Flags().GetBool("verbose")
		util.CheckError(err)

		worktrees, err := buildWorktreeStrings(verbose)
		if err != nil {
			log.Fatal(err)
		}
		for _, worktree := range worktrees {
			fmt.Println(worktree)
		}
	},
}

func buildWorktreeStrings(verbose bool) ([]string, error) {
	rawWorktrees, err := deps.Git.GetWorktrees()
	if err != nil {
		return nil, err
	}

	worktreetransformer := transformer.NewTransformer()
	worktreeObjects := worktreetransformer.TransformWorktrees(rawWorktrees)

	var worktreeBranches []string
	for _, worktree := range worktreeObjects {
		var branchDisplay string
		if verbose {
			branchDisplay = fmt.Sprintf("worktree: %s, branch: %s, fullPath: %s, commitHash: %s", worktree.Folder, worktree.BranchName, worktree.FullPath, worktree.CommitHash)
		} else {
			branchDisplay = worktree.BranchName
		}
		worktreeBranches = append(worktreeBranches, branchDisplay)
	}

	return worktreeBranches, nil
}

func init() {
	listCmd.Flags().BoolP("verbose", "v", false, "Verbose display of worktrees")
}
