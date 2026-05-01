package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/charmbracelet/log"
)

// Connect creates or connects to a tmux session at the given path
func Connect(sessionPath string) error {
	// Use a unique session name based on the path
	// For /Users/gkrohn/code/platform_work/test_zox/parent -> test_zox-parent
	// For /Users/gkrohn/code/platform_work/test_zox -> test_zox
	parts := strings.Split(strings.TrimSuffix(sessionPath, "/"), "/")
	var sessionName string
	if len(parts) >= 2 {
		// Use last two parts to make it unique
		sessionName = parts[len(parts)-2] + "-" + parts[len(parts)-1]
	} else {
		sessionName = filepath.Base(sessionPath)
	}
	
	log.Debug("Attempting to connect to tmux session", "name", sessionName, "path", sessionPath)
	
	// Check if we're inside tmux
	insideTmux := os.Getenv("TMUX") != ""
	
	// Always try to create the session - it will fail fast if it already exists
	log.Debug("Creating session (will skip if exists)")
	createSession(sessionName, sessionPath) // Ignore errors - session might already exist
	
	log.Info("Connecting to tmux session", "name", sessionName, "path", sessionPath)
	
	// Check if we're inside tmux
	if insideTmux {
		// Switch client
		switchCmd := exec.Command("tmux", "switch-client", "-t", sessionName)
		if err := switchCmd.Run(); err != nil {
			return fmt.Errorf("failed to switch tmux client: %w", err)
		}
		return nil
	}
	
	// Not inside tmux - use syscall.Exec to replace process
	tmuxPath, err := exec.LookPath("tmux")
	if err != nil {
		return fmt.Errorf("could not find tmux binary: %w", err)
	}
	
	log.Debug("Replacing process with tmux attach", "session", sessionName)
	
	args := []string{"tmux", "attach-session", "-t", sessionName}
	err = syscall.Exec(tmuxPath, args, os.Environ())
	if err != nil {
		return fmt.Errorf("failed to exec tmux: %w", err)
	}
	
	return nil
}

func sessionExists(sessionName string) bool {
	checkCmd := exec.Command("tmux", "has-session", "-t", sessionName)
	
	// Critical: Open /dev/null BEFORE creating command to avoid Cobra stdin inheritance
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		log.Warn("Failed to open /dev/null", "error", err)
		return false
	}
	defer devNull.Close()
	
	checkCmd.Stdin = devNull
	checkCmd.Stdout = devNull
	checkCmd.Stderr = devNull
	
	// Don't use process groups for simple checks
	err = checkCmd.Run()
	return err == nil
}

func createSession(sessionName, sessionPath string) error {
	log.Debug("Session doesn't exist, creating new session", "name", sessionName)
	
	// Open /dev/null BEFORE creating command to avoid Cobra stdin inheritance
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		return fmt.Errorf("failed to open /dev/null: %w", err)
	}
	defer devNull.Close()
	
	createCmd := exec.Command("tmux", "new-session", "-d", "-s", sessionName, "-c", sessionPath)
	
	// Critical: Redirect all I/O to /dev/null
	createCmd.Stdin = devNull
	createCmd.Stdout = devNull
	createCmd.Stderr = devNull
	
	// Create new process group to fully isolate
	createCmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	
	log.Debug("Creating tmux session (isolated with /dev/null)")
	err = createCmd.Run()
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			log.Debug("Session already exists")
			return nil
		}
		log.Debug("Session creation failed", "error", err)
		return fmt.Errorf("failed to create tmux session: %w", err)
	}
	
	log.Debug("Session created successfully")
	return nil
}
