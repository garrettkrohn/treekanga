package services

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/garrettkrohn/treekanga/adapters"
	"github.com/garrettkrohn/treekanga/config"
	"github.com/garrettkrohn/treekanga/connector"
	"github.com/garrettkrohn/treekanga/shell"
	"github.com/garrettkrohn/treekanga/transformer"
	util "github.com/garrettkrohn/treekanga/utility"
	"github.com/garrettkrohn/treekanga/zoxide"
)

func SetConfigForAddService(gitClient adapters.GitAdapter, cfg config.AppConfig, args []string) config.AppConfig {
	log.Info("Running configuration for add command")

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

func AddWorktree(gitClient adapters.GitAdapter, zoxide zoxide.Zoxide, connector connector.Connector, shell shell.Shell, cfg config.AppConfig) {

	err := gitClient.AddWorktree(adapters.AddWorktreeConfig{
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
	zoxidePathsToAdd := CompileZoxidePathsToAdd(cfg.ZoxideFolders, newRootDirectory)
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
