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

// resolveRepoName implements the fallback logic for determining the repo name
// 1. First tries to use the current directory name
// 2. If that doesn't exist in config, falls back to git.GetRepoName()
func resolveRepoName() string {
	// Get current working directory
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting working directory: ", err)
	}
	log.Debug("workingDir", workingDir)

	// Get directory name
	directoryName := filepath.Base(filepath.Dir(workingDir))
	log.Debug("directoryName", directoryName)

	// Check if directory name exists in viper config
	if viper.IsSet("repos." + directoryName) {
		log.Debug("Repo directory name found: ", "directory name", directoryName)
		return "repos." + directoryName
	}

	// Fallback to git.GetRepoName()
	repoName, err := deps.Git.GetRepoName(workingDir)
	if err != nil {
		log.Fatal("Error resolving repo name: ", err)
	}

	// Check if git repo name exists in viper config
	if viper.IsSet("repos." + repoName) {
		log.Debug("Repo git directory name found: ", directoryName)
		return "repos." + repoName
	}

	log.Fatal("No directory name, or git directory name found in the config")
	return ""
}

func getAddCmdConfig(cmd *cobra.Command, args []string, c *com.AddConfig) {
	addCmdFlagsAndArgs(cmd, args, c)
	setWorkingAndParentDir(c)

	// Resolve repo name only when needed (for add command)
	deps.ResolvedRepo = resolveRepoName()

	getGitConfig(c)
	getZoxideConfig(c)
	getPostScript(c)
}

func getZoxideConfig(c *com.AddConfig) {
	c.ZoxideFolders = viper.GetStringSlice(deps.ResolvedRepo + ".zoxideFolders")
	c.DirectoryReader = deps.DirectoryReader
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

	sesh, err := cmd.Flags().GetString("sesh")
	if sesh == "" {
		flags.Sesh = nil
	} else {
		flags.Sesh = &sesh
	}
	util.CheckError(err)

	pull, err := cmd.Flags().GetBool("pull")
	if err != nil {
		flags.Pull = nil
	} else {
		flags.Pull = &pull
	}
	util.CheckError(err)

	cursor, err := cmd.Flags().GetBool("cursor")
	if err != nil {
		flags.Cursor = nil
	} else {
		flags.Cursor = &cursor
	}
	util.CheckError(err)

	vscode, err := cmd.Flags().GetBool("vscode")
	if err != nil {
		flags.VsCode = nil
	} else {
		flags.VsCode = &vscode
	}
	util.CheckError(err)

	specifiedWorktreeName, err := cmd.Flags().GetString("name")
	if err != nil {
		flags.SpecifiedWorktreeName = nil
	} else {
		flags.SpecifiedWorktreeName = &specifiedWorktreeName
	}
	util.CheckError(err)

	executeScript, err := cmd.Flags().GetBool("script")
	if err != nil {
		flags.ExecuteScript = nil
	} else {
		flags.ExecuteScript = &executeScript
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
		c.GitInfo.NewBranchName = c.Args[0]
	} else {
		log.Fatal("please include new branch name as an argument")
	}

	c.GitInfo.RepoName = deps.ResolvedRepo

	if c.Flags.BaseBranch != nil {
		c.GitInfo.BaseBranchName = *c.Flags.BaseBranch
	} else {
		baseBranch = viper.GetString(deps.ResolvedRepo + ".defaultBranch")
		if baseBranch == "" {
			log.Fatal("There was no baseBranch provided, and no baseBranch in the config file")
		}
		c.GitInfo.BaseBranchName = baseBranch
	}

	t := transformer.NewTransformer()

	remoteBranches, err := deps.Git.GetRemoteBranches(c.Flags.Directory)
	util.CheckError(err)
	cleanRemoteBranches := t.RemoveOriginPrefix(remoteBranches)
	log.Debug(cleanRemoteBranches)

	localBranches, err := deps.Git.GetLocalBranches(c.Flags.Directory)
	util.CheckError(err)
	cleanLocalBranches := t.RemoveQuotes(localBranches)
	log.Debug(cleanLocalBranches)

	c.GitInfo.NewBranchExistsLocally = slices.Contains(cleanLocalBranches, c.GetNewBranchName())
	c.GitInfo.NewBranchExistsRemotely = slices.Contains(cleanRemoteBranches, c.GetNewBranchName())
	c.GitInfo.BaseBranchExistsLocally = slices.Contains(cleanLocalBranches, c.GetBaseBranchName())
	c.GitInfo.BaseBranchExistsRemotely = slices.Contains(cleanRemoteBranches, c.GetBaseBranchName())

	c.WorktreeTargetDir = resolveWorktreeTargetDir(deps.ResolvedRepo, c)
}

// resolveWorktreeTargetDir determines the target directory for the new worktree
// based on configuration and user preferences
func resolveWorktreeTargetDir(repoName string, c *com.AddConfig) string {
	// Determine the worktree name - either user specified or branch name
	worktreeName := getWorktreeName(c)

	// Check if there's a configured worktree target directory
	configWorktreeTargetDir := viper.GetString(repoName + ".worktreeTargetDir")

	if configWorktreeTargetDir != "" {
		// Use configured directory under home path
		homePath, err := os.UserHomeDir()
		if err != nil {
			log.Fatal("Error getting home directory: ", err)
		}
		return buildConfigWorktreeDir(homePath, configWorktreeTargetDir, worktreeName)
	} else {
		// Default to relative path from parent directory
		return "../" + worktreeName
	}
}

// getWorktreeName returns the name to use for the worktree directory
func getWorktreeName(c *com.AddConfig) string {
	if c.Flags.SpecifiedWorktreeName != nil && *c.Flags.SpecifiedWorktreeName != "" {
		return *c.Flags.SpecifiedWorktreeName
	}
	return c.GetNewBranchName()
}

func buildConfigWorktreeDir(homePath string, configWorktreeTargetDir string, branchName string) string {
	if configWorktreeTargetDir == "" {
		return filepath.Join(homePath, branchName)
	}
	return filepath.Join(homePath, configWorktreeTargetDir, branchName)
}

func validateConfig(c *com.AddConfig) {

	// if a path is provided, be sure it exists
	if c.Flags.Directory != nil {
		log.Debug(fmt.Sprintf("inputted path: %s ", *c.Flags.Directory))
		_, err := os.Stat(*c.Flags.Directory)
		if err != nil {
			log.Fatal("path does not exist")
		}
	}

	//baseBranch must exist
	if !c.GitInfo.BaseBranchExistsLocally && !c.GitInfo.BaseBranchExistsRemotely {
		log.Fatal("Base branch does not exist locally or remotely")
	}

}

func getPostScript(c *com.AddConfig) {
	postScript := viper.GetString(deps.ResolvedRepo + ".postScript")
	if postScript == "" {
		log.Debug("no post script found in config file")
		return
	}
	c.PostScript = postScript

	autoRunPostScript := viper.GetBool(deps.ResolvedRepo + ".autoRunPostScript")
	c.AutoRunPostScript = &autoRunPostScript

}
