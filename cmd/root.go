package cmd

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/garrettkrohn/treekanga/config"
	"github.com/garrettkrohn/treekanga/connector"
	"github.com/garrettkrohn/treekanga/directoryReader"
	"github.com/garrettkrohn/treekanga/execwrap"
	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/logger"
	"github.com/garrettkrohn/treekanga/shell"
	"github.com/garrettkrohn/treekanga/utility"
	"github.com/garrettkrohn/treekanga/zoxide"
	"github.com/spf13/cobra"
)

type Dependencies struct {
	Git             git.Git
	Zoxide          zoxide.Zoxide
	DirectoryReader directoryReader.DirectoryReader
	Connector       connector.Connector
	Shell           shell.Shell
	AppConfig       config.AppConfig
}

var (
	deps     Dependencies
	logLevel string // Variable to store the log level
)

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

			deps = Dependencies{
				Git:             git,
				Zoxide:          zoxide,
				DirectoryReader: directoryReader,
				Connector:       sesh,
				Shell:           shell,
			}

			if cmd.Name() == "completion" || cmd.HasParent() && cmd.Parent().Name() == "completion" || cmd.Name() == "clone" {
				return
			}

			// get app config
			configuration := config.NewConfig(git)
			cfg, err := configuration.GetDefaultConfig()
			utility.CheckError(err)
			deps.AppConfig = cfg

			cfg, err = configuration.ImportYamlConfigFile(cfg)

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

	options := []fang.Option{
		fang.WithVersion(version),
	}

	if err := fang.Execute(context.Background(), rootCmd, options...); err != nil {
		os.Exit(1)
	}

}

func init() {
	// Reserved for future flag and configuration settings
}
