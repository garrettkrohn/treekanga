package services

import (
	"fmt"
	"path/filepath"
	"slices"

	"github.com/charmbracelet/log"
)

func GetSeshPath(seshConnectTarget string, zoxideFolders []string, newRootDirectory string) string {
	if len(zoxideFolders) == 0 {
		return newRootDirectory
	}

	if seshConnectTarget != "" && slices.Contains(zoxideFolders, seshConnectTarget) {
		seshPath := filepath.Join(newRootDirectory, seshConnectTarget)
		log.Debug(fmt.Sprintf("specified sesh sub directory exists: %s", seshConnectTarget))
		return seshPath
	} else {
		log.Debug(fmt.Sprintf("specified sesh sub directory does not exists: %s, using new root directory: %s", seshConnectTarget, newRootDirectory))
		return newRootDirectory
	}
}
