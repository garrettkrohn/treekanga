package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/garrettkrohn/treekanga/filter"
	"github.com/garrettkrohn/treekanga/form"
	"github.com/garrettkrohn/treekanga/services"
	"github.com/garrettkrohn/treekanga/transformer"
	util "github.com/garrettkrohn/treekanga/utility"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete selected worktrees",
	Long: `Interactive deletion of worktrees with optional filtering.

    Lists all worktrees and allows you to select multiple to be deleted.
    
    Available flags:
    -s, --stale: Only show worktrees where branches don't exist on remote
    -d, --delete: CAUTION - Also delete the local branches
    -f, --force: CAUTION - Forces delete of worktree and branch`,
	Run: func(cmd *cobra.Command, args []string) {
		stale, err := cmd.Flags().GetBool("stale")
		util.CheckError(err)
		if stale {
			log.Debug("setting FilterOnlyStaleBranches = true from flags")
			deps.AppConfig.FilterOnlyStaleBranches = true
		}

		deleteBranches, err := cmd.Flags().GetBool("delete")
		util.CheckError(err)
		if deleteBranches {
			log.Debug("setting DeleteBranch = true from flags")
			deps.AppConfig.DeleteBranch = true
		}

		forceDelete, err := cmd.Flags().GetBool("force")
		util.CheckError(err)
		if forceDelete {
			log.Debug("setting ForceDelete = true from flags")
			deps.AppConfig.ForceDelete = true
		}

		numOfWorktreesRemoved, err := services.DeleteWorktrees(deps.Git,
			transformer.NewTransformer(),
			filter.NewFilter(),
			form.NewHuhForm(),
			deps.Zoxide,
			args,
			deps.AppConfig)
		if err != nil {
			cmd.PrintErrln("Error:", err)
			return
		}
		log.Info("worktrees removed", "count", numOfWorktreesRemoved)
	},
}

func init() {
	deleteCmd.Flags().BoolP("stale", "s", false, "Only show worktrees where the branches don't exist on remote")
	deleteCmd.Flags().BoolP("delete", "d", false, "CAUTION: delete the local branch")
	deleteCmd.Flags().BoolP("force", "f", false, "CAUTION: force delete the worktree and branch")
}
