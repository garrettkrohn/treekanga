/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/garrettkrohn/treekanga/adapters"
	"github.com/garrettkrohn/treekanga/config"
	"github.com/garrettkrohn/treekanga/connector"
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
	// Dependencies
	git       adapters.GitAdapter
	zoxide    adapters.Zoxide
	connector connector.Connector
	appConfig config.AppConfig
}

// NewModel creates a new TUI model with the given dependencies
func NewModel(
	table table.Model,
	spinner spinner.Model,
	git adapters.GitAdapter,
	zoxide adapters.Zoxide,
	conn connector.Connector,
	appConfig config.AppConfig,
) Model {
	return Model{
		table:      table,
		showPopup:  false,
		spinner:    spinner,
		isDeleting: false,
		theme:      DefaultTheme(),
		git:        git,
		zoxide:     zoxide,
		connector:  conn,
		appConfig:  appConfig,
	}
}
