/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/charmbracelet/log"
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

    By default, creates a new branch off of the defaultBranch defined in
    the config, or you can specify a base branch with the -b flag.

    Use --remote or --local to explicitly checkout an existing branch.`,
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

		tmux, err := cmd.Flags().GetString("tmux")
		util.CheckError(err)
		if tmux != "" {
			log.Debug(fmt.Sprintf("set TmuxConnect = %s from flags", tmux))
			deps.AppConfig.TmuxConnect = tmux
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

		remote, err := cmd.Flags().GetBool("remote")
		util.CheckError(err)
		if remote {
			log.Debug("set CheckoutRemote = true from flags")
			deps.AppConfig.CheckoutRemote = true
		}

		local, err := cmd.Flags().GetBool("local")
		util.CheckError(err)
		if local {
			log.Debug("set CheckoutLocal = true from flags")
			deps.AppConfig.CheckoutLocal = true
		}

		cfg := services.SetConfigForAddService(deps.AppConfig, args)

		services.AddWorktree(deps.Connector, deps.Shell, cfg)
	},
}

func init() {

	addCmd.Flags().BoolP("pull", "p", false, "Pull the base branch before creating new branch")
	addCmd.Flags().BoolP("cursor", "c", false, "Open up new worktree in cursor")
	addCmd.Flags().BoolP("vscode", "v", false, "Open up new worktree in vs code")
	addCmd.Flags().BoolP("script", "x", false, "Execute Custom Script")
	addCmd.Flags().BoolP("from", "f", false, "Select base branch from list of branches")
	addCmd.Flags().BoolP("remote", "r", false, "Checkout existing branch from remote")
	addCmd.Flags().BoolP("local", "L", false, "Checkout existing branch from local repository")
	addCmd.Flags().StringP("tmux", "t", "", "Connect to tmux session at subdirectory (use '.' for root)")
	addCmd.Flags().StringP("base", "b", "", "Specify the base branch for the new worktree")
	addCmd.Flags().StringP("directory", "d", "", "Specify the directory to the bare repo where the worktree will be added")
	addCmd.Flags().StringP("name", "n", "", "Specify a worktree name")
}
