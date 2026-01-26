package services

import (
	"fmt"
	"os"
	"slices"
	"sort"

	"github.com/charmbracelet/log"
	"github.com/garrettkrohn/treekanga/adapters"
	"github.com/garrettkrohn/treekanga/config"
	"github.com/garrettkrohn/treekanga/confirmer"
	"github.com/garrettkrohn/treekanga/filter"
	"github.com/garrettkrohn/treekanga/form"
	"github.com/garrettkrohn/treekanga/transformer"
	util "github.com/garrettkrohn/treekanga/utility"
	worktreeobj "github.com/garrettkrohn/treekanga/worktreeObj"
	"github.com/garrettkrohn/treekanga/zoxide"
)

func DeleteWorktrees(git adapters.GitAdapter,
	transformer *transformer.RealTransformer,
	filter filter.Filter,
	form form.Form,
	zoxide zoxide.Zoxide,
	listOfBranchesToDeleteFromArgs []string,
	cfg config.AppConfig) (int, error) {

	var selections []string
	treesToDeleteAreValid := false

	//1. get all worktrees
	worktrees := getWorktrees(git, transformer, cfg.BareRepoPath)

	//2. filter for only worktrees that don't exist on remote
	if cfg.FilterOnlyStaleBranches {
		worktrees = filterLocalBranchesOnly(git, worktrees, transformer, filter, cfg.BareRepoPath)
		if len(worktrees) == 0 {
			log.Fatal("All local branches exist on remote")
		}
	}

	// get names to display
	stringWorktrees := transformer.TransformWorktreesToBranchNames(worktrees)

	// branches can be provided via args or the form
	if len(listOfBranchesToDeleteFromArgs) > 0 {
		log.Debug("branch(es) submitted as argument(s)", "branches", listOfBranchesToDeleteFromArgs)
		treesToDeleteAreValid = validateAllBranchesToDelete(stringWorktrees, listOfBranchesToDeleteFromArgs)
		if !treesToDeleteAreValid {
			log.Error("At least one of the branches provided were not valid, please select a branch")
		} else {
			log.Info("All branches are valid")
			selections = listOfBranchesToDeleteFromArgs
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

	// get list of full paths
	worktreeFullPaths := getWorktreeFullPaths(selectedWorktreeObj)

	// remove worktrees
	removeWorktrees(worktreeFullPaths, git, zoxide, cfg.ForceDelete, cfg.BareRepoPath)

	// delete branches
	if cfg.DeleteBranch {
		log.Debug("delete branches flag true")
		deleteLocalBranches(git, selectedWorktreeObj, cfg.ForceDelete, cfg.BareRepoPath, confirmer.NewConfirmer())
	}

	return len(selectedWorktreeObj), nil
}

func getWorktreeFullPaths(worktrees []worktreeobj.WorktreeObj) []string {
	var worktreeFullPathStrings []string

	for _, worktree := range worktrees {
		worktreeFullPathStrings = append(worktreeFullPathStrings, worktree.FullPath)
	}
	return worktreeFullPathStrings

}

func deleteLocalBranches(git adapters.GitAdapter, selectedWorktreeObj []worktreeobj.WorktreeObj, forceDelete bool, bareRepoPath string, confirmer confirmer.Confirmer) {
	confirm := false

	confirmationMessage := "Are you sure you want to delete these branches: "

	for _, worktreeObj := range selectedWorktreeObj {
		confirmationMessage += worktreeObj.BranchName
	}

	confirm, err := confirmer.Confirm(confirmationMessage)
	if err != nil {
		log.Error("There was an error with the confirmation message")
	}

	if confirm {
		for _, worktreeObj := range selectedWorktreeObj {
			// Use the bare repo path if available, otherwise fall back to current directory
			dir := bareRepoPath
			if dir == "" {
				var err error
				dir, err = os.Getwd()
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
			}
			log.Debug("Deleting branch ref", "branch", worktreeObj.BranchName, "path", dir)
			git.DeleteBranchRef(worktreeObj.BranchName, dir)

			log.Debug("Deleting branch", "branch", worktreeObj.BranchName, "path", dir)
			git.DeleteBranch(worktreeObj.BranchName, dir, forceDelete)
		}
	} else {
		log.Info("No local branches were deleted")
	}

}

func validateAllBranchesToDelete(stringWorktrees []string, listOfBranchesToDelete []string) bool {
	for _, branch := range listOfBranchesToDelete {
		if !slices.Contains(stringWorktrees, branch) {
			return false
		}
	}
	return true
}

func removeWorktrees(worktreePaths []string, git adapters.GitAdapter, zoxide zoxide.Zoxide, forceDelete bool, bareRepoPath string) {
	log.Debug("removeWorktrees called", "count", len(worktreePaths))

	// Use the resolved bare repo path if available
	var path *string
	if bareRepoPath != "" {
		path = &bareRepoPath
		log.Debug("Using bare repo path for removing worktrees", "path", bareRepoPath)
	}

	for _, worktreePath := range worktreePaths {
		log.Debug("Removing worktree", "fullPath", worktreePath)
		err := git.RemoveWorktree(worktreePath, path, forceDelete)
		_ = zoxide.RemovePath(worktreePath)
		util.CheckError(err)
		log.Debug("Worktree removed successfully")
	}
}

func filterLocalBranchesOnly(git adapters.GitAdapter, worktrees []worktreeobj.WorktreeObj,
	transformer *transformer.RealTransformer,
	filter filter.Filter,
	bareRepoPath string) []worktreeobj.WorktreeObj {

	log.Info("filtering local branches only")

	// Use the resolved bare repo path if available
	var path *string
	if bareRepoPath != "" {
		path = &bareRepoPath
		log.Debug("Using bare repo path for remote branches", "path", bareRepoPath)
	}

	branches, err := git.GetRemoteBranches(path)
	util.CheckError(err)
	cleanedBranches := transformer.RemoveOriginPrefix(branches)
	worktrees = filter.GetBranchNoMatchList(cleanedBranches, worktrees)
	return worktrees
}

// TODO: remove dupilcate code here
// sortWorktreesByModTime sorts worktrees by modification time (most recent first)
func sortWorktreesByModTime(worktrees []worktreeobj.WorktreeObj) {
	sort.Slice(worktrees, func(i, j int) bool {
		statI, errI := os.Stat(worktrees[i].FullPath)
		statJ, errJ := os.Stat(worktrees[j].FullPath)

		// If there's an error accessing either path, push it to the end
		if errI != nil {
			log.Debug("Error stat'ing worktree", "path", worktrees[i].FullPath, "error", errI)
			return false
		}
		if errJ != nil {
			log.Debug("Error stat'ing worktree", "path", worktrees[j].FullPath, "error", errJ)
			return true
		}

		// Sort by modification time, most recent first
		return statI.ModTime().After(statJ.ModTime())
	})
}
