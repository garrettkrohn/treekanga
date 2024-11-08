/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/garrettkrohn/treekanga/directoryReader"
	"github.com/garrettkrohn/treekanga/execwrap"
	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/shell"
	"github.com/garrettkrohn/treekanga/zoxide"
	"github.com/spf13/cobra"
)

type Dependencies struct {
	Git             git.Git
	Zoxide          zoxide.Zoxide
	DirectoryReader directoryReader.DirectoryReader
}

var deps Dependencies

func NewRootCmd(git git.Git, zoxide zoxide.Zoxide, directoryReader directoryReader.DirectoryReader, version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "treekanga",
		Short:   "CLI application to manage git worktree",
		Long:    `CLI application to manage git worktree`,
		Version: version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			deps = Dependencies{
				Git:             git,
				Zoxide:          zoxide,
				DirectoryReader: directoryReader,
			}
		},
	}

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {

	execWrap := execwrap.NewExec()
	shell := shell.NewShell(execWrap)
	git := git.NewGit(shell)
	zoxide := zoxide.NewZoxide(shell)
	directoryReader := directoryReader.NewDirectoryReader()

	rootCmd := NewRootCmd(git, zoxide, directoryReader, version)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(cloneCmd)
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
