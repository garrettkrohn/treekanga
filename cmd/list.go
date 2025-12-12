/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/charmbracelet/log"

	"github.com/garrettkrohn/treekanga/transformer"
	util "github.com/garrettkrohn/treekanga/utility"
	worktreeobj "github.com/garrettkrohn/treekanga/worktreeObj"
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
    branch names and directory names.`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose, err := cmd.Flags().GetBool("verbose")
		util.CheckError(err)

		worktrees, err := buildWorktreeStrings(verbose)
		if err != nil {
			log.Fatal(err)
		}
		for _, worktree := range worktrees {
			fmt.Println(worktree)
		}
	},
}

func buildWorktreeStrings(verbose bool) ([]string, error) {
	var rawWorktrees []string
	var err error

	if deps.BareRepoPath != "" {
		log.Debug("Using bare repo path for worktree list", "path", deps.BareRepoPath)
		rawWorktrees, err = deps.Git.GetWorktrees(&deps.BareRepoPath)
	} else {
		log.Debug("No bare repo path set, using current directory")
		rawWorktrees, err = deps.Git.GetWorktrees(nil)
	}

	if err != nil {
		return nil, err
	}

	worktreetransformer := transformer.NewTransformer()
	worktreeObjects := worktreetransformer.TransformWorktrees(rawWorktrees)

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
	if deps.ResolvedRepo != "" {
		displayMode := viper.GetString(deps.ResolvedRepo + ".listDisplayMode")
		if displayMode == "directory" || displayMode == "folder" {
			return "directory"
		}
	}
	return "branch"
}

// getDisplayString returns the appropriate display string based on the configured mode
func getDisplayString(worktree worktreeobj.WorktreeObj, displayMode string) string {
	if displayMode == "directory" {
		return worktree.Folder
	}
	return worktree.BranchName
}

func init() {
	listCmd.Flags().BoolP("verbose", "v", false, "Verbose display of worktrees")
}
