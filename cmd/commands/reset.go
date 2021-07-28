package commands

import (
	"errors"
	"log"

	"goGit/internal/index"

	"github.com/spf13/cobra"
)

var errCommitHash = errors.New("commit hash is not specified")

func resetCmd(_ *cobra.Command, args []string) error {
	i, err := index.OpenIndex()
	if err != nil {
		log.Fatal(err)
	}

	if err := i.SetCurrentBranchCommit(args[0]); err != nil {
		log.Fatal(err)
	}

	if err := index.Deploy(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func resetCmdArgs(_ *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errCommitHash
	}

	return nil
}
