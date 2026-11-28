package fsutils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Exists(path string) (bool, os.FileInfo) {
	info, err := os.Stat(path)

	if err != nil {
		return false, nil
	}
	return true, info
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

// check for NUL byte in first 8000 bytes (same as git does)
func IsBinary(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	buf := make([]byte, 8000)

	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return false, err
	}

	return bytes.IndexByte(buf[:n], 0) != -1, nil
}

func WithAllowedRoot(target string, allowed []string, fn func(root *os.Root, relative string) error) error {
	targetAbs, err := filepath.Abs(target)
	if err != nil {
		return err
	}

	var bestRoot, bestRel string
	bestLen := -1

	for _, path := range allowed {
		allowedAbs, err := filepath.Abs(path)
		if err != nil {
			continue
		}

		rootDir := allowedAbs
		relative := ""

		info, err := os.Stat(allowedAbs)
		if err == nil && !info.IsDir() {
			// allowed file = target
			if targetAbs != allowedAbs {
				continue
			}

			rootDir = filepath.Dir(allowedAbs)
			relative = filepath.Base(allowedAbs)
		} else {
			relative, err = filepath.Rel(allowedAbs, targetAbs)
			if err != nil || !filepath.IsLocal(relative) {
				continue
			}
		}

		// pick broadest
		if bestLen == -1 || len(allowedAbs) < bestLen {
			bestRoot, bestRel, bestLen = rootDir, relative, len(allowedAbs)
		}
	}

	if bestLen == -1 {
		return fmt.Errorf("path %q is outside the allowed directories", targetAbs)
	}

	root, err := os.OpenRoot(bestRoot)
	if err != nil {
		return err
	}

	err = fn(root, bestRel)
	closeErr := root.Close()

	if err != nil {
		return err
	}
	return closeErr
}