package connector

import (
	"fmt"
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
	shortestZoxide := slices.Min(c.ZoxideConfig.FoldersToAdd)
	subFolderIsValid := slices.Contains(c.ZoxideConfig.FoldersToAdd, *c.Flags.Sesh)
	if subFolderIsValid {
		zoxidePath := c.ParentDir + "/" + c.ZoxideConfig.NewBranchName + "/" + *c.Flags.Sesh
		log.Info(fmt.Sprintf("Sesh connect to %s", zoxidePath))
		return zoxidePath
	} else {
		log.Info(fmt.Sprintf("Sesh connect to %s", shortestZoxide))
		return c.ParentDir + "/" + c.ZoxideConfig.NewBranchName
	}
}

func (r *RealConnector) SeshConnect(c *com.AddConfig) {
	zoxidePath := r.GetZoxidePath(c)
	_, err := r.shell.Cmd("sesh", "connect", zoxidePath)
	utility.CheckError(err)
}

func (r *RealConnector) VsCodeConnect(c *com.AddConfig) {
	addPath := getCodeAndCursorPath(c)
	_, err := r.shell.Cmd("code", addPath)
	utility.CheckError(err)
}

func (r *RealConnector) CursorConnect(c *com.AddConfig) {
	addPath := getCodeAndCursorPath(c)
	_, err := r.shell.Cmd("cursor", addPath)
	utility.CheckError(err)
}

func getCodeAndCursorPath(c *com.AddConfig) string {
	return c.ZoxideConfig.ParentDir + "/" + c.ZoxideConfig.NewBranchName
}
