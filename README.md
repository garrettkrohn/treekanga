# Treekanga

Treekanga is a powerful CLI tool for managing Git worktrees with ease. It simplifies the creation, management, and cleanup of worktrees, enhancing your Git workflow.

![GitHub release (latest by date)](https://img.shields.io/github/v/release/garrettkrohn/treekanga)
![License](https://img.shields.io/github/license/garrettkrohn/treekanga)

## Features

- Create new worktrees with smart branch handling
- List all worktrees in a repository
- Clean up stale worktrees
- Delete worktrees with an interactive selector
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
  # Repository name
  exampleRepository:
    # Default branch used when no base branch is specified
    defaultBranch: development
    #where should treekanga put the worktrees, assumes starting in the $HOME directory
    worktreeTargetDir: /code 
    # Folders to register with zoxide for quick navigation
    zoxideFolders:
      - frontEnd
      - frontEnd/* # adds all folders immediately within frontEnd
      - backEnd
      - backEnd/application/src/main/resources/db/migration
  
  treekanga:
    defaultBranch: main
    zoxideFolders:
      - cmd
      - git
```

## Usage

### Add a Worktree

Create a new worktree with a branch:

```bash
# Specify branch and base branch
treekanga add example_branch -b example_base_branch

# Use default base branch from config
treekanga add example_branch
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
treekanga list
```

### Clean Worktrees

Remove worktrees that no longer have a corresponding remote branch:

```bash
treekanga clean
```

### Delete Worktrees

Interactive deletion of worktrees:

```bash
treekanga delete
```

### Clone a Repository

Clone a repository as a bare worktree:

```bash
treekanga clone https://www.github.com/example/example
```

## Logging

Control log verbosity:

```bash
treekanga --log debug [command]
```

## Author

Garrett Krohn
