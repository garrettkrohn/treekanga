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
	"github.com/garrettkrohn/treekanga/transformer"
	util "github.com/garrettkrohn/treekanga/utility"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"slices"
)

var (
	newBranchName string
	baseBranch    string
)

type AddCmdFlags struct {
	Directory string
	Path      string
	Pull      bool
	Connect   string
}

type GitConfig struct {
	NewBranchName string
	BaseBrachName string
	RepoName      string
}

type TreekangaAddConfig struct {
	Flags      AddCmdFlags
	Args       []string
	GitConfig  GitConfig
	WorkingDir string
	ParentDir  string
}

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

		// addCmdFlags := getAddCmdFlags(cmd)
		c := getAddCmdConfig(cmd, args)

		validateConfig(&c)

		transformer := transformer.NewTransformer()

		// remove
		// if newBranchName == "" {
		// 	err := huh.NewInput().
		// 		Title("Input branch name").
		// 		Prompt("?").
		// 		Value(&newBranchName).
		// 		Run()
		// 	util.CheckError(err)
		// }

		// remove
		// if len(args) == 0 {
		// 	err := huh.NewInput().
		// 		Title("Input base branch (leave blank for default)").
		// 		Prompt("?").
		// 		Value(&baseBranch).
		// 		Run()
		// 	util.CheckError(err)
		// }

		// move to config
		if baseBranch == "" {
			baseBranch = viper.GetString("repos." + repoName + ".defaultBranch")
			if baseBranch == "" {
				log.Fatal("There was no baseBranch provided, and no baseBranch in the config file")
			}
		}

		// move to git config
		remoteBranches, err := deps.Git.GetRemoteBranches(addCmdFlags.Path)
		cleanRemoteBranches := transformer.RemoveOriginPrefix(remoteBranches)
		util.CheckError(err)
		localBranches, err := deps.Git.GetLocalBranches(addCmdFlags.Path)
		cleanLocalBranches := transformer.RemoveQuotes(localBranches)
		util.CheckError(err)

		// move to git config
		newBranchExistsLocally := slices.Contains(cleanLocalBranches, newBranchName)
		NewBranchExistsRemotely := slices.Contains(cleanRemoteBranches, newBranchName)
		baseBranchExistsLocally := slices.Contains(cleanLocalBranches, baseBranch)
		baseBranchExistsRemotely := slices.Contains(cleanRemoteBranches, baseBranch)

		// move to a general config log
		log.Debugf("newBranchExistsLocally: %v, newBranchExistsRemotely: %v, baseBranchExistsLocally: %v, baseBranchExistsRemotely: %v",
			newBranchExistsLocally, NewBranchExistsRemotely, baseBranchExistsLocally, baseBranchExistsRemotely)

		if !baseBranchExistsLocally && !baseBranchExistsRemotely {
			log.Fatal("Base branch does not exist locally or remotely")
		}

		folderName := "../" + newBranchName

		pull, err := cmd.Flags().GetBool("pull")

		err = deps.Git.AddWorktree(folderName, newBranchExistsLocally, NewBranchExistsRemotely,
			newBranchName, baseBranch, addCmdFlags.Path, pull, baseBranchExistsLocally, baseBranchExistsRemotely)
		util.CheckError(err)

		// if pull {
		// 	deps.Git.DeleteBranch(tempZoxideName, addCmdFlags.Path)
		// }

		log.Info(fmt.Sprintf("worktree %s created", newBranchName))

		foldersToAddFromConfig := viper.GetStringSlice("repos." + repoName + ".zoxideFolders")
		directoryReader := deps.DirectoryReader
		foldersToAdd := getListOfZoxideEntries(newBranchName, parentDir, foldersToAddFromConfig, directoryReader)

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

func getAddCmdConfig(cmd *cobra.Command, args []string) TreekangaAddConfig {
	baseConfig := getBaseConfig(cmd, args)
	getGitConfig(&baseConfig)
	return baseConfig
}

func getBaseConfig(cmd *cobra.Command, args []string) TreekangaAddConfig {
	addCmdFlags := getAddCmdFlags(cmd)

	workingDir, err := os.Getwd()
	util.CheckError(err)

	if addCmdFlags.Path != "" {
		workingDir = addCmdFlags.Path
	}

	parentDir := filepath.Dir(workingDir)

	return TreekangaAddConfig{
		Flags:      addCmdFlags,
		Args:       args,
		WorkingDir: workingDir,
		ParentDir:  parentDir,
	}
}

func getGitConfig(c *TreekangaAddConfig) {

	repoName, err := deps.Git.GetRepoName(c.WorkingDir)
	util.CheckError(err)
	c.GitConfig.RepoName = repoName

}

func getAddCmdFlags(cmd *cobra.Command) AddCmdFlags {
	directory, err := cmd.Flags().GetString("directory")
	util.CheckError(err)

	baseBranch, err := cmd.Flags().GetString("base")
	util.CheckError(err)

	connect, err := cmd.Flags().GetString("connect")
	util.CheckError(err)

	pull, err := cmd.Flags().GetBool("pull")
	util.CheckError(err)

	return AddCmdFlags{
		Directory: directory,
		Path:      baseBranch,
		Connect:   connect,
		Pull:      pull,
	}
}

func validateConfig(c *TreekangaAddConfig) {
	// make sure new branch name is included
	if len(c.Args) == 1 {
		c.GitConfig.NewBranchName = c.Args[0]
	} else {
		log.Fatal("please include news branch name as an argument")
	}

	// if a path is provided, be sure it exists
	if c.Flags.Path != "" {
		log.Debug(fmt.Sprintf("inputted path: %s ", c.Flags.Path))
		_, err := os.Stat(c.Flags.Path)
		if err != nil {
			log.Fatal("path does not exist")
		}
	}

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
