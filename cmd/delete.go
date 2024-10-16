/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/garrettkrohn/treekanga/execwrap"
	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/shell"
	worktreeobj "github.com/garrettkrohn/treekanga/worktreeObj"
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

		worktrees := getWorktrees(git)

		var stringWorktrees []string
		for _, worktreeObj := range worktrees {
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
		for _, worktreeobj := range worktrees {
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
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
