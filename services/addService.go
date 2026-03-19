package services

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/garrettkrohn/treekanga/config"
	"github.com/garrettkrohn/treekanga/connector"
	"github.com/garrettkrohn/treekanga/form"
	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/models"
	"github.com/garrettkrohn/treekanga/shell"
	"github.com/garrettkrohn/treekanga/transformer"
	"github.com/garrettkrohn/treekanga/util"
	"github.com/garrettkrohn/treekanga/utility"
)

func SetConfigForAddService(cfg config.AppConfig, args []string) config.AppConfig {
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

	remoteBranches, err := git.GetRemoteBranches(cfg.BareRepoPath)
	utility.CheckError(err)
	log.Debug("Remote branches", "branches", remoteBranches)

	localBranches, err := git.GetLocalBranches(cfg.BareRepoPath)
	utility.CheckError(err)
	log.Debug("Local branches", "branches", localBranches)

	cfg.NewBranchExistsLocally = slices.Contains(localBranches, cfg.NewBranchName)
	log.Debug(fmt.Sprintf("Setting NewBranchExistsLocally = %t from addService", cfg.NewBranchExistsLocally))

	cfg.NewBranchExistsRemotely = slices.Contains(remoteBranches, cfg.NewBranchName)
	log.Debug(fmt.Sprintf("Setting NewBranchExistsRemotely = %t from addService", cfg.NewBranchExistsRemotely))

	cfg.BaseBranchExistsLocally = slices.Contains(localBranches, cfg.BaseBranch)
	log.Debug(fmt.Sprintf("Setting BaseBranchExistsLocally = %t from addService", cfg.BaseBranchExistsLocally))

	cfg.BaseBranchExistsRemotely = slices.Contains(remoteBranches, cfg.BaseBranch)
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
	utility.CheckError(err)

	if selectedBranch == "" {
		log.Fatal("No branch selected")
	}

	log.Info("Selected base branch", "branch", selectedBranch)
	return selectedBranch
}

func AddWorktree(connector connector.Connector, shell shell.Shell, cfg config.AppConfig) {

	if cfg.UseFormToSetBaseBranch {
		worktrees, err := git.ListWorktrees(cfg.BareRepoPath)
		utility.CheckError(err)

		worktreeObjects := transformer.TransformWorktrees(worktrees)

		util.SortWorktreesByModTime(worktreeObjects)

		var branchStrings []string

		for _, wt := range worktreeObjects {
			branchStrings = append(branchStrings, wt.BranchName)
		}

		form := form.NewHuhForm()

		selectedBranch := handleFromForm(*form, branchStrings)
		cfg.BaseBranch = selectedBranch
		log.Debug(fmt.Sprintf("Set BaseBranch = %s from form selection", selectedBranch))

		// Update the BaseBranchExistsLocally flag after selection
		localBranches, err := git.GetLocalBranches(cfg.BareRepoPath)
		utility.CheckError(err)
		cfg.BaseBranchExistsLocally = slices.Contains(localBranches, cfg.BaseBranch)
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

	err := git.AddWorktree(cfg.BareRepoPath, cfg.WorktreeTargetDir, cfg.NewWorktreeName, worktreeAddArgs)
	utility.CheckError(err)

	if cfg.NewBranchExistsLocally {
		log.Info("worktree created with existing branch", "branch", cfg.NewBranchName)
	} else {
		log.Info("worktree created with new branch cut from branch",
			"newBranch", cfg.NewBranchName,
			"baseBranch", cfg.BaseBranch)
	}

	//TODO: different place for this?
	newRootDirectory := cfg.WorktreeTargetDir + "/" + cfg.NewWorktreeName

	if cfg.TmuxConnect != "" {
		connectPath := newRootDirectory
		if cfg.TmuxConnect != "." {
			connectPath = newRootDirectory + "/" + cfg.TmuxConnect
		}
		opts := models.ConnectOpts{Switch: false}
		if err := connector.Connect(connectPath, opts); err != nil {
			log.Warn("Subdirectory not found, connecting to root instead", "subdirectory", cfg.TmuxConnect)
			if err := connector.Connect(newRootDirectory, opts); err != nil {
				log.Error("Failed to connect to tmux session", "error", err)
			}
		}
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
		shell.CmdWithDir(newRootDirectory, "sh", "-c", script)
		log.Info("post script run", "command", script)
	}
}
