package utility

import (
	"log"
	"os"
)

// CheckError logs the error and terminates the program if the error is not nil.
func CheckError(err error) {
	if err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}
