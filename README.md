This is a cli application to manage git worktrees

# How to Install

`brew install garrettkrohn/treekanga/treekanga`

## Configuration

Create a yml file in this location:
`.config/treekanga/treekanga.yml`

```yaml
# example Configuration
repos:
  # This is the name of the repository
  exampleRepository:
    # This is the default branch that will be used in the add command if a baseBranch is not defined
    defaultBranch: development
    # This is a list of folders that will be added to zoxide
    zoxideFolders:
      - frontEnd
      - frontEnd/* # this will add all folders immediately within frontEnd
      - backEnd
      - backEnd/application/src/main/resources/db/migration
  treekanga:
    defaultBranch: main
    zoxideFolders:
      - cmd
      - git
```

## Commands

### Add

`treekanga add example_branch -b example_base_branch`

You can define the branch and base branch directly from the command line

`treekanga add example_branch`

logic for branche / worktrees:
- if example_branch exists locally: create a worktree with that branch
- if example_branch exists remotely: create a worktree with a new local version
  of that branch

`treekanga add example_branch -b base_branch`
error will be thrown if base_branch doesn't exist locally or remotely
- if the baseBranch exists locally: create a new example_branch off of the local
  base_branch
- if it exists locally and the pull flag (-p) is used: create a new 
  example_branch off of the remote base_branch
- if baseBranch doesn't exist locally: create new worktree with new
example_branch off of remote base_branch

`treekanga add`

You can also input via prompts

### List

`treekanga list`

List of all worktrees of repository

### Clean

`treekanga clean`

This will check what worktrees do not have a remote branch, ex. local
branches that have been merged and remote branch deleted

### Delete

`treekanga delete`

This will bring up all worktrees and allow you to select worktree(s)
to delete

### Clone

`treekanga clone https://www.github.com/example/example`

This will create a bare worktree for the given url
