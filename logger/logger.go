/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package logger

import (
	"github.com/charmbracelet/log"
	"strings" // This is required to conpare the evironment variables
)

func LoggerInit(logLevel string) {
	switch strings.ToLower(logLevel) {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}
