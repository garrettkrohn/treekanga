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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	branchName string
	baseBranch string
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
		if len(args) >= 1 {
			branchName = args[0]
		}
		if len(args) == 2 {
			baseBranch = args[1]
		}

		execWrap := execwrap.NewExec()
		shell := shell.NewShell(execWrap)
		git := git.NewGit(shell)
		filter := filter.NewFilter()
		zoxide := zoxide.NewZoxide(shell)

		//TODO: make this async for performance
		// remoteBranches, _ := git.GetRemoteBranches()
		// cleanRemoteBranches := transformer.NewWorktreeTransformer().RemoveOriginPrefix(remoteBranches)
		localBranches, _ := git.GetLocalBranches()

		if branchName == "" {
			err := huh.NewInput().
				Title("Input branch name").
				Prompt("?").
				Value(&branchName).
				Run()
			util.CheckError(err)

		}

		if len(args) == 0 {
			err := huh.NewInput().
				Title("Input base branch (leave blank for default)").
				Prompt("?").
				Value(&baseBranch).
				Run()
			util.CheckError(err)
		}

		existsLocally := filter.BranchExistsInSlice(localBranches, branchName)

		folderName := "../" + branchName

		if baseBranch == "" {
			baseBranch = "development"
		}

		action := func() { git.AddWorktree(folderName, existsLocally, branchName, baseBranch) }

		err := spinner.New().
			Title("Adding Worktree").
			Action(action).
			Run()
		util.CheckError(err)

		fmt.Printf("worktree %s created", branchName)

		addZoxideEntries(zoxide, branchName, git)

		//TODO: optional kill local session, and open it with the new branch

	},
}

func addZoxideEntries(zoxide zoxide.Zoxide, branchName string, git git.Git) {
	//TODO: zoxide entries
	workingDir, err := os.Getwd()
	util.CheckError(err)

	repoName, err := git.GetRepoName(workingDir)
	util.CheckError(err)

	folders := viper.GetStringSlice("repos." + repoName + ".zoxideFolders")

	// add base
	parentDir := filepath.Dir(workingDir)
	err = zoxide.AddPath(parentDir + "/" + branchName)

	// add all from config
	for _, folder := range folders {
		err = zoxide.AddPath(parentDir + "/" + branchName + "/" + folder)
		util.CheckError(err)
	}

}

func init() {
	rootCmd.AddCommand(addCmd)

	// Add optional arguments
	// func (f *FlagSet) StringVarP(p *string, name, shorthand string, value string, usage string) {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
