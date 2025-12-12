package utility

import (
	"github.com/charmbracelet/log"
)

// CheckError logs the error and terminates the program if the error is not nil.
func CheckError(err error) {
	if err != nil {
		log.Fatal("Error", "error", err)
	}
}
