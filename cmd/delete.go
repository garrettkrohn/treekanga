/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/garrettkrohn/treekanga/filter"
	"github.com/garrettkrohn/treekanga/transformer"
	util "github.com/garrettkrohn/treekanga/utility"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete selected worktrees",
	Long:  `List all worktrees and selected multiple to be deleted`,
	Run: func(cmd *cobra.Command, args []string) {

		transformer := transformer.NewTransformer()

		worktrees := getWorktrees(deps.Git, transformer)

		stringWorktrees := transformer.TransformWorktreesToBranchNames(worktrees)

		var selections []string

		selections = HuhMultiSelect(selections, stringWorktrees)

		//transform string selection back to worktreeobjs
		selectedWorktreeObj := filter.NewFilter().GetBranchMatchList(selections, worktrees)

		//remove worktrees
		numOfWorktreesRemoved := 0

		util.UseSpinner("Removing Worktrees", func() {
			for _, worktreeObj := range selectedWorktreeObj {
				deps.Git.RemoveWorktree(worktreeObj.Folder)
				numOfWorktreesRemoved++
			}
		})

		fmt.Printf("worktrees removed: %s", strconv.Itoa(numOfWorktreesRemoved))

	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
