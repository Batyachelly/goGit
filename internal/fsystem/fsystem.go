package fsystem

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"

	"goGit/internal/initrepo"
)

func AddFile(path string) (string, error) {
	file, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		return "", nil
	} else if err != nil {
		return "", err
	}

	fileSha, err := AddObject(bytes.NewBuffer(file))
	if err != nil {
		return "", err
	}

	return fileSha, nil
}

func AddObject(buffer *bytes.Buffer) (string, error) {
	objSha := sha1.Sum(buffer.Bytes())
	objShaHex := hex.EncodeToString(objSha[:])

	if err := ioutil.WriteFile(filepath.Join(initrepo.PathObjectsDir, objShaHex), buffer.Bytes(), 0644); err != nil {
		return "", err
	}

	return objShaHex, nil
}

func GetObject(objectHash string) ([]byte, error) {
	file, err := ioutil.ReadFile(filepath.Join(initrepo.PathObjectsDir, objectHash))
	if err != nil {
		return nil, err
	}

	return file, nil
}
