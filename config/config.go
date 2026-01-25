package config

type AppConfig struct {
	BareRepoPath string //path to the bare repo, this is where the git commnand will be run from
	RepoName     string // this is the git project name, used to find the config
}
