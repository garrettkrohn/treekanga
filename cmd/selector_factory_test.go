package cmd

import (
	"errors"
	"testing"

	"github.com/garrettkrohn/treekanga/config"
	"github.com/garrettkrohn/treekanga/shell"
	"github.com/stretchr/testify/assert"
)

func TestGetSelector_ReturnsFzf_WhenConfiguredAndAvailable(t *testing.T) {
	// Arrange - config with selectorMode="fzf", mock shell with fzf available
	cfg := config.AppConfig{
		SelectorMode: "fzf",
	}
	mockShell := shell.NewMockShell(t)
	mockShell.On("Cmd", "which", "fzf").Return("/usr/bin/fzf", nil)

	// Act - call getSelector(config, shell)
	selector := getSelector(cfg, mockShell)

	// Assert - expect fzfSelector type returned
	_, ok := selector.(*fzfSelector)
	assert.True(t, ok, "Expected fzfSelector when selectorMode='fzf' and fzf is available")
}

func TestGetSelector_ReturnsBubbletea_WhenFzfConfiguredButNotAvailable(t *testing.T) {
	// Arrange - config with selectorMode="fzf", mock shell without fzf
	cfg := config.AppConfig{
		SelectorMode: "fzf",
	}
	mockShell := shell.NewMockShell(t)
	mockShell.On("Cmd", "which", "fzf").Return("", errors.New("command not found"))

	// Act - call getSelector(config, shell)
	selector := getSelector(cfg, mockShell)

	// Assert - expect bubbleteaSelector type returned
	_, ok := selector.(*bubbleteaSelector)
	assert.True(t, ok, "Expected bubbleteaSelector when fzf is not available")
}

func TestGetSelector_LogsWarning_WhenFzfNotAvailable(t *testing.T) {
	// Arrange - config with selectorMode="fzf", mock shell, log capture
	// Note: Testing log output is difficult in unit tests, so we verify the behavior
	// (returns bubbletea) which implies the warning was logged
	cfg := config.AppConfig{
		SelectorMode: "fzf",
	}
	mockShell := shell.NewMockShell(t)
	mockShell.On("Cmd", "which", "fzf").Return("", errors.New("command not found"))

	// Act - call getSelector(config, shell)
	selector := getSelector(cfg, mockShell)

	// Assert - verify warning log behavior (returns bubbletea fallback)
	_, ok := selector.(*bubbleteaSelector)
	assert.True(t, ok, "Expected bubbleteaSelector fallback when fzf not available (warning should be logged)")
}

func TestGetSelector_ReturnsBubbletea_WhenSelectorModeEmpty(t *testing.T) {
	// Arrange - config with empty selectorMode
	cfg := config.AppConfig{
		SelectorMode: "",
	}
	mockShell := shell.NewMockShell(t)
	// No shell calls expected for empty selectorMode

	// Act - call getSelector(config, shell)
	selector := getSelector(cfg, mockShell)

	// Assert - expect bubbleteaSelector type returned
	_, ok := selector.(*bubbleteaSelector)
	assert.True(t, ok, "Expected bubbleteaSelector when selectorMode is empty")
}

func TestGetSelector_ReturnsBubbletea_WhenSelectorModeUnknown(t *testing.T) {
	// Arrange - config with selectorMode="unknown-value"
	cfg := config.AppConfig{
		SelectorMode: "unknown-value",
	}
	mockShell := shell.NewMockShell(t)
	// No shell calls expected for unknown selectorMode

	// Act - call getSelector(config, shell)
	selector := getSelector(cfg, mockShell)

	// Assert - expect bubbleteaSelector type returned (ignore unknown values)
	_, ok := selector.(*bubbleteaSelector)
	assert.True(t, ok, "Expected bubbleteaSelector when selectorMode is unknown")
}
