/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	// "fmt"
	// "github.com/charmbracelet/huh"
	// "log"
	// "os/exec"
	// "strings"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// cleanBranches := getCleanRemoteBranchNames()
		//
		// // set up branches map
		// branchesMap := make(map[string]bool)
		// for _, branch := range cleanBranches {
		// 	branchesMap[strings.TrimSpace(branch)] = true
		// }
		//
		// var branchName string
		// form := huh.NewForm(
		// 	huh.NewGroup(
		// 		huh.NewInput().
		// 			Title("Input branch name").
		// 			Prompt("?").
		// 			Value(&branchName),
		// 	),
		// )
		// formErr := form.Run()
		// if formErr != nil {
		// 	log.Fatal(formErr)
		// }
		//
		// branchExistsRemotely := !branchesMap[branchName]
		//
		// newBranch := "../" + branchName
		//
		// //TODO: need to check if the worktree exists already
		//
		// if branchExistsRemotely {
		// 	fmt.Print("doesn't exist")
		// 	cmdToRun := exec.Command("git", "worktree", "add", newBranch, "-b", branchName)
		// 	fmt.Print(cmdToRun)
		// 	_, err := cmdToRun.Output()
		// 	if err != nil {
		// 		log.Fatalf("cmd.Run() failed with %s\n", err)
		// 	}
		// } else {
		// 	fmt.Print("does exist")
		// 	cmdToRun := exec.Command("git", "worktree", "add", newBranch, branchName)
		// 	fmt.Print(cmdToRun)
		// 	_, err := cmdToRun.Output()
		// 	if err != nil {
		// 		log.Fatalf("cmd.Run() failed with %s\n", err)
		// 	}
		// }

	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
