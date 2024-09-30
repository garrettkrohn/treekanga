/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	// "os/exec"
	"strconv"
	// "strings"
	//
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/garrettkrohn/treekanga/execwrap"
	"github.com/garrettkrohn/treekanga/filter"
	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/shell"
	worktreeobj "github.com/garrettkrohn/treekanga/worktreeObj"
	"github.com/garrettkrohn/treekanga/worktreeTransformer"
	"github.com/spf13/cobra"
)

// cleanCmd repesents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean unused worktrees",
	Long: `Compare all local worktree branches with the remote branches,
    allow the user to select which worktrees they would like to delete.`,
	Run: func(cmd *cobra.Command, args []string) {

		execWrap := execwrap.NewExec()
		shell := shell.NewShell(execWrap)
		git := git.NewGit(shell)

		branches, error := git.GetRemoteBranches()
		if error != nil {
			log.Fatal(error)
		}

		worktreeStrings, wError := git.GetWorktrees()
		if wError != nil {
			log.Fatal(wError)
		}

		worktreeTransformer := worktreetransformer.NewWorktreeTransformer()
		worktrees := worktreeTransformer.TransformWorktrees(worktreeStrings)

		filter := filter.NewFilter()
		noMatchList := filter.GetBranchNoMatchList(branches, worktrees)

		// transform worktreeobj into strings for selection
		var stringWorktrees []string
		for _, worktreeObj := range noMatchList {
			stringWorktrees = append(stringWorktrees, worktreeObj.BranchName)
		}

		var selections []string

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Value(&selections).
					OptionsFunc(func() []huh.Option[string] {
						return huh.NewOptions(stringWorktrees...)
					}, &stringWorktrees).
					Title("Local endpoints that do not exist on remote").
					Height(25),
			),
		)

		formErr := form.Run()
		if formErr != nil {
			log.Fatal(formErr)
		}

		//transform string selection back to worktreeobjs
		var selectedWorktreeObj []worktreeobj.WorktreeObj
		for _, worktreeobj := range noMatchList {
			for _, str := range selections {
				if worktreeobj.BranchName == str {
					selectedWorktreeObj = append(selectedWorktreeObj, worktreeobj)
					break
				}
			}
		}

		//remove worktrees

		numOfWorktreesRemoved := 0

		action := func() {

			for _, worktreeObj := range selectedWorktreeObj {
				git.RemoveWorktree(worktreeObj.Folder)
				numOfWorktreesRemoved++
			}
		}
		err := spinner.New().
			Title("Removing Worktrees").
			Action(action).
			Run()

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("worktrees removed: %s", strconv.Itoa(numOfWorktreesRemoved))

	},
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
