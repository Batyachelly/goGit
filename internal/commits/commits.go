package commits

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"goGit/internal/fsystem"
	"goGit/internal/initrepo"
)

type ICommit interface {
	Save() (string, error)
	ToIndexObjects() (map[string]string, error)

	Tree() string
	Parent() (ICommit, error)
	Name() string
	Hash() string
}

type commit struct {
	tree, name, hash string
	parentHash       string
}

type ITree interface {
	ToObject() (string, error)
}

type tree struct {
	trees map[string]*tree
	blobs map[string]string
}

func OpenCommit(commitHash string) (ICommit, error) {
	file, err := fsystem.GetObject(commitHash)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(file), "\n")

	return &commit{
		tree:       strings.Split(lines[0], " ")[1],
		parentHash: strings.Split(lines[1], " ")[1],
		name:       strings.Split(lines[2], " ")[1],
		hash:       commitHash,
	}, nil
}

func NewCommit(tree string, parent ICommit, name string) (ICommit, error) {
	var parentHash string

	if parent == nil {
		parentHash = "nil"
	} else {
		parentHash = parent.Hash()
	}

	c := commit{
		tree:       tree,
		parentHash: parentHash,
		name:       name,
	}

	commitHash, err := c.Save()
	if err != nil {
		return nil, err
	}

	c.hash = commitHash

	return &c, nil
}

func (c commit) Save() (string, error) {
	var b bytes.Buffer

	currentTime := time.Now()

	b.WriteString(fmt.Sprintf("tree %s\n", c.tree))
	b.WriteString(fmt.Sprintf("parent %s\n", c.parentHash))
	b.WriteString(fmt.Sprintf("name %s\n", c.name))
	b.WriteString(fmt.Sprintf("%s\n", currentTime.Format(time.UnixDate)))

	return fsystem.AddObject(&b)
}

// Превращает древо объектов в список объектов индекса
func (c *commit) ToIndexObjects() (map[string]string, error) {
	objects := make(map[string]string)

	if err := treeToObjects(c.Tree(), "", objects); err != nil {
		return nil, err
	}

	return objects, nil
}

func (c commit) Tree() string {
	return c.tree
}

func (c commit) Parent() (ICommit, error) {
	p, err := OpenCommit(c.parentHash)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return p, nil
}

func (c commit) Name() string {
	return c.name
}

func (c commit) Hash() string {
	return c.hash
}

func treeToObjects(treeHash string, path string, objects map[string]string) error {
	file, err := ioutil.ReadFile(filepath.Join(initrepo.PathObjectsDir, treeHash))
	if err != nil {
		return err
	}

	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		lineParts := strings.Split(line, " ")

		switch lineParts[0] {
		case "tree":
			if err := treeToObjects(lineParts[1], filepath.Join(path, lineParts[2]), objects); err != nil {
				return err
			}
		case "blob":
			objects[filepath.Join(path, lineParts[2])] = lineParts[1]
		}
	}

	return nil
}

// Превращает список [(хэш_файла, путь_к_файлу), ...] в граф состоящий из структур типа ITree.
func MakeGraphTree(objects map[string]string) ITree {
	t := &tree{blobs: make(map[string]string), trees: make(map[string]*tree)}

	for filePath, fileHash := range objects {
		dir, file := filepath.Split(filePath)

		if dir == "" {
			t.blobs[filePath] = fileHash

			continue
		}

		prevT := t
		for _, d := range strings.Split(filepath.Clean(dir), string(filepath.Separator)) {
			if t, ok := prevT.trees[d]; ok {
				prevT = t
			} else {
				prevT.trees[d] = &tree{blobs: map[string]string{}, trees: map[string]*tree{}}
				prevT = prevT.trees[d]
			}
		}

		prevT.blobs[file] = fileHash
	}

	return t
}

// Превратить в объект текущее древо.
// Предварительно рекурсивно получаются объекты всех деревьев на которые указывает текущее дерево.
func (n *tree) ToObject() (string, error) {
	var treeObject bytes.Buffer

	for dir, t := range n.trees { // каждое дерево на которое указывает текущее дерево превращаем в объект
		treeHash, err := t.ToObject()
		if err != nil {
			return "", err
		}

		treeObject.WriteString(fmt.Sprintf("tree %s %s\n", treeHash, dir))
	}

	for file, blobHash := range n.blobs {
		treeObject.WriteString(fmt.Sprintf("blob %s %s\n", blobHash, file))
	}

	return fsystem.AddObject(&treeObject)
}
