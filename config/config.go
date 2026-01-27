package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

type AppConfig struct {
	BareRepoPath               string   // path to the bare repo, this is where the git commnand will be run from
	RepoNameForConfig          string   // this is the git project name, used to find the config
	ParentDirOfBareRepo        string   // this is an option for configuration to allow the user to have multiple configs for multiple instances of one project
	BaseBranch                 string   // default base branch
	WorktreeTargetDir          string   // this is where the added worktree will be
	ListDisplayMode            string   // branch or directory
	ZoxideFolders              []string // list of zoxide folders to be added
	PostScriptPath             string   // path to the post script to be run
	RunPostScript              bool     // run the post script without the execute flag
	PullBeforeCuttingNewBranch bool     // pull before cutting new branch

	// DELETE COMMAND
	FilterOnlyStaleBranches bool // only show branches that don't exist on remote
	DeleteBranch            bool // in addition to the worktree, delete the branch as well
	ForceDelete             bool // use --force when deleting

	// ADD COMMAND
	SeshConnect              string
	CursorConnect            bool
	VsCodeConnect            bool
	NewWorktreeName          string
	NewBranchName            string
	UseFormToSetBaseBranch   bool
	NewBranchExistsLocally   bool
	NewBranchExistsRemotely  bool
	BaseBranchExistsLocally  bool
	BaseBranchExistsRemotely bool
}

type Config interface {
	GetDefaultConfig(bareRepoPath string, projectName string) (AppConfig, error)
	ImportYamlConfigFile(cfg AppConfig) (AppConfig, error)
}

type ConfigInstance struct {
}

func NewConfig() Config {
	return &ConfigInstance{}
}

func (c *ConfigInstance) GetDefaultConfig(bareRepoPath string, projectName string) (AppConfig, error) {

	return AppConfig{
		BareRepoPath:               bareRepoPath,
		RepoNameForConfig:          projectName,
		ParentDirOfBareRepo:        filepath.Base(filepath.Dir(bareRepoPath)), // this produces just the base of the parent dir `/Users/gkrohn/code/treekanga_work/.bare` => `treekanga_work`
		BaseBranch:                 "development",
		WorktreeTargetDir:          "~",
		ListDisplayMode:            "branch",
		ZoxideFolders:              []string{},
		PostScriptPath:             "",
		RunPostScript:              false,
		PullBeforeCuttingNewBranch: false,
		FilterOnlyStaleBranches:    false,
		DeleteBranch:               false,
		ForceDelete:                false,
	}, nil
}

func getRepoConfigPrefix(repoNameForConfig string, parentDirOfBareRepo string) string {
	//1. check for a config with the project name
	repoConfig := viper.GetStringMap("repos." + repoNameForConfig)

	if len(repoConfig) > 0 {
		log.Debug(fmt.Sprintf("configuration found under repo name: %s", repoNameForConfig))
		return "repos." + repoNameForConfig + "."
	}

	//2. check for a config with parent of the bare repo
	repoConfig = viper.GetStringMap("repos." + parentDirOfBareRepo)

	if len(repoConfig) > 0 {
		log.Debug(fmt.Sprintf("configuration found under parent of bare directory name %s", parentDirOfBareRepo))

		return "repos." + parentDirOfBareRepo + "."
	}
	log.Fatal(fmt.Sprintf("no configuration could be found by repo name: %s or parent of bare directory name: %s", repoNameForConfig, parentDirOfBareRepo))

	return ""
}

func (c *ConfigInstance) ImportYamlConfigFile(cfg AppConfig) (AppConfig, error) {

	repoconfig := viper.GetStringMap("repos")

	if repoconfig == nil {
		log.Fatal("could not find configuration file")
	}

	viperRepoPrefix := getRepoConfigPrefix(cfg.RepoNameForConfig, cfg.ParentDirOfBareRepo)

	if viperRepoPrefix == "" {
		log.Fatal("error loading config")
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
			cfg.BaseBranch = defaultBranch
		}
	}

	if viper.IsSet(viperRepoPrefix + "worktreeTargetDir") {
		worktreeTargetDir := viper.GetString(viperRepoPrefix + "worktreeTargetDir")
		if worktreeTargetDir != "" {
			// Expand tilde to home directory
			homeDir, err := os.UserHomeDir()
			if err == nil {
				worktreeTargetDir = filepath.Join(homeDir, strings.TrimPrefix(worktreeTargetDir, "~/"))
			}
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
			log.Debug("setting autoRunPostScript: true from config")
			cfg.RunPostScript = autoRunPostScript
		}
	}

	return cfg, nil
}

func (cfg *AppConfig) Print() {
	log.Info("=== AppConfig ===")
	log.Info(fmt.Sprintf("BareRepoPath: %s", cfg.BareRepoPath))
	log.Info(fmt.Sprintf("RepoNameForConfig: %s", cfg.RepoNameForConfig))
	log.Info(fmt.Sprintf("DefaultBaseBranch: %s", cfg.BaseBranch))
	log.Info(fmt.Sprintf("WorktreeTargetDir: %s", cfg.WorktreeTargetDir))
	log.Info(fmt.Sprintf("ListDisplayMode: %s", cfg.ListDisplayMode))
	log.Info(fmt.Sprintf("ZoxideFolders: %v", cfg.ZoxideFolders))
	log.Info(fmt.Sprintf("PostScriptPath: %s", cfg.PostScriptPath))
	log.Info(fmt.Sprintf("AutoRunPostScript: %t", cfg.RunPostScript))
	log.Info(fmt.Sprintf("PullBeforeCuttingNewBranch: %t", cfg.PullBeforeCuttingNewBranch))
	log.Info(fmt.Sprintf("FilterOnlyStaleBranches: %t", cfg.FilterOnlyStaleBranches))
	log.Info(fmt.Sprintf("DeleteBranch: %t", cfg.DeleteBranch))
	log.Info(fmt.Sprintf("ForceDelete: %t", cfg.ForceDelete))
	log.Info("================")
}
