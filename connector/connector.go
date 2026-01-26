package connector

import (
	"slices"

	"github.com/charmbracelet/log"
	com "github.com/garrettkrohn/treekanga/common"
	"github.com/garrettkrohn/treekanga/shell"
	"github.com/garrettkrohn/treekanga/utility"
)

type Connector interface {
	SeshConnect(seshPath string)
	VsCodeConnect(newRootPath string)
	CursorConnect(newRootPath string)
	GetZoxidePath(c *com.AddConfig) string
}

type RealConnector struct {
	shell shell.Shell
}

func NewConnector(shell shell.Shell) Connector {
	return &RealConnector{shell}
}

func (r *RealConnector) GetZoxidePath(c *com.AddConfig) string {
	if len(c.ZoxideFolders) == 0 {
		return c.GetZoxideBasePath()
	}

	seshTarget := c.GetSeshTarget()

	if seshTarget != "" && slices.Contains(c.ZoxideFolders, seshTarget) {
		zoxidePath := c.GetZoxidePath(seshTarget)
		log.Info("Sesh connect", "path", zoxidePath)
		return zoxidePath
	} else {
		basePath := c.GetZoxideBasePath()
		log.Info("Sesh connect", "path", basePath)
		return basePath
	}
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
