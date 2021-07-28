package commands

import (
	"goGit/internal/initrepo"

	"github.com/spf13/cobra"
)

func initCmd(_ *cobra.Command, _ []string) error {
	return initrepo.InitRepo()
}
