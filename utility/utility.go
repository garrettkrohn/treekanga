package utility

import (
	"log"
)

// CheckError logs the error and terminates the program if the error is not nil.
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
