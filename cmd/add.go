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
	"github.com/charmbracelet/log"
	"github.com/garrettkrohn/treekanga/directoryReader"
	"github.com/garrettkrohn/treekanga/filter"
	"github.com/garrettkrohn/treekanga/transformer"
	util "github.com/garrettkrohn/treekanga/utility"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"slices"
)

var (
	branchName string
	baseBranch string
)

const tempZoxideName = "temp_treekanga_worktree"

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a git worktree",
	Long: `You may use this command with zero arguments, and you
    will be prompeted to input the branch name and base branch.

    Alternatively, you may the branch name as an argument, 
    treekanga will create this branch off of the defaultBranch 
    defined in the config, or use the current branch.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) >= 1 {
			branchName = args[0]
		}

		path, err := cmd.Flags().GetString("directory")
		log.Debug(path)

		baseBranch, err := cmd.Flags().GetString("base")
		util.CheckError(err)

		filter := filter.NewFilter()
		transformer := transformer.NewTransformer()

		workingDir, err := os.Getwd()
		util.CheckError(err)
		if path != "" {
			workingDir = path
		}

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

		// action := func() {
		remoteBranches, err := deps.Git.GetRemoteBranches(path)
		cleanRemoteBranches := transformer.RemoveOriginPrefix(remoteBranches)
		util.CheckError(err)
		localBranches, err := deps.Git.GetLocalBranches(path)
		cleanLocalBranches := transformer.RemoveQuotes(localBranches)
		util.CheckError(err)

		existsLocally := filter.BranchExistsInSlice(cleanLocalBranches, branchName)
		existsRemotely := filter.BranchExistsInSlice(cleanRemoteBranches, branchName)

		if existsRemotely {
			deps.Git.FetchOrigin(branchName, path)
			log.Debug("Branch exists remotely:", "branch name", branchName)
		} else {
			log.Debug("Branch does not exist remotely:", "branch name", branchName)
		}

		if existsLocally {
			log.Debug("Branch exists locally:", "branch name", branchName)
		} else {
			log.Debug("Branch does not exist locally:", "branch name", branchName)
		}

		folderName := "../" + branchName

		if baseBranch == "" {
			baseBranch = viper.GetString("repos." + repoName + ".defaultBranch")
			if baseBranch == "" {
				log.Fatal("There was no baseBranch provided, and no baseBranch in the config file")
			}
		}

		pull, err := cmd.Flags().GetBool("pull")
		if pull {
			log.Info("pulling base branch before creating worktree", "base branch", baseBranch)
			deps.Git.FetchOrigin(baseBranch, path)
			deps.Git.CreateTempBranch(path)
			baseBranch = tempZoxideName
		}

		err = deps.Git.AddWorktree(folderName, existsLocally, branchName, baseBranch, path)
		util.CheckError(err)

		if pull {
			deps.Git.DeleteBranch(tempZoxideName, path)
		}
		// }

		// err = spinner.New().
		// 	Title("Adding Worktree").
		// 	Action(action).
		// 	Run()
		// util.CheckError(err)

		log.Info(fmt.Sprintf("worktree %s created", branchName))

		foldersToAddFromConfig := viper.GetStringSlice("repos." + repoName + ".zoxideFolders")
		directoryReader := deps.DirectoryReader
		foldersToAdd := getListOfZoxideEntries(branchName, parentDir, foldersToAddFromConfig, directoryReader)

		addZoxideEntries(foldersToAdd)

		if cmd.Flags().Changed("connect") {
			connect, err := cmd.Flags().GetString("connect")
			log.Debug(connect)
			util.CheckError(err)

			// shortestZoxide := findShortestString(foldersToAdd)
			shortestZoxide := slices.Min(foldersToAdd)
			subFolderIsValid := slices.Contains(foldersToAddFromConfig, connect)
			if connect != "" && subFolderIsValid {
				zoxidePath := shortestZoxide + "/" + connect
				log.Info(fmt.Sprintf("Sesh connect to %s", zoxidePath))
				deps.Sesh.SeshConnect(zoxidePath)
			} else {
				log.Info(fmt.Sprintf("Sesh connect to %s", shortestZoxide))
				deps.Sesh.SeshConnect(shortestZoxide)
			}
		}

	},
}

func getListOfZoxideEntries(branchName string, parentDir string, foldersToAddFromConfig []string, directoryReader directoryReader.DirectoryReader) []string {
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
	addCmd.Flags().BoolP("pull", "p", false, "Pull the base branch before creating new branch")
	addCmd.Flags().StringP("connect", "c", "", "Automatically connect to a sesh upon creation")
	addCmd.Flags().StringP("base", "b", "", "Specify the base branch for the new worktree")
	addCmd.Flags().StringP("directory", "d", "", "Specify the directory to the bare repo where the worktree will be added")
}
