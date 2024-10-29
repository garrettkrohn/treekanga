package cmd

import (
	"testing"

	"github.com/garrettkrohn/treekanga/directoryReader"
	"github.com/stretchr/testify/assert"
)

func TestZoxideAddEntriesNoConfig(t *testing.T) {
	mockDirectoryReader := directoryReader.NewMockDirectoryReader(t)
	actualOutput := getListOfZoxideEntries("baseBranch", "repoName", "parentDir", nil, mockDirectoryReader)

	expectedOutput := []string{"parentDir/baseBranch"}

	assert.Equal(t, expectedOutput, actualOutput)
}

func TestZoxideAddEntriesWithConfigNoWildcard(t *testing.T) {
	mockDirectoryReader := directoryReader.NewMockDirectoryReader(t)
	actualOutput := getListOfZoxideEntries("baseBranch", "repoName", "parentDir", []string{"test"}, mockDirectoryReader)

	expectedOutput := []string{"parentDir/baseBranch", "parentDir/baseBranch/test"}

	assert.Equal(t, expectedOutput, actualOutput)
}

func TestZoxideAddEntriesWithConfigWithWildcard(t *testing.T) {
	mockDirectoryReader := directoryReader.NewMockDirectoryReader(t)
	mockDirectoryReader.On("GetFoldersInDirectory", "parentDir/baseBranch/test").Return([]string{
		"folder1",
		"folder2",
	}, nil)
	actualOutput := getListOfZoxideEntries("baseBranch", "repoName", "parentDir", []string{"test", "test/*"}, mockDirectoryReader)

	expectedOutput := []string{
		"parentDir/baseBranch",
		"parentDir/baseBranch/test",
		"parentDir/baseBranch/test/folder1",
		"parentDir/baseBranch/test/folder2"}

	assert.Equal(t, expectedOutput, actualOutput)
}
