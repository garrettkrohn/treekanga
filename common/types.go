package common

import "github.com/garrettkrohn/treekanga/directoryReader"

type TreekangaAddConfig struct {
	Flags        AddCmdFlags
	Args         []string
	GitConfig    GitConfig
	WorkingDir   string
	ParentDir    string
	ZoxideConfig ZoxideConfig
}

type AddCmdFlags struct {
	Directory  *string
	BaseBranch *string
	Pull       *bool
	Connect    *string
}

type GitConfig struct {
	NewBranchName            string
	BaseBranchName           string
	RepoName                 string
	NumOfRemoteBranches      int
	NumOfLocalBranches       int
	NewBranchExistsLocally   bool
	NewBranchExistsRemotely  bool
	BaseBranchExistsLocally  bool
	BaseBranchExistsRemotely bool
	FolderPath               string
}

type ZoxideConfig struct {
	NewBranchName   string
	ParentDir       string
	FoldersToAdd    []string
	DirectoryReader directoryReader.DirectoryReader
}
