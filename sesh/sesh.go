package sesh

import (
	"github.com/garrettkrohn/treekanga/shell"
)

type Sesh interface {
	SeshConnect(seshName string) error
}

type RealSesh struct {
	shell shell.Shell
}

func NewSesh(shell shell.Shell) Sesh {
	return &RealSesh{shell}
}

func (r *RealSesh) SeshConnect(seshName string) error {
	_, err := r.shell.Cmd("sesh", "connect", seshName)
	return err
}
