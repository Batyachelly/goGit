package initrepo

import (
	"io/ioutil"
	"os"
)

const (
	PathIndexFile  = ".goGit/index"
	PathHeadFile   = ".goGit/HEAD"
	PathObjectsDir = ".goGit/objects/"
	PathRefDir     = ".goGit/refs/"
	defaultBranch  = "ref master"
)

func InitRepo() error {
	if err := os.MkdirAll(PathObjectsDir, 0771); err != nil {
		return err
	}

	if err := os.MkdirAll(PathRefDir, 0771); err != nil {
		return err
	}

	return ioutil.WriteFile(PathHeadFile, []byte(defaultBranch), 0771)
}
