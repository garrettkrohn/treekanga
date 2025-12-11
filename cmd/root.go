package cmd

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/garrettkrohn/treekanga/connector"
	"github.com/garrettkrohn/treekanga/directoryReader"
	"github.com/garrettkrohn/treekanga/execwrap"
	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/logger"
	"github.com/garrettkrohn/treekanga/shell"
	"github.com/garrettkrohn/treekanga/zoxide"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Dependencies struct {
	Git             git.Git
	Zoxide          zoxide.Zoxide
	DirectoryReader directoryReader.DirectoryReader
	Connector       connector.Connector
	Shell           shell.Shell
	ResolvedRepo    string
}

var (
	deps     Dependencies
	logLevel string // Variable to store the log level
)

// resolveRepoName implements the fallback logic for determining the repo name
// 1. First tries to use the current directory name
// 2. If that doesn't exist in config, falls back to git.GetRepoName()
func resolveRepoName(git git.Git) string {
	// Get current working directory
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting working directory: ", err)
	}
	log.Debug("workingDir", workingDir)

	// Get directory name
	directoryName := filepath.Base(filepath.Dir(workingDir))
	log.Debug("directoryName", directoryName)

	// Check if directory name exists in viper config
	if viper.IsSet("repos." + directoryName) {
		log.Debug("Repo directory name found: ", "directory name", directoryName)
		return "repos." + directoryName
	}

	// Fallback to git.GetRepoName()
	repoName, err := git.GetRepoName(workingDir)
	if err != nil {
		log.Fatal("Error resolving repo name: ", err)
	}

	// Check if git repo name exists in viper config
	if viper.IsSet("repos." + repoName) {
		log.Debug("Repo git directory name found: ", directoryName)
		return "repos." + repoName
	}

	log.Fatal("No directory name, or git directory name found in the config")
	return ""
}

func NewRootCmd(git git.Git,
	zoxide zoxide.Zoxide,
	directoryReader directoryReader.DirectoryReader,
	sesh connector.Connector,
	shell shell.Shell,
	version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "treekanga",
		Short:   "CLI application to manage git worktree",
		Long:    `CLI application to manage git worktree`,
		Version: version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logger.LoggerInit(logLevel)

			resolvedRepo := resolveRepoName(git)

			deps = Dependencies{
				Git:             git,
				Zoxide:          zoxide,
				DirectoryReader: directoryReader,
				Connector:       sesh,
				Shell:           shell,
				ResolvedRepo:    resolvedRepo,
			}

		},
	}

	// Add the log level flag
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log", "l", "", "Set the log level (e.g., debug, info, warn, error)")

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {

	execWrap := execwrap.NewExec()
	shell := shell.NewShell(execWrap)
	git := git.NewGit(shell)
	zoxide := zoxide.NewZoxide(shell)
	connector := connector.NewConnector(shell)
	directoryReader := directoryReader.NewDirectoryReader()

	rootCmd := NewRootCmd(git, zoxide, directoryReader, connector, shell, version)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(cloneCmd)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Reserved for future flag and configuration settings
}
