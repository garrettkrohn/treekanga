package services

import (
	"github.com/charmbracelet/log"
	"github.com/garrettkrohn/treekanga/config"
	"github.com/garrettkrohn/treekanga/connector"
	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/shell"
	util "github.com/garrettkrohn/treekanga/utility"
	"github.com/garrettkrohn/treekanga/zoxide"
)

func AddWorktree(gitClient git.GitAdapter, zoxide zoxide.Zoxide, connector connector.Connector, shell shell.Shell, cfg config.AppConfig) {

	log.Info("Validating configuration")
	// c := com.AddConfig{}
	// getAddCmdConfig(cmd, args, &c)

	// validateConfig(&c)

	// log.Debug("Adding worktree with config:")
	// PrintConfig(c)
	err := gitClient.AddWorktree(git.AddWorktreeConfig{
		BareRepoPath:               cfg.BareRepoPath,
		TargetDirectory:            cfg.TargetDirectory,
		NewBranchExistsLocally:     cfg.NewBranchExistsLocally,
		NewBranchExistsRemotely:    cfg.NewBranchExistsRemotely,
		BaseBranchExistsLocally:    cfg.BaseBranchExistsLocally,
		NewBranchName:              cfg.NewBranchName,
		PullBeforeCuttingNewBranch: cfg.PullBeforeCuttingNewBranch,
		BaseBranch:                 cfg.BaseBranch,
	})
	util.CheckError(err)

	if c.GitInfo.NewBranchExistsLocally {
		log.Info("worktree created with existing branch", "branch", c.GetNewBranchName())
	} else {
		log.Info("worktree created with new branch cut from branch",
			"newBranch", c.GetNewBranchName(),
			"baseBranch", c.GitInfo.BaseBranchName)
	}

	log.Info("adding zoxide entries")
	zoxide.AddZoxideEntries(&c)

	if c.HasSeshTarget() {
		connector.SeshConnect(&c)
	}

	if c.ShouldOpenCursor() {
		connector.CursorConnect(&c)
	}

	if c.ShouldOpenVSCode() {
		connector.VsCodeConnect(&c)
	}

	if cfg.RunPostScript {
		log.Info("Runnning post script")
		script := c.GetPostScript()
		shell.CmdWithDir(c.WorktreeTargetDir, "sh", "-c", script)
		log.Info("post script run", "command", script)
	}
}
