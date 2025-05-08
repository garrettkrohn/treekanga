/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	com "github.com/garrettkrohn/treekanga/common"
	"github.com/garrettkrohn/treekanga/directoryReader"
	"github.com/garrettkrohn/treekanga/transformer"
	util "github.com/garrettkrohn/treekanga/utility"

	"slices"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	newBranchName string
	baseBranch    string
)

// const tempZoxideName = "temp_treekanga_worktree"

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

		c := com.AddConfig{}
		getAddCmdConfig(cmd, args, &c)

		validateConfig(&c)

		log.Debug("Adding worktree with config:")
		PrintConfig(c)
		err := deps.Git.AddWorktree(&c)
		util.CheckError(err)

		log.Info(fmt.Sprintf("worktree %s created", newBranchName))

		foldersToAdd := getListOfZoxideEntries(&c.ZoxideConfig)
		addZoxideEntries(foldersToAdd)

	},
}

func connectSesh(c *com.AddConfig) {
	if c.Flags.Connect == nil {
		return
	}

	log.Debug("connecting to: %s", c.Flags.Connect)

	// shortestZoxide := findShortestString(foldersToAdd)
	shortestZoxide := slices.Min(c.ZoxideConfig.FoldersToAdd)
	subFolderIsValid := slices.Contains(c.ZoxideConfig.FoldersToAdd, *c.Flags.Connect)
	if subFolderIsValid {
		zoxidePath := shortestZoxide + "/" + *c.Flags.Connect
		log.Info(fmt.Sprintf("Sesh connect to %s", zoxidePath))
		deps.Sesh.SeshConnect(zoxidePath)
	} else {
		log.Info(fmt.Sprintf("Sesh connect to %s", shortestZoxide))
		deps.Sesh.SeshConnect(shortestZoxide)
	}
}

func getAddCmdConfig(cmd *cobra.Command, args []string, c *com.AddConfig) {
	addCmdFlagsAndArgs(cmd, args, c)
	setWorkingAndParentDir(c)
	getGitConfig(c)
	getZoxideConfig(c)
}

func getZoxideConfig(c *com.AddConfig) {
	c.ZoxideConfig = com.ZoxideConfig{

		NewBranchName:   c.GitConfig.NewBranchName,
		ParentDir:       c.ParentDir,
		FoldersToAdd:    viper.GetStringSlice("repos." + c.GitConfig.RepoName + ".zoxideFolders"),
		DirectoryReader: deps.DirectoryReader,
	}

}

func addCmdFlagsAndArgs(cmd *cobra.Command, args []string, c *com.AddConfig) {
	flags := com.AddCmdFlags{}
	directory, err := cmd.Flags().GetString("directory")
	if directory == "" {
		flags.Directory = nil
	} else {
		flags.Directory = &directory
	}
	util.CheckError(err)

	baseBranch, err := cmd.Flags().GetString("base")
	if baseBranch == "" {
		flags.BaseBranch = nil
	} else {
		flags.BaseBranch = &baseBranch
	}
	util.CheckError(err)

	// Refactor connect
	connect, err := cmd.Flags().GetString("connect")
	if connect == "" {
		flags.Connect = nil
	} else {
		flags.Connect = &connect
	}
	util.CheckError(err)

	// Refactor pull
	pull, err := cmd.Flags().GetBool("pull")
	if err != nil {
		flags.Pull = nil
	} else {
		flags.Pull = &pull
	}
	util.CheckError(err)

	c.Flags = flags
	c.Args = args
}

func setWorkingAndParentDir(c *com.AddConfig) {
	// working dir
	workingDir, err := os.Getwd()
	util.CheckError(err)
	if c.Flags.Directory != nil {
		workingDir = *c.Flags.Directory
	}

	//parent dir
	parentDir := filepath.Dir(workingDir)

	c.WorkingDir = workingDir
	c.ParentDir = parentDir

}

func getGitConfig(c *com.AddConfig) {

	if len(c.Args) == 1 {
		c.GitConfig.NewBranchName = c.Args[0]
	} else {
		log.Fatal("please include news branch name as an argument")
	}

	repoName, err := deps.Git.GetRepoName(c.WorkingDir)
	util.CheckError(err)
	c.GitConfig.RepoName = repoName

	if c.Flags.BaseBranch != nil {
		c.GitConfig.BaseBranchName = *c.Flags.BaseBranch
	} else {
		baseBranch = viper.GetString("repos." + repoName + ".defaultBranch")
		if baseBranch == "" {
			log.Fatal("There was no baseBranch provided, and no baseBranch in the config file")
		}
		c.GitConfig.BaseBranchName = baseBranch
	}

	t := transformer.NewTransformer()

	remoteBranches, err := deps.Git.GetRemoteBranches(c.Flags.Directory)
	util.CheckError(err)
	cleanRemoteBranches := t.RemoveOriginPrefix(remoteBranches)
	log.Debug(cleanRemoteBranches)
	c.GitConfig.NumOfRemoteBranches = len(cleanRemoteBranches)

	localBranches, err := deps.Git.GetLocalBranches(c.Flags.Directory)
	util.CheckError(err)
	cleanLocalBranches := t.RemoveQuotes(localBranches)
	log.Debug(cleanLocalBranches)
	c.GitConfig.NumOfLocalBranches = len(cleanLocalBranches)

	c.GitConfig.NewBranchExistsLocally = slices.Contains(cleanLocalBranches, c.GitConfig.NewBranchName)
	c.GitConfig.NewBranchExistsRemotely = slices.Contains(cleanRemoteBranches, c.GitConfig.NewBranchName)
	c.GitConfig.BaseBranchExistsLocally = slices.Contains(cleanLocalBranches, c.GitConfig.BaseBranchName)
	c.GitConfig.BaseBranchExistsRemotely = slices.Contains(cleanRemoteBranches, c.GitConfig.BaseBranchName)

	c.GitConfig.FolderPath = "../" + c.GitConfig.NewBranchName

}

func validateConfig(c *com.AddConfig) {

	// if a path is provided, be sure it exists
	if c.Flags.Directory != nil {
		log.Debug(fmt.Sprintf("inputted path: %s ", c.Flags.Directory))
		_, err := os.Stat(*c.Flags.Directory)
		if err != nil {
			log.Fatal("path does not exist")
		}
	}

	//baseBranch must exist
	if !c.GitConfig.BaseBranchExistsLocally && !c.GitConfig.BaseBranchExistsRemotely {
		log.Fatal("Base branch does not exist locally or remotely")
	}

}

func getListOfZoxideEntries(c *com.ZoxideConfig) []string {
	baseName := c.ParentDir + "/" + c.NewBranchName

	var foldersToAdd []string
	foldersToAdd = append(foldersToAdd, baseName)

	foldersToAdd = addConfigFolders(foldersToAdd, c.FoldersToAdd, baseName, c.DirectoryReader)

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

	addCmd.Flags().BoolP("pull", "p", false, "Pull the base branch before creating new branch")
	addCmd.Flags().StringP("connect", "c", "", "Automatically connect to a sesh upon creation")
	addCmd.Flags().StringP("base", "b", "", "Specify the base branch for the new worktree")
	addCmd.Flags().StringP("directory", "d", "", "Specify the directory to the bare repo where the worktree will be added")
}
