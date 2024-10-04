package zoxide

import (
	"fmt"

	"github.com/garrettkrohn/treekanga/shell"
)

type Zoxide interface {
	AddPath(path string) error
}

type RealZoxide struct {
	shell shell.Shell
}

func NewZoxide(shell shell.Shell) Zoxide {
	return &RealZoxide{shell}
}

func (r *RealZoxide) AddPath(path string) error {
	fmt.Print(path)
	_, err := r.shell.Cmd("zoxide", "add", path)
	return err
}
