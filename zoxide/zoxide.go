package zoxide

import (
	"fmt"
	"strings"

	"github.com/garrettkrohn/treekanga/shell"
	util "github.com/garrettkrohn/treekanga/utility"
)

type Zoxide interface {
	AddPath(path string) error
	RemovePath(path string) error
	AddZoxideEntries(zoxideFolders []string)
	QueryScore(path string) (float64, error)
}

type RealZoxide struct {
	shell shell.Shell
}

func NewZoxide(shell shell.Shell) Zoxide {
	return &RealZoxide{shell}
}

func (r *RealZoxide) AddPath(path string) error {
	_, err := r.shell.Cmd("zoxide", "add", path)
	return err
}

func (r *RealZoxide) RemovePath(path string) error {
	_, err := r.shell.Cmd("zoxide", "remove", path)
	return err
}

func (r *RealZoxide) AddZoxideEntries(zoxideFolders []string) {
	for _, folder := range zoxideFolders {
		err := r.AddPath(folder)
		util.CheckError(err)
	}
}

func (r *RealZoxide) QueryScore(path string) (float64, error) {
	output, err := r.shell.Cmd("zoxide", "query", "--score", path)
	if err != nil {
		// If zoxide doesn't have this path, return 0
		return 0, nil
	}
	// Parse the output which should be in format "score path"
	parts := strings.Fields(output)
	if len(parts) < 1 {
		return 0, nil
	}
	// Try to parse the score
	var score float64
	_, parseErr := fmt.Sscanf(parts[0], "%f", &score)
	if parseErr != nil {
		return 0, nil
	}
	return score, nil
}
