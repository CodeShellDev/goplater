package fsutils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func Exists(path string) (bool, os.FileInfo) {
	info, err := os.Stat(path)

	return !errors.Is(err, os.ErrNotExist), info
}

func IsDir(path string) bool {
	exists, info := Exists(path)

	if exists {
		return info.IsDir()
	}

	return false
}

func IsFile(path string) bool {
	exists, info := Exists(path)

	if exists {
		return !info.IsDir()
	}

	return false
}

func SafeRelative(base, path string) (string, error) {
	baseAbs, err := filepath.Abs(base)
	if err != nil {
		return "", err
	}

	pathAbs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	relPath, err := filepath.Rel(baseAbs, pathAbs)
	if err != nil {
		return "", err
	}

	return relPath, nil
}

func IsInside(file, dir string) bool {
	file, err := filepath.Abs(file)
	if err != nil {
		return false
	}

	dir, err = filepath.Abs(dir)
	if err != nil {
		return false
	}

	relPath, err := filepath.Rel(dir, file)
	if err != nil {
		return false
	}

	return relPath == "." || (relPath != ".." && !strings.HasPrefix(relPath, ".." + string(filepath.Separator)))
}