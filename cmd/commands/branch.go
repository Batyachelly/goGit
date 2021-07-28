package commands

import (
	"log"

	"goGit/internal/index"

	"github.com/spf13/cobra"
)

func branchCmd(_ *cobra.Command, _ []string) error {
	commitHash, b, err := index.GetHead()
	if err != nil {
		return err
	}

	if b == nil {
		log.Println(commitHash)
	}

	log.Println(b.Name())

	return nil
}
