package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/garrettkrohn/treekanga/git"
	spinner "github.com/garrettkrohn/treekanga/spinnerHuh"
	util "github.com/garrettkrohn/treekanga/utility"
	"github.com/spf13/cobra"
)

var (
	url        string
	folderName string
)

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone a repository as a bare repo",
	Long: `Clone a repository as a bare repository for worktree management.

    Bare repositories are ideal for worktree workflows as they don't have 
    a working directory, allowing you to create multiple worktrees from 
    the same repository.
    
    Usage:
      treekanga clone <repository_url> [folder_name]
    
    If no folder name is provided, it will use the repository name 
    with "_bare" suffix.`,
	Run: func(cmd *cobra.Command, args []string) {
		CloneBareRepo(deps.Git, spinner.NewRealHuhSpinner(), args)
	},
}

func CloneBareRepo(git git.Git, spinner spinner.HuhSpinner, args []string) {
	if len(args) == 0 {
		fmt.Print("must include url to clone, folder name can be included optionally")
	}

	url = args[0]

	if len(args) == 2 {
		folderName = args[1]
	} else {
		folderName = getProjectName(url)
		folderName = fmt.Sprintf("%s_bare", folderName)
	}

	spinner.Title("Cloning bare repo")
	spinner.Action(func() {
		err := deps.Git.CloneBare(url, folderName)
		util.CheckError(err)

		workingDir, err := os.Getwd()
		util.CheckError(err)

		barePath := workingDir + "/" + folderName

		deps.Git.ConfigureGitBare(barePath)
	})
	spinner.Run()
}

func getProjectName(url string) string {
	lastSlashIndex := strings.LastIndex(url, "/")
	if lastSlashIndex == -1 {
		return url
	}
	return url[lastSlashIndex+1:]

}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
