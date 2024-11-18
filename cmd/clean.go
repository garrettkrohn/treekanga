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

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean unused worktrees",
	Long: `Compare all local worktree branches with the remote branches,
    allow the user to select which worktrees they would like to delete.`,
	Run: func(cmd *cobra.Command, args []string) {
		numOfWorktreesRemoved, err := cleanWorktrees(deps.Git, transformer.NewTransformer(), filter.NewFilter(), spinner.NewRealHuhSpinner(), form.NewHuhForm())
		if err != nil {
			cmd.PrintErrln("Error:", err)
			return
		}
		fmt.Printf("worktrees removed: %s", strconv.Itoa(numOfWorktreesRemoved))
	},
}

// cleanWorktrees performs the core logic of cleaning worktrees
func cleanWorktrees(git git.Git, transformer *transformer.RealTransformer, filter filter.Filter, spinner spinner.HuhSpinner, form form.Form) (int, error) {
	var worktrees []worktreeobj.WorktreeObj
	spinner.Title("Fetching Worktrees")
	spinner.Action(func() {
		worktrees = getWorktrees(git, transformer)
	})
	spinner.Run()

	var cleanedBranches []string
	spinner.Title("Fetching Remote Branches")
	spinner.Action(func() {
		cleanedBranches = getRemoteBranches(git, transformer)
	})
	spinner.Run()

	noMatchList := filter.GetBranchNoMatchList(cleanedBranches, worktrees)

	if len(noMatchList) == 0 {
		fmt.Println("All local branches exist on remote")
		os.Exit(1)
	}

	// Transform worktree objects into strings for selection
	stringWorktrees := transformer.TransformWorktreesToBranchNames(noMatchList)

	var selections []string
	form.SetSelections(&selections)
	form.SetOptions(stringWorktrees)
	err := form.Run()
	util.CheckError(err)

	// Transform string selection back to worktree objects
	selectedWorktreeObj := filter.GetBranchMatchList(selections, noMatchList)

	removeWorktrees(selectedWorktreeObj, spinner)

	return len(selectedWorktreeObj), nil
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
