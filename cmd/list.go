/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os/exec"
	"regexp"
	"strings"
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
		cmdToRun := exec.Command("git", "worktree", "list")
		output, err := cmdToRun.Output()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}

		lines := strings.Split(string(output), "\n")
		worktrees := make([]Worktree, 0, len(lines))
		for _, line := range lines {
			parts := strings.SplitN(line, " ", 2)
			if len(parts) < 2 {
				continue
			}
			worktrees = append(worktrees, Worktree{Path: parts[0], Head: parts[1]})
		}

		// Now worktrees contains the results of the command
		for _, wt := range worktrees {
			// fmt.Printf("Path: %s, Head: %s\n", wt.Path, wt.Head)
			splitPath := strings.Split(wt.Path, "/")
			fmt.Printf("Folder: %s\n", splitPath[5])

			branch := ExtractTextInBrackets(wt.Head)
			fmt.Printf("Branch: %s\n", branch)
		}
	},
}

func ExtractTextInBrackets(s string) string {
	re := regexp.MustCompile(`\[(.*?)\]`)
	matches := re.FindAllStringSubmatch(s, -1)
	var results []string
	for _, match := range matches {
		results = append(results, match[1])
	}
	return strings.Join(results, ", ")
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
