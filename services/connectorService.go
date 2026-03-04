package services

import (
	"github.com/charmbracelet/log"
)

func GetSeshPath(seshConnectTarget string, newRootDirectory string) string {
	log.Debug("using new root directory for sesh connection", "path", newRootDirectory)
	return newRootDirectory
}
