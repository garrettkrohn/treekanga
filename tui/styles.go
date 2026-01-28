/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package tui

import "github.com/charmbracelet/lipgloss"

// baseTableStyle returns the base style for the table
func (m Model) baseTableStyle() lipgloss.Style {
	// Change border color based on focus
	borderColor := m.theme.BorderDim
	if !m.logsFocused {
		borderColor = m.theme.Accent
	}
	
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1).
		Width(m.termWidth - 2)
}
