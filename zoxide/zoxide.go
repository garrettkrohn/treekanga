package zoxide

import (
	"fmt"
	"strings"

	"github.com/garrettkrohn/treekanga/common"
	"github.com/garrettkrohn/treekanga/directoryReader"
	"github.com/garrettkrohn/treekanga/shell"
	util "github.com/garrettkrohn/treekanga/utility"
)

type Zoxide interface {
	AddPath(path string) error
	RemovePath(path string) error
	AddZoxideEntries(c *common.AddConfig)
	QueryScore(path string) (float64, error)
	GetQueryList(pathSearch string) ([]string, error)
}

type RealZoxide struct {
	shell shell.Shell
}

func NewZoxide(shell shell.Shell) Zoxide {
	return &RealZoxide{shell}
}

func (r *RealZoxide) AddPath(path string) error {
	_, err := r.shell.Cmd("zoxide", "add", path)
	return err
}

func (r *RealZoxide) RemovePath(path string) error {
	_, err := r.shell.Cmd("zoxide", "remove", path)
	return err
}

func (r *RealZoxide) AddZoxideEntries(c *common.AddConfig) {
	baseName := c.GetZoxideBasePath()

	var foldersToAdd []string
	foldersToAdd = append(foldersToAdd, baseName)

	foldersToAdd = addConfigFolders(foldersToAdd, c.ZoxideFolders, baseName, c.DirectoryReader)

	for _, folder := range foldersToAdd {
		err := r.AddPath(folder)
		util.CheckError(err)
	}
}

func (r *RealZoxide) QueryScore(path string) (float64, error) {
	output, err := r.shell.Cmd("zoxide", "query", "--score", path)
	if err != nil {
		// If zoxide doesn't have this path, return 0
		return 0, nil
	}
	// Parse the output which should be in format "score path"
	parts := strings.Fields(output)
	if len(parts) < 1 {
		return 0, nil
	}
	// Try to parse the score
	var score float64
	_, parseErr := fmt.Sscanf(parts[0], "%f", &score)
	if parseErr != nil {
		return 0, nil
	}
	return score, nil
}

func (r *RealZoxide) GetQueryList(pathSearch string) ([]string, error) {
	// Get all zoxide entries
	output, err := r.shell.Cmd("zoxide", "query", "--list")
	if err != nil {
		return []string{}, err
	}
	if output == "" {
		return []string{}, nil
	}

	// Split into lines and filter for entries that start with pathSearch
	allEntries := strings.Split(output, "\n")
	var filteredEntries []string
	for _, entry := range allEntries {
		entry = strings.TrimSpace(entry)
		if entry != "" && strings.HasPrefix(entry, pathSearch) {
			filteredEntries = append(filteredEntries, entry)
		}
	}

	return filteredEntries, nil
}

func addConfigFolders(foldersToAdd []string, foldersToAddFromConfig []string, baseName string, directoryReader directoryReader.DirectoryReader) []string {
	for _, folder := range foldersToAddFromConfig {
		if !isLastCharWildcard(folder) {
			newFolderFromConfig := baseName + "/" + folder
			foldersToAdd = append(foldersToAdd, newFolderFromConfig)
		} else {
			pathUpTillWildcard := getPathUntilLastSlash(folder)
			baseFolderToSearch := baseName + "/" + pathUpTillWildcard
			configFolders, err := directoryReader.GetFoldersInDirectory(baseFolderToSearch)

			for _, configFolder := range configFolders {
				newConfigFolder := baseFolderToSearch + "/" + configFolder
				foldersToAdd = append(foldersToAdd, newConfigFolder)
			}
			util.CheckError(err)
		}
	}
	return foldersToAdd
}

func getPathUntilLastSlash(input string) string {
	parts := strings.Split(input, "/")
	if len(parts) > 1 {
		return strings.Join(parts[:len(parts)-1], "/")
	}
	return ""
}

func isLastCharWildcard(input string) bool {
	parts := strings.Split(input, "/")
	lastSegment := parts[len(parts)-1]
	return strings.HasSuffix(lastSegment, "*")
}
