/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package cmd

import (
	"strings"

	"github.com/charmbracelet/log"
	"github.com/garrettkrohn/treekanga/models"
	"github.com/spf13/cobra"
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
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
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
}
