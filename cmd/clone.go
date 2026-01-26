package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/garrettkrohn/treekanga/adapters"
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

func CloneBareRepo(git adapters.GitAdapter, spinner spinner.HuhSpinner, args []string) {
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

	// Clone with streaming output so user can see git progress
	err := git.CloneBare(url, folderName)
	util.CheckError(err)

	workingDir, err := os.Getwd()
	util.CheckError(err)

	barePath := workingDir + "/" + folderName

	git.ConfigureGitBare(barePath)

	fmt.Printf("\n✓ Successfully cloned %s\n", folderName)
}

func getProjectName(url string) string {
	lastSlashIndex := strings.LastIndex(url, "/")
	if lastSlashIndex == -1 {
		return url
	}
	return url[lastSlashIndex+1:]

}

func init() {
	// No additional flags needed for clone command
}
