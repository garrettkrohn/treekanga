/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"github.com/charmbracelet/log"

	"github.com/garrettkrohn/treekanga/models"
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

		expand, err := cmd.Flags().GetBool("expand")
		utilpkg.CheckError(err)

		global, err := cmd.Flags().GetBool("global")
		utilpkg.CheckError(err)

		worktrees, err := buildWorktreeStrings(verbose, global, expand)
		if err != nil {
			log.Fatal(err)
		}
		for _, worktree := range worktrees {
			fmt.Println(worktree)
		}
	},
}

func buildWorktreeStrings(verbose, global, expand bool) ([]string, error) {
	// get fetcher
	fetcher := getFetcher(global)
	rawWorktrees, err := fetcher.fetch()
	log.Debug("rawWorktrees", rawWorktrees)
	utilpkg.CheckError(err)

	// get lister
	lister := getLister(verbose, global, expand)
	worktreeStrings, err := lister.list()
	utilpkg.CheckError(err)

	log.Debug(worktreeStrings)
	return worktreeStrings, nil
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
	listCmd.Flags().BoolP("expand", "e", false, "Expand the root with all defined sub folders")
	listCmd.Flags().BoolP("global", "g", false, "Show all worktrees for every repo in the config file")
}
