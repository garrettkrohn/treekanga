package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/garrettkrohn/treekanga/adapters"
	"github.com/garrettkrohn/treekanga/directoryReader"
	util "github.com/garrettkrohn/treekanga/utility"
)

func CompileZoxidePathsToAdd(foldersToAddFromConfig []string, newRootDirectory string, directoryReader directoryReader.DirectoryReader) []string {
	var foldersToAdd []string
	// add the root
	foldersToAdd = append(foldersToAdd, newRootDirectory)

	for _, folder := range foldersToAddFromConfig {
		if !isLastCharWildcard(folder) {
			newFolderFromConfig := filepath.Join(newRootDirectory, folder)
			// validate
			if checkPath(newFolderFromConfig) {
				foldersToAdd = append(foldersToAdd, newFolderFromConfig)
			} else {
				log.Error(fmt.Sprintf("zoxide path %s does not exist", newFolderFromConfig))
			}
		} else {
			pathUpTillWildcard := getPathUntilLastSlash(folder)
			baseFolderToSearch := filepath.Join(newRootDirectory, pathUpTillWildcard)
			configFolders, err := directoryReader.GetFoldersInDirectory(baseFolderToSearch)

			for _, configFolder := range configFolders {
				newConfigFolder := filepath.Join(baseFolderToSearch, configFolder)
				// validate
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

func GetQueryList(zoxideAdapter adapters.Zoxide, pathSearch string) ([]string, error) {
	allEntries, err := zoxideAdapter.QueryList(pathSearch)
	if err != nil {
		log.Error(err)
		return []string{}, err
	}

	// Filter for entries that start with pathSearch
	var filteredEntries []string
	for _, entry := range allEntries {
		entry = strings.TrimSpace(entry)
		if entry != "" && strings.HasPrefix(entry, pathSearch) {
			filteredEntries = append(filteredEntries, entry)
		}
	}

	return filteredEntries, nil
}
