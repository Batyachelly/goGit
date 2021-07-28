package branches

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"goGit/internal/commits"
	"goGit/internal/initrepo"
)

type IBranch interface {
	SetCommit(commits.ICommit) error
	GetCommit() commits.ICommit
	Name() string
}

type branch struct {
	name   string
	commit commits.ICommit
}

func (b branch) SetCommit(c commits.ICommit) error {
	branchPath := filepath.Join(initrepo.PathRefDir, b.name)
	return ioutil.WriteFile(branchPath, []byte(c.Hash()), 0644)
}

func NewBranch(name string, commit commits.ICommit) (IBranch, error) {
	branchPath := filepath.Join(initrepo.PathRefDir, name)
	if err := ioutil.WriteFile(branchPath, []byte(commit.Hash()), 0644); err != nil {
		return nil, err
	}

	return &branch{
		name:   name,
		commit: commit,
	}, nil
}

func GetBranch(name string) (IBranch, error) {
	branchPath := filepath.Join(initrepo.PathRefDir, name)
	commitHash, err := ioutil.ReadFile(branchPath)

	if os.IsNotExist(err) {
		return &branch{
			name:   name,
			commit: nil,
		}, nil
	} else if err != nil {
		return nil, err
	}

	c, err := commits.OpenCommit(string(commitHash))
	if err != nil {
		return nil, err
	}

	return &branch{
		name:   name,
		commit: c,
	}, nil
}

func (b branch) GetCommit() commits.ICommit {
	return b.commit
}

func (b branch) Name() string {
	return b.name
}
