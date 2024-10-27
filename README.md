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
      - backEnd
      - backeEnd/application/src/main/resources/db/migration
  treekanga:
    defaultBranch: main
    zoxideFolders:
      - cmd
      - git
```
