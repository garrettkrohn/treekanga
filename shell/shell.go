package shell

import (
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
