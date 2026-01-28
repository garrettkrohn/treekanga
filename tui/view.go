/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// View renders the TUI based on the current model state
func (m Model) View() string {
	// If terminal size not set yet, return empty
	if m.termWidth == 0 || m.termHeight == 0 {
		return ""
	}

	// Show spinner popup if adding (no logs in background)
	if m.isAdding {
		return m.renderAddSpinnerPopup()
	}

	// Show add input prompt
	if m.showAddInput {
		baseView := m.renderSplitView()
		return m.renderAddInputPopup(baseView)
	}

	// Show spinner popup if deleting (no logs in background)
	if m.isDeleting {
		return m.renderSpinnerPopup()
	}

	// Show delete confirmation dialog
	if m.showDeleteConfirm {
		baseView := m.renderSplitView()
		return m.renderConfirmPopup(baseView)
	}

	// If popup is showing, render it with visual separation
	if m.showPopup {
		return m.renderModalPopup()
	}

	// Render the split view with table and logs
	return m.renderSplitView()
}

// renderSplitView renders the main view with table on top and logs on bottom
func (m Model) renderSplitView() string {
	// Render table
	tableView := m.baseTableStyle().Render(m.table.View())

	// Render logs pane
	logsView := m.renderLogsPane()

	// Combine vertically
	splitView := lipgloss.JoinVertical(
		lipgloss.Left,
		tableView,
		logsView,
	)

	// Add help text at the bottom
	return splitView + "\n" + helpText(m) + "\n"
}

// renderLogsPane renders the logs section as a pane
func (m Model) renderLogsPane() string {
	// Style for the logs header
	headerStyle := lipgloss.NewStyle().
		Foreground(m.theme.Accent).
		Bold(true).
		Padding(0, 1)

	// Different border color if focused
	borderColor := m.theme.BorderDim
	if m.logsFocused {
		borderColor = m.theme.Accent
	}

	header := headerStyle.Render("Operation Logs")

	// Style for the logs viewport container
	logsStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1).
		Width(m.termWidth - 2).
		Height(m.logsViewport.Height + 2)

	return logsStyle.Render(header + "\n" + m.logsViewport.View())
}

// renderSpinnerPopup shows a small centered popup with spinner
func (m Model) renderSpinnerPopup() string {
	spinnerStyle := lipgloss.NewStyle().Foreground(m.theme.Accent).Bold(true)
	messageStyle := lipgloss.NewStyle().Foreground(m.theme.TextFg)

	content := fmt.Sprintf("\n  %s  %s\n",
		spinnerStyle.Render(m.spinner.View()),
		messageStyle.Render(fmt.Sprintf("Deleting worktree: %s...", m.deletingName)))

	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Accent).
		Padding(1, 3).
		Align(lipgloss.Center)

	popup := popupStyle.Render(content)

	// Place popup and fill whitespace with spaces
	return lipgloss.Place(
		m.termWidth,
		m.termHeight,
		lipgloss.Center,
		lipgloss.Center,
		popup,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.NoColor{}),
	)
}

// renderConfirmPopup shows the delete confirmation dialog as a popup
func (m Model) renderConfirmPopup(background string) string {
	errorStyle := lipgloss.NewStyle().
		Foreground(m.theme.ErrorFg).
		Bold(true)

	messageStyle := lipgloss.NewStyle().
		Foreground(m.theme.TextFg)

	errorMsg := errorStyle.Render("⚠ Error deleting worktree")
	message := fmt.Sprintf("\n%s\n\n%s\n\nWorktree '%s' contains uncommitted changes.\n\nForce delete and discard all changes?",
		errorMsg,
		messageStyle.Render(m.deleteConfirmError),
		messageStyle.Bold(true).Render(m.pendingDeleteName))

	hintStyle := lipgloss.NewStyle().
		Foreground(m.theme.MutedFg).
		Italic(true).
		Align(lipgloss.Center)

	hint := hintStyle.Render("\nPress [Y]es to force delete • [N]o to cancel")

	content := message + "\n" + hint

	popupWidth := m.termWidth * 2 / 3
	if popupWidth > 60 {
		popupWidth = 60
	}

	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.ErrorFg).
		Padding(1, 2).
		Width(popupWidth).
		Align(lipgloss.Left)

	popup := popupStyle.Render(content)

	return lipgloss.Place(
		m.termWidth,
		m.termHeight,
		lipgloss.Center,
		lipgloss.Center,
		popup,
	)
}

// renderModalPopup shows a prominent centered popup
func (m Model) renderModalPopup() string {
	// Create popup (60% width, 70% height to leave visible margins)
	popupWidth := (m.termWidth * 3) / 5
	popupHeight := (m.termHeight * 7) / 10

	// Ensure minimum size
	if popupWidth < 45 {
		popupWidth = 45
	}
	if popupHeight < 12 {
		popupHeight = 12
	}

	// Add padding around popup to create "floating" effect
	marginSize := 2

	// Create the popup content area
	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(m.theme.Accent).
		Padding(1, 2).
		Width(popupWidth - (marginSize * 2)).
		Height(popupHeight - (marginSize * 2))

	popupContent := popupStyle.Render(m.popupList.View())

	// Add hint text above popup
	hintStyle := lipgloss.NewStyle().
		Foreground(m.theme.MutedFg).
		Italic(true).
		Align(lipgloss.Center)

	hint := hintStyle.Render("↑↓ to navigate • Enter/o to select • ESC/q to close")

	// Combine hint and popup
	fullPopup := lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		hint,
		"",
		popupContent,
	)

	// Center the popup with margins to show it's floating
	return lipgloss.Place(
		m.termWidth,
		m.termHeight,
		lipgloss.Center,
		lipgloss.Center,
		fullPopup,
	)
}

// renderAddInputPopup shows the add worktree input prompt as a popup
func (m Model) renderAddInputPopup(background string) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(m.theme.Accent).
		Bold(true)

	messageStyle := lipgloss.NewStyle().
		Foreground(m.theme.TextFg)

	title := titleStyle.Render("Add Worktree")
	prompt := messageStyle.Render("Enter command (e.g., branch_name -p -s client-ui):")

	// Show error if present
	errorMsg := ""
	if m.addError != "" {
		errorStyle := lipgloss.NewStyle().
			Foreground(m.theme.ErrorFg).
			Bold(true)
		errorMsg = "\n" + errorStyle.Render("⚠ Error: "+m.addError) + "\n"
	}

	content := fmt.Sprintf("\n%s\n\n%s\n\n%s\n%s",
		title,
		prompt,
		m.addInput.View(),
		errorMsg)

	hintStyle := lipgloss.NewStyle().
		Foreground(m.theme.MutedFg).
		Italic(true).
		Align(lipgloss.Center)

	hint := hintStyle.Render("\nPress Enter to add • ESC to cancel")

	fullContent := content + hint

	popupWidth := m.termWidth * 2 / 3
	if popupWidth > 80 {
		popupWidth = 80
	}

	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Accent).
		Padding(1, 2).
		Width(popupWidth).
		Align(lipgloss.Left)

	popup := popupStyle.Render(fullContent)

	return lipgloss.Place(
		m.termWidth,
		m.termHeight,
		lipgloss.Center,
		lipgloss.Center,
		popup,
	)
}

// renderAddSpinnerPopup shows a spinner while adding worktree
func (m Model) renderAddSpinnerPopup() string {
	spinnerStyle := lipgloss.NewStyle().Foreground(m.theme.Accent).Bold(true)
	messageStyle := lipgloss.NewStyle().Foreground(m.theme.TextFg)

	content := fmt.Sprintf("\n  %s  %s\n",
		spinnerStyle.Render(m.spinner.View()),
		messageStyle.Render(fmt.Sprintf("Adding worktree: %s...", m.addingBranchName)))

	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Accent).
		Padding(1, 3).
		Align(lipgloss.Center)

	popup := popupStyle.Render(content)

	// Place popup and fill whitespace with spaces
	return lipgloss.Place(
		m.termWidth,
		m.termHeight,
		lipgloss.Center,
		lipgloss.Center,
		popup,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.NoColor{}),
	)
}

// helpText renders the help text with keymaps using styled key hints
func helpText(m Model) string {
	footerStyle := lipgloss.NewStyle().
		Foreground(m.theme.TextFg).
		Padding(0, 1)

	var hints []string
	if m.logsFocused {
		// Show log navigation hints when logs are focused
		hints = []string{
			m.renderKeyHint("j/k", "Scroll"),
			m.renderKeyHint("d/u", "Half page"),
			m.renderKeyHint("g/G", "Top/Bottom"),
			m.renderKeyHint("h", "Focus table"),
			m.renderKeyHint("q", "Quit"),
		}
	} else {
		// Show table navigation hints when table is focused
		hints = []string{
			m.renderKeyHint("↑/↓", "Navigate"),
			m.renderKeyHint("a", "Add"),
			m.renderKeyHint("o", "Open"),
			m.renderKeyHint("d", "Delete"),
			m.renderKeyHint("D", "Delete+Branch"),
			m.renderKeyHint("l", "Focus logs"),
			m.renderKeyHint("q", "Quit"),
		}
	}

	footerContent := "  " + lipgloss.JoinHorizontal(lipgloss.Left, hints...)
	return footerStyle.Render(footerContent)
}

// renderKeyHint renders a single key hint with pill/badge styling
func (m Model) renderKeyHint(key, label string) string {
	keyStyle := lipgloss.NewStyle().
		Foreground(m.theme.AccentFg).
		Background(m.theme.Accent).
		Bold(true).
		Padding(0, 1)
	labelStyle := lipgloss.NewStyle().Foreground(m.theme.Accent)
	return fmt.Sprintf("%s %s  ", keyStyle.Render(key), labelStyle.Render(label))
}
