/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/charmbracelet/log"
	com "github.com/garrettkrohn/treekanga/common"
	"github.com/garrettkrohn/treekanga/services"
	util "github.com/garrettkrohn/treekanga/utility"

	"github.com/spf13/cobra"
)

var (
	baseBranch string
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a git worktree",
	Long: `Create a new worktree with the specified branch name.

    The branch name is required as an argument. Treekanga will create 
    this branch off of the defaultBranch defined in the config, or 
    you can specify a base branch with the -b flag.

    Available flags:
    -b, --base: Specify the base branch for the new worktree
    -f, --from: Select base branch from list of existing worktrees (sorted by recent use)
    -p, --pull: Pull the base branch before creating new branch
    -c, --cursor: Open the new worktree in Cursor
    -v, --vscode: Open the new worktree in VS Code
    -s, --sesh: Connect to a sesh session after creation
    -d, --directory: Specify the directory to the bare repo`,
	Run: func(cmd *cobra.Command, args []string) {

		directory, err := cmd.Flags().GetString("directory")
		util.CheckError(err)
		if directory != "" {
			log.Debug(fmt.Sprintf("set Directory = %s by flags", directory))
			deps.AppConfig.WorktreeTargetDir = directory
		}

		baseBranch, err := cmd.Flags().GetString("base")
		util.CheckError(err)
		if baseBranch != "" {
			log.Debug(fmt.Sprintf("set baseBranch = %s by flags", baseBranch))
			deps.AppConfig.BaseBranch = baseBranch
		}

		sesh, err := cmd.Flags().GetString("sesh")
		util.CheckError(err)
		if sesh != "" {
			log.Debug("set SeshConnect = true from flags")
			deps.AppConfig.SeshConnect = sesh
		}

		pull, err := cmd.Flags().GetBool("pull")
		util.CheckError(err)
		if pull {
			log.Debug("set PullBeforeCuttingNewBranch = true from flags")
			deps.AppConfig.PullBeforeCuttingNewBranch = true
		}

		cursor, err := cmd.Flags().GetBool("cursor")
		util.CheckError(err)
		if cursor {
			log.Debug("set CursorConnect = true from flags")
			deps.AppConfig.CursorConnect = true
		}

		vscode, err := cmd.Flags().GetBool("vscode")
		util.CheckError(err)
		if vscode {
			log.Debug("set VsCodeConnect = true from flags")
			deps.AppConfig.VsCodeConnect = true
		}

		specifiedWorktreeName, err := cmd.Flags().GetString("name")
		util.CheckError(err)
		if specifiedWorktreeName != "" {
			deps.AppConfig.NewWorktreeName = specifiedWorktreeName
		}

		executeScript, err := cmd.Flags().GetBool("script")
		util.CheckError(err)
		if executeScript {
			log.Debug("set RunPostScript = true from flags")
			deps.AppConfig.RunPostScript = true
		}

		from, err := cmd.Flags().GetBool("from")
		util.CheckError(err)
		if from {
			log.Debug("set UseFormToSetBaseBranch = true from flags")
			deps.AppConfig.UseFormToSetBaseBranch = true
		}

		cfg := services.SetConfigForAddService(deps.Git, deps.AppConfig, args)

		services.AddWorktree(deps.Git, deps.Zoxide, deps.Connector, deps.Shell, cfg)
	},
}

func runPostScript(c com.AddConfig) bool {
	if !c.HasPostScript() {
		return false
	}

	if *c.Flags.ExecuteScript || *c.AutoRunPostScript {
		return true
	}

	return false
}

func init() {

	addCmd.Flags().BoolP("pull", "p", false, "Pull the base branch before creating new branch")
	addCmd.Flags().BoolP("cursor", "c", false, "Open up new worktree in cursor")
	addCmd.Flags().BoolP("vscode", "v", false, "Open up new worktree in vs code")
	addCmd.Flags().BoolP("script", "x", false, "Execute Custom Script")
	addCmd.Flags().BoolP("from", "f", false, "Select base branch from list of branches")
	addCmd.Flags().StringP("sesh", "s", "", "Automatically connect to a sesh upon creation")
	addCmd.Flags().StringP("base", "b", "", "Specify the base branch for the new worktree")
	addCmd.Flags().StringP("directory", "d", "", "Specify the directory to the bare repo where the worktree will be added")
	addCmd.Flags().StringP("name", "n", "", "Specify a worktree name")
}
