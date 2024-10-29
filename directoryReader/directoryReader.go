package directoryReader

import (
	"fmt"
	"os"

	"github.com/garrettkrohn/treekanga/utility"
)

type DirectoryReader interface {
	GetFoldersInDirectory(dirPath string) ([]string, error)
}

type DirectoryReaderImpl struct{}

func NewDirectoryReader() DirectoryReader {
	return &DirectoryReaderImpl{}
}

func (d DirectoryReaderImpl) GetFoldersInDirectory(dirPath string) ([]string, error) {
	var folders []string

	fmt.Println(dirPath)
	entries, err := os.ReadDir(dirPath)
	utility.CheckError(err)
	fmt.Println("entries", entries)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			folders = append(folders, entry.Name())
		}
	}

	return folders, nil
}
