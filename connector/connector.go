package connector

import (
	"slices"

	"github.com/charmbracelet/log"
	com "github.com/garrettkrohn/treekanga/common"
	"github.com/garrettkrohn/treekanga/shell"
	"github.com/garrettkrohn/treekanga/utility"
)

type Connector interface {
	SeshConnect(c *com.AddConfig)
	VsCodeConnect(c *com.AddConfig)
	CursorConnect(c *com.AddConfig)
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

func (r *RealConnector) SeshConnect(c *com.AddConfig) {
	zoxidePath := r.GetZoxidePath(c)
	_, err := r.shell.Cmd("sesh", "connect", zoxidePath)
	utility.CheckError(err)
}

func (r *RealConnector) VsCodeConnect(c *com.AddConfig) {
	addPath := c.GetWorktreePath()
	_, err := r.shell.Cmd("code", addPath)
	utility.CheckError(err)
}

func (r *RealConnector) CursorConnect(c *com.AddConfig) {
	addPath := c.GetWorktreePath()
	_, err := r.shell.Cmd("cursor", addPath)
	utility.CheckError(err)
}
