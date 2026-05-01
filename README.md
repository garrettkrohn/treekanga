# Treekanga

Treekanga is a powerful CLI tool for managing Git worktrees with ease. It simplifies the creation, management, and cleanup of worktrees, enhancing your Git workflow.

![GitHub release (latest by date)](https://img.shields.io/github/v/release/garrettkrohn/treekanga)
![License](https://img.shields.io/github/license/garrettkrohn/treekanga)

## Features

- Create new worktrees with smart branch handling
- List all worktrees in a repository
- Delete worktrees with stale branch filtering and interactive selector
- Clone repositories as bare worktrees
- Simple YAML configuration

## Installation

### Homebrew (macOS and Linux)

```bash
brew install garrettkrohn/treekanga/treekanga
```

## Configuration

Create a YAML configuration file at `~/.config/treekanga/treekanga.yml`:

```yaml
# Example Configuration

# Global selector mode (optional)
# Use "fzf" for fzf selector, or omit for built-in selector
selectorMode: fzf

repos:
  # Repository name or the parent of the bare repo
  exampleRepository:
    # Default branch used when no base branch is specified
    defaultBranch: development
    #where should treekanga put the worktrees, assumes starting in the $HOME directory
    worktreeTargetDir: /code 
    # Display mode for the list command: "branch" (default) or "directory"/"folder"
    # "branch" shows branch names, "directory" shows directory names
    listDisplayMode: branch
    # Folders to show with --all flag (subdirectories within worktrees)
    zoxideFolders:
      - frontEnd
      - frontEnd/* # adds all folders immediately within frontEnd
      - backEnd
      - backEnd/application/src/main/resources/db/migration
    postScript: ~/dotfiles/scripts/test_script.sh
    autoRunPostScript: false
    tuiTheme: catppuccin-mocha
  
  treekanga:
    bareRepoName: treekanga_bare
    defaultBranch: main
    listDisplayMode: directory
    zoxideFolders:
      - cmd
      - adapters
```

## Deprecated config options
```yaml
    bareRepoName: .bare # this was used to specify the name of the bare repo,
but now I am using git commands to find the bare repo, so the user does not need
to define it in the config.
```

## Usage

### Add a Worktree

Create a new worktree with a branch:

```bash
# Specify branch and base branch
treekanga add example_branch -b example_base_branch

# Use default base branch from config
treekanga add example_branch

# Pull the base branch before creating new branch
treekanga add example_branch -p

# Open in Cursor after creation
treekanga add example_branch -c

# Open in VS Code after creation
treekanga add example_branch -v

# Connect to tmux session at subdirectory (or use '.' for root)
treekanga add example_branch -t frontend
treekanga add example_branch -t .

# Specify custom directory for bare repo
treekanga add example_branch -d /path/to/bare/repo
```

Branch handling logic:
- If `example_branch` exists locally: Create a worktree with that branch
- If `example_branch` exists remotely: Create a worktree with a new local version of that branch
- With base branch specified:
  - If base branch exists locally: Create a new branch off the local base branch
  - If using the pull flag (`-p`): Create a new branch off the remote base branch
  - If base branch doesn't exist locally: Create new worktree with new branch off remote base branch

### List Worktrees

Display all worktrees in the current repository:

```bash
# List worktrees (display format based on config)
treekanga list

# Verbose output showing all details
treekanga list -v

# Show all worktrees plus subdirectories from zoxideFolders config
treekanga list --all
```

By default, the list command displays branch names. You can configure it to display directory names instead using the `listDisplayMode` configuration option:

- `branch` (default): Display branch names
- `directory` or `folder`: Display directory names

Example configuration:
```yaml
repos:
  myrepo:
    listDisplayMode: directory
```

The verbose flag (`-v`) will always show all details including both branch names and directory names, regardless of the configured display mode.

#### List All with Subdirectories

The `--all` or `-a` flag expands the list to include subdirectories within each worktree based on the `zoxideFolders` configuration. This is useful when you have a monorepo structure and want to quickly connect to specific subdirectories.

Example configuration:
```yaml
repos:
  platform:
    zoxideFolders:
      - parent
      - ui
      - backend/*  # Wildcard to include all folders in backend
```

With this configuration, `treekanga list --all` would show:
- `/code/platform_work` (worktree root)
- `/code/platform_work/parent`
- `/code/platform_work/ui`
- `/code/platform_work/backend/api`
- `/code/platform_work/backend/services`
- etc.

### Delete Worktrees

Interactive deletion of worktrees:

```bash
# Delete any worktrees interactively
treekanga delete

# Only show worktrees where branches don't exist on remote (stale worktrees)
treekanga delete --stale

# Also delete the local branches (use with caution)
treekanga delete --delete
```

### Clone a Repository

Clone a repository as a bare worktree:

```bash
treekanga clone https://www.github.com/example/example
```

### Connect to a Session

Connect to a tmux session using various strategies:

```bash
# Connect to an existing tmux session by name
treekanga connect my-session

# Connect to a worktree by name
treekanga connect feature-branch

# Connect to a directory (absolute or relative path)
treekanga connect ~/code/myproject
treekanga connect ./my-worktree

# Switch to a session when already inside tmux
treekanga connect my-session --switch
```

The connect command will automatically:
1. Check for an existing tmux session with the given name
2. Look for a worktree matching the name or path
3. Check if the input is a valid directory path
4. Create a new tmux session if none exists

### Interactive Selection Mode

The `connect` command supports interactive selection with the `--select` flag:

#### Flat Mode (All Worktrees)
```bash
treekanga connect --select
```
Shows all worktrees from all repos in your config. Select one to connect.

#### Hierarchical Mode (Pick Repo, Then Worktree)
```bash
treekanga connect --select --by-repo
```
First select a repo, then select a worktree within that repo.

#### Bare Repo Mode
```bash
treekanga connect --select --bare
```
Shows all bare repos and connects to the selected one.

### Selector Configuration

By default, treekanga uses a built-in interactive selector. If you prefer `fzf`, add this to your config:

```yaml
selectorMode: fzf
```

**Requirements**: fzf version 0.20.0 or higher in your PATH.

**Fallback**: If fzf is configured but not found, treekanga will warn and use the built-in selector.

### Troubleshooting

**"fzf not found in PATH" warning**:
- Install fzf: `brew install fzf` (macOS) or see [fzf installation](https://github.com/junegunn/fzf#installation)
- Or remove `selectorMode: fzf` from config to use built-in selector

### TUI (In Beta)

```bash
treekanga tui
```

### TUI Available Themes
```bash
"dracula"
"dracula-light"
"narna"
"clean-light"
"solarized-dark"
"solarized-light"
"gruvbox-dark"
"gruvbox-light"
"nord"
"monokai"
"catppuccin-mocha"
"catppuccin-latte"
"rose-pine-dawn"
"one-light"
"everforest-light"
"everforest-dark"
"modern"
"tokyo-night"
"one-dark"
"rose-pine"
"ayu-mirage"
"kanagawa"         
```

## Tmux Integration

Treekanga works great with tmux for quick worktree/bare repo switching. Add these keybinds to your `~/.config/tmux/tmux.conf`:

```bash
# Select worktree from all repos (flat view)
bind-key "w" run-shell "tmux popup -E -w 80% -h 90% 'treekanga connect --select --switch'"

# Select repo first, then worktree (hierarchical)
bind-key "W" run-shell "tmux popup -E -w 80% -h 90% 'treekanga connect --select --by-repo --switch'"

# Select bare repo
bind-key "b" run-shell "tmux popup -E -w 80% -h 90% 'treekanga connect --select --bare --switch'"
```

Then reload your config:
```bash
tmux source ~/.config/tmux/tmux.conf
```

**Usage**:
- Press `Ctrl-a w` to quickly switch worktrees
- Press `Ctrl-a W` for hierarchical selection
- Press `Ctrl-a b` to navigate to a bare repo

**Note**: Adjust `Ctrl-a` to your tmux prefix if different.

## Logging

Control log verbosity:

```bash
treekanga --log debug [command]
```

## Author

Garrett Krohn
