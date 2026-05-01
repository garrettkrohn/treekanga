package models

// Connection represents an established connection to a session
type Connection struct {
	Session Session
	Found   bool // Whether the connection was found
	New     bool // Whether the session is new
}

// Session represents a tmux session or directory
type Session struct {
	Name string // The display name
	Path string // The absolute directory path
	Src  string // The source of the session (tmux, worktree, dir)
}

// ConnectOpts represents options for connecting to a session
type ConnectOpts struct {
	Switch bool // Whether to switch to the session (rather than attach)
}
