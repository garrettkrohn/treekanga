/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package tui

import (
	"bytes"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/garrettkrohn/treekanga/config"
	"github.com/garrettkrohn/treekanga/services"
	"github.com/garrettkrohn/treekanga/utility"
)

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update handles all events and updates the model state
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle window size changes
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height

		// Split the screen: 60% for table, 40% for logs
		tableHeight := (msg.Height * 6) / 10
		logsHeight := msg.Height - tableHeight - 8 // Account for borders and help text

		if logsHeight < 5 {
			logsHeight = 5
		}

		m.table.SetHeight(tableHeight)
		m.logsViewport.Width = msg.Width - 6 // Account for borders and padding
		m.logsViewport.Height = logsHeight

		return m, nil
	case deleteCompleteMsg:
		m.isDeleting = false
		// Log the success
		m.addOperationLog(OperationLog{
			Timestamp: time.Now(),
			Operation: "delete",
			Target:    msg.worktreeName,
			Command:   msg.worktreeName,
			Status:    "success",
			Message:   msg.output,
		})
		// Rebuild the table with updated data
		rows, err := BuildWorktreeTableRows(m.git, m.appConfig)
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
		// Log the error (will be updated if force delete succeeds)
		m.addOperationLog(OperationLog{
			Timestamp: time.Now(),
			Operation: "delete",
			Target:    msg.worktreeName,
			Command:   msg.worktreeName,
			Status:    "error",
			Message:   msg.output,
		})
		return m, nil
	case addCompleteMsg:
		m.isAdding = false
		if msg.err != nil {
			m.addError = msg.err.Error()
			m.showAddInput = true
			// Log the error
			m.addOperationLog(OperationLog{
				Timestamp: time.Now(),
				Operation: "add",
				Target:    msg.branchName,
				Command:   m.addingCommand,
				Status:    "error",
				Message:   msg.output,
			})
			return m, nil
		}
		// Log the success
		m.addOperationLog(OperationLog{
			Timestamp: time.Now(),
			Operation: "add",
			Target:    msg.branchName,
			Command:   m.addingCommand,
			Status:    "success",
			Message:   msg.output,
		})
		// Rebuild the table with updated data
		rows, err := BuildWorktreeTableRows(m.git, m.appConfig)
		if err != nil {
			return m, tea.Printf("Error refreshing worktrees: %v", err)
		}
		m.table.SetRows(rows)
		return m, tea.Printf("Added worktree: %s", msg.branchName)
	case addErrorMsg:
		m.isAdding = false
		m.addError = msg.err.Error()
		m.showAddInput = true
		// Log the error
		m.addOperationLog(OperationLog{
			Timestamp: time.Now(),
			Operation: "add",
			Target:    msg.branchName,
			Command:   m.addingCommand,
			Status:    "error",
			Message:   msg.output,
		})
		return m, nil
	}

	// If add input is showing, handle it first
	if m.showAddInput {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				// Parse and execute add command
				input := strings.TrimSpace(m.addInput.Value())
				if input == "" {
					return m, nil
				}
				m.showAddInput = false
				m.isAdding = true
				m.addingBranchName = parseFirstArg(input)
				m.addingCommand = input
				m.addError = ""
				m.addInput.SetValue("") // Clear input for next time
				return m, tea.Batch(m.performAdd(input), m.spinner.Tick)
			case "esc", "ctrl+c":
				// User cancelled
				m.showAddInput = false
				m.addError = ""
				m.addInput.SetValue("")
				return m, nil
			}
		}
		m.addInput, cmd = m.addInput.Update(msg)
		return m, cmd
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
				return m, tea.Batch(m.performDelete(m.pendingDeletePath, m.pendingDeleteName, m.pendingBranchName, true, false), m.spinner.Tick)
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
					m.connector.SeshConnect(item.title)
					return m, tea.Quit
				}
				m.showPopup = false
				return m, nil
			}
		}
		m.popupList, cmd = m.popupList.Update(msg)
		return m, cmd
	}

	// Handle focus switching and navigation
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "h":
			// Focus table (up/left)
			if m.logsFocused {
				m.logsFocused = false
				m.table.Focus()
				return m, nil
			}
		case "l":
			// Focus logs (down/right)
			if !m.logsFocused {
				m.logsFocused = true
				m.table.Blur()
				return m, nil
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}

		// If logs are focused, handle log navigation
		if m.logsFocused {
			switch msg.String() {
			case "j", "down":
				m.logsViewport.LineDown(1)
				return m, nil
			case "k", "up":
				m.logsViewport.LineUp(1)
				return m, nil
			case "d", "ctrl+d":
				m.logsViewport.HalfViewDown()
				return m, nil
			case "u", "ctrl+u":
				m.logsViewport.HalfViewUp()
				return m, nil
			case "g":
				m.logsViewport.GotoTop()
				return m, nil
			case "G":
				m.logsViewport.GotoBottom()
				return m, nil
			}
			return m, nil
		}

		// Table is focused - handle table operations
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "a":
			// Show add input prompt
			m.showAddInput = true
			m.addInput.Focus()
			return m, nil
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
			return m, tea.Batch(m.performDelete(worktreePath, worktreeName, branchName, false, false), m.spinner.Tick)
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
			return m, tea.Batch(m.performDelete(worktreePath, worktreeName, branchName, false, true), m.spinner.Tick)
		case "o":
			selectedRow := m.table.SelectedRow()
			if len(selectedRow) < 3 {
				return m, tea.Printf("No worktree selected")
			}

			zoxideEntries, err := services.GetQueryList(m.zoxide, selectedRow[2])
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

	// Update spinner if deleting or adding
	if m.isDeleting || m.isAdding {
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// performDelete performs the deletion in the background
func (m Model) performDelete(worktreePath, worktreeName, branchName string, force bool, deleteBranch bool) tea.Cmd {
	return func() tea.Msg {
		// Add a minimum display time for the spinner
		startTime := time.Now()
		minDisplayTime := 1 * time.Second

		// Capture log output - write ONLY to buffer, not to stderr
		var logBuffer bytes.Buffer
		log.SetOutput(&logBuffer)

		log.Debug("Removing worktree", "fullPath", worktreePath, "force", force)

		var err error
		if force {
			err = m.git.RemoveWorktree(worktreeName, &worktreePath, true)
		} else {
			err = m.git.RemoveWorktree(worktreeName, &worktreePath, false)
		}

		if err != nil {
			log.SetOutput(os.Stderr)
			output := logBuffer.String()
			// Ensure spinner shows for at least minDisplayTime before showing error
			elapsed := time.Since(startTime)
			if elapsed < minDisplayTime {
				time.Sleep(minDisplayTime - elapsed)
			}
			return deleteErrorMsg{
				err:          err,
				worktreePath: worktreePath,
				worktreeName: worktreeName,
				branchName:   branchName,
				output:       output,
			}
		}

		_ = m.zoxide.RemovePath(worktreePath)
		log.Debug("Worktree removed successfully")

		if deleteBranch {
			log.Debug("Deleting branch", "branchName", branchName)
			err = m.git.DeleteBranchRef(branchName, m.appConfig.BareRepoPath)
			if force {
				err = m.git.DeleteBranch(branchName, m.appConfig.BareRepoPath, true)
			} else {
				err = m.git.DeleteBranch(branchName, m.appConfig.BareRepoPath, false)
			}
			if err != nil {
				log.Warn("Failed to delete branch", "branchName", branchName, "error", err)
			}
		}

		// Restore stderr as log output
		log.SetOutput(os.Stderr)
		output := logBuffer.String()

		// Ensure spinner shows for at least minDisplayTime
		elapsed := time.Since(startTime)
		if elapsed < minDisplayTime {
			time.Sleep(minDisplayTime - elapsed)
		}

		return deleteCompleteMsg{
			err:          nil,
			worktreeName: worktreeName,
			output:       output,
		}
	}
}

// performAdd performs the add worktree operation in the background
func (m Model) performAdd(input string) tea.Cmd {
	return func() tea.Msg {
		// Add a minimum display time for the spinner
		startTime := time.Now()
		minDisplayTime := 1 * time.Second

		// Parse the input string into args and flags
		args, cfg, err := parseAddCommand(input, m.appConfig)
		if err != nil {
			return addCompleteMsg{err: err, branchName: parseFirstArg(input)}
		}

		if len(args) == 0 {
			return addCompleteMsg{err: nil, branchName: ""}
		}

		log.Debug("Adding worktree", "input", input, "branch", args[0])

		// Configure the add service
		cfg = services.SetConfigForAddService(m.git, cfg, args)

		// Capture log output - write ONLY to buffer, not to stderr
		var logBuffer bytes.Buffer
		log.SetOutput(&logBuffer)

		// Call the add service with panic recovery
		var addErr error
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Error("Panic during add worktree", "error", r)
					addErr = err
				}
			}()
			services.AddWorktree(m.git, m.zoxide, m.connector, m.shell, cfg)
		}()

		// Restore stderr as log output
		log.SetOutput(os.Stderr)

		// Capture the output
		output := logBuffer.String()

		log.Debug("Worktree added successfully")

		// Ensure spinner shows for at least minDisplayTime
		elapsed := time.Since(startTime)
		if elapsed < minDisplayTime {
			time.Sleep(minDisplayTime - elapsed)
		}

		return addCompleteMsg{
			err:        addErr,
			branchName: cfg.NewBranchName,
			output:     output,
		}
	}
}

// parseFirstArg extracts the first non-flag argument from input
func parseFirstArg(input string) string {
	parts := strings.Fields(input)
	for _, part := range parts {
		if !strings.HasPrefix(part, "-") {
			return part
		}
	}
	return ""
}

// parseAddCommand parses the input string and builds config based on flags
func parseAddCommand(input string, baseConfig config.AppConfig) ([]string, config.AppConfig, error) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil, baseConfig, nil
	}

	// Create a copy of the config to modify
	cfg := baseConfig

	// Extract the branch name (first non-flag argument)
	var branchName string
	var args []string
	i := 0

	for i < len(parts) {
		part := parts[i]

		// Handle flags
		if strings.HasPrefix(part, "-") {
			switch part {
			case "-p", "--pull":
				cfg.PullBeforeCuttingNewBranch = true
			case "-c", "--cursor":
				cfg.CursorConnect = true
			case "-v", "--vscode":
				cfg.VsCodeConnect = true
			case "-x", "--script":
				cfg.RunPostScript = true
			case "-f", "--from":
				cfg.UseFormToSetBaseBranch = true
			case "-s", "--sesh":
				// Next part should be the sesh value
				if i+1 < len(parts) && !strings.HasPrefix(parts[i+1], "-") {
					cfg.SeshConnect = parts[i+1]
					i++
				}
			case "-b", "--base":
				// Next part should be the base branch
				if i+1 < len(parts) && !strings.HasPrefix(parts[i+1], "-") {
					cfg.BaseBranch = parts[i+1]
					i++
				}
			case "-d", "--directory":
				// Next part should be the directory
				if i+1 < len(parts) && !strings.HasPrefix(parts[i+1], "-") {
					cfg.WorktreeTargetDir = parts[i+1]
					i++
				}
			case "-n", "--name":
				// Next part should be the worktree name
				if i+1 < len(parts) && !strings.HasPrefix(parts[i+1], "-") {
					cfg.NewWorktreeName = parts[i+1]
					i++
				}
			}
		} else if branchName == "" {
			// First non-flag argument is the branch name
			branchName = part
		}
		i++
	}

	if branchName == "" {
		return nil, cfg, nil
	}

	args = []string{branchName}
	return args, cfg, nil
}

// addOperationLog adds a log entry to the operation history (max 100 entries)
func (m *Model) addOperationLog(log OperationLog) {
	m.operationLogs = append(m.operationLogs, log)
	// Keep only the last 100 logs
	if len(m.operationLogs) > 100 {
		m.operationLogs = m.operationLogs[1:]
	}
	// Update the viewport content
	m.updateLogsViewport()
}

// updateLogsViewport updates the viewport content with current logs
func (m *Model) updateLogsViewport() {
	if len(m.operationLogs) == 0 {
		m.logsViewport.SetContent("No operations logged yet.")
		return
	}

	var content strings.Builder

	// Show logs in reverse order (most recent first)
	for i := len(m.operationLogs) - 1; i >= 0; i-- {
		log := m.operationLogs[i]

		// Status indicator
		statusIcon := "✓"
		if log.Status == "error" {
			statusIcon = "✗"
		}

		// Format timestamp
		timestamp := log.Timestamp.Format("15:04:05")

		// Build command line: timestamp + status + treekanga operation [command or target]
		commandLine := timestamp + " " + statusIcon + " treekanga " + log.Operation

		// Use command if available (includes flags), otherwise use target
		if log.Command != "" {
			commandLine += " " + log.Command
		} else if log.Target != "" {
			commandLine += " " + log.Target
		}

		content.WriteString(commandLine)
		content.WriteString("\n")

		// Output section
		if log.Message != "" {
			cleanMsg := strings.TrimSpace(log.Message)
			if cleanMsg != "" {
				content.WriteString(cleanMsg)
				content.WriteString("\n")
			}
		}

		// Separator between entries
		content.WriteString("\n")
	}

	m.logsViewport.SetContent(content.String())
}
