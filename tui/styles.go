/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package tui

import "github.com/charmbracelet/lipgloss"

// baseTableStyle returns the base style for the table
func (m Model) baseTableStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Accent).
		Padding(0, 1)
}
