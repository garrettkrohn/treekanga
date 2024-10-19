package zoxide

import (
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
	_, err := r.shell.Cmd("zoxide", "add", path)
	return err
}
