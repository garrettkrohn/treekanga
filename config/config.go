package config

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/garrettkrohn/treekanga/git"
	"github.com/garrettkrohn/treekanga/utility"
	"github.com/spf13/viper"
)

type AppConfig struct {
	BareRepoPath               string   // path to the bare repo, this is where the git commnand will be run from
	RepoNameForConfig          string   // this is the git project name, used to find the config
	DefaultBaseBranch          string   // default base branch
	WorktreeTargetDir          string   // this is where the added worktree will be
	ListDisplayMode            string   // branch or directory
	ZoxideFolders              []string // list of zoxide folders to be added
	PostScriptPath             string   // path to the post script to be run
	AutoRunPostScript          bool     // run the post script without the execute flag
	PullBeforeCuttingNewBranch bool     // pull before cutting new branch
}

type Config interface {
	GetDefaultConfig() (AppConfig, error)
	ImportYamlConfigFile(cfg AppConfig) (AppConfig, error)
}

type ConfigInstance struct {
	git git.Git
}

func NewConfig(git git.Git) Config {
	return &ConfigInstance{git}
}

func (c *ConfigInstance) GetDefaultConfig() (AppConfig, error) {

	bareRepoPath, err := c.git.GetBareRepoPath()
	utility.CheckError(err)

	projectName, err := c.git.GetProjectName()
	utility.CheckError(err)

	return AppConfig{
		BareRepoPath:               bareRepoPath,
		RepoNameForConfig:          projectName,
		DefaultBaseBranch:          "development",
		WorktreeTargetDir:          "~",
		ListDisplayMode:            "branch",
		ZoxideFolders:              []string{},
		PostScriptPath:             "",
		AutoRunPostScript:          false,
		PullBeforeCuttingNewBranch: false,
	}, nil
}

// TODO: somehow the log level is not set
func (c *ConfigInstance) ImportYamlConfigFile(cfg AppConfig) (AppConfig, error) {

	var viperRepoPrefix string

	repoconfig := viper.GetStringMap("repos")

	if repoconfig == nil {
		// return error
	}

	repoConfig := viper.GetStringMap("repos." + cfg.RepoNameForConfig)

	if len(repoConfig) == 0 {
		log.Error("could not find repo config in configuration yaml")
		//TODO: need to support parent directory name as the config
		// try it here
	} else {
		log.Info(fmt.Sprintf("configuration found under repo name: %s", cfg.RepoNameForConfig))
		viperRepoPrefix = "repos." + cfg.RepoNameForConfig + "."
	}

	log.Info(fmt.Sprintf("All keys in repoconfig: %v", repoConfig))
	for key := range repoConfig {
		log.Info(fmt.Sprintf("Key: '%s' (len: %d)", key, len(key)))
	}

	if viper.IsSet(viperRepoPrefix + "autoPull") {
		autoPull := viper.GetBool(viperRepoPrefix + "autoPull")
		if autoPull {
			log.Debug("setting PullBeforeCuttingNewBranch = true from config")
			cfg.PullBeforeCuttingNewBranch = true
		}
	}

	if viper.IsSet(viperRepoPrefix + "defaultBranch") {
		defaultBranch := viper.GetString(viperRepoPrefix + "defaultBranch")
		if defaultBranch != "" {
			log.Debug(fmt.Sprintf("setting defaultBranch: %s from config", defaultBranch))
			cfg.DefaultBaseBranch = defaultBranch
		}
	}

	if viper.IsSet(viperRepoPrefix + "worktreeTargetDir") {
		worktreeTargetDir := viper.GetString(viperRepoPrefix + "worktreeTargetDir")
		if worktreeTargetDir != "" {
			log.Debug(fmt.Sprintf("setting worktreeTargetDir: %s from config", worktreeTargetDir))
			cfg.WorktreeTargetDir = worktreeTargetDir
		}
	}

	if viper.IsSet(viperRepoPrefix + "listDisplayMode") {
		listDisplayMode := viper.GetString(viperRepoPrefix + "listDisplayMode")
		if listDisplayMode != "" {
			log.Debug(fmt.Sprintf("setting listDisplayMode: %s from config", listDisplayMode))
			cfg.ListDisplayMode = listDisplayMode
		}
	}

	if viper.IsSet(viperRepoPrefix + "zoxideFolders") {
		zoxideFolders := viper.GetStringSlice(viperRepoPrefix + "zoxideFolders")
		if len(zoxideFolders) > 0 {
			log.Debug(fmt.Sprintf("setting zoxideFolders: %s from config", zoxideFolders))
			cfg.ZoxideFolders = zoxideFolders
		}
	}

	if viper.IsSet(viperRepoPrefix + "postScript") {
		postScript := viper.GetString(viperRepoPrefix + "postScript")
		if postScript != "" {
			log.Debug(fmt.Sprintf("setting postScript: %s from config", postScript))
			cfg.PostScriptPath = postScript
		}
	}

	if viper.IsSet(viperRepoPrefix + "autoRunPostScript") {
		autoRunPostScript := viper.GetBool(viperRepoPrefix + "autoRunPostScript")
		if autoRunPostScript {
			log.Debug(fmt.Sprintf("setting autoRunPostScript: %s from config", autoRunPostScript))
			cfg.AutoRunPostScript = autoRunPostScript
		}
	}

	return cfg, nil

}
