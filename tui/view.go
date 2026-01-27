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
	// Render the base view (table + help text)
	baseView := m.baseTableStyle().Render(m.table.View()) + "\n" + helpText(m) + "\n"

	// Show spinner popup if deleting
	if m.isDeleting {
		return m.renderSpinnerPopup(baseView)
	}

	// Show delete confirmation dialog
	if m.showDeleteConfirm {
		return m.renderConfirmPopup(baseView)
	}

	// If popup is showing, render it with visual separation
	if m.showPopup {
		return m.renderModalPopup()
	}

	return baseView
}

// renderSpinnerPopup shows a small centered popup with spinner
func (m Model) renderSpinnerPopup(background string) string {
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

	return lipgloss.Place(
		m.termWidth,
		m.termHeight,
		lipgloss.Center,
		lipgloss.Center,
		popup,
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

// helpText renders the help text with keymaps using styled key hints
func helpText(m Model) string {
	footerStyle := lipgloss.NewStyle().
		Foreground(m.theme.TextFg).
		Padding(0, 1)

	hints := []string{
		m.renderKeyHint("↑/↓", "Navigate"),
		m.renderKeyHint("o", "Open"),
		m.renderKeyHint("d", "Delete"),
		m.renderKeyHint("D", "Delete+Branch"),
		m.renderKeyHint("enter", "Select"),
		m.renderKeyHint("q", "Quit"),
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
