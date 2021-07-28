package commands

import (
	"errors"
	"path/filepath"

	"goGit/internal/index"

	"github.com/spf13/cobra"
)

var errFilePath = errors.New("nothing specified, nothing added")

func addCmd(_ *cobra.Command, args []string) error {
	basePath, err := filepath.Abs(".")
	if err != nil {
		return err
	}

	filePath, err := filepath.Abs(args[0])
	if err != nil {
		return err
	}

	i, err := index.OpenIndex()
	if err != nil {
		return err
	}

	return i.AddFile(filePath[len(basePath)+1:])
}

func addCmdArgs(_ *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errFilePath
	}

	return nil
}
