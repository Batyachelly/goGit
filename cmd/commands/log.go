package commands

import (
	"log"

	"goGit/internal/index"

	"github.com/spf13/cobra"
)

func logCmd(_ *cobra.Command, _ []string) error {
	buffer, err := index.Log()
	if err != nil {
		return err
	}

	log.Println(buffer.String())

	return nil
}
