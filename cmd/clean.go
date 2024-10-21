/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/garrettkrohn/treekanga/execwrap"
	"github.com/garrettkrohn/treekanga/filter"
	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/shell"
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

		//TODO: add spinner for all these calls
		execWrap := execwrap.NewExec()
		shell := shell.NewShell(execWrap)
		git := git.NewGit(shell)
		transformer := transformer.NewTransformer()

		var worktrees []worktreeobj.WorktreeObj
		util.UseSpinner("Fetching Worktrees", func() {
			worktrees = getWorktrees(git, transformer)
		})

		var cleanedBranches []string
		util.UseSpinner("Fetching Remote Branches", func() {
			cleanedBranches = getRemoteBranches(git, transformer)
		})

		filter := filter.NewFilter()
		noMatchList := filter.GetBranchNoMatchList(cleanedBranches, worktrees)

		if len(noMatchList) == 0 {
			fmt.Println("All local branches exist on remote")
			os.Exit(1)
		}

		// transform worktreeobj into strings for selection
		stringWorktrees := transformer.TransformWorktreesToBranchNames(noMatchList)

		var selections []string
		selections = HuhMultiSelect(selections, stringWorktrees)

		//transform string selection back to worktreeobjs
		selectedWorktreeObj := filter.GetBranchMatchList(selections, noMatchList)

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

func HuhMultiSelect(selections []string, stringOptions []string) []string {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Value(&selections).
				OptionsFunc(func() []huh.Option[string] {
					return huh.NewOptions(stringOptions...)
				}, &stringOptions).
				Title("Local endpoints that do not exist on remote").
				Height(25),
		),
	)

	formErr := form.Run()
	if formErr != nil {
		log.Fatal(formErr)
	}
	return selections
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
