/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/charmbracelet/log"

	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/models"
	"github.com/garrettkrohn/treekanga/services"
	"github.com/garrettkrohn/treekanga/util"
	utilpkg "github.com/garrettkrohn/treekanga/utility"
)

type Worktree struct {
	Path string
	Head string
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all worktrees",
	Long: `Display all worktrees in the current repository.

    By default, shows the branch name for each worktree.
    You can configure the display mode using the 'listDisplayMode' 
    configuration option:
      - "branch" (default): Display branch names
      - "directory" or "folder": Display directory names
    
    Configuration example:
      repos:
        myrepo:
          listDisplayMode: directory
    
    Use the -v/--verbose flag to show all details including both
    branch names and directory names.
    
    Use the -a/--all flag to show all worktrees plus subdirectories
    defined in the zoxideFolders configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose, err := cmd.Flags().GetBool("verbose")
		utilpkg.CheckError(err)

		all, err := cmd.Flags().GetBool("all")
		utilpkg.CheckError(err)

		worktrees, err := buildWorktreeStrings(verbose, all)
		if err != nil {
			log.Fatal(err)
		}
		for _, worktree := range worktrees {
			fmt.Println(worktree)
		}
	},
}

func buildWorktreeStrings(verbose bool, all bool) ([]string, error) {
	rawWorktrees, err := git.ListWorktrees(deps.AppConfig.BareRepoPath)
	if err != nil {
		return nil, err
	}

	worktreeObjects := util.ParseWorktrees(rawWorktrees)

	// Sort worktrees by most recently modified
	sortWorktreesByModTime(worktreeObjects)

	// If --all flag is set, expand with zoxide folders
	if all {
		log.Debug("Expanding worktrees with zoxide folders", "zoxideFolders", deps.AppConfig.ZoxideFolders)
		allPaths := services.ExpandWorktreesWithZoxideFolders(worktreeObjects, deps.AppConfig.ZoxideFolders, deps.DirectoryReader)
		return allPaths, nil
	}

	// Get the display mode from config (default to "branch" for backward compatibility)
	displayMode := getListDisplayMode()
	log.Debug("List display mode", "mode", displayMode)

	var worktreeBranches []string
	for _, worktree := range worktreeObjects {
		var branchDisplay string
		if verbose {
			branchDisplay = fmt.Sprintf("worktree: %s, branch: %s, fullPath: %s, commitHash: %s", worktree.Folder, worktree.BranchName, worktree.FullPath, worktree.CommitHash)
		} else {
			branchDisplay = getDisplayString(worktree, displayMode)
		}
		worktreeBranches = append(worktreeBranches, branchDisplay)
	}

	return worktreeBranches, nil
}

// getListDisplayMode retrieves the configured display mode for list command
// Returns "branch" or "directory" (default: "branch")
func getListDisplayMode() string {
	if deps.AppConfig.RepoNameForConfig != "" {
		displayMode := viper.GetString("repos." + deps.AppConfig.RepoNameForConfig + ".listDisplayMode")
		if displayMode == "directory" || displayMode == "folder" {
			return "directory"
		}
	}
	return "branch"
}

// getDisplayString returns the appropriate display string based on the configured mode
func getDisplayString(worktree models.Worktree, displayMode string) string {
	if displayMode == "directory" {
		return worktree.Folder
	}
	return worktree.BranchName
}

// sortWorktreesByModTime sorts worktrees by modification time (most recent first)
func sortWorktreesByModTime(worktrees []models.Worktree) {
	sort.Slice(worktrees, func(i, j int) bool {
		statI, errI := os.Stat(worktrees[i].FullPath)
		statJ, errJ := os.Stat(worktrees[j].FullPath)

		// If there's an error accessing either path, push it to the end
		if errI != nil {
			log.Debug("Error stat'ing worktree", "path", worktrees[i].FullPath, "error", errI)
			return false
		}
		if errJ != nil {
			log.Debug("Error stat'ing worktree", "path", worktrees[j].FullPath, "error", errJ)
			return true
		}

		// Sort by modification time, most recent first
		return statI.ModTime().After(statJ.ModTime())
	})
}

func init() {
	listCmd.Flags().BoolP("verbose", "v", false, "Verbose display of worktrees")
	listCmd.Flags().BoolP("all", "a", false, "Show all worktrees plus subdirectories from zoxideFolders config")
}
