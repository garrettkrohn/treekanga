/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package cmd

import (
	"github.com/charmbracelet/log"
	com "github.com/garrettkrohn/treekanga/common"
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
    -p, --pull: Pull the base branch before creating new branch
    -c, --cursor: Open the new worktree in Cursor
    -v, --vscode: Open the new worktree in VS Code
    -s, --sesh: Connect to a sesh session after creation
    -d, --directory: Specify the directory to the bare repo`,
	Run: func(cmd *cobra.Command, args []string) {

		c := com.AddConfig{}
		getAddCmdConfig(cmd, args, &c)

		validateConfig(&c)

		log.Debug("Adding worktree with config:")
		PrintConfig(c)
		err := deps.Git.AddWorktree(&c)
		util.CheckError(err)

		log.Info("worktree created", "branch", c.GetNewBranchName())

		deps.Zoxide.AddZoxideEntries(&c)

		if c.HasSeshTarget() {
			deps.Connector.SeshConnect(&c)
		}

		if c.ShouldOpenCursor() {
			deps.Connector.CursorConnect(&c)
		}

		if c.ShouldOpenVSCode() {
			deps.Connector.VsCodeConnect(&c)
		}

		if c.HasPostScript() && *c.Flags.ExecuteScript {
			script := c.GetPostScript()
			deps.Shell.CmdWithDir(c.WorktreeTargetDir, "sh", "-c", script)
			log.Info("post script run", "command", script)
		}
	},
}

func init() {

	addCmd.Flags().BoolP("pull", "p", false, "Pull the base branch before creating new branch")
	addCmd.Flags().BoolP("cursor", "c", false, "Open up new worktree in cursor")
	addCmd.Flags().BoolP("vscode", "v", false, "Open up new worktree in vs code")
	addCmd.Flags().BoolP("script", "x", false, "Execute Custom Script")
	addCmd.Flags().StringP("sesh", "s", "", "Automatically connect to a sesh upon creation")
	addCmd.Flags().StringP("base", "b", "", "Specify the base branch for the new worktree")
	addCmd.Flags().StringP("directory", "d", "", "Specify the directory to the bare repo where the worktree will be added")
	addCmd.Flags().StringP("name", "n", "", "Specify a worktree name")
}
