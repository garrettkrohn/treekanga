package connector

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/garrettkrohn/treekanga/adapters"
	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/models"
	"github.com/garrettkrohn/treekanga/shell"
	"github.com/garrettkrohn/treekanga/transformer"
	"github.com/garrettkrohn/treekanga/utility"
)

type Connector interface {
	Connect(name string, opts models.ConnectOpts) error
	VsCodeConnect(newRootPath string)
	CursorConnect(newRootPath string)
}

type RealConnector struct {
	shell shell.Shell
	tmux  adapters.Tmux
}

func NewConnector(shell shell.Shell) Connector {
	return &RealConnector{
		shell: shell,
		tmux:  adapters.NewTmux(shell),
	}
}

// Connect attempts to connect to a session using various strategies
func (r *RealConnector) Connect(name string, opts models.ConnectOpts) error {
	strategies := []func(string) (models.Connection, error){
		r.tmuxStrategy,
		r.worktreeStrategy,
		r.dirStrategy,
	}

	for _, strategy := range strategies {
		connection, err := strategy(name)
		if err != nil {
			return fmt.Errorf("connection strategy error: %w", err)
		}
		if connection.Found {
			return r.connectToTmux(connection, opts)
		}
	}

	return fmt.Errorf("no connection found for '%s'", name)
}

// tmuxStrategy checks if a tmux session with the given name exists
func (r *RealConnector) tmuxStrategy(name string) (models.Connection, error) {
	session, exists := r.tmux.FindSession(name)
	if !exists {
		return models.Connection{Found: false}, nil
	}
	return models.Connection{
		Found:   true,
		Session: session,
		New:     false,
	}, nil
}

// worktreeStrategy checks if the name matches a worktree path
func (r *RealConnector) worktreeStrategy(name string) (models.Connection, error) {
	// Try to get bare repo path
	bareRepoPath, err := git.GetBareRepoPath("")
	if err != nil {
		// Not in a git repo, skip this strategy
		return models.Connection{Found: false}, nil
	}

	worktrees, err := git.ListWorktrees(bareRepoPath)
	if err != nil {
		return models.Connection{Found: false}, nil
	}

	// Parse worktrees and check if name matches any worktree path or name
	worktreeObjects := transformer.TransformWorktrees(worktrees)
	for _, wt := range worktreeObjects {
		// Check if name matches the full path or the directory name
		if wt.FullPath == name || wt.Folder == name {
			sessionName := r.generateSessionName(wt.FullPath)
			return models.Connection{
				Found: true,
				New:   true,
				Session: models.Session{
					Name: sessionName,
					Path: wt.FullPath,
					Src:  "worktree",
				},
			}, nil
		}
	}

	return models.Connection{Found: false}, nil
}

// dirStrategy checks if the name is a valid directory path
func (r *RealConnector) dirStrategy(name string) (models.Connection, error) {
	// Expand home directory if needed
	path := name
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return models.Connection{Found: false}, nil
		}
		path = filepath.Join(homeDir, strings.TrimPrefix(path, "~/"))
	}

	// Check if it's an absolute path
	if !filepath.IsAbs(path) {
		// Try to make it absolute
		absPath, err := filepath.Abs(path)
		if err != nil {
			return models.Connection{Found: false}, nil
		}
		path = absPath
	}

	// Check if directory exists
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return models.Connection{Found: false}, nil
	}

	sessionName := r.generateSessionName(path)
	return models.Connection{
		Found: true,
		New:   true,
		Session: models.Session{
			Name: sessionName,
			Path: path,
			Src:  "dir",
		},
	}, nil
}

// generateSessionName creates a valid tmux session name from a path
func (r *RealConnector) generateSessionName(path string) string {
	// Use the basename of the path
	name := filepath.Base(path)
	// Replace dots and other special characters with underscores
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ReplaceAll(name, " ", "_")
	return name
}

// connectToTmux handles the actual connection to tmux
func (r *RealConnector) connectToTmux(connection models.Connection, opts models.ConnectOpts) error {
	if connection.New {
		// Create new session
		if err := r.tmux.NewSession(connection.Session.Name, connection.Session.Path); err != nil {
			return fmt.Errorf("failed to create tmux session: %w", err)
		}
	}

	// Switch or attach to the session
	return r.tmux.SwitchOrAttach(connection.Session.Name, opts)
}

func (r *RealConnector) VsCodeConnect(newRootPath string) {
	_, err := r.shell.Cmd("code", newRootPath)
	utility.CheckError(err)
}

func (r *RealConnector) CursorConnect(newRootPath string) {
	_, err := r.shell.Cmd("cursor", newRootPath)
	utility.CheckError(err)
}
