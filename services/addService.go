package services

import (
	"fmt"
	"os"
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

	// Sanitize worktree name by replacing slashes with dashes to prevent nested directory issues
	cfg.NewWorktreeName = strings.ReplaceAll(cfg.NewWorktreeName, "/", "-")

	// When checking out an existing remote branch, fetch it first so a branch
	// pushed after the last fetch is still found. A fetch failure here (e.g.
	// the branch doesn't exist at all) isn't fatal on its own - the existence
	// check further down will catch a genuinely missing branch.
	if cfg.CheckoutRemote {
		if err := git.Fetch(cfg.BareRepoPath, cfg.NewBranchName); err != nil {
			log.Debug("Failed to fetch target branch from remote, falling back to cached remote-tracking refs", "branch", cfg.NewBranchName, "error", err)
		}
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
	CheckoutRemote             bool
	CheckoutLocal              bool
	NewBranchExistsLocally     bool
	NewBranchExistsRemotely    bool
	BaseBranchExistsLocally    bool
	BaseBranchExistsRemotely   bool
	NewBranchName              string
	PullBeforeCuttingNewBranch bool
	BaseBranch                 string
	NewWorktreeName            string
}

func GetAddWorktreeArguements(params AddWorktreeConfig) []string {
	// Case 1: Checkout existing remote branch
	if params.CheckoutRemote {
		return []string{params.NewBranchName}
	}

	// Case 2: Checkout existing local branch
	if params.CheckoutLocal {
		return []string{params.NewBranchName}
	}

	// Case 3: Default mode - create new branch from base branch
	// Base branch exists locally
	if params.BaseBranchExistsLocally {
		if params.PullBeforeCuttingNewBranch && params.BaseBranchExistsRemotely {
			// Create new branch from remote version of base branch (only if it exists remotely)
			return []string{"-b", params.NewBranchName, "--no-track", "origin/" + params.BaseBranch}
		} else {
			// Create new branch from local version of base branch
			return []string{"-b", params.NewBranchName, "--no-track", params.BaseBranch}
		}
	}

	// Base branch only exists remotely
	return []string{"-b", params.NewBranchName, "--no-track", "origin/" + params.BaseBranch}
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

	// Validation: Check mode and branch existence constraints
	if cfg.CheckoutRemote {
		// --remote mode: branch must exist remotely
		if !cfg.NewBranchExistsRemotely {
			log.Fatal(fmt.Sprintf("Branch '%s' not found on remote", cfg.NewBranchName))
		}
		log.Debug("Checkout mode: remote - ignoring -b and -p flags if set")
	} else if cfg.CheckoutLocal {
		// --local mode: branch must exist locally
		if !cfg.NewBranchExistsLocally {
			log.Fatal(fmt.Sprintf("Branch '%s' not found locally", cfg.NewBranchName))
		}
		log.Debug("Checkout mode: local - ignoring -b and -p flags if set")
	} else {
		// Default mode: branch must NOT exist
		if cfg.NewBranchExistsLocally || cfg.NewBranchExistsRemotely {
			log.Fatal(fmt.Sprintf("Branch '%s' already exists. Use --remote or --local to checkout existing branch.", cfg.NewBranchName))
		}
		log.Debug("Default mode: creating new branch")
	}

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

		// Update the BaseBranchExists flags after selection
		localBranches, err := git.GetLocalBranches(cfg.BareRepoPath)
		utility.CheckError(err)
		cfg.BaseBranchExistsLocally = slices.Contains(localBranches, cfg.BaseBranch)
		log.Debug(fmt.Sprintf("Updated BaseBranchExistsLocally = %t after form selection", cfg.BaseBranchExistsLocally))

		remoteBranches, err := git.GetRemoteBranches(cfg.BareRepoPath)
		utility.CheckError(err)
		cfg.BaseBranchExistsRemotely = slices.Contains(remoteBranches, cfg.BaseBranch)
		log.Debug(fmt.Sprintf("Updated BaseBranchExistsRemotely = %t after form selection", cfg.BaseBranchExistsRemotely))
	}

	// Fetch the latest state of base branch if pull flag is set
	if cfg.PullBeforeCuttingNewBranch && cfg.BaseBranchExistsRemotely {
		err := git.Fetch(cfg.BareRepoPath, cfg.BaseBranch)
		utility.CheckError(err)
		log.Debug(fmt.Sprintf("Fetched latest state of %s from remote", cfg.BaseBranch))
	}

	worktreeAddArgs := GetAddWorktreeArguements(AddWorktreeConfig{
		BareRepoPath:               cfg.BareRepoPath,
		WorktreeTargetDirectory:    cfg.WorktreeTargetDir,
		CheckoutRemote:             cfg.CheckoutRemote,
		CheckoutLocal:              cfg.CheckoutLocal,
		NewBranchExistsLocally:     cfg.NewBranchExistsLocally,
		NewBranchExistsRemotely:    cfg.NewBranchExistsRemotely,
		BaseBranchExistsLocally:    cfg.BaseBranchExistsLocally,
		BaseBranchExistsRemotely:   cfg.BaseBranchExistsRemotely,
		NewBranchName:              cfg.NewBranchName,
		PullBeforeCuttingNewBranch: cfg.PullBeforeCuttingNewBranch,
		BaseBranch:                 cfg.BaseBranch,
		NewWorktreeName:            cfg.NewWorktreeName,
	})

	err := git.AddWorktree(cfg.BareRepoPath, cfg.WorktreeTargetDir, cfg.NewWorktreeName, worktreeAddArgs)
	utility.CheckError(err)

	//TODO: different place for this?
	newRootDirectory := cfg.WorktreeTargetDir + "/" + cfg.NewWorktreeName

	// Set upstream for new branches (not existing ones)
	if !cfg.CheckoutRemote && !cfg.CheckoutLocal {
		err = git.SetUpstream(newRootDirectory, cfg.NewBranchName)
		if err != nil {
			log.Warn("Failed to set upstream branch", "branch", cfg.NewBranchName, "error", err)
		}
	}

	if cfg.CheckoutRemote {
		log.Info("worktree created with remote branch", "branch", cfg.NewBranchName)
	} else if cfg.CheckoutLocal {
		log.Info("worktree created with local branch", "branch", cfg.NewBranchName)
	} else {
		log.Info("worktree created with new branch cut from branch",
			"newBranch", cfg.NewBranchName,
			"baseBranch", cfg.BaseBranch)
	}

	if cfg.TmuxConnect != "" {
		connectPath := newRootDirectory
		if cfg.TmuxConnect != "." {
			connectPath = newRootDirectory + "/" + cfg.TmuxConnect
		}
		opts := models.ConnectOpts{Switch: false}
		if err := connector.ConnectWithConfig(connectPath, opts, cfg.PostScriptPath, cfg.RunPostScript); err != nil {
			log.Warn("Subdirectory not found, connecting to root instead", "subdirectory", cfg.TmuxConnect)
			if err := connector.ConnectWithConfig(newRootDirectory, opts, cfg.PostScriptPath, cfg.RunPostScript); err != nil {
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

	// Post-script execution is handled by ConnectWithConfig when using tmux connect flag
	// If not connecting to a new session, run the script in the current context
	if cfg.RunPostScript && cfg.TmuxConnect == "" && !cfg.CursorConnect && !cfg.VsCodeConnect {
		log.Info("Running post script in current session")
		script := cfg.PostScriptPath
		// Expand tilde in script path
		expandedPath := script
		if len(script) > 2 && script[0] == '~' && script[1] == '/' {
			homeDir, err := os.UserHomeDir()
			if err == nil {
				expandedPath = homeDir + script[1:]
			}
		}
		// Run the script in a subshell so the user stays in their current directory
		command := fmt.Sprintf("(cd %s && sh %s)", newRootDirectory, expandedPath)
		shell.Cmd("tmux", "send-keys", "-t", ".", command, "Enter")
		log.Info("Post script command sent to current session")
	}
}
