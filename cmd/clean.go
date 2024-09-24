/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	// "github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"log"
	"os/exec"
	"strings"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		getFetch := exec.Command("git", "fetch", "origin")
		getFetch.Run()

		getAllBranchesCmd := exec.Command("git", "branch", "-r")
		allBranches, err := getAllBranchesCmd.Output()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}

		branches := strings.Split(string(allBranches), "\n")
		for _, branch := range branches {
			cleanBranch := strings.Replace(branch, "origin/", "", 1)
			fmt.Println(cleanBranch)
		}

		cmdToRun := exec.Command("git", "worktree", "list")
		allWorktrees, err := cmdToRun.Output()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}

		lines := strings.Split(string(allWorktrees), "\n")
		worktrees := make([]Worktree, 0, len(lines))
		for _, line := range lines {
			parts := strings.SplitN(line, " ", 2)
			if len(parts) < 2 {
				continue
			}
			worktrees = append(worktrees, Worktree{Path: parts[0], Head: parts[1]})
		}

		branchesMap := make(map[string]bool)
		for _, branch := range branches {
			branchesMap[branch] = true
		}

		var match []Worktree
		var noMatch []Worktree
		for _, worktree := range worktrees {

			branch := ExtractTextInBrackets(worktree.Head)
			if branchesMap[branch] {
				match = append(match, worktree)
			} else {
				noMatch = append(noMatch, worktree)
			}
		}

		for _, tree := range noMatch {
			fmt.Printf(tree.Head)
		}
		//
		// for _, wt := range worktrees {
		// 	// fmt.Printf("Path: %s, Head: %s\n", wt.Path, wt.Head)
		// 	// splitPath := strings.Split(wt.Path, "/")
		// 	// fmt.Printf("Folder: %s\n", splitPath[5])
		//
		// 	branch := ExtractTextInBrackets(wt.Head)
		// 	// fmt.Printf("Branch: %s\n", branch)
		//
		// 	exists, err := branchExistsRemotely(branch, "origin")
		// 	if err != nil {
		// 		// handle error
		// 	}
		// 	confirm := false
		// 	huh.NewConfirm().
		// 		Title("Do you want to delete %, branch").
		// 		Affirmative("Yes").
		// 		Negative("No").
		// 		Value(&confirm)
		//
		// 	if exists {
		// 		fmt.Println("Branch exists remotely\n")
		// 	} else {
		// 		fmt.Println("Branch does not exist remotely\n")
		// 	}
		//
		// }
	},
}

func branchExistsRemotely(branchName string, remoteName string) (bool, error) {
	cmd := exec.Command("git", "ls-remote", "--heads", remoteName, branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}

	// If the output is empty, the branch does not exist remotely
	return strings.TrimSpace(string(output)) != "", nil
}
func init() {
	rootCmd.AddCommand(cleanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
