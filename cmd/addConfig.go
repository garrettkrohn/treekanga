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

// resolveRepoNameAndPath implements the fallback logic for determining the repo name and bare repo path
// 1. First tries to use the current directory name
// 2. Checks sibling directories for matching bareRepoName configs
// 3. Checks if parent directory matches any repo's bareRepoName config
// 4. If that doesn't exist in config, falls back to git.GetRepoName()
// Returns: (repoName, bareRepoPath)
func resolveRepoNameAndPath() (string, string) {
	log.Debug("=== Starting bare repo resolution ===")

	// Get current working directory
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting working directory: ", err)
	}
	log.Debug("Current working directory", "path", workingDir)

	// Get directory name (parent directory of current working directory)
	parentDir := filepath.Dir(workingDir)
	directoryName := filepath.Base(parentDir)
	currentDirName := filepath.Base(workingDir)
	log.Debug("Directory info", "current", currentDirName, "parent", directoryName, "parentPath", parentDir)

	// Check if directory name exists in viper config
	log.Debug("Step 1: Checking if parent directory name exists in config", "checking", "repos."+directoryName)
	if viper.IsSet("repos." + directoryName) {
		log.Debug("✓ Repo directory name found in config", "directory name", directoryName)
		// Parent directory is likely the bare repo
		return "repos." + directoryName, parentDir
	}
	log.Debug("✗ Parent directory name not found in config")

	// Check sibling directories to find a matching bareRepoName
	// This handles cases where we're in a worktree and the bare repo is a sibling
	log.Debug("Step 2: Checking sibling directories for matching bareRepoName", "parentDir", parentDir)
	repoKey, bareRepoPath := findRepoByBareRepoInSiblings(parentDir)
	if repoKey != "" {
		log.Debug("✓ Repo found by checking sibling directories", "repo", repoKey, "bareRepoPath", bareRepoPath)
		return repoKey, bareRepoPath
	}
	log.Debug("✗ No matching bareRepoName found in siblings")

	// Check if the parent directory matches any repo's bareRepoName config
	log.Debug("Step 3: Checking if parent directory matches any bareRepoName config", "directoryName", directoryName)
	repoKey = findRepoByBareRepoName(directoryName)
	if repoKey != "" {
		log.Debug("✓ Repo found by bareRepoName match", "repo", repoKey, "bareRepoName", directoryName)
		// Parent directory is the bare repo
		return repoKey, parentDir
	}
	log.Debug("✗ Parent directory doesn't match any bareRepoName")

	// Check if we're in a nested structure (e.g., project/.bare/worktree)
	// Try going up one more level
	grandparentPath := filepath.Dir(parentDir)
	grandparentDir := filepath.Base(grandparentPath)
	log.Debug("Step 4: Checking grandparent directory", "grandparent", grandparentDir, "checking", "repos."+grandparentDir)
	if viper.IsSet("repos." + grandparentDir) {
		log.Debug("✓ Repo found by grandparent directory", "directory name", grandparentDir)
		return "repos." + grandparentDir, grandparentPath
	}
	log.Debug("✗ Grandparent directory not found in config")

	// Fallback to git.GetRepoName()
	log.Debug("Step 5: Falling back to git.GetRepoName()")
	repoName, err := deps.Git.GetRepoName(workingDir)
	if err != nil {
		log.Error("Error resolving repo name via git", "error", err)
		log.Fatal("Error resolving repo name: ", err)
	}
	log.Debug("Git repo name resolved", "repoName", repoName)

	// Check if git repo name exists in viper config
	log.Debug("Step 6: Checking if git repo name exists in config", "checking", "repos."+repoName)
	if viper.IsSet("repos." + repoName) {
		log.Debug("✓ Repo git directory name found in config", "repoName", repoName)
		// Try to find the actual bare repo path
		bareRepoPath = determineBareRepoPath(repoName, workingDir)
		return "repos." + repoName, bareRepoPath
	}
	log.Debug("✗ Git repo name not found in config")

	log.Error("Failed to resolve repo name through all methods")
	log.Fatal("No directory name, or git directory name found in the config")
	return "", ""
}

// determineBareRepoPath tries to determine the path to the bare repository
// when we've resolved the repo name but don't have a specific path
func determineBareRepoPath(repoName string, workingDir string) string {
	log.Debug("  → Attempting to determine bare repo path", "repoName", repoName, "workingDir", workingDir)

	// Check if there's a configured bareRepoName
	configuredBareRepoName := viper.GetString(repoName + ".bareRepoName")
	if configuredBareRepoName != "" {
		log.Debug("  → Found configured bareRepoName", "bareRepoName", configuredBareRepoName)
		// Look for this directory as a sibling
		parentDir := filepath.Dir(workingDir)
		bareRepoPath := filepath.Join(parentDir, configuredBareRepoName)
		if _, err := os.Stat(bareRepoPath); err == nil {
			log.Debug("  → ✓ Found bare repo at expected location", "path", bareRepoPath)
			return bareRepoPath
		}
	}

	// Default to working directory (might be the bare repo itself)
	log.Debug("  → Defaulting to working directory", "path", workingDir)
	return workingDir
}

// findRepoByBareRepoName searches all repo configs to find one with a matching bareRepoName
func findRepoByBareRepoName(bareRepoName string) string {
	log.Debug("  → Searching for bareRepoName in configs", "looking for", bareRepoName)
	repos := viper.GetStringMap("repos")
	log.Debug("  → Found repos in config", "count", len(repos), "repos", repos)

	for repoName := range repos {
		configuredBareRepoName := viper.GetString(fmt.Sprintf("repos.%s.bareRepoName", repoName))
		log.Debug("  → Checking repo", "repo", repoName, "configuredBareRepoName", configuredBareRepoName, "looking for", bareRepoName)

		if configuredBareRepoName != "" && configuredBareRepoName == bareRepoName {
			log.Debug("  → ✓ Match found!", "repo", repoName, "bareRepoName", configuredBareRepoName)
			return "repos." + repoName
		}
	}
	log.Debug("  → ✗ No matching bareRepoName found in any repo config")
	return ""
}

// findRepoByBareRepoInSiblings checks sibling directories to see if any match a configured bareRepoName
// This is useful when we're in a worktree and need to find the bare repo which is a sibling directory
// Returns: (repoKey, bareRepoPath)
func findRepoByBareRepoInSiblings(parentDir string) (string, string) {
	log.Debug("  → Checking sibling directories in parent", "parentDir", parentDir)

	// Read the parent directory to get all siblings
	entries, err := os.ReadDir(parentDir)
	if err != nil {
		log.Debug("  → ✗ Could not read parent directory", "error", err)
		return "", ""
	}

	log.Debug("  → Found entries in parent directory", "count", len(entries))

	// Check each sibling directory to see if it matches a configured bareRepoName
	for _, entry := range entries {
		if entry.IsDir() {
			dirName := entry.Name()
			log.Debug("  → Checking sibling directory", "name", dirName)
			repoKey := findRepoByBareRepoName(dirName)
			if repoKey != "" {
				bareRepoPath := filepath.Join(parentDir, dirName)
				log.Debug("  → ✓ Found bare repo in sibling directory!", "directory", dirName, "repo", repoKey, "path", bareRepoPath)
				return repoKey, bareRepoPath
			}
		} else {
			log.Debug("  → Skipping non-directory entry", "name", entry.Name())
		}
	}

	log.Debug("  → ✗ No matching bare repo found in sibling directories")
	return "", ""
}

func getAddCmdConfig(cmd *cobra.Command, args []string, c *com.AddConfig) {
	addCmdFlagsAndArgs(cmd, args, c)

	// Resolve repo name and bare repo path early
	repoName, bareRepoPath := resolveRepoNameAndPath()
	deps.ResolvedRepo = repoName
	deps.BareRepoPath = bareRepoPath

	// If user didn't provide -d flag, use the resolved bare repo path for git operations
	if c.Flags.Directory == nil && bareRepoPath != "" {
		log.Debug("Using resolved bare repo path for git operations", "path", bareRepoPath)
		c.Flags.Directory = &bareRepoPath
	}

	setWorkingAndParentDir(c)

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
