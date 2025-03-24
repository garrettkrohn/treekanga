package cmd

import (
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/garrettkrohn/treekanga/filter"
	"github.com/garrettkrohn/treekanga/form"
	"github.com/garrettkrohn/treekanga/git"
	spinner "github.com/garrettkrohn/treekanga/spinnerHuh"
	"github.com/garrettkrohn/treekanga/transformer"
	util "github.com/garrettkrohn/treekanga/utility"
	worktreeobj "github.com/garrettkrohn/treekanga/worktreeObj"
	"github.com/garrettkrohn/treekanga/zoxide"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete selected worktrees",
	Long:  `List all worktrees and selected multiple to be deleted`,
	Run: func(cmd *cobra.Command, args []string) {
		numOfWorktreesRemoved, err := deleteWorktrees(deps.Git,
			transformer.NewTransformer(),
			filter.NewFilter(),
			spinner.NewRealHuhSpinner(),
			form.NewHuhForm(),
			deps.Zoxide)
		if err != nil {
			cmd.PrintErrln("Error:", err)
			return
		}
		log.Info("worktrees removed: %s", strconv.Itoa(numOfWorktreesRemoved))
	},
}

// deleteWorktrees performs the core logic of deleting worktrees
func deleteWorktrees(git git.Git,
	transformer *transformer.RealTransformer,
	filter filter.Filter,
	spinner spinner.HuhSpinner,
	form form.Form,
	zoxide zoxide.Zoxide) (int, error) {
	worktrees := getWorktrees(git, transformer)

	stringWorktrees := transformer.TransformWorktreesToBranchNames(worktrees)

	var selections []string

	form.SetSelections(&selections)
	form.SetOptions(stringWorktrees)
	err := form.Run()
	util.CheckError(err)

	// Transform string selection back to worktree objects
	selectedWorktreeObj := filter.GetBranchMatchList(selections, worktrees)

	removeWorktrees(selectedWorktreeObj, spinner, git, zoxide)

	return len(selectedWorktreeObj), nil
}

func removeWorktrees(worktrees []worktreeobj.WorktreeObj, spinner spinner.HuhSpinner, git git.Git, zoxide zoxide.Zoxide) {
	spinner.Title("Deleting Worktrees")
	spinner.Action(func() {
		for _, worktreeObj := range worktrees {
			_, err := git.RemoveWorktree(worktreeObj.Folder)
			_ = zoxide.RemovePath(worktreeObj.FullPath)
			util.CheckError(err)
		}
	})
	spinner.Run()

}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
