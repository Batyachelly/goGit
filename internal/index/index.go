package index

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"goGit/internal/branches"
	"goGit/internal/commits"
	"goGit/internal/fsystem"
	"goGit/internal/initrepo"
)

type IIndex interface {
	Checkout(branchName string) error
	SetCurrentBranchCommit(commitHash string) error
	AddFile(path string) error
	Commit(name string) error
	Status() *bytes.Buffer
}

type Index struct {
	Objects map[string]string
}

func OpenIndex() (IIndex, error) {
	if _, err := os.Stat(initrepo.PathIndexFile); os.IsNotExist(err) {
		i := &Index{Objects: make(map[string]string)}
		if err := i.saveIndex(); err != nil {
			return nil, err
		}

		return i, nil
	}

	indexFile, err := ioutil.ReadFile(initrepo.PathIndexFile)
	if err != nil {
		return nil, err
	}

	dec := gob.NewDecoder(bytes.NewReader(indexFile))

	i := &Index{}
	if err := dec.Decode(i); err != nil {
		return nil, err
	}

	return i, nil
}

func (i *Index) Checkout(branchName string) error {
	commit, _, err := GetHead()
	if err != nil {
		return err
	}

	if _, err := os.Stat(filepath.Join(initrepo.PathRefDir, branchName)); os.IsNotExist(err) {
		_, err := branches.NewBranch(branchName, commit)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	ref := fmt.Sprintf("ref %s", branchName)
	if err := ioutil.WriteFile(initrepo.PathHeadFile, []byte(ref), 0644); err != nil {
		return err
	}

	return i.setCommit(commit)
}

func (i *Index) saveIndex() error {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	if err := enc.Encode(i); err != nil {
		return err
	}

	return ioutil.WriteFile(initrepo.PathIndexFile, buffer.Bytes(), 0644)
}

func (i *Index) SetCurrentBranchCommit(commitHash string) error {
	_, branch, err := GetHead()
	if err != nil {
		return err
	}

	if branch == nil {
		return nil
	}

	commit, err := commits.OpenCommit(commitHash)
	if err != nil {
		return err
	}

	if err := i.setCommit(commit); err != nil {
		return err
	}

	return branch.SetCommit(commit)
}

func (i *Index) setCommit(c commits.ICommit) error {
	objects, err := c.ToIndexObjects()
	if err != nil {
		return err
	}

	i.Objects = objects

	return i.saveIndex()
}

func (i *Index) AddFile(path string) error {
	fileShaHex, err := fsystem.AddFile(path)
	if err != nil {
		return err
	}

	for filePath := range i.Objects {
		if filePath != path {
			continue
		}

		if fileShaHex == "" {
			// Удаляем объект если такой есть в списке и его хэш пуст(т.е. файл отсутствует)
			delete(i.Objects, filePath)
		} else {
			i.Objects[filePath] = fileShaHex
		}

		return i.saveIndex()
	}

	i.Objects[path] = fileShaHex

	return i.saveIndex()
}

// Получить коммит и ветку(если она существует) на которую указывает HEAD.
func GetHead() (commits.ICommit, branches.IBranch, error) {
	target, err := ioutil.ReadFile(initrepo.PathHeadFile)
	if err != nil {
		return nil, nil, err
	}

	if strings.Contains(string(target), "ref ") {
		branchName := strings.Split(string(target), " ")[1]

		branch, err := branches.GetBranch(branchName)
		if err != nil {
			return nil, nil, err
		}

		return branch.GetCommit(), branch, nil
	}

	commit, err := commits.OpenCommit(string(target))
	if err != nil {
		return nil, nil, err
	}

	return commit, nil, nil
}

func Log() (*bytes.Buffer, error) {
	commit, _, err := GetHead()
	if err != nil {
		return nil, err
	}

	if commit == nil {
		return &bytes.Buffer{}, nil
	}

	var buffer bytes.Buffer

	for {
		buffer.WriteString(fmt.Sprintf("%s : %s\n", commit.Name(), commit.Hash()))

		commit, err = commit.Parent()
		if err != nil {
			return nil, err
		} else if commit == nil {
			break
		}
	}

	return &buffer, nil
}

func Deploy() error {
	commit, _, err := GetHead()
	if err != nil {
		return err
	}

	objects, err := commit.ToIndexObjects()
	if err != nil {
		return err
	}

	for filePath, fileHash := range objects {
		file, err := fsystem.GetObject(fileHash)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(filePath, file, 0644); err != nil {
			return err
		}
	}

	return nil
}

func (i *Index) Commit(name string) error {

	return nil
}

func (i *Index) Status() *bytes.Buffer {
	var b bytes.Buffer
	for filePath, fileHash := range i.Objects {
		b.WriteString(fmt.Sprintf("%s %s\n", fileHash[:8], filePath))
	}

	return &b
}
