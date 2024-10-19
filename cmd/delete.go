/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/garrettkrohn/treekanga/execwrap"
	"github.com/garrettkrohn/treekanga/filter"
	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/shell"
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

		execWrap := execwrap.NewExec()
		shell := shell.NewShell(execWrap)
		git := git.NewGit(shell)
		transformer := transformer.NewTransformer()

		worktrees := getWorktrees(git, transformer)

		stringWorktrees := transformer.TransformWorktreesToBranchNames(worktrees)

		var selections []string

		selections = HuhMultiSelect(selections, stringWorktrees)

		//transform string selection back to worktreeobjs
		selectedWorktreeObj := filter.NewFilter().GetBranchMatchList(selections, worktrees)

		//remove worktrees
		numOfWorktreesRemoved := 0

		util.UseSpinner("Removing Worktrees", func() {
			for _, worktreeObj := range selectedWorktreeObj {
				git.RemoveWorktree(worktreeObj.Folder)
				numOfWorktreesRemoved++
			}
		})

		fmt.Printf("worktrees removed: %s", strconv.Itoa(numOfWorktreesRemoved))

	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
