package commands

import (
	"fmt"
	"log"

	"goGit/internal/index"

	"github.com/spf13/cobra"
)

func statusCmd(_ *cobra.Command, _ []string) error {
	i, err := index.OpenIndex()
	if err != nil {
		log.Fatal(err)
	}
	buffer := i.Status()

	fmt.Print(buffer.String())

	return nil
}
