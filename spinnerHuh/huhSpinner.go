package spinnerhuh

import (
	"github.com/charmbracelet/huh/spinner"
)

// HuhSpinner interface defines the methods for a spinner.
type HuhSpinner interface {
	Title(string) HuhSpinner
	Action(func()) HuhSpinner
	Run() error
}

// RealHuhSpinner is the concrete implementation of the Spinner interface.
type RealHuhSpinner struct {
	spinner *spinner.Spinner
}

// NewRealHuhSpinner creates and returns a new RealSpinner instance.
func NewRealHuhSpinner() *RealHuhSpinner {
	return &RealHuhSpinner{spinner: spinner.New()}
}

// Title sets the title of the spinner.
func (rs *RealHuhSpinner) Title(title string) HuhSpinner {
	rs.spinner.Title(title)
	return rs
}

// Action sets the action of the spinner.
func (rs *RealHuhSpinner) Action(action func()) HuhSpinner {
	rs.spinner.Action(action)
	return rs
}

// Run starts the spinner.
func (rs *RealHuhSpinner) Run() error {
	return rs.spinner.Run()
}
