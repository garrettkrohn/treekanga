/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/garrettkrohn/treekanga/adapters"
	"github.com/garrettkrohn/treekanga/config"
	"github.com/garrettkrohn/treekanga/connector"
	"github.com/garrettkrohn/treekanga/shell"
)

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
	theme              *Theme
	// Add command state
	showAddInput       bool
	addInput           textinput.Model
	isAdding           bool
	addingBranchName   string
	addError           string
	// Dependencies
	git       adapters.GitAdapter
	zoxide    adapters.Zoxide
	connector connector.Connector
	shell     shell.Shell
	appConfig config.AppConfig
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

	return Model{
		table:        table,
		showPopup:    false,
		spinner:      spinner,
		isDeleting:   false,
		showAddInput: false,
		addInput:     ti,
		isAdding:     false,
		theme:        DefaultTheme(),
		git:          git,
		zoxide:       zoxide,
		connector:    conn,
		shell:        shell,
		appConfig:    appConfig,
	}
}
