/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	// "os/exec"
	"strconv"
	// "strings"
	//
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/garrettkrohn/treekanga/execwrap"
	"github.com/garrettkrohn/treekanga/filter"
	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/shell"
	worktreeobj "github.com/garrettkrohn/treekanga/worktreeObj"
	"github.com/garrettkrohn/treekanga/worktreeTransformer"
	"github.com/spf13/cobra"
)

// cleanCmd repesents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean unused worktrees",
	Long: `Compare all local worktree branches with the remote branches,
    allow the user to select which worktrees they would like to delete.`,
	Run: func(cmd *cobra.Command, args []string) {

		execWrap := execwrap.NewExec()
		shell := shell.NewShell(execWrap)
		git := git.NewGit(shell)

		branches, error := git.GetRemoteBranches()
		if error != nil {
			log.Fatal(error)
		}

		worktreeStrings, wError := git.GetWorktrees()
		if wError != nil {
			log.Fatal(wError)
		}

		worktreeTransformer := worktreetransformer.NewWorktreeTransformer()
		worktrees := worktreeTransformer.TransformWorktrees(worktreeStrings)

		filter := filter.NewFilter()
		noMatchList := filter.GetBranchNoMatchList(branches, worktrees)
		fmt.Print(noMatchList)

		// transform worktreeobj into strings for selection
		var stringWorktrees []string
		for _, worktreeObj := range noMatchList {
			stringWorktrees = append(stringWorktrees, worktreeObj.BranchName)
		}

		var selections []string

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Value(&selections).
					OptionsFunc(func() []huh.Option[string] {
						return huh.NewOptions(stringWorktrees...)
					}, &stringWorktrees).
					Title("Local endpoints that do not exist on remote").
					Height(25),
			),
		)

		formErr := form.Run()
		if formErr != nil {
			log.Fatal(formErr)
		}

		//transform string selection back to worktreeobjs
		var selectedWorktreeObj []worktreeobj.WorktreeObj
		for _, worktreeobj := range noMatchList {
			for _, str := range selections {
				if worktreeobj.BranchName == str {
					selectedWorktreeObj = append(selectedWorktreeObj, worktreeobj)
					break
				}
			}
		}

		//remove worktrees

		numOfWorktreesRemoved := 0

		action := func() {

			for _, worktreeObj := range selectedWorktreeObj {
				git.RemoveWorktree(worktreeObj.Folder)
				numOfWorktreesRemoved++
			}
		}
		err := spinner.New().
			Title("Removing Worktrees").
			Action(action).
			Run()

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("worktrees removed: %s", strconv.Itoa(numOfWorktreesRemoved))

	},
}

// func getWorktrees() []Worktree {
//
// 	//get all worktrees
// 	cmdToRun := exec.Command("git", "worktree", "list")
// 	allWorktrees, err := cmdToRun.Output()
// 	if err != nil {
// 		log.Fatalf("cmd.Run() failed with %s\n", err)
// 	}
//
// 	//clean worktrees
// 	lines := strings.Split(string(allWorktrees), "\n")
// 	worktrees := make([]Worktree, 0, len(lines))
// 	for _, line := range lines {
// 		parts := strings.SplitN(line, " ", 2)
// 		if len(parts) < 2 {
// 			continue
// 		}
// 		worktrees = append(worktrees, Worktree{Path: parts[0], Head: ExtractTextInBrackets(parts[1])})
// 	}
//
// 	return worktrees
// }
//
// func getCleanRemoteBranchNames() []string {
// 	// fetch
// 	getFetch := exec.Command("git", "fetch", "origin")
// 	getFetch.Run()
//
// 	//get all branches
// 	getAllBranchesCmd := exec.Command("git", "branch", "-r")
// 	allBranches, err := getAllBranchesCmd.Output()
// 	if err != nil {
// 		log.Fatalf("cmd.Run() failed with %s\n", err)
// 	}
//
// 	//clean branch names
// 	branches := strings.Split(string(allBranches), "\n")
// 	var cleanBranches []string
// 	for _, branch := range branches {
// 		cleanBranch := strings.Replace(branch, "origin/", "", 1)
// 		cleanBranches = append(cleanBranches, cleanBranch)
// 	}
// 	return cleanBranches
// }
//
// func branchExistsRemotely(branchName string, remoteName string) (bool, error) {
// 	cmd := exec.Command("git", "ls-remote", "--heads", remoteName, branchName)
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return false, err
// 	}
//
// 	// If the output is empty, the branch does not exist remotely
// 	return strings.TrimSpace(string(output)) != "", nil
// }

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
