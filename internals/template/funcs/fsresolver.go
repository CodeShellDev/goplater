package funcs

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/utils/fsutils"
)

type FsResolver struct {}

func NewFsResolver() *FsResolver {
	return &FsResolver{}
}

func (r *FsResolver) Trusted() bool {
	return true
}

func (r *FsResolver) CanResolve(path string) bool {
	return strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../") || strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "/")
}

func (r *FsResolver) Resolve(ctx *gplt.Context, path string) (string, error) {
	tmplContext, err := context.Extract(ctx)

	if err != nil {
		return "", err
	}

	filePathAbs, err := ResolvePath(tmplContext, path + "." + gplt.IMPORT_FILE_EXTENSION)

	if err != nil {
		return "", err
	}

	var info os.FileInfo

	err = fsutils.WithAllowedRoot(filePathAbs, tmplContext.Options.AllowedInputs, func(root *os.Root, relative string) error {
		i, err := root.Stat(relative)
		info = i
		return err
	})

	if err != nil {
		return "", err
	}

	if info.IsDir() {
		return "", errors.New("import is not a valid file")
	}

	return SafeRead(filePathAbs, tmplContext)
}

func SafeRead(path string, tmplContext *context.TemplateContext) (string, error) {
	var result string

	err := fsutils.WithAllowedRoot(path, tmplContext.Options.AllowedInputs,
		func(root *os.Root, relative string) error {
			data, err := root.ReadFile(relative)
			if err != nil {
				return err
			}

			result = string(data)
			return nil
		},
	)

	if err != nil {
		return "", err
	}

	return result, nil
}

func (r *FsResolver) DeriveName(path string) (string, error) {
	name := filepath.Base(path)
	name = strings.TrimSuffix(name, filepath.Ext(name))
	return name, nil
}

func (r *FsResolver) CleanPath(ctx *gplt.Context, path string) (string, error) {
	tmplContext, err := context.Extract(ctx)

	if err != nil {
		return "", err
	}

	return ResolvePath(tmplContext, path + "." + gplt.IMPORT_FILE_EXTENSION)
}

func ResolvePath(tc *context.TemplateContext, path string) (string, error) {
	switch {
	case strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../"):

		inv, err := filepath.Abs(tc.Local.Invoker)
		if err != nil {
			return "", err
		}
		return GetAbsPathWithSource(path, filepath.Dir(inv))

	case strings.HasPrefix(path, "~/"):
		return GetAbsPathWithSource("./" + strings.TrimPrefix(path, "~/"), tc.Options.Source)

	default:
		return filepath.Abs(path)
	}
}

func GetAbsPathWithSource(path, source string) (string, error) {
	source, err := filepath.Abs(source)
	if err != nil {
		return "", err
	}

	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}

	return filepath.Join(source, path), nil
}