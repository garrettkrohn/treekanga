package cmd

import (
	"errors"
	"testing"

	"github.com/garrettkrohn/treekanga/shell"
	"github.com/stretchr/testify/assert"
)

func TestFzfSelector_IsAvailable_WhenFzfInPath(t *testing.T) {
	// Arrange - mock shell with fzf available
	mockShell := shell.NewMockShell(t)
	mockShell.On("Cmd", "which", "fzf").Return("/usr/bin/fzf", nil)

	selector := &fzfSelector{shell: mockShell}

	// Act - call IsAvailable()
	result := selector.IsAvailable()

	// Assert - expect true
	assert.True(t, result, "Expected IsAvailable to return true when fzf is in PATH")
}

func TestFzfSelector_IsAvailable_WhenFzfNotInPath(t *testing.T) {
	// Arrange - mock shell without fzf
	mockShell := shell.NewMockShell(t)
	mockShell.On("Cmd", "which", "fzf").Return("", errors.New("command not found"))

	selector := &fzfSelector{shell: mockShell}

	// Act - call IsAvailable()
	result := selector.IsAvailable()

	// Assert - expect false
	assert.False(t, result, "Expected IsAvailable to return false when fzf is not in PATH")
}

func TestFzfSelector_Select_ReturnsSelectedItem(t *testing.T) {
	// Arrange - mock shell to return "item2" from fzf command
	mockShell := shell.NewMockShell(t)
	items := []string{"item1", "item2", "item3"}

	// Mock the exec function to return "item2"
	selector := &fzfSelector{
		shell: mockShell,
		execFunc: func(items []string, prompt string) (string, error) {
			return "item2", nil
		},
	}

	// Act - call Select() with items
	result, err := selector.Select(items, "Select: ")

	// Assert - expect "item2" returned
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, "item2", result, "Expected selected item to be 'item2'")
}

func TestFzfSelector_Select_ReturnsErrorOnFzfFailure(t *testing.T) {
	// Arrange - mock shell to fail fzf execution
	mockShell := shell.NewMockShell(t)
	items := []string{"item1"}

	// Mock the exec function to return error
	selector := &fzfSelector{
		shell: mockShell,
		execFunc: func(items []string, prompt string) (string, error) {
			return "", errors.New("fzf execution failed")
		},
	}

	// Act - call Select() with items
	result, err := selector.Select(items, "Choose: ")

	// Assert - expect error returned
	assert.Error(t, err, "Expected error when fzf execution fails")
	assert.Empty(t, result, "Expected empty result on error")
}

func TestFzfSelector_Select_ReturnsErrorOnUserCancel(t *testing.T) {
	// Arrange - mock shell to return exit code 130 (user cancelled)
	mockShell := shell.NewMockShell(t)
	items := []string{"item1"}

	// Mock the exec function to return cancellation error
	selector := &fzfSelector{
		shell: mockShell,
		execFunc: func(items []string, prompt string) (string, error) {
			return "", errors.New("selection cancelled by user")
		},
	}

	// Act - call Select() with items
	result, err := selector.Select(items, "Pick: ")

	// Assert - expect error with "cancelled" message
	assert.Error(t, err, "Expected error when user cancels")
	assert.Contains(t, err.Error(), "cancel", "Expected error message to contain 'cancel'")
	assert.Empty(t, result, "Expected empty result on cancellation")
}
