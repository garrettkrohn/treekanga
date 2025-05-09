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
}

type RealSesh struct {
	shell shell.Shell
}

func NewSesh(shell shell.Shell) Sesh {
	return &RealSesh{shell}
}

func (r *RealSesh) SeshConnect(c *com.AddConfig) {

	log.Debug("connecting to: %s", c.Flags.Connect)

	shortestZoxide := slices.Min(c.ZoxideConfig.FoldersToAdd)
	subFolderIsValid := slices.Contains(c.ZoxideConfig.FoldersToAdd, *c.Flags.Connect)
	if subFolderIsValid {
		zoxidePath := shortestZoxide + "/" + *c.Flags.Connect
		log.Info(fmt.Sprintf("Sesh connect to %s", zoxidePath))
		_, err := r.shell.Cmd("sesh", "connect", zoxidePath)
		utility.CheckError(err)
		// deps.Sesh.SeshConnect(zoxidePath)
	} else {
		log.Info(fmt.Sprintf("Sesh connect to %s", shortestZoxide))
		_, err := r.shell.Cmd("sesh", "connect", shortestZoxide)
		utility.CheckError(err)
		// deps.Sesh.SeshConnect(shortestZoxide)
	}
}
