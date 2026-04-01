package cmd

import (
	"github.com/garrettkrohn/treekanga/confirmer"
	"github.com/garrettkrohn/treekanga/services"
	"github.com/spf13/cobra"
)

var renameCmd = &cobra.Command{
	Use:   "rename <new-branch-name>",
	Short: "Rename the current worktree and branch",
	Long: `Rename the current worktree's branch and folder structure.

    This command renames both the git branch and the worktree folder.
    Branch names can contain slashes (e.g., feature/new-feature), which
    will be converted to dashes in the folder name (feature-new-feature).

    Example usage:
      treekanga rename feature/new-feature
      treekanga rename bugfix/issue-123
      treekanga rename feature/new-feature -s  # auto-switch tmux session

    Flags:
    -s, --switch: Automatically switch to new tmux session (skip prompt)

    Important notes:
    - Only works from within a worktree (not from the bare repository)
    - The new branch name must not already exist locally or remotely
    - After rename, you'll need to cd to the new folder path
    - Your shell will be in an invalid directory after the rename`,
	Run: func(cmd *cobra.Command, args []string) {
		autoSwitch, err := cmd.Flags().GetBool("switch")
		if err != nil {
			cmd.PrintErrln("Error:", err)
			return
		}

		err = services.ExecuteRename(
			deps.AppConfig,
			args,
			deps.Connector,
			confirmer.NewConfirmer(),
			autoSwitch,
		)
		if err != nil {
			cmd.PrintErrln("Error:", err)
			return
		}
	},
}

func init() {
	renameCmd.Flags().BoolP("switch", "s", false, "Automatically switch to new tmux session without prompting")
}
