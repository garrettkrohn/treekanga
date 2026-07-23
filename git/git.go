package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

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

// SetUpstream configures the upstream branch using git config for the current branch in a worktree
func SetUpstream(worktreePath, branchName string) error {
	// Set remote for the branch
	remoteArgs := []string{"-C", worktreePath, "config", "branch." + branchName + ".remote", "origin"}
	err := runCommand("git", remoteArgs...)
	if err != nil {
		log.Debug("Failed to set remote for branch", "branch", branchName, "error", err)
		return err
	}

	// Set merge ref for the branch
	mergeArgs := []string{"-C", worktreePath, "config", "branch." + branchName + ".merge", "refs/heads/" + branchName}
	err = runCommand("git", mergeArgs...)
	if err != nil {
		log.Debug("Failed to set merge ref for branch", "branch", branchName, "error", err)
		return err
	}

	log.Debug("Set upstream config for branch", "branch", branchName, "remote", "origin", "merge", "refs/heads/"+branchName)
	return nil
}

// UnsetUpstream removes the upstream tracking config for a branch in a worktree
func UnsetUpstream(worktreePath, branchName string) error {
	for _, key := range []string{"branch." + branchName + ".remote", "branch." + branchName + ".merge"} {
		err := runCommand("git", "-C", worktreePath, "config", "--unset", key)
		if err != nil {
			// exit code 5 means the key didn't exist — not an error
			if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 5 {
				continue
			}
			return fmt.Errorf("failed to unset %s: %w", key, err)
		}
	}
	log.Debug("Unset upstream config for branch", "branch", branchName)
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

// Fetch fetches updates for a specific branch from remote
func Fetch(bareRepoPath, branch string) error {
	args := []string{"-C", bareRepoPath, "fetch", "origin", branch}
	err := runCommand("git", args...)
	if err != nil {
		return fmt.Errorf("failed to fetch branch %s: %w", branch, err)
	}
	log.Debug("Fetched latest state from remote", "branch", branch)
	return nil
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
// If forceSubmodules is true, manually moves the directory and updates git's worktree config
func MoveWorktree(bareRepoPath, oldPath, newPath string, forceSubmodules bool) error {
	// If forceSubmodules is enabled, skip the git command and go straight to manual move
	if forceSubmodules {
		log.Info("Force submodules enabled, using manual move workaround")
		return moveWorktreeManually(bareRepoPath, oldPath, newPath)
	}

	args := []string{"-C", bareRepoPath, "worktree", "move", oldPath, newPath}
	err := runCommand("git", args...)
	if err != nil {
		return fmt.Errorf("failed to move worktree from %s to %s: %w", oldPath, newPath, err)
	}
	return nil
}

// moveWorktreeManually manually moves a worktree directory and updates git's internal references
// This is a workaround for git's limitation with submodules
func moveWorktreeManually(bareRepoPath, oldPath, newPath string) error {
	// Step 1: Move the directory
	log.Debug("Manually moving directory", "from", oldPath, "to", newPath)
	err := os.Rename(oldPath, newPath)
	if err != nil {
		return fmt.Errorf("failed to move directory: %w", err)
	}

	// Step 2: Update the worktree's gitdir file to point to the correct location
	gitdirPath := filepath.Join(newPath, ".git")
	gitdirContent, err := os.ReadFile(gitdirPath)
	if err != nil {
		// Try to rollback the directory move
		os.Rename(newPath, oldPath)
		return fmt.Errorf("failed to read .git file: %w", err)
	}

	// The .git file contains: "gitdir: /path/to/bare/.git/worktrees/name"
	// We need to extract the worktree name and verify it's correct
	gitdirLine := strings.TrimSpace(string(gitdirContent))
	if !strings.HasPrefix(gitdirLine, "gitdir: ") {
		os.Rename(newPath, oldPath)
		return fmt.Errorf("unexpected .git file format: %s", gitdirLine)
	}

	worktreeGitDir := strings.TrimPrefix(gitdirLine, "gitdir: ")

	// Step 3: Update the worktree's gitdir location in the bare repo
	gitdirLocationFile := filepath.Join(worktreeGitDir, "gitdir")
	log.Debug("Updating gitdir location", "file", gitdirLocationFile, "newPath", newPath)

	err = os.WriteFile(gitdirLocationFile, []byte(filepath.Join(newPath, ".git")), 0644)
	if err != nil {
		// Try to rollback
		os.Rename(newPath, oldPath)
		return fmt.Errorf("failed to update gitdir location: %w", err)
	}

	log.Info("Successfully moved worktree manually", "from", oldPath, "to", newPath)
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

// GetWorkingTreeStatus reports whether a worktree has staged, modified
// (unstaged), or untracked changes.
func GetWorkingTreeStatus(worktreePath string) (staged, modified, untracked bool, err error) {
	output, err := runCommandOutput("git", "-C", worktreePath, "status", "--porcelain=v1", "--untracked-files=all")
	if err != nil {
		return false, false, false, fmt.Errorf("failed to get status for %s: %w", worktreePath, err)
	}

	if output == "" {
		return false, false, false, nil
	}

	for _, line := range strings.Split(output, "\n") {
		if len(line) < 2 {
			continue
		}
		x, y := line[0], line[1]
		if x == '?' && y == '?' {
			untracked = true
			continue
		}
		if x != ' ' {
			staged = true
		}
		if y != ' ' {
			modified = true
		}
	}

	return staged, modified, untracked, nil
}

// GetAheadBehind returns how many commits HEAD is ahead/behind compareRef
// in a given worktree.
func GetAheadBehind(worktreePath, compareRef string) (ahead, behind int, err error) {
	output, err := runCommandOutput("git", "-C", worktreePath, "rev-list", "--left-right", "--count", "HEAD..."+compareRef)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to compute ahead/behind for %s against %s: %w", worktreePath, compareRef, err)
	}

	fields := strings.Fields(output)
	if len(fields) != 2 {
		return 0, 0, fmt.Errorf("unexpected rev-list output for %s: %q", worktreePath, output)
	}

	ahead, err = strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse ahead count: %w", err)
	}
	behind, err = strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse behind count: %w", err)
	}

	return ahead, behind, nil
}

// GetUpstreamBranch returns the upstream (remote-tracking) branch for a
// worktree's current branch, or "" if no upstream is configured.
func GetUpstreamBranch(worktreePath string) (string, error) {
	output, err := runCommandOutput("git", "-C", worktreePath, "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	if err != nil {
		// No upstream configured - not an error condition for callers.
		return "", nil
	}
	return strings.TrimSpace(output), nil
}

// IsMerged reports whether branchName's content is already present in
// targetRef: either as a literal ancestor, or via a squash-merge content
// match (the branch's aggregate diff since its merge-base matches the
// patch-id of some commit in targetRef since that same merge-base).
func IsMerged(worktreePath, branchName, targetRef string) (bool, error) {
	isAncestor, err := isAncestor(worktreePath, branchName, targetRef)
	if err != nil {
		return false, err
	}
	if isAncestor {
		return true, nil
	}

	base, err := mergeBase(worktreePath, branchName, targetRef)
	if err != nil {
		// No common ancestor - branch and target share no history.
		return false, nil
	}

	branchPatchID, err := patchID(worktreePath, base, branchName)
	if err != nil {
		return false, err
	}
	if branchPatchID == "" {
		// Branch has no changes relative to the merge-base.
		return false, nil
	}

	commits, err := commitsSince(worktreePath, base, targetRef)
	if err != nil {
		return false, err
	}

	for _, commit := range commits {
		pid, err := patchID(worktreePath, commit+"^", commit)
		if err != nil {
			continue
		}
		if pid != "" && pid == branchPatchID {
			return true, nil
		}
	}

	return false, nil
}

func isAncestor(worktreePath, ancestorRef, ref string) (bool, error) {
	command := exec.Command("git", "-C", worktreePath, "merge-base", "--is-ancestor", ancestorRef, ref)
	command.SysProcAttr = setSysProcAttr()
	err := command.Run()
	if err == nil {
		return true, nil
	}
	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
		return false, nil
	}
	return false, fmt.Errorf("failed to check ancestry of %s in %s: %w", ancestorRef, ref, err)
}

func mergeBase(worktreePath, refA, refB string) (string, error) {
	output, err := runCommandOutput("git", "-C", worktreePath, "merge-base", refA, refB)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

func commitsSince(worktreePath, base, ref string) ([]string, error) {
	output, err := runCommandOutput("git", "-C", worktreePath, "log", "--format=%H", base+".."+ref)
	if err != nil {
		return nil, fmt.Errorf("failed to list commits between %s and %s: %w", base, ref, err)
	}
	if output == "" {
		return nil, nil
	}
	return strings.Split(output, "\n"), nil
}

// patchID returns the stable patch-id for the diff between fromRef and
// toRef, or "" if the diff is empty.
func patchID(worktreePath, fromRef, toRef string) (string, error) {
	diffCmd := exec.Command("git", "-C", worktreePath, "diff", fromRef, toRef)
	diffCmd.SysProcAttr = setSysProcAttr()
	diffOutput, err := diffCmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to diff %s..%s: %w", fromRef, toRef, err)
	}
	if len(strings.TrimSpace(string(diffOutput))) == 0 {
		return "", nil
	}

	patchIDCmd := exec.Command("git", "-C", worktreePath, "patch-id", "--stable")
	patchIDCmd.SysProcAttr = setSysProcAttr()
	patchIDCmd.Stdin = strings.NewReader(string(diffOutput))
	patchIDOutput, err := patchIDCmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to compute patch-id: %w", err)
	}

	fields := strings.Fields(string(patchIDOutput))
	if len(fields) == 0 {
		return "", nil
	}
	return fields[0], nil
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
	command.SysProcAttr = setSysProcAttr()

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

	command.SysProcAttr = setSysProcAttr()

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
