package cmd

import (
	"fmt"
	"os"
	"slices"

	"github.com/charmbracelet/huh"
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

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete selected worktrees",
	Long: `Interactive deletion of worktrees with optional filtering.

    Lists all worktrees and allows you to select multiple to be deleted.
    
    Available flags:
    -s, --stale: Only show worktrees where branches don't exist on remote
    -d, --delete: CAUTION - Also delete the local branches`,
	Run: func(cmd *cobra.Command, args []string) {
		stale, err := cmd.Flags().GetBool("stale")
		util.CheckError(err)
		deleteBranches, err := cmd.Flags().GetBool("delete")
		util.CheckError(err)

		numOfWorktreesRemoved, err := deleteWorktrees(deps.Git,
			transformer.NewTransformer(),
			filter.NewFilter(),
			spinner.NewRealHuhSpinner(),
			form.NewHuhForm(),
			deps.Zoxide,
			args,
			stale,
			deleteBranches)
		if err != nil {
			cmd.PrintErrln("Error:", err)
			return
		}
		log.Info("worktrees removed", "count", numOfWorktreesRemoved)
	},
}

func deleteWorktrees(git git.Git,
	transformer *transformer.RealTransformer,
	filter filter.Filter,
	spinner spinner.HuhSpinner,
	form form.Form,
	zoxide zoxide.Zoxide,
	listOfBranchesToDelete []string,
	applyStaleFilter bool,
	deleteBranches bool) (int, error) {

	var selections []string
	treesToDeleteAreValid := false

	worktrees := getWorktrees(git, transformer)

	if applyStaleFilter {
		worktrees = filterLocalBranchesOnly(worktrees, transformer, filter)
		if len(worktrees) == 0 {
			log.Fatal("All local branches exist on remote")
		}
	}

	// get names to display
	stringWorktrees := transformer.TransformWorktreesToBranchNames(worktrees)

	// branches can be provided via args or the form
	if len(listOfBranchesToDelete) > 0 {
		log.Debug("branch(es) submitted as argument(s)", "branches", listOfBranchesToDelete)
		treesToDeleteAreValid = validateAllBranchesToDelete(stringWorktrees, listOfBranchesToDelete)
		if !treesToDeleteAreValid {
			log.Error("At least one of the branches provided were not valid, please select a branch")
		} else {
			log.Info("All branches are valid")
			selections = listOfBranchesToDelete
		}
	}

	// need to make this cleaner
	if !treesToDeleteAreValid {
		log.Debug("activating selection form")
		form.SetSelections(&selections)
		form.SetOptions(stringWorktrees)
		err := form.Run()
		util.CheckError(err)
	}

	// transform selection back into worktreeObj
	selectedWorktreeObj := filter.GetBranchMatchList(selections, worktrees)

	// remove worktrees
	removeWorktrees(selectedWorktreeObj, spinner, git, zoxide)

	// delete branches
	if deleteBranches {
		log.Debug("delete branches flag true")
		deleteLocalBranches(selectedWorktreeObj)
	}

	return len(selectedWorktreeObj), nil
}

func deleteLocalBranches(selectedWorktreeObj []worktreeobj.WorktreeObj) {
	confirm := false

	confirmationMessage := "Are you sure you want to delete these branches: "

	for _, worktreeObj := range selectedWorktreeObj {
		confirmationMessage += worktreeObj.BranchName
	}

	confirmDialog := huh.NewConfirm().
		Title(confirmationMessage).
		Affirmative("Yes!").
		Negative("No.").
		Value(&confirm)

	confirmDialog.Run()

	if confirm {
		for _, worktreeObj := range selectedWorktreeObj {
			// Use the bare repo path if available, otherwise fall back to current directory
			dir := deps.BareRepoPath
			if dir == "" {
				var err error
				dir, err = os.Getwd()
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
			}
			log.Debug("Deleting branch ref", "branch", worktreeObj.BranchName, "path", dir)
			deps.Git.DeleteBranchRef(worktreeObj.BranchName, dir)

			log.Debug("Deleting branch", "branch", worktreeObj.BranchName, "path", dir)
			deps.Git.DeleteBranch(worktreeObj.BranchName, dir)
		}
	} else {
		log.Info("No local branches were deleted")
	}

}

func getWorktrees(git git.Git, transformer *transformer.RealTransformer) []worktreeobj.WorktreeObj {
	var worktreeStrings []string
	var wError error

	if deps.BareRepoPath != "" {
		worktreeStrings, wError = git.GetWorktrees(&deps.BareRepoPath)
	} else {
		worktreeStrings, wError = git.GetWorktrees(nil)
	}

	if wError != nil {
		log.Fatal(wError)
	}

	worktrees := transformer.TransformWorktrees(worktreeStrings)

	// Sort worktrees by most recently modified
	sortWorktreesByModTime(worktrees)

	return worktrees
}

func validateAllBranchesToDelete(stringWorktrees []string, listOfBranchesToDelete []string) bool {
	for _, branch := range listOfBranchesToDelete {
		if !slices.Contains(stringWorktrees, branch) {
			return false
		}
	}
	return true
}

func removeWorktrees(worktrees []worktreeobj.WorktreeObj, spinner spinner.HuhSpinner, git git.Git, zoxide zoxide.Zoxide) {
	log.Info("removing %d worktrees", len(worktrees))

	// Use the resolved bare repo path if available
	var path *string
	if deps.BareRepoPath != "" {
		path = &deps.BareRepoPath
		log.Debug("Using bare repo path for removing worktrees", "path", deps.BareRepoPath)
	}

	for _, worktreeObj := range worktrees {
		log.Info("Removing worktree", "fullPath", worktreeObj.FullPath, "folder", worktreeObj.Folder, "branch", worktreeObj.BranchName)
		output, err := git.RemoveWorktree(worktreeObj.FullPath, path)
		log.Debug("RemoveWorktree returned", "output", output, "error", err)
		_ = zoxide.RemovePath(worktreeObj.FullPath)
		util.CheckError(err)
		log.Debug("Worktree removed successfully")
	}
}

func filterLocalBranchesOnly(worktrees []worktreeobj.WorktreeObj,
	transformer *transformer.RealTransformer,
	filter filter.Filter) []worktreeobj.WorktreeObj {

	log.Info("filtering local branches only")

	// Use the resolved bare repo path if available
	var path *string
	if deps.BareRepoPath != "" {
		path = &deps.BareRepoPath
		log.Debug("Using bare repo path for remote branches", "path", deps.BareRepoPath)
	}

	branches, err := deps.Git.GetRemoteBranches(path)
	util.CheckError(err)
	cleanedBranches := transformer.RemoveOriginPrefix(branches)
	worktrees = filter.GetBranchNoMatchList(cleanedBranches, worktrees)
	return worktrees
}

func init() {
	deleteCmd.Flags().BoolP("stale", "s", false, "Only show worktrees where the branches don't exist on remote")
	deleteCmd.Flags().BoolP("delete", "d", false, "CAUTION: delete the local branch")
}
