package sesh

import (
	"fmt"
	"slices"

	"github.com/charmbracelet/log"
	com "github.com/garrettkrohn/treekanga/common"
	"github.com/garrettkrohn/treekanga/shell"
	"github.com/garrettkrohn/treekanga/utility"
)

type Sesh interface {
	SeshConnect(c *com.AddConfig)
	GetZoxidePath(c *com.AddConfig) string
}

type RealSesh struct {
	shell shell.Shell
}

func NewSesh(shell shell.Shell) Sesh {
	return &RealSesh{shell}
}

func (r *RealSesh) GetZoxidePath(c *com.AddConfig) string {
	shortestZoxide := slices.Min(c.ZoxideConfig.FoldersToAdd)
	subFolderIsValid := slices.Contains(c.ZoxideConfig.FoldersToAdd, *c.Flags.Connect)
	if subFolderIsValid {
		zoxidePath := c.ParentDir + "/" + c.ZoxideConfig.NewBranchName + "/" + *c.Flags.Connect
		log.Info(fmt.Sprintf("Sesh connect to %s", zoxidePath))
		return zoxidePath
	} else {
		log.Info(fmt.Sprintf("Sesh connect to %s", shortestZoxide))
		return c.ParentDir + "/" + c.ZoxideConfig.NewBranchName
	}
}

func (r *RealSesh) SeshConnect(c *com.AddConfig) {
	zoxidePath := r.GetZoxidePath(c)
	_, err := r.shell.Cmd("sesh", "connect", zoxidePath)
	utility.CheckError(err)
}
