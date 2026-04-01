package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/charmbracelet/log"
)

// AddWorktree creates a new worktree with flexible arguments
func AddWorktree(bareRepoPath, worktreeTargetDir, worktreeName string, worktreeArgs []string) error {
	args := []string{"-C", bareRepoPath, "worktree", "add", filepath.Join(worktreeTargetDir, worktreeName)}
	args = append(args, worktreeArgs...)

	fullCommand := strings.Join(append([]string{"git"}, args...), " ")
	log.Debug("Executing git worktree command", "command", fullCommand)

	err := runCommand("git", args...)
	if err != nil {
		return fmt.Errorf("failed to add worktree: %v\nCommand: %s", err, fullCommand)
	}

	return nil
}

// RemoveWorktree removes a worktree (worktreePath can be name or path)
func RemoveWorktree(bareRepoPath, worktreePath string, force bool) error {
	args := []string{"-C", bareRepoPath, "worktree", "remove", worktreePath}
	if force {
		args = append(args, "--force")
	}

	err := runCommand("git", args...)
	if err != nil {
		log.Debug(fmt.Errorf("failed to remove worktree %s: %w", worktreePath, err))
		return err
	}
	return nil
}

// ListWorktrees returns raw worktree list output
func ListWorktrees(bareRepoPath string) ([]string, error) {
	args := []string{"-C", bareRepoPath, "worktree", "list"}
	output, err := runCommandOutput("git", args...)
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(output), "\n"), nil
}

// GetRemoteBranches lists remote branches (without fetching)
func GetRemoteBranches(bareRepoPath string) ([]string, error) {
	args := []string{"-C", bareRepoPath, "branch", "-r", "--format=%(refname:short)"}
	output, err := runCommandOutput("git", args...)
	if err != nil {
		return nil, err
	}
	branches := strings.Split(strings.TrimSpace(output), "\n")

	// Remove "origin/" prefix
	cleaned := make([]string, 0, len(branches))
	for _, branch := range branches {
		if branch != "" && branch != "origin" {
			cleaned = append(cleaned, strings.TrimPrefix(branch, "origin/"))
		}
	}
	return cleaned, nil
}

// GetLocalBranches lists local branches
func GetLocalBranches(bareRepoPath string) ([]string, error) {
	args := []string{"-C", bareRepoPath, "branch", "--format=%(refname:short)"}
	output, err := runCommandOutput("git", args...)
	if err != nil {
		return nil, err
	}
	branches := strings.Split(strings.TrimSpace(output), "\n")

	// Remove quotes if present
	cleaned := make([]string, 0, len(branches))
	for _, branch := range branches {
		branch = strings.Trim(branch, "'\"")
		if branch != "" {
			cleaned = append(cleaned, branch)
		}
	}
	return cleaned, nil
}

// DeleteBranch deletes a local branch
func DeleteBranch(bareRepoPath, branch string, force bool) error {
	args := []string{"-C", bareRepoPath, "branch"}
	if force {
		args = append(args, "-D", branch)
	} else {
		args = append(args, "-d", branch)
	}
	return runCommand("git", args...)
}

// RenameBranch renames a local branch
func RenameBranch(bareRepoPath, oldName, newName string) error {
	args := []string{"-C", bareRepoPath, "branch", "-m", oldName, newName}
	err := runCommand("git", args...)
	if err != nil {
		return fmt.Errorf("failed to rename branch from %s to %s: %w", oldName, newName, err)
	}
	return nil
}

// MoveWorktree moves a worktree to a new location
func MoveWorktree(bareRepoPath, oldPath, newPath string) error {
	args := []string{"-C", bareRepoPath, "worktree", "move", oldPath, newPath}
	err := runCommand("git", args...)
	if err != nil {
		return fmt.Errorf("failed to move worktree from %s to %s: %w", oldPath, newPath, err)
	}
	return nil
}

// GetCurrentBranch returns the current branch name for a given directory
func GetCurrentBranch(dir string) (string, error) {
	command := exec.Command("git", "branch", "--show-current")
	command.Dir = dir
	output, err := command.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	branchName := strings.TrimSpace(string(output))
	return branchName, nil
}

// CloneBare clones a repository as bare
func CloneBare(url, folderName string) error {
	return runCommand("git", "clone", "--progress", "--bare", url, folderName)
}

// ConfigureBare configures a bare repository for worktree usage
func ConfigureBare(bareRepoPath string) error {
	_, err := runCommandOutput("git", "-C", bareRepoPath, "config", "remote.origin.fetch", "+refs/heads/*:refs/remotes/origin/*")
	if err != nil {
		return err
	}

	// Fetch to populate remote-tracking branches
	_, err = runCommandOutput("git", "-C", bareRepoPath, "fetch", "origin")
	if err != nil {
		log.Debug("Warning: fetch after bare config failed", "error", err)
	}

	return nil
}

// GetBareRepoPath returns the path to the bare repository
func GetBareRepoPath(dir string) (string, error) {
	if dir != "" {
		command := exec.Command("git", "rev-parse", "--git-common-dir")
		command.Dir = dir
		output, err := command.Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(output)), nil
	}
	return runCommandOutput("git", "rev-parse", "--git-common-dir")
}

// GetProjectName returns the project name from git config
func GetProjectName() (string, error) {
	url, err := runCommandOutput("git", "config", "--get", "remote.origin.url")
	if err != nil {
		return "", err
	}

	output, err := runCommandOutput("basename", "-s", ".git", strings.TrimSpace(url))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

// Helper functions

func runCommand(cmd string, args ...string) error {
	log.Debug(cmd, "args", args)

	command := exec.Command(cmd, args...)
	command.Stdin = nil

	// Set environment to prevent git from using pagers or editors
	command.Env = append(os.Environ(),
		"GIT_PAGER=cat",
		"GIT_EDITOR=true",
		"EDITOR=true",
		"VISUAL=true",
	)

	// Create new process group
	command.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	output, err := command.CombinedOutput()

	// Log output
	if len(output) > 0 {
		for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
			if line != "" {
				log.Info(line)
			}
		}
	}

	return err
}

func runCommandOutput(cmd string, args ...string) (string, error) {
	log.Debug(cmd, "args", args)

	command := exec.Command(cmd, args...)
	command.Stdin = nil

	command.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	output, err := command.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			errString := strings.TrimSpace(string(exitErr.Stderr))
			if strings.HasPrefix(errString, "no server running on") {
				return "", nil
			}
		}
		return "", err
	}

	return strings.TrimSuffix(string(output), "\n"), nil
}
