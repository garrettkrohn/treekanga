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
	AddWorktree(string, bool, bool, string, string, string, bool, bool, bool) error
	GetRepoName(path string) (string, error)
	FetchOrigin(branch string, path string) error
	CloneBare(string, string) error
	PullBranch(url string) error
	CreateTempBranch(path string) error
	DeleteBranch(branch string, path string) error
	DeleteBranchRef(branch string, path string) error
	ConfigureGitBare() error
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
	branchCmd = append(branchCmd, "branch", "-r", "--format=%(refname:short)")
	list, err := g.shell.ListCmd("git", branchCmd...)
	return list, err
}

func (g *RealGit) GetLocalBranches(path string) ([]string, error) {
	gitCmd := getBaseCommandWithOrWithoutPath(path)
	gitCmd = append(gitCmd, "branch", "--format=%(refname:short)")
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

func (g *RealGit) AddWorktree(folderName string,
	newBranchExistsLocally bool, newBranchExistRemotely bool,
	branchName string, baseBranch string, path string, pull bool,
	baseBranchExistsLocally bool, baseBrachExistsRemotely bool) error {
	gitCommand := getBaseCommandWithOrWithoutPath(path)
	gitCommand = append(gitCommand, "worktree", "add", folderName)

	// create worktree off of local branch
	if newBranchExistsLocally || newBranchExistRemotely {
		gitCommand = append(gitCommand, branchName)
	} else if baseBranchExistsLocally {
		if pull {
			gitCommand = append(gitCommand, "-b", branchName, "origin/"+baseBranch, "--no-track")
		} else {
			gitCommand = append(gitCommand, "-b", branchName, baseBranch)

		}
	} else {
		gitCommand = append(gitCommand, "-b", branchName, "origin/"+baseBranch, "--no-track")
	}

	output, err := g.shell.Cmd("git", gitCommand...)

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

func (g *RealGit) DeleteBranchRef(branch string, path string) error {
	gitCmd := fmt.Sprintf("%s/refs/heads/%s", path, branch)
	_, err := g.shell.Cmd("update-ref", "-d", gitCmd)
	if err != nil {
		return err
	}
	return nil
}

func (g *RealGit) ConfigureGitBare() error {
	_, err := g.shell.Cmd("git", "config", "remote.origin.fetch", "+refs/heads/*:refs/remotes/origin/*")
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
