package services

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/garrettkrohn/treekanga/directoryReader"
	util "github.com/garrettkrohn/treekanga/utility"
)

func CompileZoxidePathsToAdd(foldersToAddFromConfig []string, newRootDirectory string) []string {

	directoryReader := directoryReader.NewDirectoryReader()
	var foldersToAdd []string

	//add the root
	foldersToAdd = append(foldersToAdd, newRootDirectory)

	for _, folder := range foldersToAddFromConfig {
		if !isLastCharWildcard(folder) {
			newFolderFromConfig := newRootDirectory + "/" + folder
			//validate
			if checkPath(newFolderFromConfig) {
				foldersToAdd = append(foldersToAdd, newFolderFromConfig)
			} else {
				log.Error(fmt.Sprintf("zoxide path %s does not exist", newFolderFromConfig))
			}
		} else {
			pathUpTillWildcard := getPathUntilLastSlash(folder)
			baseFolderToSearch := newRootDirectory + "/" + pathUpTillWildcard
			configFolders, err := directoryReader.GetFoldersInDirectory(baseFolderToSearch)

			for _, configFolder := range configFolders {
				newConfigFolder := baseFolderToSearch + "/" + configFolder
				//validate
				if checkPath(newConfigFolder) {
					foldersToAdd = append(foldersToAdd, newConfigFolder)
				} else {
					log.Error(fmt.Sprintf("zoxide path %s does not exist", newConfigFolder))
				}
			}
			util.CheckError(err)
		}
	}

	return foldersToAdd
}

func checkPath(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func isLastCharWildcard(input string) bool {
	parts := strings.Split(input, "/")
	lastSegment := parts[len(parts)-1]
	return strings.HasSuffix(lastSegment, "*")
}

func getPathUntilLastSlash(input string) string {
	parts := strings.Split(input, "/")
	if len(parts) > 1 {
		return strings.Join(parts[:len(parts)-1], "/")
	}
	return ""
}
