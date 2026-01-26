package adapters

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"

	"github.com/garrettkrohn/treekanga/shell"
)

type GitAdapter interface {
	GetRemoteBranches(*string) ([]string, error)
	GetLocalBranches(*string) ([]string, error)
	GetWorktrees(path *string) ([]string, error)
	RemoveWorktree(worktreeName string, path *string, forceDelete bool) error
	AddWorktree(params AddWorktreeConfig) error
	GetRepoName(path string) (string, error)
	CloneBare(string, string) error
	DeleteBranchRef(branch string, path string) error
	DeleteBranch(branch string, path string, forceDelete bool) error
	ConfigureGitBare(path string) error
	GetBareRepoPath() (string, error)
	GetProjectName() (string, error)
}

type RealGitAdapter struct {
	shell shell.Shell
}

func NewGitAdapter(shell shell.Shell) GitAdapter {
	return &RealGitAdapter{shell}
}

func (g *RealGitAdapter) GetRemoteBranches(path *string) ([]string, error) {
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

func (g *RealGitAdapter) GetLocalBranches(path *string) ([]string, error) {
	gitCmd := getBaseArguementsWithOrWithoutPath(path)
	gitCmd = append(gitCmd, "branch", "--format='%(refname:short)'")
	branches, err := g.shell.ListCmd("git", gitCmd...)
	if err != nil {
		return nil, err
	}
	return branches, nil
}

func (g *RealGitAdapter) GetWorktrees(path *string) ([]string, error) {
	gitCmd := getBaseArguementsWithOrWithoutPath(path)
	gitCmd = append(gitCmd, "worktree", "list")
	out, err := g.shell.ListCmd("git", gitCmd...)
	if err != nil {
		log.Fatal(err)
	}
	return out, nil
}

func (g *RealGitAdapter) RemoveWorktree(worktreeName string, path *string, forceDelete bool) error {
	gitCmd := getBaseArguementsWithOrWithoutPath(path)
	gitCmd = append(gitCmd, "worktree", "remove", worktreeName)
	if forceDelete {
		gitCmd = append(gitCmd, "--force")
	}
	err := g.shell.CmdWithStreaming("git", gitCmd...)
	if err != nil {
		log.Debug(fmt.Errorf("failed to remove worktree %s: %w", worktreeName, err))
		return err
	}
	return nil
}

// fixWorktreeGitFile fixes the .git file in a worktree to point to the correct bare repo location
func (g *RealGitAdapter) fixWorktreeGitFile(worktreePath string, bareRepoPath string) error {
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

func (g *RealGitAdapter) AddWorktree(params AddWorktreeConfig) error {
	// Build base command
	gitCommand := getBaseArguementsWithOrWithoutPath(&params.BareRepoPath)
	gitCommand = append(gitCommand, "worktree", "add", params.WorktreeTargetDirectory+"/"+params.NewWorktreeName)

	// Add branch-specific arguments
	branchArgs := g.determineBranchArguments(params)
	gitCommand = append(gitCommand, branchArgs...)

	// Log the full command for debugging
	fullCommand := strings.Join(append([]string{"git"}, gitCommand...), " ")
	log.Debug("Executing git worktree command", "command", fullCommand)

	err := g.shell.CmdWithStreaming("git", gitCommand...)
	if err != nil {
		return fmt.Errorf("failed to add worktree: %v\nCommand: %s", err, fullCommand)
	}

	return nil
}

func (g *RealGitAdapter) determineBranchArguments(params AddWorktreeConfig) []string {
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

// Note: path is figured out in add.go
func (g *RealGitAdapter) GetRepoName(path string) (string, error) {
	out, err := g.shell.Cmd("git", "-C", path, "config", "--get", "remote.origin.url")
	if err != nil {
		return "", err
	}
	repoName := strings.TrimSuffix(filepath.Base(out), filepath.Ext(out))
	return repoName, nil
}

func (g *RealGitAdapter) CloneBare(url string, folderName string) error {
	err := g.shell.CmdWithStreaming("git", "clone", "--progress", "--bare", url, folderName)
	if err != nil {
		return err
	}
	return nil
}

// NOTE: I this can be removed
// func (g *RealGitAdapter) CreateTempBranch(path string) error {
// 	gitCmd := getBaseArguementsWithOrWithoutPath(&path)
// 	gitCmd = append(gitCmd, "branch", tempZoxideName, "FETCH_HEAD")
// 	_, err := g.shell.Cmd("git", gitCmd...)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (g *RealGitAdapter) DeleteBranchRef(branch string, path string) error {
	gitCmd := fmt.Sprintf("%s/refs/heads/%s", path, branch)
	err := g.shell.CmdWithStreaming("update-ref", "-d", gitCmd)
	if err != nil {
		return err
	}

	return nil
}

func (g *RealGitAdapter) DeleteBranch(branch string, path string, forceDelete bool) error {
	gitCmd := getBaseArguementsWithOrWithoutPath(&path)
	if forceDelete {
		gitCmd = append(gitCmd, "branch", "-D", branch)
	} else {
		gitCmd = append(gitCmd, "branch", "-d", branch)

	}
	err := g.shell.CmdWithStreaming("git", gitCmd...)
	if err != nil {
		return err
	}

	return nil
}

func (g *RealGitAdapter) ConfigureGitBare(path string) error {
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

func (g *RealGitAdapter) GetBareRepoPath() (string, error) {
	return g.shell.Cmd("git", "rev-parse", "--git-common-dir")
}

func (g *RealGitAdapter) GetProjectName() (string, error) {
	// First get the remote URL
	url, err := g.shell.Cmd("git", "config", "--get", "remote.origin.url")
	if err != nil {
		return "", err
	}

	// Then extract the basename
	return g.shell.Cmd("basename", "-s", ".git", strings.TrimSpace(url))
}

func getBaseArguementsWithOrWithoutPath(path *string) []string {
	gitCommand := make([]string, 0)

	if path != nil {
		gitCommand = append(gitCommand, "-C", *path)
	}

	return gitCommand
}
