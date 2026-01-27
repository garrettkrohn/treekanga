package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSeshPathNoMatch(t *testing.T) {
	//Arrange
	seshConnectTarget := "x"
	zoxideFolders := []string{"cmd", "git"}
	newRootDirectory := "/Users/gkrohn/code/treekanga_work/test"

	//Act
	res := GetSeshPath(seshConnectTarget, zoxideFolders, newRootDirectory)

	//Assert
	assert.Equal(t, "/Users/gkrohn/code/treekanga_work/test", res)
}

func TestGetSeshPathMatch(t *testing.T) {
	//Arrange
	seshConnectTarget := "cmd"
	zoxideFolders := []string{"cmd", "git"}
	newRootDirectory := "/Users/gkrohn/code/treekanga_work/test"

	//Act
	res := GetSeshPath(seshConnectTarget, zoxideFolders, newRootDirectory)

	//Assert
	assert.Equal(t, "/Users/gkrohn/code/treekanga_work/test/cmd", res)
}

func TestGetSeshPathEmptyZoxideFolders(t *testing.T) {
	//Arrange
	seshConnectTarget := "cmd"
	zoxideFolders := []string{}
	newRootDirectory := "/Users/gkrohn/code/treekanga_work/test"

	//Act
	res := GetSeshPath(seshConnectTarget, zoxideFolders, newRootDirectory)

	//Assert
	assert.Equal(t, "/Users/gkrohn/code/treekanga_work/test", res)
}

func TestGetSeshPathEmptyTarget(t *testing.T) {
	//Arrange
	seshConnectTarget := ""
	zoxideFolders := []string{"cmd", "git"}
	newRootDirectory := "/Users/gkrohn/code/treekanga_work/test"

	//Act
	res := GetSeshPath(seshConnectTarget, zoxideFolders, newRootDirectory)

	//Assert
	assert.Equal(t, "/Users/gkrohn/code/treekanga_work/test", res)
}
