package common

import (
	"path/filepath"

	"github.com/garrettkrohn/treekanga/directoryReader"
)

type AddConfig struct {
	// Input from user
	Args  []string
	Flags AddCmdFlags

	// Resolved paths and directories
	WorkingDir        string
	ParentDir         string
	WorktreeTargetDir string

	// Git repository information
	GitInfo GitInfo

	// External tool configurations
	ZoxideFolders   []string
	DirectoryReader directoryReader.DirectoryReader

	// Custom scripts
	PostScript        string
	AutoRunPostScript *bool
}

// Helper methods for the AddConfig struct
func (c *AddConfig) GetNewBranchName() string {
	return c.GitInfo.NewBranchName
}

func (c *AddConfig) GetBaseBranchName() string {
	return c.GitInfo.BaseBranchName
}

func (c *AddConfig) GetRepoName() string {
	return c.GitInfo.RepoName
}

func (c *AddConfig) GetWorktreePath() string {
	return c.WorktreeTargetDir
}

func (c *AddConfig) GetZoxideBasePath() string {
	return c.WorktreeTargetDir
}

func (c *AddConfig) GetZoxidePath(subFolder string) string {
	if subFolder != "" {
		return filepath.Join(c.WorktreeTargetDir, subFolder)
	}
	return c.WorktreeTargetDir
}

func (c *AddConfig) ShouldPull() bool {
	return c.Flags.Pull != nil && *c.Flags.Pull
}

func (c *AddConfig) ShouldOpenCursor() bool {
	return c.Flags.Cursor != nil && *c.Flags.Cursor
}

func (c *AddConfig) ShouldOpenVSCode() bool {
	return c.Flags.VsCode != nil && *c.Flags.VsCode
}

func (c *AddConfig) GetSeshTarget() string {
	if c.Flags.Sesh != nil {
		return *c.Flags.Sesh
	}
	return ""
}

func (c *AddConfig) HasSeshTarget() bool {
	return c.Flags.Sesh != nil && *c.Flags.Sesh != ""
}

func (c *AddConfig) HasPostScript() bool {
	if c.PostScript != "" {
		return true
	}
	return false
}

func (c *AddConfig) GetPostScript() string {
	return c.PostScript
}

type AddCmdFlags struct {
	Directory             *string
	BaseBranch            *string
	Pull                  *bool
	Sesh                  *string
	Cursor                *bool
	VsCode                *bool
	SpecifiedWorktreeName *string
	ExecuteScript         *bool
}

type GitInfo struct {
	NewBranchName            string
	BaseBranchName           string
	RepoName                 string
	NewBranchExistsLocally   bool
	NewBranchExistsRemotely  bool
	BaseBranchExistsLocally  bool
	BaseBranchExistsRemotely bool
}

// Legacy type alias for backward compatibility during transition
type GitConfig = GitInfo
type ZoxideConfig struct {
	NewBranchName   string
	ParentDir       string
	FoldersToAdd    []string
	DirectoryReader directoryReader.DirectoryReader
}
