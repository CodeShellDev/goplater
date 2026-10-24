package outputresolver

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/codeshelldev/goplater/utils/fsutils"
)
var (
	ErrDirToFile      	= errors.New("directory source needs directory output")
	ErrOutsideRoot    	= errors.New("source is outside root")
	ErrNoSharedParent 	= errors.New("source and output share no parent, use flatten or set root")
	ErrEmptyRoot      	= errors.New("root is empty")
)

func hasTrailingSep(p string) bool {
	return strings.HasSuffix(p, "/") || strings.HasSuffix(p, string(filepath.Separator))
}

// trailing slash or exists as dir
func IsDirPath(p string) bool {
	if hasTrailingSep(p) || p == "." || p == ".." {
		return true
	}

	return fsutils.IsDir(p)
}

func isOutside(rel string) bool {
	return rel == ".." || strings.HasPrefix(rel, ".." + string(filepath.Separator))
}

func isWithin(parent, child string) bool {
	rel, err := filepath.Rel(parent, child)
	return err == nil && !isOutside(rel)
}

func isVolumeRoot(p string) bool {
	return filepath.Dir(p) == p
}

// walk up from a-Path until dir contains b ("" if none)
func CommonParent(a, b string) string {
	for p := a; ; p = filepath.Dir(p) {
		if isWithin(p, b) {
			return p
		}

		if isVolumeRoot(p) {
			return ""
		}
	}
}

// doesn't preserve structure ("dumps")
func ResolveOutput(source, output string) string {
	if !IsDirPath(output) {
		return output
	}

	return filepath.Join(output, filepath.Base(filepath.Clean(source)))
}

// preserve output judging by root
func ResolveOutputPreserved(root, source, output string) (string, error) {
	isDirSource := IsDirPath(source)

	// file output => rename or overwrite
	if !IsDirPath(output) {
		if isDirSource {
			return "", ErrDirToFile
		}

		return filepath.Abs(output)
	}

	if root == "" {
		return "", ErrEmptyRoot
	}

	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}

	sourceAbs, err := filepath.Abs(source)
	if err != nil {
		return "", err
	}

	rel, err := filepath.Rel(rootAbs, sourceAbs)
	if err != nil {
		return "", err
	}

	if isOutside(rel) {
		return "", errors.New(ErrOutsideRoot.Error() + ": " + source + " not in " + root)
	}

	// source is root and a file: keep basename
	if rel == "." && !isDirSource {
		rel = filepath.Base(sourceAbs)
	}

	return filepath.Abs(filepath.Join(output, rel))
}

// infer root by common parent of source, output
func InferRoot(source, output string) (string, error) {
	sourceAbs, err := filepath.Abs(source)
	if err != nil {
		return "", err
	}

	outputAbs, err := filepath.Abs(output)
	if err != nil {
		return "", err
	}

	common := CommonParent(sourceAbs, outputAbs)
	if common == "" || isVolumeRoot(common) {
		return "", ErrNoSharedParent
	}

	rel, err := filepath.Rel(common, sourceAbs)
	if err != nil {
		return "", err
	}

	// source is common parent (or one segment only) => keep
	if rel == "." || !strings.Contains(rel, string(filepath.Separator)) {
		return common, nil
	}

	first, _, _ := strings.Cut(rel, string(filepath.Separator))

	return filepath.Join(common, first), nil
}

// preserve structure by infering root
func ResolveOutputInferred(source, output string) (string, error) {
	// file output => no root needed
	if !IsDirPath(output) {
		if IsDirPath(source) {
			return "", ErrDirToFile
		}

		return filepath.Abs(output)
	}

	root, err := InferRoot(source, output)
	if err != nil {
		return "", err
	}

	return ResolveOutputPreserved(root, source, output)
}