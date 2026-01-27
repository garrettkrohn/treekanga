package connector

import (
	"github.com/garrettkrohn/treekanga/shell"
	"github.com/garrettkrohn/treekanga/utility"
)

type Connector interface {
	SeshConnect(seshPath string)
	VsCodeConnect(newRootPath string)
	CursorConnect(newRootPath string)
}

type RealConnector struct {
	shell shell.Shell
}

func NewConnector(shell shell.Shell) Connector {
	return &RealConnector{shell}
}

func (r *RealConnector) SeshConnect(seshPath string) {
	_, err := r.shell.Cmd("sesh", "connect", seshPath)
	utility.CheckError(err)
}

func (r *RealConnector) VsCodeConnect(newRootPath string) {
	_, err := r.shell.Cmd("code", newRootPath)
	utility.CheckError(err)
}

func (r *RealConnector) CursorConnect(newRootPath string) {
	_, err := r.shell.Cmd("cursor", newRootPath)
	utility.CheckError(err)
}
