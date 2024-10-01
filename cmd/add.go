/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/garrettkrohn/treekanga/execwrap"
	"github.com/garrettkrohn/treekanga/filter"
	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/shell"
	// "github.com/garrettkrohn/treekanga/transformer"
	"github.com/spf13/cobra"
	// "log"
	// "os/exec"
	// "strings"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		execWrap := execwrap.NewExec()
		shell := shell.NewShell(execWrap)
		git := git.NewGit(shell)
		filter := filter.NewFilter()

		//TODO: make this async for performance
		// remoteBranches, _ := git.GetRemoteBranches()
		// cleanRemoteBranches := transformer.NewWorktreeTransformer().RemoveOriginPrefix(remoteBranches)
		localBranches, _ := git.GetLocalBranches()

		var branchName string
		var folderName string
		huh.NewInput().
			Title("Input branch name").
			Prompt("?").
			Value(&branchName).
			Run()

		huh.NewInput().
			Title("Input folder name (leave blank for same as folder)").
			Prompt("?").
			Value(&folderName).
			Run()

		// existsOnRemote := filter.BranchExistsInSlice(cleanRemoteBranches, branchName)
		existsLocally := filter.BranchExistsInSlice(localBranches, branchName)

		if folderName == "" {
			folderName = branchName
		}

		folderName = "../" + folderName

		action := func() { git.AddWorktree(folderName, existsLocally, branchName) }

		err := spinner.New().
			Title("Adding Worktree").
			Action(action).
			Run()
		if err != nil {
			fmt.Print(err)
		}

		fmt.Printf("worktree %s created", branchName)

	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
