/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/garrettkrohn/treekanga/filter"
	util "github.com/garrettkrohn/treekanga/utility"

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
	Short: "Add a git worktree",
	Long: `You may use this command with zero arguments, and you
    will be prompeted to input the branch name and base branch.

    Alternatively, you may the branch name as an argument, 
    treekanga will create this branch off of the defaultBranch 
    defined in the config, or use the current branch.

    You may also pass in the new branch and the base branch as
    arguments.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 1 {
			branchName = args[0]
		}
		if len(args) == 2 {
			baseBranch = args[1]
		}

		filter := filter.NewFilter()

		workingDir, err := os.Getwd()
		util.CheckError(err)

		repoName, err := deps.Git.GetRepoName(workingDir)
		util.CheckError(err)

		parentDir := filepath.Dir(workingDir)

		localBranches, _ := deps.Git.GetLocalBranches()

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
			baseBranch = viper.GetString("repos." + repoName + ".defaultBranch")
		}

		action := func() { deps.Git.AddWorktree(folderName, existsLocally, branchName, baseBranch) }

		err = spinner.New().
			Title("Adding Worktree").
			Action(action).
			Run()
		util.CheckError(err)

		fmt.Printf("worktree %s created", branchName)

		addZoxideEntries(branchName, repoName, parentDir)

		//TODO: optional kill local session, and open it with the new branch

	},
}

func addZoxideEntries(branchName string, repoName string, parentDir string) {
	folders := viper.GetStringSlice("repos." + repoName + ".zoxideFolders")

	// add base
	err := deps.Zoxide.AddPath(parentDir + "/" + branchName)
	util.CheckError(err)

	// add all from config
	for _, folder := range folders {
		err = deps.Zoxide.AddPath(parentDir + "/" + branchName + "/" + folder)
		util.CheckError(err)
	}

}

func init() {

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
