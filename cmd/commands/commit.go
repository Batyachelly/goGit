package commands

import (
	"errors"

	"goGit/internal/index"

	"github.com/spf13/cobra"
)

var errCommitMessage = errors.New("commit message is not specified")

func commitCmd(_ *cobra.Command, args []string) error {
	i, err := index.OpenIndex()
	if err != nil {
		return err
	}

	return i.Commit(args[0])
}

func commitCmdArgs(_ *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errCommitMessage
	}

	return nil
}
