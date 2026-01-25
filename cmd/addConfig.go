package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	com "github.com/garrettkrohn/treekanga/common"
	"github.com/garrettkrohn/treekanga/transformer"
	"github.com/garrettkrohn/treekanga/utility"
	util "github.com/garrettkrohn/treekanga/utility"
	worktreeobj "github.com/garrettkrohn/treekanga/worktreeObj"
)

// resolveRepoNameAndPath implements the fallback logic for determining the repo name and bare repo path
// 1. Checks if parent directory matches any repo's bareRepoName config
// 2. Tries to use the parent directory name
// 3. If that doesn't exist in config, falls back to git.GetRepoName()
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

	bareRepoPath, err := deps.Git.GetBareRepoPath()
	utility.CheckError(err)

	projectName, err := deps.Git.GetProjectName()
	utility.CheckError(err)

	if bareRepoPath != "" && projectName != "" {
		return projectName, bareRepoPath
	}

	// Check if the parent directory matches any repo's bareRepoName config
	log.Debug("Step 1: Checking if parent directory matches any bareRepoName config", "directoryName", directoryName)
	repoKey := findRepoByBareRepoName(directoryName, parentDir)
	if repoKey != "" {
		log.Debug("✓ Repo found by bareRepoName match", "repo", repoKey, "bareRepoName", directoryName)
		// Parent directory is the bare repo
		return repoKey, parentDir
	}
	log.Debug("✗ Parent directory doesn't match any bareRepoName")

	// Check if directory name exists in viper config
	log.Debug("Step 2: Checking if parent directory name exists in config", "checking", "repos."+directoryName)
	if viper.IsSet("repos." + directoryName) {
		log.Debug("✓ Repo directory name found in config", "directory name", directoryName)
		bareRepoName := viper.GetString("repos." + directoryName + ".bareRepoName")
		log.Debug(fmt.Sprintf("bareRepoName %s", bareRepoName))

		return "repos." + directoryName, parentDir + "/" + bareRepoName
	}
	log.Debug("✗ Parent directory name not found in config")

	// Check if we're in a nested structure (e.g., project/.bare/worktree)
	// Try going up one more level
	grandparentPath := filepath.Dir(parentDir)
	grandparentDir := filepath.Base(grandparentPath)
	log.Debug("Step 3: Checking grandparent directory", "grandparent", grandparentDir, "checking", "repos."+grandparentDir)
	if viper.IsSet("repos." + grandparentDir) {
		log.Debug("✓ Repo found by grandparent directory", "directory name", grandparentDir)
		return "repos." + grandparentDir, grandparentPath
	}
	log.Debug("✗ Grandparent directory not found in config")

	// Fallback to git.GetRepoName()
	log.Debug("Step 4: Falling back to git.GetRepoName()")
	repoName, err := deps.Git.GetRepoName(workingDir)
	if err != nil {
		log.Error("Error resolving repo name via git", "error", err)
		log.Fatal("Error resolving repo name: ", err)
	}
	log.Debug("Git repo name resolved", "repoName", repoName)

	// Check if git repo name exists in viper config
	log.Debug("Step 5: Checking if git repo name exists in config", "checking", "repos."+repoName)
	if viper.IsSet("repos." + repoName) {
		log.Debug("✓ Repo git directory name found in config", "repoName", repoName)
		// Try to find the actual bare repo path
		bareRepoPath := determineBareRepoPath(repoName, workingDir)
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
	// If we're in a worktree, read the .git file to find the bare repo
	gitFilePath := filepath.Join(workingDir, ".git")
	if content, err := os.ReadFile(gitFilePath); err == nil {
		gitFileContent := string(content)
		// Check if this is a worktree (file format: "gitdir: /path/to/.bare/worktrees/worktree_name")
		if strings.HasPrefix(gitFileContent, "gitdir: ") {
			gitdir := strings.TrimSpace(strings.TrimPrefix(gitFileContent, "gitdir: "))
			log.Debug("  → Found .git file pointing to", "gitdir", gitdir)
			// Extract bare repo path from gitdir
			// gitdir format: /path/to/.bare/worktrees/worktree_name
			// We need to go up two levels: worktree_name -> worktrees -> .bare
			if gitdir != "" {
				bareRepoPath := filepath.Dir(filepath.Dir(gitdir))
				if _, err := os.Stat(bareRepoPath); err == nil {
					log.Debug("  → ✓ Found bare repo from .git file", "path", bareRepoPath)
					return bareRepoPath
				}
			}
		}
	}

	// Default to working directory (might be the bare repo itself)
	log.Debug("  → Defaulting to working directory", "path", workingDir)
	return workingDir
}

// findRepoByBareRepoName searches all repo configs to find one with a matching bareRepoName
func findRepoByBareRepoName(bareRepoName string, parentDir ...string) string {
	log.Debug("  → Searching for bareRepoName in configs", "looking for", bareRepoName)
	repos := viper.GetStringMap("repos")
	log.Debug("  → Found repos in config", "count", len(repos), "repos", repos)

	// Collect all matching repos
	var matches []string
	for repoName := range repos {
		configuredBareRepoName := viper.GetString(fmt.Sprintf("repos.%s.bareRepoName", repoName))
		log.Debug("  → Checking repo", "repo", repoName, "configuredBareRepoName", configuredBareRepoName, "looking for", bareRepoName)

		if configuredBareRepoName != "" && configuredBareRepoName == bareRepoName {
			log.Debug("  → ✓ Match found!", "repo", repoName, "bareRepoName", configuredBareRepoName)
			matches = append(matches, repoName)
		}
	}

	// If no matches, return empty
	if len(matches) == 0 {
		log.Debug("  → ✗ No matching bareRepoName found in any repo config")
		return ""
	}

	// If only one match, return it
	if len(matches) == 1 {
		return "repos." + matches[0]
	}

	// Multiple matches found - try to disambiguate using worktreetargetdir and parentDir
	log.Debug("  → Multiple matches found, disambiguating", "matches", matches, "count", len(matches))

	if len(parentDir) > 0 {
		currentParent := parentDir[0]
		log.Debug("  → Using parent directory for disambiguation", "parentDir", currentParent)

		// Check each match for worktreetargetdir
		for _, repoName := range matches {
			worktreeTargetDir := viper.GetString(fmt.Sprintf("repos.%s.worktreeTargetDir", repoName))
			log.Debug("  → Checking worktreetargetdir", "repo", repoName, "worktreeTargetDir", worktreeTargetDir)

			if worktreeTargetDir != "" {
				// Expand the worktreeTargetDir if it starts with ~
				if len(worktreeTargetDir) > 0 && worktreeTargetDir[0] == '~' {
					homeDir, err := os.UserHomeDir()
					if err == nil {
						worktreeTargetDir = filepath.Join(homeDir, worktreeTargetDir[1:])
					}
				}

				// Check if the parentDir matches or contains the base of worktreeTargetDir
				worktreeBase := filepath.Base(worktreeTargetDir)
				parentBase := filepath.Base(currentParent)
				log.Debug("  → Comparing paths", "parentBase", parentBase, "worktreeBase", worktreeBase)

				if parentBase == worktreeBase || currentParent == worktreeTargetDir {
					log.Debug("  → ✓ Disambiguated by worktreetargetdir match!", "repo", repoName)
					return "repos." + repoName
				}
			}
		}
	}

	// If we still can't disambiguate, return the first match but warn
	log.Warn("  → ⚠ Multiple repos with same bareRepoName, using first match", "bareRepoName", bareRepoName, "matches", matches, "selected", matches[0])
	return "repos." + matches[0]
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

	// Handle the --from flag to select base branch from worktrees
	if c.Flags.From != nil && *c.Flags.From {
		handleFromFlag(c)
	}

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

	from, err := cmd.Flags().GetBool("from")
	if err != nil {
		flags.From = nil
	} else {
		flags.From = &from
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

	repoName := deps.ResolvedRepo

	c.GitInfo.RepoName = repoName

	if c.Flags.BaseBranch != nil {
		c.GitInfo.BaseBranchName = *c.Flags.BaseBranch
	} else {
		baseBranch = viper.GetString("repos." + deps.ResolvedRepo + ".defaultBranch")
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

	c.WorktreeTargetDir = resolveWorktreeTargetDir(repoName, c)
	autoPull := viper.GetBool("repos."+repoName+".autoPull") == true
	if autoPull == true {
		c.AutoPull = true
	}
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
		log.Debug("inputted path", "path", *c.Flags.Directory)
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

// handleFromFlag prompts the user to select a base branch from existing worktrees
// sorted by most recent use (via zoxide scores)
func handleFromFlag(c *com.AddConfig) {
	log.Debug("Handling --from flag to select base branch from worktrees")

	// Get all worktrees
	var rawWorktrees []string
	var err error

	if c.Flags.Directory != nil {
		log.Debug("Using provided directory for worktree list", "path", *c.Flags.Directory)
		rawWorktrees, err = deps.Git.GetWorktrees(c.Flags.Directory)
	} else {
		log.Debug("No directory provided, using current directory")
		rawWorktrees, err = deps.Git.GetWorktrees(nil)
	}

	if err != nil {
		log.Fatal("Error getting worktrees:", err)
	}

	// Transform worktrees to objects
	worktreetransformer := transformer.NewTransformer()
	worktreeObjects := worktreetransformer.TransformWorktrees(rawWorktrees)

	if len(worktreeObjects) == 0 {
		log.Fatal("No worktrees found")
	}

	// Sort worktrees by zoxide score (most recent first)
	sortedWorktrees := sortWorktreesByZoxideScore(worktreeObjects)

	// Create options for selection (branch names)
	var options []string
	for _, wt := range sortedWorktrees {
		options = append(options, wt.BranchName)
	}

	// Present selection interface
	var selectedBranch string
	err = selectBranchForm(&selectedBranch, options)
	util.CheckError(err)

	if selectedBranch == "" {
		log.Fatal("No branch selected")
	}

	log.Info("Selected base branch", "branch", selectedBranch)

	// Set the selected branch as the base branch
	c.Flags.BaseBranch = &selectedBranch
}

// sortWorktreesByZoxideScore sorts worktrees by their zoxide scores (most recent first)
func sortWorktreesByZoxideScore(worktrees []worktreeobj.WorktreeObj) []worktreeobj.WorktreeObj {
	type worktreeWithScore struct {
		worktree worktreeobj.WorktreeObj
		score    float64
	}

	// Get scores for all worktrees
	var worktreesWithScores []worktreeWithScore
	for _, wt := range worktrees {
		score, err := deps.Zoxide.QueryScore(wt.FullPath)
		if err != nil {
			log.Debug("Error querying zoxide score", "path", wt.FullPath, "error", err)
			score = 0
		}
		log.Debug("Zoxide score", "branch", wt.BranchName, "path", wt.FullPath, "score", score)
		worktreesWithScores = append(worktreesWithScores, worktreeWithScore{
			worktree: wt,
			score:    score,
		})
	}

	// Sort by score (highest first)
	slices.SortFunc(worktreesWithScores, func(a, b worktreeWithScore) int {
		if a.score > b.score {
			return -1
		}
		if a.score < b.score {
			return 1
		}
		return 0
	})

	// Extract sorted worktrees
	var sorted []worktreeobj.WorktreeObj
	for _, wts := range worktreesWithScores {
		sorted = append(sorted, wts.worktree)
	}

	return sorted
}

// selectBranchForm presents a selection interface for choosing a branch
func selectBranchForm(selection *string, options []string) error {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select base branch for new worktree:").
				Options(huh.NewOptions(options...)...).
				Value(selection),
		),
	)
	return form.Run()
}
