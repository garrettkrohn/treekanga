package services

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/garrettkrohn/treekanga/directoryReader"
	"github.com/garrettkrohn/treekanga/models"
)

// ExpandWorktreesWithZoxideFolders takes a list of worktrees and expands them
// with subdirectories based on the zoxideFolders configuration
func ExpandWorktreesWithZoxideFolders(worktrees []models.Worktree, zoxideFolders []string, dirReader directoryReader.DirectoryReader) []string {
	var allPaths []string

	for _, worktree := range worktrees {
		// Always add the worktree root
		allPaths = append(allPaths, worktree.FullPath)

		// Add subdirectories based on zoxideFolders config
		for _, folder := range zoxideFolders {
			expandedPaths := ExpandZoxideFolder(worktree.FullPath, folder, dirReader)
			allPaths = append(allPaths, expandedPaths...)
		}
	}

	return allPaths
}

// ExpandZoxideFolder expands a zoxide folder pattern for a given worktree root
func ExpandZoxideFolder(worktreeRoot string, folder string, dirReader directoryReader.DirectoryReader) []string {
	var paths []string

	// Check if folder has wildcard
	if !isLastCharWildcard(folder) {
		// No wildcard - just append the folder to the worktree root
		fullPath := filepath.Join(worktreeRoot, folder)
		if checkPath(fullPath) {
			paths = append(paths, fullPath)
		} else {
			log.Debug("Path does not exist", "path", fullPath)
		}
	} else {
		// Has wildcard - expand all matching directories
		pathUpTillWildcard := getPathUntilLastSlash(folder)
		baseFolderToSearch := filepath.Join(worktreeRoot, pathUpTillWildcard)
		
		folders, err := dirReader.GetFoldersInDirectory(baseFolderToSearch)
		if err != nil {
			log.Debug("Could not read directory", "path", baseFolderToSearch, "error", err)
			return paths
		}

		for _, subFolder := range folders {
			fullPath := filepath.Join(baseFolderToSearch, subFolder)
			if checkPath(fullPath) {
				paths = append(paths, fullPath)
			}
		}
	}

	return paths
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
