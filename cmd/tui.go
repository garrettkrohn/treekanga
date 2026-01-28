/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/garrettkrohn/treekanga/tui"
	"github.com/garrettkrohn/treekanga/utility"
)

// tuiCmd represents the tui command
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Interactive TUI for managing worktrees",
	Long: `Launch an interactive terminal user interface (TUI) for managing worktrees.

    The TUI provides:
    - Interactive table view of all worktrees
    - Real-time operation logs in the bottom pane
    - Add worktrees with the 'a' key
    - Delete worktrees with the 'd' key
    - Connect to worktrees with the 'o' key
    - Switch focus between panes with 'h' (table) and 'l' (logs)
    - Navigate with arrow keys or j/k (vim-style)
    - Press 'q' to quit`,
	Run: func(cmd *cobra.Command, args []string) {
		columns := []table.Column{
			{Title: "Name", Width: 45},
			{Title: "Branch", Width: 45},
			{Title: "fullPath", Width: 70},
			{Title: "CommitHash", Width: 25},
		}

		rows, err := tui.BuildWorktreeTableRows(deps.Git, deps.AppConfig)
		utility.CheckError(err)

		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(25),
		)

		// Apply Catppuccin theme colours to table
		theme := tui.Catppuccin()
		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(theme.BorderDim).
			BorderBottom(true).
			Foreground(theme.Cyan).
			Bold(true)
		s.Selected = s.Selected.
			Foreground(theme.AccentFg).
			Background(theme.Accent).
			Bold(true)
		t.SetStyles(s)

		// Initialize spinner with theme color
		sp := spinner.New()
		sp.Spinner = spinner.Dot
		sp.Style = lipgloss.NewStyle().Foreground(theme.Accent)

		m := tui.NewModel(t, sp, deps.Git, deps.Zoxide, deps.Connector, deps.Shell, deps.AppConfig)
		p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
		if _, err := p.Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	},
}

func init() {
	// Add flags here if needed
}
