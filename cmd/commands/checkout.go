package commands

import (
	"errors"

	"goGit/internal/index"

	"github.com/spf13/cobra"
)

var errBranchName = errors.New("branch name is not specified")

func checkoutCmd(_ *cobra.Command, args []string) error {
	i, err := index.OpenIndex()
	if err != nil {
		return err
	}

	if err := i.Checkout(args[0]); err != nil {
		return err
	}

	return index.Deploy()
}

func checkoutCmdArgs(_ *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errBranchName
	}

	return nil
}
