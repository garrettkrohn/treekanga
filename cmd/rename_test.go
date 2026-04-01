package cmd

import (
	"testing"

	"github.com/garrettkrohn/treekanga/services"
	"github.com/stretchr/testify/assert"
)

func TestValidateRenameArgs(t *testing.T) {
	t.Run("valid single argument", func(t *testing.T) {
		args := []string{"new-branch"}
		branchName, err := services.ValidateRenameArgs(args)
		assert.NoError(t, err)
		assert.Equal(t, "new-branch", branchName)
	})

	t.Run("valid argument with slashes", func(t *testing.T) {
		args := []string{"feature/new-branch"}
		branchName, err := services.ValidateRenameArgs(args)
		assert.NoError(t, err)
		assert.Equal(t, "feature/new-branch", branchName)
	})

	t.Run("error when no arguments", func(t *testing.T) {
		args := []string{}
		_, err := services.ValidateRenameArgs(args)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "provide new branch name")
	})

	t.Run("error when too many arguments", func(t *testing.T) {
		args := []string{"branch1", "branch2"}
		_, err := services.ValidateRenameArgs(args)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "too many arguments")
	})

	t.Run("error when empty string", func(t *testing.T) {
		args := []string{"  "}
		_, err := services.ValidateRenameArgs(args)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be empty")
	})
}
