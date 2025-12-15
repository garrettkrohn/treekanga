package git

import (
	"fmt"
	"os"
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
	GetWorktrees(path *string) ([]string, error)
	RemoveWorktree(worktreeName string, path *string) (string, error)
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

func (g *RealGit) GetWorktrees(path *string) ([]string, error) {
	gitCmd := getBaseArguementsWithOrWithoutPath(path)
	gitCmd = append(gitCmd, "worktree", "list")
	out, err := g.shell.ListCmd("git", gitCmd...)
	if err != nil {
		log.Fatal(err)
	}
	return out, nil
}

func (g *RealGit) RemoveWorktree(worktreeName string, path *string) (string, error) {
	// When using a bare repo, convert absolute worktree path to relative path
	worktreePath := worktreeName
	if path != nil && filepath.IsAbs(worktreeName) {
		// Fix the .git file if it points to the wrong location (from old bare repo)
		err := g.fixWorktreeGitFile(worktreeName, *path)
		if err != nil {
			log.Debug("Could not fix .git file", "error", err)
		}
		
		relPath, err := filepath.Rel(*path, worktreeName)
		if err == nil {
			worktreePath = relPath
			log.Debug("Using relative path for worktree removal", "absolute", worktreeName, "relative", relPath)
		}
	}
	
	gitCmd := getBaseArguementsWithOrWithoutPath(path)
	gitCmd = append(gitCmd, "worktree", "remove", worktreePath, "--force")
	log.Debug("git args", "args", gitCmd)
	out, err := g.shell.Cmd("git", gitCmd...)
	if err != nil {
		return "", fmt.Errorf("failed to remove worktree %s: %w", worktreeName, err)
	}
	return out, nil
}

// fixWorktreeGitFile fixes the .git file in a worktree to point to the correct bare repo location
func (g *RealGit) fixWorktreeGitFile(worktreePath string, bareRepoPath string) error {
	gitFilePath := filepath.Join(worktreePath, ".git")
	worktreeName := filepath.Base(worktreePath)
	expectedGitDir := filepath.Join(bareRepoPath, "worktrees", worktreeName)
	
	// Write the corrected .git file using Go's file I/O
	content := fmt.Sprintf("gitdir: %s\n", expectedGitDir)
	err := os.WriteFile(gitFilePath, []byte(content), 0644)
	if err != nil {
		log.Debug("Failed to fix .git file", "error", err, "path", gitFilePath)
		return err
	}
	log.Debug("Fixed .git file", "path", gitFilePath, "points to", expectedGitDir)
	return nil
}

func (g *RealGit) AddWorktree(c *com.AddConfig) error {
	// Build base command
	gitCommand := getBaseArguementsWithOrWithoutPath(c.Flags.Directory)
	gitCommand = append(gitCommand, "worktree", "add", c.GetWorktreePath())

	// Add branch-specific arguments
	branchArgs := g.determineBranchArguments(c)
	gitCommand = append(gitCommand, branchArgs...)

	// Log the full command for debugging
	fullCommand := strings.Join(append([]string{"git"}, gitCommand...), " ")
	log.Debug("Executing git worktree command", "command", fullCommand)

	output, err := g.shell.Cmd("git", gitCommand...)
	if err != nil {
		return fmt.Errorf("failed to add worktree: %v\nCommand: %s\nOutput: %s", err, fullCommand, output)
	}

	return nil
}

func (g *RealGit) determineBranchArguments(c *com.AddConfig) []string {
	// Case 1: Branch already exists (locally or remotely) - just checkout
	if c.GitInfo.NewBranchExistsLocally || c.GitInfo.NewBranchExistsRemotely {
		return []string{c.GetNewBranchName()}
	}

	// Case 2: Base branch exists locally
	if c.GitInfo.BaseBranchExistsLocally {
		if c.ShouldPull() || c.AutoPull {
			// Create new branch from remote version of base branch
			return []string{"-b", c.GetNewBranchName(), "origin/" + c.GetBaseBranchName(), "--no-track"}
		} else {
			// Create new branch from local version of base branch
			return []string{"-b", c.GetNewBranchName(), c.GetBaseBranchName()}
		}
	}

	// Case 3: Base branch only exists remotely
	return []string{"-b", c.GetNewBranchName(), "origin/" + c.GetBaseBranchName(), "--no-track"}
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

// NOTE: I this can be removed
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
	
	// After configuring the fetch refspec, we need to fetch to populate remote-tracking branches
	// In a bare clone, branches are initially in refs/heads/, but we want them in refs/remotes/origin/
	_, err = g.shell.Cmd("git", "-C", path, "fetch", "origin")
	if err != nil {
		log.Debug("Warning: fetch after bare config failed", "error", err)
		// Don't return error as the repo might still be usable
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
