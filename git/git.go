package git

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/garrettkrohn/treekanga/shell"
)

type Git interface {
	ShowTopLevel(name string) (bool, string, error)
	GitCommonDir(name string) (bool, string, error)
	Clone(name string) (string, error)
	GetRemoteBranches() ([]string, error)
	GetLocalBranches() ([]string, error)
	GetWorktrees() ([]string, error)
	RemoveWorktree(string) (string, error)
	AddWorktree(string, bool, string, string) error
	GetRepoName(path string) (string, error)
	FetchOrigin(branch string) error
	CloneBare(string, string) error
}

type RealGit struct {
	shell shell.Shell
}

func NewGit(shell shell.Shell) Git {
	return &RealGit{shell}
}

func (g *RealGit) ShowTopLevel(path string) (bool, string, error) {
	out, err := g.shell.Cmd("git", "-C", path, "rev-parse", "--show-toplevel")
	if err != nil {
		return false, "", err
	}
	return true, out, nil
}

func (g *RealGit) GitCommonDir(path string) (bool, string, error) {
	out, err := g.shell.Cmd("git", "-C", path, "rev-parse", "--git-common-dir")
	if err != nil {
		return false, "", err
	}
	return true, out, nil
}

func (g *RealGit) Clone(name string) (string, error) {
	out, err := g.shell.Cmd("git", "clone", name)
	if err != nil {
		return "", err
	}
	return out, nil
}

func (g *RealGit) GetRemoteBranches() ([]string, error) {
	// prune
	g.shell.Cmd("git", "fetch", "--prune")

	// fetch
	g.shell.Cmd("git", "fetch", "origin")

	//get all branches
	return g.shell.ListCmd("git", "branch", "-r")
}

func (g *RealGit) GetLocalBranches() ([]string, error) {
	branches, err := g.shell.ListCmd("git", "branch")
	if err != nil {
		return nil, err
	}
	return branches, nil
}

func (g *RealGit) GetWorktrees() ([]string, error) {
	out, err := g.shell.ListCmd("git", "worktree", "list")
	if err != nil {
		log.Fatal(err)
	}
	return out, nil
}

func (g *RealGit) RemoveWorktree(worktreeName string) (string, error) {
	out, err := g.shell.Cmd("git", "worktree", "remove", worktreeName, "--force")
	if err != nil {
		log.Fatal(err)
	}
	return out, nil
}

func (g *RealGit) AddWorktree(folderName string, existsOnRemote bool,
	branchName string, baseBranch string) error {
	if existsOnRemote {
		_, err := g.shell.Cmd("git", "worktree", "add", folderName, branchName)
		return err
	} else {
		_, err := g.shell.Cmd("git", "worktree", "add", folderName, "-b", branchName, baseBranch)
		return err
	}
}

func (g *RealGit) GetRepoName(path string) (string, error) {
	out, err := g.shell.Cmd("git", "-C", path, "config", "--get", "remote.origin.url")
	if err != nil {
		return "", err
	}
	repoName := strings.TrimSuffix(filepath.Base(out), filepath.Ext(out))
	return repoName, nil
}

func (g *RealGit) FetchOrigin(branch string) error {
	_, err := g.shell.Cmd("git", "fetch", "origin", branch)
	if err != nil {
		return err
	}
	return nil
}

func (g *RealGit) CloneBare(url string, folderName string) error {
	_, err := g.shell.Cmd("git", "clone", "--bare", url, folderName)
	if err != nil {
		return err
	}
	return nil

}
