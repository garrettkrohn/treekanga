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
	Short: "List worktrees",
	Long:  `List worktrees`,
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
