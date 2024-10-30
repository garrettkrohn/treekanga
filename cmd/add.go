/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/garrettkrohn/treekanga/directoryReader"
	"github.com/garrettkrohn/treekanga/filter"
	"github.com/garrettkrohn/treekanga/transformer"
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

		action := func() {
			remoteBranches, err := deps.Git.GetRemoteBranches()
			cleanRemoteBranches := transformer.NewTransformer().RemoveOriginPrefix(remoteBranches)
			util.CheckError(err)
			existsRemotely := filter.BranchExistsInSlice(cleanRemoteBranches, branchName)

			if existsRemotely {
				deps.Git.FetchOrigin(branchName)
			}

			folderName := "../" + branchName

			if baseBranch == "" {
				baseBranch = viper.GetString("repos." + repoName + ".defaultBranch")
			}

			deps.Git.AddWorktree(folderName, existsRemotely, branchName, baseBranch)
		}

		err = spinner.New().
			Title("Adding Worktree").
			Action(action).
			Run()
		util.CheckError(err)

		fmt.Printf("worktree %s created", branchName)

		foldersToAddFromConfig := viper.GetStringSlice("repos." + repoName + ".zoxideFolders")
		directoryReader := deps.DirectoryReader
		foldersToAdd := getListOfZoxideEntries(branchName, repoName, parentDir, foldersToAddFromConfig, directoryReader)

		addZoxideEntries(foldersToAdd)

		//TODO: optional kill local session, and open it with the new branch

	},
}

func getListOfZoxideEntries(branchName string, repoName string, parentDir string, foldersToAddFromConfig []string, directoryReader directoryReader.DirectoryReader) []string {
	baseName := parentDir + "/" + branchName

	var foldersToAdd []string
	foldersToAdd = append(foldersToAdd, baseName)

	foldersToAdd = addConfigFolders(foldersToAdd, foldersToAddFromConfig, baseName, directoryReader)

	return foldersToAdd
}

func addConfigFolders(foldersToAdd []string, foldersToAddFromConfig []string, baseName string, directoryReader directoryReader.DirectoryReader) []string {
	for _, folder := range foldersToAddFromConfig {
		if !isLastCharWildcard(folder) {
			newFolderFromConfig := baseName + "/" + folder
			foldersToAdd = append(foldersToAdd, newFolderFromConfig)
		} else {
			pathUpTillWildcard := getPathUntilLastSlash(folder)
			baseFolderToSearch := baseName + "/" + pathUpTillWildcard
			configFolders, err := directoryReader.GetFoldersInDirectory(baseFolderToSearch)

			for _, configFolder := range configFolders {
				newConfigFolder := baseFolderToSearch + "/" + configFolder
				foldersToAdd = append(foldersToAdd, newConfigFolder)
			}
			util.CheckError(err)
		}
	}
	return foldersToAdd
}

func isLastCharWildcard(input string) bool {
	parts := strings.Split(input, "/")
	lastSegment := parts[len(parts)-1]
	return strings.HasSuffix(lastSegment, "*")
}

func getPathUntilLastSlash(input string) string {
	parts := strings.Split(input, "/")
	if len(parts) > 1 {
		return strings.Join(parts[:len(parts)-1], "/")
	}
	return ""
}

func addZoxideEntries(folders []string) {
	for _, folder := range folders {
		err := deps.Zoxide.AddPath(folder)
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
