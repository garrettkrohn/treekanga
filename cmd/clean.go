/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/garrettkrohn/treekanga/execwrap"
	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/shell"
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
		fmt.Print(branches)
		fmt.Print(error)

		cleanBranches := getCleanRemoteBranchNames()

		worktrees := getWorktrees()

		// set up branches map
		branchesMap := make(map[string]bool)
		for _, branch := range cleanBranches {
			branchesMap[strings.TrimSpace(branch)] = true
		}

		var match []Worktree
		var noMatch []Worktree
		for _, worktree := range worktrees {

			branch := worktree.Head
			branch = strings.TrimSpace(branch)

			if branch == "" {
				continue
			}

			if !branchesMap[branch] {
				match = append(match, worktree)
			} else {
				noMatch = append(noMatch, worktree)
			}
		}

		var selection []string

		for _, tree := range noMatch {
			localString := tree.Head
			selection = append(selection, localString)
		}

		var choices []string

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Value(&choices).
					OptionsFunc(func() []huh.Option[string] {
						return huh.NewOptions(selection...)
					}, &selection).
					Title("Local endpoints that do not exist on remote").
					Height(25),
			),
		)

		formErr := form.Run()
		if formErr != nil {
			log.Fatal(formErr)
		}

		numOfWorktreesRemoved := 0

		action := func() {

			for _, branch := range choices {
				deleteWorktree := exec.Command("git", "worktree", "remove", branch, "--force")
				dwErr := deleteWorktree.Run()
				numOfWorktreesRemoved++
				if dwErr != nil {
					log.Fatal(dwErr)
				}

			}
		}
		err := spinner.New().
			Title("Removing Worktrees").
			Action(action).
			Run()

		if err != nil {
			log.Fatal(err)
		}

		if numOfWorktreesRemoved > 0 {
			fmt.Printf("%s Worktrees removed!", strconv.Itoa(numOfWorktreesRemoved))
		}

	},
}

func getWorktrees() []Worktree {

	//get all worktrees
	cmdToRun := exec.Command("git", "worktree", "list")
	allWorktrees, err := cmdToRun.Output()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	//clean worktrees
	lines := strings.Split(string(allWorktrees), "\n")
	worktrees := make([]Worktree, 0, len(lines))
	for _, line := range lines {
		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}
		worktrees = append(worktrees, Worktree{Path: parts[0], Head: ExtractTextInBrackets(parts[1])})
	}

	return worktrees
}

func getCleanRemoteBranchNames() []string {
	// fetch
	getFetch := exec.Command("git", "fetch", "origin")
	getFetch.Run()

	//get all branches
	getAllBranchesCmd := exec.Command("git", "branch", "-r")
	allBranches, err := getAllBranchesCmd.Output()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	//clean branch names
	branches := strings.Split(string(allBranches), "\n")
	var cleanBranches []string
	for _, branch := range branches {
		cleanBranch := strings.Replace(branch, "origin/", "", 1)
		cleanBranches = append(cleanBranches, cleanBranch)
	}
	return cleanBranches
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
