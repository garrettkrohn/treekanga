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
	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/models"
	"github.com/garrettkrohn/treekanga/transformer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:     "connect [session-name]",
	Aliases: []string{"cn"},
	Short:   "Connect to a tmux session",
	Long: `Connect to a tmux session by name, worktree path, or directory path.

The connect command will try to find a session using the following strategies:
1. Existing tmux session with the given name
2. Worktree matching the given name or path
3. Directory path (absolute or relative)

If a session doesn't exist, it will be created automatically.

Examples:
  # Connect to an existing tmux session
  treekanga connect my-session

  # Connect to a worktree by name
  treekanga connect feature-branch

  # Connect to a directory
  treekanga connect ~/code/myproject

  # Switch to a session (when already in tmux)
  treekanga connect my-session --switch`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		selectFlag, err := cmd.Flags().GetBool("select")
		if err != nil {
			log.Fatal(err)
			return
		}

		// Handle --select mode
		if selectFlag {
			handleSelectMode(cmd)
			return
		}

		// Existing direct connect logic
		if len(args) == 0 {
			log.Fatal("please provide a session name or path")
			return
		}

		name := strings.Join(args, " ")
		if name == "" {
			return
		}

		switchFlag, err := cmd.Flags().GetBool("switch")
		if err != nil {
			log.Fatal(err)
			return
		}

		opts := models.ConnectOpts{
			Switch: switchFlag,
		}

		log.Debug("Attempting to connect", "name", name, "switch", switchFlag)

		if err := deps.Connector.Connect(name, opts); err != nil {
			log.Fatal(err)
			return
		}

		log.Info("Connected successfully", "session", name)
	},
}

func init() {
	connectCmd.Flags().BoolP("switch", "s", false, "Switch to the session (rather than attach). Useful when already inside tmux.")
	connectCmd.Flags().Bool("select", false, "Interactive selection mode")
	connectCmd.Flags().Bool("bare", false, "Select from bare repos (use with --select)")
	connectCmd.Flags().Bool("by-repo", false, "Select repo first, then worktree (use with --select)")
}

// listBareRepos returns a list of bare repo paths from the config that exist on disk
func listBareRepos() ([]string, error) {
	// Initialize with empty slice to avoid nil return
	bareRepos := []string{}

	for _, worktreeDir := range deps.AppConfig.AllBareRepoPaths {
		bareRepoPath := filepath.Join(worktreeDir, ".bare")

		// Check if the bare repo path exists
		if _, err := os.Stat(bareRepoPath); err == nil {
			bareRepos = append(bareRepos, bareRepoPath)
		} else {
			log.Debug("Skipping non-existent bare repo", "path", bareRepoPath)
		}
	}

	return bareRepos, nil
}

// formatBareRepoForDisplay formats a bare repo path for user-friendly display
func formatBareRepoForDisplay(bareRepoPath string) string {
	// Extract repo name from path
	// Example: /Users/gkrohn/code/cal_work/.bare -> "cal"
	parentDir := filepath.Dir(bareRepoPath)
	repoName := filepath.Base(parentDir)

	// Strip _work suffix if present
	repoName = strings.TrimSuffix(repoName, "_work")

	return fmt.Sprintf("%s -> %s", repoName, bareRepoPath)
}

// listAllWorktrees returns all worktrees from all repos using globalFetcher
func listAllWorktrees() ([]models.Worktree, error) {
	fetcher := &globalFetcher{}
	worktrees, err := fetcher.fetch()
	if err != nil {
		return nil, err
	}

	// Ensure we return an empty slice instead of nil
	if worktrees == nil {
		worktrees = []models.Worktree{}
	}

	// Already sorted by mod time in globalFetcher
	return worktrees, nil
}

// listReposForSelection returns a list of repo names from the config
func listReposForSelection() ([]string, error) {
	repoNames := []string{}

	// Get repos map from viper config
	repoconfig := viper.GetStringMap("repos")
	if repoconfig == nil {
		return repoNames, nil
	}

	// Extract repo names (keys from the map)
	for repoName := range repoconfig {
		repoNames = append(repoNames, repoName)
	}

	return repoNames, nil
}

// listWorktreesForRepo returns worktrees for a specific repo
func listWorktreesForRepo(repoName string) ([]models.Worktree, error) {
	// Get the worktreeTargetDir for this repo
	worktreeTargetDir := viper.GetString("repos." + repoName + ".worktreeTargetDir")
	if worktreeTargetDir == "" {
		return nil, fmt.Errorf("repo %s not found in config", repoName)
	}

	// Expand tilde to home directory
	homeDir, err := os.UserHomeDir()
	if err == nil {
		worktreeTargetDir = filepath.Join(homeDir, strings.TrimPrefix(worktreeTargetDir, "~/"))
	}

	// Find the bare repo path
	bareRepoPath, err := findBareRepoFromWorktreeDir(worktreeTargetDir)
	if err != nil {
		return nil, fmt.Errorf("could not find bare repo for %s: %w", repoName, err)
	}

	// List worktrees
	rawWorktrees, err := git.ListWorktrees(bareRepoPath)
	if err != nil {
		return nil, fmt.Errorf("could not list worktrees for %s: %w", repoName, err)
	}

	// Transform to models.Worktree
	worktreeObjects := transformer.TransformWorktrees(rawWorktrees)

	return worktreeObjects, nil
}

// formatWorktreeForDisplay formats a worktree for user-friendly display
func formatWorktreeForDisplay(wt models.Worktree) string {
	// Extract repo name from the full path
	// Example: /Users/gkrohn/code/cal_work/feature-branch -> "cal"
	parentDir := filepath.Dir(wt.FullPath)
	repoName := filepath.Base(parentDir)

	// Strip _work suffix if present
	repoName = strings.TrimSuffix(repoName, "_work")

	return fmt.Sprintf("%s - %s", repoName, wt.BranchName)
}

// formatRepoForDisplay formats a repo name with worktree count
func formatRepoForDisplay(repoName string, count int) string {
	return fmt.Sprintf("%s (%d worktrees)", repoName, count)
}

// handleSelectMode handles the --select flag logic
func handleSelectMode(cmd *cobra.Command) {
	bareFlag, _ := cmd.Flags().GetBool("bare")
	byRepoFlag, _ := cmd.Flags().GetBool("by-repo")
	switchFlag, _ := cmd.Flags().GetBool("switch")

	selector := getSelector(deps.AppConfig, deps.Shell)
	opts := models.ConnectOpts{Switch: switchFlag}

	if bareFlag {
		handleBareRepoSelection(selector, opts)
	} else if byRepoFlag {
		handleHierarchicalSelection(selector, opts)
	} else {
		handleFlatSelection(selector, opts)
	}
}

// handleBareRepoSelection handles bare repo selection mode
func handleBareRepoSelection(selector Selector, opts models.ConnectOpts) {
	bareRepos, err := listBareRepos()
	if err != nil {
		log.Fatal("Error listing bare repos", "error", err)
		return
	}

	if len(bareRepos) == 0 {
		log.Fatal("No bare repos found in config")
		return
	}

	// Format for display
	displayItems := make([]string, len(bareRepos))
	for i, repo := range bareRepos {
		displayItems[i] = formatBareRepoForDisplay(repo)
	}

	// Show selector
	selected, err := selector.Select(displayItems, "Select bare repo: ")
	if err != nil {
		log.Fatal("Selection cancelled", "error", err)
		return
	}

	// Parse selection to extract path (format: "repo -> /path/.bare")
	parts := strings.Split(selected, " -> ")
	if len(parts) != 2 {
		log.Fatal("Invalid selection format")
		return
	}
	bareRepoPath := parts[1]

	// Connect to bare repo
	if err := deps.Connector.Connect(bareRepoPath, opts); err != nil {
		log.Fatal("Failed to connect", "error", err)
		return
	}

	log.Info("Connected successfully to bare repo", "path", bareRepoPath)
}

// handleHierarchicalSelection handles hierarchical selection (repo first, then worktree)
func handleHierarchicalSelection(selector Selector, opts models.ConnectOpts) {
	// Step 1: Select repo
	repos, err := listReposForSelection()
	if err != nil {
		log.Fatal("Error listing repos", "error", err)
		return
	}

	if len(repos) == 0 {
		log.Fatal("No repos found in config")
		return
	}

	// Get worktree counts for display
	displayItems := make([]string, len(repos))
	for i, repoName := range repos {
		worktrees, _ := listWorktreesForRepo(repoName)
		displayItems[i] = formatRepoForDisplay(repoName, len(worktrees))
	}

	selectedRepo, err := selector.Select(displayItems, "Select repo: ")
	if err != nil {
		log.Fatal("Selection cancelled", "error", err)
		return
	}

	// Parse repo name from selection (format: "repo-name (X worktrees)")
	repoName := strings.Split(selectedRepo, " (")[0]

	// Step 2: Select worktree in that repo
	worktrees, err := listWorktreesForRepo(repoName)
	if err != nil {
		log.Fatal("Error listing worktrees", "error", err, "repo", repoName)
		return
	}

	if len(worktrees) == 0 {
		log.Fatal("No worktrees found for repo", "repo", repoName)
		return
	}

	// Format for display
	displayItems = make([]string, len(worktrees))
	for i, wt := range worktrees {
		displayItems[i] = formatWorktreeForDisplay(wt)
	}

	selectedWorktree, err := selector.Select(displayItems, "Select worktree: ")
	if err != nil {
		log.Fatal("Selection cancelled", "error", err)
		return
	}

	// Find the worktree path from the selection
	var selectedPath string
	for _, wt := range worktrees {
		if formatWorktreeForDisplay(wt) == selectedWorktree {
			selectedPath = wt.FullPath
			break
		}
	}

	if selectedPath == "" {
		log.Fatal("Could not find selected worktree")
		return
	}

	// Connect to selected worktree
	if err := deps.Connector.Connect(selectedPath, opts); err != nil {
		log.Fatal("Failed to connect", "error", err)
		return
	}

	log.Info("Connected successfully", "worktree", selectedPath)
}

// handleFlatSelection handles flat selection (all worktrees from all repos)
func handleFlatSelection(selector Selector, opts models.ConnectOpts) {
	worktrees, err := listAllWorktrees()
	if err != nil {
		log.Fatal("Error listing worktrees", "error", err)
		return
	}

	if len(worktrees) == 0 {
		log.Fatal("No worktrees found")
		return
	}

	// Format for display
	displayItems := make([]string, len(worktrees))
	for i, wt := range worktrees {
		displayItems[i] = formatWorktreeForDisplay(wt)
	}

	// Show selector
	selected, err := selector.Select(displayItems, "Select worktree: ")
	if err != nil {
		log.Fatal("Selection cancelled", "error", err)
		return
	}

	// Find the worktree path from the selection
	var selectedPath string
	for _, wt := range worktrees {
		if formatWorktreeForDisplay(wt) == selected {
			selectedPath = wt.FullPath
			break
		}
	}

	if selectedPath == "" {
		log.Fatal("Could not find selected worktree")
		return
	}

	// Connect to selected worktree
	if err := deps.Connector.Connect(selectedPath, opts); err != nil {
		log.Fatal("Failed to connect", "error", err)
		return
	}

	log.Info("Connected successfully", "worktree", selectedPath)
}
