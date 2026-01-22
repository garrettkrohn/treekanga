/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/charmbracelet/log"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/garrettkrohn/treekanga/transformer"
	"github.com/garrettkrohn/treekanga/utility"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type tuiModel struct {
	table              table.Model
	showPopup          bool
	popupList          list.Model
	termWidth          int
	termHeight         int
	spinner            spinner.Model
	isDeleting         bool
	deletingName       string
	showDeleteConfirm  bool
	deleteConfirmError string
	pendingDeletePath  string
	pendingDeleteName  string
	pendingBranchName  string
}

// popupItem represents an item in the popup list
type popupItem struct {
	title string
	desc  string
}

func (i popupItem) Title() string       { return i.title }
func (i popupItem) Description() string { return i.desc }
func (i popupItem) FilterValue() string { return i.title }

// deleteCompleteMsg is sent when deletion is complete
type deleteCompleteMsg struct {
	err          error
	worktreeName string
}

// deleteErrorMsg is sent when deletion fails
type deleteErrorMsg struct {
	err          error
	worktreePath string
	worktreeName string
	branchName   string
}

// performDelete performs the deletion in the background
func performDelete(worktreePath, worktreeName, branchName string, force bool, deleteBranch bool) tea.Cmd {
	return func() tea.Msg {
		// Add a minimum display time for the spinner
		startTime := time.Now()
		minDisplayTime := 1 * time.Second

		log.Debug("Removing worktree", "fullPath", worktreePath, "force", force)

		var err error
		if force {
			err = deps.Git.RemoveWorktreeForce(worktreePath, &worktreePath)
		} else {
			err = deps.Git.RemoveWorktree(worktreePath, &worktreePath)
		}

		if err != nil {
			// Ensure spinner shows for at least minDisplayTime before showing error
			elapsed := time.Since(startTime)
			if elapsed < minDisplayTime {
				time.Sleep(minDisplayTime - elapsed)
			}
			return deleteErrorMsg{err: err, worktreePath: worktreePath, worktreeName: worktreeName, branchName: branchName}
		}

		_ = deps.Zoxide.RemovePath(worktreePath)
		log.Debug("Worktree removed successfully")

		if deleteBranch {
			log.Debug("Deleting branch", "branchName", branchName)
			err = deps.Git.DeleteBranch(branchName, deps.BareRepoPath)
			if err != nil {
				log.Warn("Failed to delete branch", "branchName", branchName, "error", err)
			}
		}

		// Ensure spinner shows for at least minDisplayTime
		elapsed := time.Since(startTime)
		if elapsed < minDisplayTime {
			time.Sleep(minDisplayTime - elapsed)
		}

		return deleteCompleteMsg{err: nil, worktreeName: worktreeName}
	}
}

func (m tuiModel) Init() tea.Cmd { return m.spinner.Tick }

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle window size changes
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		m.table.SetHeight(msg.Height - 4) // Leave room for borders/padding
		return m, nil
	case deleteCompleteMsg:
		m.isDeleting = false
		// Rebuild the table with updated data
		rows, err := buildWorktreeTableRows()
		if err != nil {
			return m, tea.Printf("Error refreshing worktrees: %v", err)
		}
		m.table.SetRows(rows)
		return m, tea.Printf("Deleted worktree: %s", msg.worktreeName)
	case deleteErrorMsg:
		m.isDeleting = false
		m.showDeleteConfirm = true
		m.deleteConfirmError = msg.err.Error()
		m.pendingDeletePath = msg.worktreePath
		m.pendingDeleteName = msg.worktreeName
		m.pendingBranchName = msg.branchName
		return m, nil
	}

	// If delete confirmation is showing, handle it first
	if m.showDeleteConfirm {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "y", "Y":
				// User confirmed force delete
				m.showDeleteConfirm = false
				m.isDeleting = true
				m.deletingName = m.pendingDeleteName
				return m, performDelete(m.pendingDeletePath, m.pendingDeleteName, m.pendingBranchName, true, false)
			case "n", "N", "esc", "q":
				// User cancelled
				m.showDeleteConfirm = false
				m.deleteConfirmError = ""
				m.pendingDeletePath = ""
				m.pendingDeleteName = ""
				m.pendingBranchName = ""
				return m, nil
			}
		}
		return m, nil
	}

	// If popup is open, handle its inputs first
	if m.showPopup {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc", "q":
				m.showPopup = false
				return m, nil
			case "enter", "o":
				// Handle selection from popup
				selected := m.popupList.SelectedItem()
				if item, ok := selected.(popupItem); ok {
					m.showPopup = false
					// Connect to the selected path using sesh
					deps.Connector.SeshConnectWithString(item.title)
					return m, tea.Quit
				}
				m.showPopup = false
				return m, nil
			}
		}
		m.popupList, cmd = m.popupList.Update(msg)
		return m, cmd
	}

	// Normal table handling
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "d":
			selectedRow := m.table.SelectedRow()
			if len(selectedRow) < 3 {
				return m, tea.Printf("No worktree selected")
			}
			worktreePath := selectedRow[2]
			worktreeName := selectedRow[0]
			branchName := selectedRow[1]

			// Start the deletion process with spinner
			m.isDeleting = true
			m.deletingName = worktreeName
			return m, performDelete(worktreePath, worktreeName, branchName, false, false)
		case "D":
			selectedRow := m.table.SelectedRow()
			if len(selectedRow) < 3 {
				return m, tea.Printf("No worktree selected")
			}
			worktreePath := selectedRow[2]
			worktreeName := selectedRow[0]
			branchName := selectedRow[1]

			// Start the deletion process with spinner
			m.isDeleting = true
			m.deletingName = worktreeName
			return m, performDelete(worktreePath, worktreeName, branchName, false, true)
		case "o":
			selectedRow := m.table.SelectedRow()
			if len(selectedRow) < 3 {
				return m, tea.Printf("No worktree selected")
			}

			zoxideEntries, err := deps.Zoxide.GetQueryList(selectedRow[2])
			utility.CheckError(err)

			log.Info(zoxideEntries)

			items := getPopupItems(zoxideEntries)
			delegate := list.NewDefaultDelegate()
			delegate.SetSpacing(0)          // Remove spacing between items
			popupHeight := m.termHeight - 4 // Use most of the terminal height
			m.popupList = list.New(items, delegate, m.termWidth, popupHeight)
			m.popupList.Title = "Select a sesh to connect to"
			m.popupList.SetShowStatusBar(true)
			m.popupList.SetFilteringEnabled(false)
			m.showPopup = true
			return m, nil
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}

	// Update spinner if deleting
	if m.isDeleting {
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m tuiModel) View() string {
	// Show spinner if deleting
	if m.isDeleting {
		return fmt.Sprintf("\n\n   %s Deleting worktree: %s...\n\n", m.spinner.View(), m.deletingName)
	}

	// Show delete confirmation dialog
	if m.showDeleteConfirm {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

		confirmStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("196")).
			Padding(1, 2).
			Width(m.termWidth - 4)

		errorMsg := errorStyle.Render("Error deleting worktree:")
		confirmMsg := fmt.Sprintf("\n%s\n\n%s\n\nWorktree '%s' contains uncommitted changes.\nForce delete and discard all changes?\n\n[Y]es / [N]o",
			errorMsg,
			m.deleteConfirmError,
			m.pendingDeleteName)

		return "\n" + confirmStyle.Render(confirmMsg) + "\n"
	}

	if m.showPopup {
		// Render popup overlay with transparent background
		popupHeight := m.termHeight - 4
		popupStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2).
			Width(m.termWidth).
			Height(popupHeight)

		popup := popupStyle.Render(m.popupList.View())

		// Show just the popup, overlaying the table
		tableView := baseStyle.Render(m.table.View())

		// Create an overlay effect
		return tableView + "\n\n" + popup + "\n"
	}

	// Render help text with keymaps
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Padding(1, 0, 0, 2)

	helpText := helpStyle.Render(
		"Keymaps: [↑/↓] Navigate • [o] Open • [d] Delete Worktree • [D] Delete Worktree and Branch • [enter] Select • [q] Quit",
	)

	return baseStyle.Render(m.table.View()) + "\n" + helpText + "\n"
}

// getPopupItems returns the list of items to display in the popup
// You can customize this function to populate from any source
func getPopupItems(zoxideEntries []string) []list.Item {

	var returnItems []list.Item
	for _, item := range zoxideEntries { // iterate over your actual data
		returnItems = append(returnItems, popupItem{
			title: item,
			desc:  "", // add description if needed
		})
	}
	return returnItems
}

// tuiCmd represents the tui command
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Interactive TUI for managing worktrees",
	Long: `Launch an interactive terminal user interface (TUI) for managing worktrees.

    The TUI provides:
    - Interactive table view of all worktrees
    - Delete worktrees with the 'd' key
    - Connect to worktrees with the 'o' key
    - Navigate with arrow keys
    - Press 'q' to quit`,
	Run: func(cmd *cobra.Command, args []string) {
		columns := []table.Column{
			{Title: "Name", Width: 45},
			{Title: "Branch", Width: 45},
			{Title: "fullPath", Width: 70},
			{Title: "CommitHash", Width: 25},
		}

		rows, err := buildWorktreeTableRows()
		utility.CheckError(err)

		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(25),
		)

		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false)
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(false)
		t.SetStyles(s)

		// Initialize spinner
		sp := spinner.New()
		sp.Spinner = spinner.Dot
		sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

		m := tuiModel{
			table:      t,
			showPopup:  false,
			spinner:    sp,
			isDeleting: false,
		}
		p := tea.NewProgram(m, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	},
}

func buildWorktreeTableRows() ([]table.Row, error) {
	var rawWorktrees []string
	var err error

	if deps.BareRepoPath != "" {
		log.Debug("Using bare repo path for worktree list", "path", deps.BareRepoPath)
		rawWorktrees, err = deps.Git.GetWorktrees(&deps.BareRepoPath)
	} else {
		log.Debug("No bare repo path set, using current directory")
		rawWorktrees, err = deps.Git.GetWorktrees(nil)
	}

	if err != nil {
		return nil, err
	}

	worktreetransformer := transformer.NewTransformer()
	worktreeObjects := worktreetransformer.TransformWorktrees(rawWorktrees)

	// Sort worktrees by most recently modified
	sortWorktreesByModTime(worktreeObjects)

	var worktreeBranches []table.Row
	for _, worktree := range worktreeObjects {
		worktreeBranches = append(worktreeBranches, table.Row{worktree.Folder, worktree.BranchName, worktree.FullPath, worktree.CommitHash})
	}

	return worktreeBranches, nil
}

func init() {
	// Add flags here if needed
}
