package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	com "github.com/garrettkrohn/treekanga/common"
	"github.com/garrettkrohn/treekanga/transformer"
	util "github.com/garrettkrohn/treekanga/utility"
)

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
