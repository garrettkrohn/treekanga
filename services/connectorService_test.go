package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSeshPath(t *testing.T) {
	//Arrange
	seshConnectTarget := "x"
	newRootDirectory := "/Users/gkrohn/code/treekanga_work/test"

	//Act
	res := GetSeshPath(seshConnectTarget, newRootDirectory)

	//Assert
	assert.Equal(t, "/Users/gkrohn/code/treekanga_work/test", res)
}
