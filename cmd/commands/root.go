package commands

import (
	"log"

	"github.com/spf13/cobra"
)

func initServerCmd(root *cobra.Command) {
	root.AddCommand(
		&cobra.Command{
			Use:   "add [file path]",
			Short: "add object to index",
			RunE:  addCmd,
			Args:  addCmdArgs,
		},
		&cobra.Command{
			Use:   "branch",
			Short: "show current branch",
			RunE:  branchCmd,
		},
		&cobra.Command{
			Use:   "checkout [branch]",
			Short: "Runs publisher",
			RunE:  checkoutCmd,
			Args:  checkoutCmdArgs,
		},
		&cobra.Command{
			Use:   "commit [commit message]",
			Short: "commit index",
			RunE:  commitCmd,
			Args:  commitCmdArgs,
		},
		&cobra.Command{
			Use:   "init",
			Short: "init goGit repo",
			RunE:  initCmd,
		},
		&cobra.Command{
			Use:   "log",
			Short: "show log",
			RunE:  logCmd,
		},
		&cobra.Command{
			Use:   "reset [commit message]",
			Short: "reset branch to commit",
			RunE:  resetCmd,
			Args:  resetCmdArgs,
		},
		&cobra.Command{
			Use:   "status",
			Short: "show index status",
			RunE:  statusCmd,
		},
	)
}

func Execute() {
	rootCmd := &cobra.Command{
		Use: "goGit",
	}

	initServerCmd(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
