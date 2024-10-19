/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"log"
	// "regexp"
	// "strings"

	"github.com/garrettkrohn/treekanga/execwrap"
	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/shell"
	"github.com/garrettkrohn/treekanga/transformer"
)

type Worktree struct {
	Path string
	Head string
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		shell := shell.NewShell(execwrap.NewExec())
		git := git.NewGit(shell)

		rawWorktrees, err := git.GetWorktrees()

		if err != nil {
			log.Fatal(err)
		}

		worktreetransformer := transformer.NewTransformer()
		worktreeObjects := worktreetransformer.TransformWorktrees(rawWorktrees)

		for _, worktree := range worktreeObjects {
			fmt.Println(worktree.BranchName)
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
