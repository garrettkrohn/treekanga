/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/garrettkrohn/treekanga/filter"
	"github.com/garrettkrohn/treekanga/form"
	"github.com/garrettkrohn/treekanga/git"
	spinner "github.com/garrettkrohn/treekanga/spinnerHuh"
	"github.com/garrettkrohn/treekanga/transformer"
	util "github.com/garrettkrohn/treekanga/utility"
	worktreeobj "github.com/garrettkrohn/treekanga/worktreeObj"
	"github.com/spf13/cobra"
)

// cleanCmd repesents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean unused worktrees",
	Long: `Compare all local worktree branches with the remote branches,
    allow the user to select which worktrees they would like to delete.`,
	Run: func(cmd *cobra.Command, args []string) {

		transformer := transformer.NewTransformer()

		var worktrees []worktreeobj.WorktreeObj
		spinner := spinner.NewRealHuhSpinner()
		spinner.Title("Fetching Worktrees")
		spinner.Action(func() {
			worktrees = getWorktrees(deps.Git, transformer)
		})
		spinner.Run()

		var cleanedBranches []string
		spinner.Title("Fetching Remote Branches")
		spinner.Action(func() {
			cleanedBranches = getRemoteBranches(deps.Git, transformer)
		})
		spinner.Run()

		filter := filter.NewFilter()
		noMatchList := filter.GetBranchNoMatchList(cleanedBranches, worktrees)

		if len(noMatchList) == 0 {
			fmt.Println("All local branches exist on remote")
			os.Exit(1)
		}

		// transform worktreeobj into strings for selection
		stringWorktrees := transformer.TransformWorktreesToBranchNames(noMatchList)

		var selections []string
		form := form.NewHuhForm()
		form.SetSelections(&selections)
		form.SetOptions(stringWorktrees)
		err := form.Run()
		util.CheckError(err)

		//transform string selection back to worktreeobjs
		selectedWorktreeObj := filter.GetBranchMatchList(selections, noMatchList)

		//remove worktrees
		spinner.Title("Removing Worktrees")
		spinner.Action(func() {
			for _, worktreeObj := range selectedWorktreeObj {
				deps.Git.RemoveWorktree(worktreeObj.Folder)
			}
		})
		spinner.Run()

		fmt.Printf("worktrees removed: %s", strconv.Itoa(len(selectedWorktreeObj)))

	},
}

func getWorktrees(git git.Git, transformer *transformer.RealTransformer) []worktreeobj.WorktreeObj {
	worktreeStrings, wError := git.GetWorktrees()
	if wError != nil {
		log.Fatal(wError)
	}

	worktrees := transformer.TransformWorktrees(worktreeStrings)

	return worktrees
}

func getRemoteBranches(git git.Git, transformer *transformer.RealTransformer) []string {
	branches, error := git.GetRemoteBranches()
	if error != nil {
		log.Fatal(error)
	}
	cleanedBranches := transformer.RemoveOriginPrefix(branches)

	return cleanedBranches
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
