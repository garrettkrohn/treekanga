package connector

import (
	"testing"

	com "github.com/garrettkrohn/treekanga/common"
	"github.com/garrettkrohn/treekanga/shell"
	"github.com/stretchr/testify/assert"
)

func TestSeshConnectNoMatch(t *testing.T) {
	//Arrange
	ms := shell.NewMockShell(t)
	c := getTestAddConfig()
	s := NewConnector(ms)

	connect := "x"
	c.Flags.Sesh = &connect

	//Act
	res := s.GetZoxidePath(&c)

	//Assert
	assert.Equal(t, "/Users/gkrohn/code/treekanga_work/test", res)

}

func TestSeshConnectMatch(t *testing.T) {
	//Arrange
	ms := shell.NewMockShell(t)
	c := getTestAddConfig()
	s := NewConnector(ms)

	connect := "cmd"
	c.Flags.Sesh = &connect

	//Act
	res := s.GetZoxidePath(&c)

	//Assert
	assert.Equal(t, "/Users/gkrohn/code/treekanga_work/test/cmd", res)

}

func getTestAddConfig() com.AddConfig {
	return com.AddConfig{
		Flags: com.AddCmdFlags{
			Directory:  nil,
			BaseBranch: nil,
			Pull:       nil,
			Sesh:       nil,
		},
		Args:         []string{"test"},
		GitInfo:      com.GitInfo{NewBranchName: "test"},
		WorkingDir:   "/Users/gkrohn/code/treekanga_work/treekanga_bare",
		ParentDir:    "/Users/gkrohn/code/treekanga_work",
		WorktreeTargetDir: "/Users/gkrohn/code/treekanga_work/test",
		ZoxideFolders: []string{
			"cmd",
			"git",
		},
		DirectoryReader: nil, // Replace with actual implementation if needed
	}
}
