/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/garrettkrohn/treekanga/adapters"
	"github.com/garrettkrohn/treekanga/config"
	"github.com/garrettkrohn/treekanga/connector"
	"github.com/garrettkrohn/treekanga/models"
	"github.com/garrettkrohn/treekanga/shell"
)

// OperationLog represents a single operation in the history
type OperationLog struct {
	Timestamp time.Time
	Operation string // "add", "delete", etc.
	Target    string // branch name, worktree name, etc.
	Command   string // full command or details
	Status    string // "success", "error"
	Message   string // result or error message
}

// Model represents the TUI model
type Model struct {
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
	// Add command state
	showAddInput        bool
	addInput            textinput.Model
	isAdding            bool
	addingBranchName    string
	addingCommand       string
	addError            string
	showBranchSelection bool
	pendingAddInput     string
	pendingAddArgs      []string
	pendingAddConfig    config.AppConfig
	// Log viewer state
	logsFocused   bool
	logsViewport  viewport.Model
	operationLogs []OperationLog
	// Dependencies
	git       adapters.GitAdapter
	zoxide    adapters.Zoxide
	connector connector.Connector
	shell     shell.Shell
	appConfig config.AppConfig
}

// theme returns the theme from the app config
func (m Model) theme() *models.Theme {
	return m.appConfig.Theme
}

// NewModel creates a new TUI model with the given dependencies
func NewModel(
	table table.Model,
	spinner spinner.Model,
	git adapters.GitAdapter,
	zoxide adapters.Zoxide,
	conn connector.Connector,
	shell shell.Shell,
	appConfig config.AppConfig,
) Model {
	// Initialize text input for add command
	ti := textinput.New()
	ti.Placeholder = "branch_name -p -s client-ui"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 80

	// Initialize viewport for logs
	vp := viewport.New(80, 10)
	vp.SetContent("No operations logged yet.")

	return Model{
		table:               table,
		showPopup:           false,
		spinner:             spinner,
		isDeleting:          false,
		showAddInput:        false,
		addInput:            ti,
		isAdding:            false,
		showBranchSelection: false,
		logsFocused:         false,
		logsViewport:        vp,
		operationLogs:       []OperationLog{},
		git:                 git,
		zoxide:              zoxide,
		connector:           conn,
		shell:               shell,
		appConfig:           appConfig,
	}
}
