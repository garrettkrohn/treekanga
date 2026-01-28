# Treekanga

Treekanga is a powerful CLI tool for managing Git worktrees with ease. It simplifies the creation, management, and cleanup of worktrees, enhancing your Git workflow.

![GitHub release (latest by date)](https://img.shields.io/github/v/release/garrettkrohn/treekanga)
![License](https://img.shields.io/github/license/garrettkrohn/treekanga)

## Features

- Create new worktrees with smart branch handling
- List all worktrees in a repository
- Delete worktrees with stale branch filtering and interactive selector
- Clone repositories as bare worktrees
- Zoxide integration for quick navigation
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
    # Folders to register with zoxide for quick navigation
    zoxideFolders:
      - frontEnd
      - frontEnd/* # adds all folders immediately within frontEnd
      - backEnd
      - backEnd/application/src/main/resources/db/migration
    postScript: ~/dotfiles/scripts/test_script.sh
    autoRunPostScript: false
  
  treekanga:
    bareRepoName: treekanga_bare
    defaultBranch: main
    listDisplayMode: directory
    zoxideFolders:
      - cmd
      - git
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

# Connect to a sesh session after creation
treekanga add example_branch -s session_name

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

### TUI (In Beta)

```bash
treekanga tui
```

## Logging

Control log verbosity:

```bash
treekanga --log debug [command]
```

## Author

Garrett Krohn
