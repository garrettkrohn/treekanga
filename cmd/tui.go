/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package cmd

import (
	// "github.com/charmbracelet/log"
	// tui "github.com/garrettkrohn/treekanga/tui"
	"github.com/spf13/cobra"
)

// tuiCmd represents the tui command
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// path := "*"

		// worktrees, err := deps.Git.GetWorktrees(&path)
		// if err != nil {
		// 	log.Info("error", err)
		// }

		// tui.Main(worktrees)
	},
}

func init() {
	// Add flags here if needed
}
