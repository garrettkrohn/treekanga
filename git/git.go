package git

import (
	"fmt"
	"github.com/charmbracelet/log"
	"path/filepath"
	"strings"

	"github.com/garrettkrohn/treekanga/shell"
)

const tempZoxideName = "temp_treekanga_worktree"

// TODO: make a a function to add the directory
type Git interface {
	ShowTopLevel(name string) (bool, string, error)
	GitCommonDir(name string) (bool, string, error)
	Clone(name string) (string, error)
	GetRemoteBranches(string) ([]string, error)
	GetLocalBranches(string) ([]string, error)
	GetWorktrees() ([]string, error)
	RemoveWorktree(string) (string, error)
	AddWorktree(string, bool, string, string, string) error
	GetRepoName(path string) (string, error)
	FetchOrigin(branch string, path string) error
	CloneBare(string, string) error
	PullBranch(url string) error
	CreateTempBranch(path string) error
	DeleteBranch(branch string, path string) error
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

func (g *RealGit) GetRemoteBranches(path string) ([]string, error) {
	// prune
	fetchCmd := getBaseCommandWithOrWithoutPath(path)
	fetchCmd = append(fetchCmd, "fetch", "--prune")
	g.shell.Cmd("git", fetchCmd...)

	// fetch
	fetchCmd2 := getBaseCommandWithOrWithoutPath(path)
	fetchCmd2 = append(fetchCmd2, "fetch", "origin")
	g.shell.Cmd("git", fetchCmd2...)

	//get all branches
	branchCmd := getBaseCommandWithOrWithoutPath(path)
	branchCmd = append(branchCmd, "branch", "-r", "--format=\"%(refname:short)\"")
	list, err := g.shell.ListCmd("git", branchCmd...)
	return list, err
}

func (g *RealGit) GetLocalBranches(path string) ([]string, error) {
	gitCmd := getBaseCommandWithOrWithoutPath(path)
	gitCmd = append(gitCmd, "branch", "--format='%(refname:short)'")
	branches, err := g.shell.ListCmd("git", gitCmd...)
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

func (g *RealGit) AddWorktree(folderName string, existsLocally bool,
	branchName string, baseBranch string, path string) error {
	gitCommand := getBaseCommandWithOrWithoutPath(path)
	gitCommand = append(gitCommand, "worktree", "add", folderName)

	if existsLocally {
		gitCommand = append(gitCommand, branchName)
	} else {
		gitCommand = append(gitCommand, "-b", branchName, baseBranch)
	}

	output, err := g.shell.Cmd("git", gitCommand...)

	// var err error
	// var output string
	//
	// if existsLocally {
	// 	log.Debug("branch exists on remote")
	// 	output, err = g.shell.Cmd("git", "worktree", "add", folderName, branchName)
	// } else {
	// 	log.Debug("branch does not exist on remote")
	// 	output, err = g.shell.Cmd("git", "worktree", "add", folderName, "-b", branchName, baseBranch)
	// }

	if err != nil {
		return fmt.Errorf("failed to add worktree: %v, %s", err, output)
	}

	return nil
}

// Note: path is figured out in add.go
func (g *RealGit) GetRepoName(path string) (string, error) {
	out, err := g.shell.Cmd("git", "-C", path, "config", "--get", "remote.origin.url")
	if err != nil {
		return "", err
	}
	repoName := strings.TrimSuffix(filepath.Base(out), filepath.Ext(out))
	return repoName, nil
}

func (g *RealGit) FetchOrigin(branch string, path string) error {
	gitCmd := getBaseCommandWithOrWithoutPath(path)
	gitCmd = append(gitCmd, "fetch", "origin", branch)
	_, err := g.shell.Cmd("git", gitCmd...)
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

func (g *RealGit) PullBranch(url string) error {
	//TODO: need to make this configurable to be run from a worktree
	_, err := g.shell.Cmd("git", "-c", "/Users/gkrohn/code/platform_work/development", "pull", url)
	if err != nil {
		return err
	}
	return nil
}

func (g *RealGit) CreateTempBranch(path string) error {
	gitCmd := getBaseCommandWithOrWithoutPath(path)
	gitCmd = append(gitCmd, "branch", tempZoxideName, "FETCH_HEAD")
	_, err := g.shell.Cmd("git", gitCmd...)
	if err != nil {
		return err
	}
	return nil
}

func (g *RealGit) DeleteBranch(branch string, path string) error {
	gitCmd := getBaseCommandWithOrWithoutPath(path)
	gitCmd = append(gitCmd, "checkout", "-d", branch)
	_, err := g.shell.Cmd("git", gitCmd...)
	if err != nil {
		return err
	}
	return nil
}

func getBaseCommandWithOrWithoutPath(path string) []string {
	gitCommand := make([]string, 0)

	if path != "" {
		gitCommand = append(gitCommand, "-C", path)
	}

	return gitCommand
}
