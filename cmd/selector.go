package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/garrettkrohn/treekanga/config"
	"github.com/garrettkrohn/treekanga/shell"
)

// Selector provides an interface for interactive item selection
type Selector interface {
	IsAvailable() bool
	Select(items []string, prompt string) (string, error)
}

// fzfSelector implements Selector using the fzf command-line tool
type fzfSelector struct {
	shell    shell.Shell
	execFunc func(items []string, prompt string) (string, error) // For testing
}

// IsAvailable checks if fzf is available in the PATH
func (f *fzfSelector) IsAvailable() bool {
	_, err := f.shell.Cmd("which", "fzf")
	return err == nil
}

// Select presents items to the user via fzf and returns the selected item
func (f *fzfSelector) Select(items []string, prompt string) (string, error) {
	if len(items) == 0 {
		return "", errors.New("no items to select from")
	}

	// Use injected exec function for testing, or default implementation
	if f.execFunc != nil {
		return f.execFunc(items, prompt)
	}

	return f.execFzf(items, prompt)
}

// execFzf is the real fzf execution implementation
func (f *fzfSelector) execFzf(items []string, prompt string) (string, error) {
	// Pipe items to fzf stdin
	input := strings.Join(items, "\n")
	cmd := exec.Command("fzf", fmt.Sprintf("--prompt=%s", prompt))
	cmd.Stdin = strings.NewReader(input)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// Check if this is a user cancellation (exit code 130)
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 130 {
				return "", errors.New("selection cancelled by user")
			}
		}
		return "", fmt.Errorf("fzf execution failed: %w", err)
	}

	// Get the selected item and trim whitespace
	selected := strings.TrimSpace(stdout.String())
	if selected == "" {
		return "", errors.New("no item selected")
	}

	return selected, nil
}

// bubbleteaSelector implements Selector using bubbletea (always available)
type bubbleteaSelector struct{}

// IsAvailable always returns true since bubbletea is built-in
func (b *bubbleteaSelector) IsAvailable() bool {
	return true
}

// Select presents items to the user via bubbletea TUI and returns the selected item
func (b *bubbleteaSelector) Select(items []string, prompt string) (string, error) {
	if len(items) == 0 {
		return "", errors.New("no items to select from")
	}

	// Create the model
	model := selectorModel{
		items:  items,
		cursor: 0,
		prompt: prompt,
	}

	// Run the bubbletea program
	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("error running selector: %w", err)
	}

	// Extract result from final model
	result := finalModel.(selectorModel)
	if result.err != nil {
		return "", result.err
	}

	return result.selected, nil
}

// selectorModel is the bubbletea model for item selection
type selectorModel struct {
	items    []string
	cursor   int
	selected string
	err      error
	prompt   string
}

// Init initializes the model (no command needed)
func (m selectorModel) Init() tea.Cmd {
	return nil
}

// Update handles keyboard input and updates the model
func (m selectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			// Move cursor up
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyDown:
			// Move cursor down
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case tea.KeyEnter:
			// Select current item
			m.selected = m.items[m.cursor]
			return m, tea.Quit
		case tea.KeyEsc:
			// Cancel selection
			m.err = errors.New("selection cancelled by user")
			return m, tea.Quit
		case tea.KeyRunes:
			// Handle 'q' for quit
			if len(msg.Runes) > 0 && msg.Runes[0] == 'q' {
				m.err = errors.New("selection cancelled by user")
				return m, tea.Quit
			}
			// Handle 'j' and 'k' for vim-style navigation
			if len(msg.Runes) > 0 && msg.Runes[0] == 'j' {
				if m.cursor < len(m.items)-1 {
					m.cursor++
				}
			}
			if len(msg.Runes) > 0 && msg.Runes[0] == 'k' {
				if m.cursor > 0 {
					m.cursor--
				}
			}
		case tea.KeyCtrlC:
			// Cancel with Ctrl+C
			m.err = errors.New("selection cancelled by user")
			return m, tea.Quit
		}
	}

	return m, nil
}

// View renders the selector UI
func (m selectorModel) View() string {
	var b strings.Builder

	// Show prompt
	if m.prompt != "" {
		b.WriteString(m.prompt)
		b.WriteString("\n\n")
	}

	// Render items with cursor
	for i, item := range m.items {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		b.WriteString(fmt.Sprintf("%s %s\n", cursor, item))
	}

	// Show help text
	b.WriteString("\n")
	b.WriteString("↑/k up • ↓/j down • enter select • q/esc cancel\n")

	return b.String()
}

// getSelector returns the appropriate Selector based on config
func getSelector(cfg config.AppConfig, sh shell.Shell) Selector {
	// Check if fzf is requested
	if cfg.SelectorMode == "fzf" {
		fzf := &fzfSelector{shell: sh}
		if fzf.IsAvailable() {
			return fzf
		}
		log.Warn("fzf not found in PATH, using built-in selector")
	}

	// Default to bubbletea selector for empty or unknown values
	return &bubbleteaSelector{}
}
