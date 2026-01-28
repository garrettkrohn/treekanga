package services

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/garrettkrohn/treekanga/adapters"
	"github.com/garrettkrohn/treekanga/config"
	"github.com/garrettkrohn/treekanga/connector"
	"github.com/garrettkrohn/treekanga/directoryReader"
	"github.com/garrettkrohn/treekanga/form"
	"github.com/garrettkrohn/treekanga/shell"
	"github.com/garrettkrohn/treekanga/transformer"
	"github.com/garrettkrohn/treekanga/utility"
	util "github.com/garrettkrohn/treekanga/utility"
)

func SetConfigForAddService(gitClient adapters.GitAdapter, cfg config.AppConfig, args []string) config.AppConfig {
	log.Debug("Running configuration for add command")

	if len(args) == 1 {
		cfg.NewBranchName = strings.TrimSpace(args[0])
		log.Debug(fmt.Sprintf("Setting newBranchName = %s in addService", args[0]))
	} else {
		log.Fatal("please include new branch name as an argument")
	}

	if cfg.NewWorktreeName == "" {
		cfg.NewWorktreeName = cfg.NewBranchName
		log.Debug(fmt.Sprintf("No worktree name specified in flags, so defaults to new branch name: %s", cfg.NewWorktreeName))
	}

	t := transformer.NewTransformer()

	remoteBranches, err := gitClient.GetRemoteBranches(&cfg.BareRepoPath)
	util.CheckError(err)
	cleanRemoteBranches := t.RemoveOriginPrefix(remoteBranches)
	log.Debug(cleanRemoteBranches)

	localBranches, err := gitClient.GetLocalBranches(&cfg.BareRepoPath)
	util.CheckError(err)
	cleanLocalBranches := t.RemoveQuotes(localBranches)
	log.Debug(cleanLocalBranches)

	cfg.NewBranchExistsLocally = slices.Contains(cleanLocalBranches, cfg.NewBranchName)
	log.Debug(fmt.Sprintf("Setting NewBranchExistsLocally = %t from addService", cfg.NewBranchExistsLocally))

	cfg.NewBranchExistsRemotely = slices.Contains(cleanRemoteBranches, cfg.NewBranchName)
	log.Debug(fmt.Sprintf("Setting NewBranchExistsRemotely = %t from addService", cfg.NewBranchExistsRemotely))

	cfg.BaseBranchExistsLocally = slices.Contains(cleanLocalBranches, cfg.BaseBranch)
	log.Debug(fmt.Sprintf("Setting BaseBranchExistsLocally = %t from addService", cfg.BaseBranchExistsLocally))

	cfg.BaseBranchExistsRemotely = slices.Contains(cleanRemoteBranches, cfg.BaseBranch)
	log.Debug(fmt.Sprintf("Setting BaseBranchExistsRemotely = %t from addService", cfg.BaseBranchExistsRemotely))

	return cfg
}

type AddWorktreeConfig struct {
	BareRepoPath               string
	WorktreeTargetDirectory    string
	NewBranchExistsLocally     bool
	NewBranchExistsRemotely    bool
	BaseBranchExistsLocally    bool
	NewBranchName              string
	PullBeforeCuttingNewBranch bool
	BaseBranch                 string
	NewWorktreeName            string
}

func GetAddWorktreeArguements(params AddWorktreeConfig) []string {
	// Case 1: Branch already exists (locally or remotely) - just checkout
	if params.NewBranchExistsLocally || params.NewBranchExistsRemotely {
		return []string{params.NewBranchName}
	}

	// Case 2: Base branch exists locally
	if params.BaseBranchExistsLocally {
		if params.PullBeforeCuttingNewBranch {
			// Create new branch from remote version of base branch
			return []string{"-b", params.NewBranchName, "origin/" + params.BaseBranch, "--no-track"}
		} else {
			// Create new branch from local version of base branch
			return []string{"-b", params.NewBranchName, params.BaseBranch}
		}
	}

	// Case 3: Base branch only exists remotely
	return []string{"-b", params.NewBranchName, "origin/" + params.BaseBranch, "--no-track"}
}

func handleFromForm(form form.HuhForm, worktrees []string) string {
	// Present selection interface
	var selectedBranch string
	form.SetSingleSelection(&selectedBranch)
	form.SetOptions(worktrees)
	form.SetTitle("Select base branch for new worktree:")
	err := form.Run()
	util.CheckError(err)

	if selectedBranch == "" {
		log.Fatal("No branch selected")
	}

	log.Info("Selected base branch", "branch", selectedBranch)
	return selectedBranch
}

func AddWorktree(gitClient adapters.GitAdapter, zoxide adapters.Zoxide, connector connector.Connector, shell shell.Shell, cfg config.AppConfig) {

	if cfg.UseFormToSetBaseBranch {
		worktrees, err := gitClient.GetWorktrees(&cfg.BareRepoPath)
		utility.CheckError(err)

		worktreeObjects := transformer.NewTransformer().TransformWorktrees(worktrees)

		SortWorktreesByModTime(worktreeObjects)

		var branchStrings []string

		for _, wt := range worktreeObjects {
			branchStrings = append(branchStrings, wt.BranchName)
		}

		form := form.NewHuhForm()

		selectedBranch := handleFromForm(*form, branchStrings)
		cfg.BaseBranch = selectedBranch
		log.Debug(fmt.Sprintf("Set BaseBranch = %s from form selection", selectedBranch))
		
		// Update the BaseBranchExistsLocally flag after selection
		t := transformer.NewTransformer()
		localBranches, err := gitClient.GetLocalBranches(&cfg.BareRepoPath)
		utility.CheckError(err)
		cleanLocalBranches := t.RemoveQuotes(localBranches)
		cfg.BaseBranchExistsLocally = slices.Contains(cleanLocalBranches, cfg.BaseBranch)
		log.Debug(fmt.Sprintf("Updated BaseBranchExistsLocally = %t after form selection", cfg.BaseBranchExistsLocally))
	}

	worktreeAddArgs := GetAddWorktreeArguements(AddWorktreeConfig{
		BareRepoPath:               cfg.BareRepoPath,
		WorktreeTargetDirectory:    cfg.WorktreeTargetDir,
		NewBranchExistsLocally:     cfg.NewBranchExistsLocally,
		NewBranchExistsRemotely:    cfg.NewBranchExistsRemotely,
		BaseBranchExistsLocally:    cfg.BaseBranchExistsLocally,
		NewBranchName:              cfg.NewBranchName,
		PullBeforeCuttingNewBranch: cfg.PullBeforeCuttingNewBranch,
		BaseBranch:                 cfg.BaseBranch,
		NewWorktreeName:            cfg.NewWorktreeName,
	})

	err := gitClient.AddWorktree(cfg.BareRepoPath, cfg.WorktreeTargetDir, cfg.NewWorktreeName, worktreeAddArgs)
	util.CheckError(err)

	if cfg.NewBranchExistsLocally {
		log.Info("worktree created with existing branch", "branch", cfg.NewBranchName)
	} else {
		log.Info("worktree created with new branch cut from branch",
			"newBranch", cfg.NewBranchName,
			"baseBranch", cfg.BaseBranch)
	}

	//TODO: different place for this?
	newRootDirectory := cfg.WorktreeTargetDir + "/" + cfg.NewWorktreeName

	// always add the root, add folders if included
	log.Info("adding zoxide entries")
	directoryReader := directoryReader.NewDirectoryReader()
	zoxidePathsToAdd := CompileZoxidePathsToAdd(cfg.ZoxideFolders, newRootDirectory, directoryReader)
	zoxide.AddZoxideEntries(zoxidePathsToAdd)

	if cfg.SeshConnect != "" {
		seshConnectPath := GetSeshPath(cfg.SeshConnect, cfg.ZoxideFolders, newRootDirectory)
		connector.SeshConnect(seshConnectPath)
	}

	if cfg.CursorConnect {
		connector.CursorConnect(newRootDirectory)
	}

	if cfg.VsCodeConnect {
		connector.VsCodeConnect(newRootDirectory)
	}

	//TODO: allow the user to set where this is run?
	if cfg.RunPostScript {
		log.Info("Runnning post script")
		script := cfg.PostScriptPath
		shell.CmdWithDir(cfg.WorktreeTargetDir, "sh", "-c", script)
		log.Info("post script run", "command", script)
	}
}
