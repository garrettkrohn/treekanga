/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package main

import (
	"fmt"
	"log"

	"github.com/garrettkrohn/treekanga/cmd"
	"github.com/spf13/viper"

	"io" // This is used to enable the multiwritter and be able to write to the log file and the console at the same time
	"log/slog"
	"os"
	"path"    // This is used to create the path where the log files will be stored
	"strings" // This is required to conpare the evironment variables
	"time"    // This is used to get the current date and create the log file
)

var version = "dev"

func initConfig() {
	viper.SetConfigName("treekanga")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/treekanga")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}
func loggerInit() {
	var f *os.File
	var err error
	fileOnly := false

	if f, err = createLoggerFile(); err != nil {
		slog.Error("Unable to create logger file", "error", err)
		os.Exit(1)
	}

	env := os.Getenv("ENV")
	handlerOptions := &slog.HandlerOptions{}

	switch strings.ToLower(env) {
	case "debug":
		handlerOptions.Level = slog.LevelDebug
	case "info":
		handlerOptions.Level = slog.LevelInfo
	case "error":
		handlerOptions.Level = slog.LevelError
	default:
		handlerOptions.Level = slog.LevelWarn
		fileOnly = true
	}

	var loggerHandler *slog.JSONHandler
	if !fileOnly {
		multiWriter := io.MultiWriter(os.Stdout, f)
		loggerHandler = slog.NewJSONHandler(multiWriter, handlerOptions)
	} else {
		loggerHandler = slog.NewJSONHandler(f, handlerOptions)
	}
	slog.SetDefault(slog.New(loggerHandler))
}

func createLoggerFile() (*os.File, error) {
	now := time.Now()
	date := fmt.Sprintf("%s.log", now.Format("2006-01-02"))

	// TempDir returns the default directory to use for temporary files.
	//
	// On Unix systems, it returns $TMPDIR if non-empty, else /tmp.
	// On Windows, it uses GetTempPath, returning the first non-empty
	// value from %TMP%, %TEMP%, %USERPROFILE%, or the Windows directory.
	// On Plan 9, it returns /tmp.
	userTempDir := os.TempDir()
	slog.Debug("createLoggerFile:", "userTempDir", userTempDir)

	if err := os.MkdirAll(path.Join(userTempDir, "treekanga"), 0755); err != nil {
		return nil, err
	}

	fileFullPath := path.Join(userTempDir, "treekanga", date)
	file, err := os.OpenFile(fileFullPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func main() {

	initConfig()

	loggerInit()

	cmd.Execute(version)
}
