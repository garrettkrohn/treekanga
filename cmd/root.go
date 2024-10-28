/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/garrettkrohn/treekanga/execwrap"
	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/shell"
	"github.com/garrettkrohn/treekanga/zoxide"
	"github.com/spf13/cobra"
)

type Dependencies struct {
	Git    git.Git
	Zoxide zoxide.Zoxide
}

var deps Dependencies

func NewRootCmd(git git.Git, zoxide zoxide.Zoxide) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "treekanga",
		Short:   "CLI application to manage git worktree",
		Long:    `CLI application to manage git worktree`,
		Version: `v0.1.0-beta`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			deps = Dependencies{
				Git:    git,
				Zoxide: zoxide,
			}
		},
	}

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	execWrap := execwrap.NewExec()
	shell := shell.NewShell(execWrap)
	git := git.NewGit(shell)
	zoxide := zoxide.NewZoxide(shell)

	rootCmd := NewRootCmd(git, zoxide)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(deleteCmd)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.treekanga.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
