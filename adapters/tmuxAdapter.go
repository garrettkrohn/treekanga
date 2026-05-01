package adapters

import (
	"os"
	"strings"

	"github.com/garrettkrohn/treekanga/models"
	"github.com/garrettkrohn/treekanga/shell"
)

type Tmux interface {
	ListSessions() ([]models.Session, error)
	NewSession(sessionName string, startDir string) error
	IsAttached() bool
	AttachSession(targetSession string) error
	SwitchClient(targetSession string) error
	SwitchOrAttach(name string, opts models.ConnectOpts) error
	FindSession(name string) (models.Session, bool)
	KillSession(sessionName string) error
	GetCurrentSessionName() (string, error)
}

type RealTmux struct {
	shell shell.Shell
}

func NewTmux(shell shell.Shell) Tmux {
	return &RealTmux{shell}
}

func (t *RealTmux) ListSessions() ([]models.Session, error) {
	output, err := t.shell.Cmd("tmux", "list-sessions", "-F", "#{session_name}:#{session_path}")
	if err != nil {
		// If tmux server isn't running, return empty list
		if strings.Contains(err.Error(), "no server running") {
			return []models.Session{}, nil
		}
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	sessions := make([]models.Session, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			sessions = append(sessions, models.Session{
				Name: parts[0],
				Path: parts[1],
				Src:  "tmux",
			})
		}
	}

	return sessions, nil
}

func (t *RealTmux) FindSession(name string) (models.Session, bool) {
	sessions, err := t.ListSessions()
	if err != nil {
		return models.Session{}, false
	}

	for _, session := range sessions {
		if session.Name == name {
			return session, true
		}
	}

	return models.Session{}, false
}

func (t *RealTmux) NewSession(sessionName string, startDir string) error {
	_, err := t.shell.Cmd("tmux", "new-session", "-d", "-s", sessionName, "-c", startDir)
	return err
}

func (t *RealTmux) AttachSession(targetSession string) error {
	_, err := t.shell.Cmd("tmux", "attach-session", "-t", targetSession)
	return err
}

func (t *RealTmux) SwitchClient(targetSession string) error {
	_, err := t.shell.Cmd("tmux", "switch-client", "-t", targetSession)
	return err
}

func (t *RealTmux) IsAttached() bool {
	return len(os.Getenv("TMUX")) > 0
}

func (t *RealTmux) SwitchOrAttach(name string, opts models.ConnectOpts) error {
	if opts.Switch || t.IsAttached() {
		return t.SwitchClient(name)
	}
	return t.AttachSession(name)
}

func (t *RealTmux) KillSession(sessionName string) error {
	_, err := t.shell.Cmd("tmux", "kill-session", "-t", sessionName)
	return err
}

func (t *RealTmux) GetCurrentSessionName() (string, error) {
	output, err := t.shell.Cmd("tmux", "display-message", "-p", "#{session_name}")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}
