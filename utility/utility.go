package util

import (
	"github.com/charmbracelet/huh/spinner"
	"log"
)

// CheckError logs the error and terminates the program if the error is not nil.
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func UseSpinner(title string, action func()) {
	err := spinner.New().
		Title(title).
		Action(action).
		Run()

	if err != nil {
		log.Fatal(err)
	}
}
