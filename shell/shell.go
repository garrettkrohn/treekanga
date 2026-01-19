package shell

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/log"

	"github.com/garrettkrohn/treekanga/execwrap"
)

type Shell interface {
	Cmd(cmd string, arg ...string) (string, error)
	ListCmd(cmd string, arg ...string) ([]string, error)
	CmdWithDir(dir string, cmd string, args ...string) (string, error)
	CmdWithStreaming(cmd string, args ...string) error
}

type RealShell struct {
	exec execwrap.Exec
}

func NewShell(exec execwrap.Exec) Shell {
	return &RealShell{exec}
}

// NOTE: This was refactored to allow for the wrapper to potentially
// set a working directory.  I would like to adopt this pattern in
// the git wrapper, instead of using the -c command
func (c *RealShell) Cmd(cmd string, args ...string) (string, error) {
	log.Debug(cmd, "args", args)

	foundCmd, err := c.exec.LookPath(cmd)
	if err != nil {
		return "", err
	}
	command := exec.Command(foundCmd, args...)
	string, err := runCommand(command)
	return string, err
}

func runCommand(command *exec.Cmd) (string, error) {
	var stdout, stderr bytes.Buffer
	command.Stdin = os.Stdin
	command.Stdout = &stdout
	command.Stderr = os.Stderr
	command.Stderr = &stderr
	if err := command.Start(); err != nil {
		return "", err
	}
	if err := command.Wait(); err != nil {
		errString := strings.TrimSpace(stderr.String())
		if strings.HasPrefix(errString, "no server running on") {
			return "", nil
		}
		return "", err
	}
	trimmedOutput := strings.TrimSuffix(string(stdout.String()), "\n")
	return trimmedOutput, nil

}

func (c *RealShell) CmdWithDir(dir string, cmd string, args ...string) (string, error) {
	log.Debug(cmd, "workingdir:", dir, "args", args)

	foundCmd, err := c.exec.LookPath(cmd)
	if err != nil {
		return "", err
	}
	command := exec.Command(foundCmd, args...)
	command.Dir = dir
	string, err := runCommand(command)
	return string, err
}

func (c *RealShell) ListCmd(cmd string, arg ...string) ([]string, error) {
	log.Debug(cmd, "args", arg)
	command := c.exec.Command(cmd, arg...)
	output, err := command.Output()
	return strings.Split(string(output), "\n"), err
}

// CmdWithStreaming executes a command and streams output through the logger
func (c *RealShell) CmdWithStreaming(cmd string, args ...string) error {
	log.Debug(cmd, "args", args)

	foundCmd, err := c.exec.LookPath(cmd)
	if err != nil {
		return err
	}
	command := exec.Command(foundCmd, args...)
	command.Stdin = os.Stdin

	// Create pipes for stdout and stderr
	stdoutPipe, err := command.StdoutPipe()
	if err != nil {
		return err
	}
	stderrPipe, err := command.StderrPipe()
	if err != nil {
		return err
	}

	// Start the command
	if err := command.Start(); err != nil {
		return err
	}

	// Create channels to signal when reading is done and collect stderr
	done := make(chan bool, 2)
	var stderrBuffer bytes.Buffer

	// Read stdout and log it
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			log.Info(scanner.Text())
		}
		done <- true
	}()

	// Read stderr, log it, and capture it
	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			log.Info(line)
			stderrBuffer.WriteString(line)
			stderrBuffer.WriteString("\n")
		}
		done <- true
	}()

	// Wait for both readers to finish
	<-done
	<-done

	// Wait for command to complete
	err = command.Wait()
	if err != nil {
		// Include stderr output in the error message
		stderrOutput := strings.TrimSpace(stderrBuffer.String())
		if stderrOutput != "" {
			return &ExecError{
				Err:    err,
				Stderr: stderrOutput,
			}
		}
		return err
	}
	return nil
}

// ExecError wraps an execution error with stderr output
type ExecError struct {
	Err    error
	Stderr string
}

func (e *ExecError) Error() string {
	return e.Stderr
}

func (e *ExecError) Unwrap() error {
	return e.Err
}
