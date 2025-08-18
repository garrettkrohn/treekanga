package worktreeobj

type WorktreeObj struct {
	FullPath   string
	Folder     string // This is also the worktree name
	BranchName string
	CommitHash string
}
