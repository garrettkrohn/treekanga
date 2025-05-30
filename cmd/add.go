/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/charmbracelet/log"
	com "github.com/garrettkrohn/treekanga/common"
	util "github.com/garrettkrohn/treekanga/utility"

	"github.com/spf13/cobra"
)

var (
	newBranchName string
	baseBranch    string
)

// const tempZoxideName = "temp_treekanga_worktree"

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a git worktree",
	Long: `You may use this command with zero arguments, and you
    will be prompeted to input the branch name and base branch.

    Alternatively, you may the branch name as an argument, 
    treekanga will create this branch off of the defaultBranch 
    defined in the config, or use the current branch.`,
	Run: func(cmd *cobra.Command, args []string) {

		c := com.AddConfig{}
		getAddCmdConfig(cmd, args, &c)

		validateConfig(&c)

		log.Debug("Adding worktree with config:")
		PrintConfig(c)
		err := deps.Git.AddWorktree(&c)
		util.CheckError(err)

		log.Info(fmt.Sprintf("worktree %s created", c.GetNewBranchName()))

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
	},
}

func init() {

	addCmd.Flags().BoolP("pull", "p", false, "Pull the base branch before creating new branch")
	addCmd.Flags().BoolP("cursor", "c", false, "Open up new worktree in cursor")
	addCmd.Flags().BoolP("vscode", "v", false, "Open up new worktree in vs code")
	addCmd.Flags().StringP("sesh", "s", "", "Automatically connect to a sesh upon creation")
	addCmd.Flags().StringP("base", "b", "", "Specify the base branch for the new worktree")
	addCmd.Flags().StringP("directory", "d", "", "Specify the directory to the bare repo where the worktree will be added")
}
