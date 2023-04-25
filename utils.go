package vfsgo

import (
	"errors"
	"os"
	"path/filepath"
)

func getProjRoot() (string, error) {
	projRoot, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(projRoot, "go.mod")); err == nil {
			break
		}

		if projRoot == filepath.Dir(projRoot) {
			return "", errors.New("can't find proj root")
		}

		projRoot = filepath.Dir(projRoot)
	}

	return projRoot, nil
}
