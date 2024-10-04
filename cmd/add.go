/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/garrettkrohn/treekanga/execwrap"
	"github.com/garrettkrohn/treekanga/filter"
	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/shell"
	util "github.com/garrettkrohn/treekanga/utility"
	"github.com/garrettkrohn/treekanga/zoxide"

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
		zoxide := zoxide.NewZoxide(shell)

		//TODO: make this async for performance
		// remoteBranches, _ := git.GetRemoteBranches()
		// cleanRemoteBranches := transformer.NewWorktreeTransformer().RemoveOriginPrefix(remoteBranches)
		localBranches, _ := git.GetLocalBranches()

		var branchName string
		err := huh.NewInput().
			Title("Input branch name").
			Prompt("?").
			Value(&branchName).
			Run()
		util.CheckError(err)

		var baseBranch string
		err = huh.NewInput().
			Title("Input base branch (leave blank for default)").
			Prompt("?").
			Value(&baseBranch).
			Run()
		util.CheckError(err)

		// existsOnRemote := filter.BranchExistsInSlice(cleanRemoteBranches, branchName)
		existsLocally := filter.BranchExistsInSlice(localBranches, branchName)

		folderName := "../" + branchName

		if baseBranch == "" {
			baseBranch = "development"
		}

		action := func() { git.AddWorktree(folderName, existsLocally, branchName, baseBranch) }

		err = spinner.New().
			Title("Adding Worktree").
			Action(action).
			Run()
		util.CheckError(err)

		fmt.Printf("worktree %s created", branchName)

		addZoxideEntries(zoxide, branchName)

		//TODO: optional kill local session, and open it with the new branch

	},
}

func addZoxideEntries(zoxide zoxide.Zoxide, branchName string) {
	//TODO: zoxide entries
	workingDir, err := os.Getwd()
	util.CheckError(err)

	parentDir := filepath.Dir(workingDir)
	err = zoxide.AddPath(parentDir + "/" + branchName)
	util.CheckError(err)

	err = zoxide.AddPath(parentDir + "/" + branchName + "/ui")
	util.CheckError(err)

	err = zoxide.AddPath(parentDir + "/" + branchName + "/parent")
	util.CheckError(err)

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
