package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZoxideAddEntriesNoConfig(t *testing.T) {
	actualOutput := getListOfZoxideEntries("baseBranch", "repoName", "parentDir", nil)

	expectedOutput := []string{"parentDir/baseBranch"}

	assert.Equal(t, expectedOutput, actualOutput)
}

func TestZoxideAddEntriesWithConfigNoWildcard(t *testing.T) {
	actualOutput := getListOfZoxideEntries("baseBranch", "repoName", "parentDir", []string{"test"})

	expectedOutput := []string{"parentDir/baseBranch", "parentDir/baseBranch/test"}

	assert.Equal(t, expectedOutput, actualOutput)
}
