package git

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"

	com "github.com/garrettkrohn/treekanga/common"
	"github.com/garrettkrohn/treekanga/shell"
)

const tempZoxideName = "temp_treekanga_worktree"

// TODO: make a a function to add the directory
type Git interface {
	GetRemoteBranches(*string) ([]string, error)
	GetLocalBranches(*string) ([]string, error)
	GetWorktrees() ([]string, error)
	RemoveWorktree(string) (string, error)
	AddWorktree(c *com.AddConfig) error
	GetRepoName(path string) (string, error)
	CloneBare(string, string) error
	DeleteBranchRef(branch string, path string) error
	ConfigureGitBare(path string) error
}

type RealGit struct {
	shell shell.Shell
}

func NewGit(shell shell.Shell) Git {
	return &RealGit{shell}
}

func (g *RealGit) GetRemoteBranches(path *string) ([]string, error) {
	// prune
	fetchCmd := getBaseArguementsWithOrWithoutPath(path)
	fetchCmd = append(fetchCmd, "fetch", "--prune")
	g.shell.Cmd("git", fetchCmd...)

	// fetch
	fetchCmd2 := getBaseArguementsWithOrWithoutPath(path)
	fetchCmd2 = append(fetchCmd2, "fetch", "origin")
	g.shell.Cmd("git", fetchCmd2...)

	//get all branches
	branchCmd := getBaseArguementsWithOrWithoutPath(path)
	branchCmd = append(branchCmd, "branch", "-r", "--format=%(refname:short)")
	list, err := g.shell.ListCmd("git", branchCmd...)
	return list, err
}

func (g *RealGit) GetLocalBranches(path *string) ([]string, error) {
	gitCmd := getBaseArguementsWithOrWithoutPath(path)
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

func (g *RealGit) AddWorktree(c *com.AddConfig) error {
	gitCommand := getBaseArguementsWithOrWithoutPath(c.Flags.Directory)
	// Use the WorktreeTargetDir directly as it already contains the full path
	worktreePath := c.GetWorktreePath()
	gitCommand = append(gitCommand, "worktree", "add", worktreePath)

	// create worktree off of local branch
	if c.GitInfo.NewBranchExistsLocally || c.GitInfo.NewBranchExistsRemotely {
		gitCommand = append(gitCommand, c.GetNewBranchName())
	} else if c.GitInfo.BaseBranchExistsLocally {
		if c.ShouldPull() {
			gitCommand = append(gitCommand, "-b", c.GetNewBranchName(), "origin/"+c.GetBaseBranchName(), "--no-track")
		} else {
			gitCommand = append(gitCommand, "-b", c.GetNewBranchName(), c.GetBaseBranchName())

		}
	} else {
		gitCommand = append(gitCommand, "-b", c.GetNewBranchName(), "origin/"+c.GetBaseBranchName(), "--no-track")
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

func (g *RealGit) CloneBare(url string, folderName string) error {
	_, err := g.shell.Cmd("git", "clone", "--bare", url, folderName)
	if err != nil {
		return err
	}
	return nil
}

func (g *RealGit) CreateTempBranch(path string) error {
	gitCmd := getBaseArguementsWithOrWithoutPath(&path)
	gitCmd = append(gitCmd, "branch", tempZoxideName, "FETCH_HEAD")
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

func (g *RealGit) ConfigureGitBare(path string) error {
	//This is really ugly, but necessary to set up the bare repo correctly.  The issue was trying
	//to get the shell to keep the "" around the +refs...
	_, err := g.shell.Cmd("sh", "-c", fmt.Sprintf(`git -C %s config remote.origin.fetch "+refs/heads/*:refs/remotes/origin/*"`, path))
	if err != nil {
		return err
	}
	return nil
}

func getBaseArguementsWithOrWithoutPath(path *string) []string {
	gitCommand := make([]string, 0)

	if path != nil {
		gitCommand = append(gitCommand, "-C", *path)
	}

	return gitCommand
}
