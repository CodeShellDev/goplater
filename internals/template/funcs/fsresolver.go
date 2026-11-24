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
	tmplContext, ok := ctx.Get(context.TemplateContextKey).(context.TemplateContext)

	if !ok {
		return "", errors.New("missing template context")
	}

	filePathAbs := ResolvePath(tmplContext, path + "." + gplt.IMPORT_FILE_EXTENSION)

	var info os.FileInfo

	err := fsutils.WithAllowedRoot(filePathAbs, tmplContext.Options.AllowedInputs, func(root *os.Root, relative string) error {
		i, err := root.Stat(relative)
		info = i
		return err
	})

	if err != nil {
		return "", errors.New("file does not exist")
	}

	if info.IsDir() {
		return "", errors.New("import is not a valid file")
	}

	return SafeRead(filePathAbs, &tmplContext)
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
	tmplContext, _ := ctx.Get(context.TemplateContextKey).(context.TemplateContext)

	return ResolvePath(tmplContext, path + "." + gplt.IMPORT_FILE_EXTENSION), nil
}

func ResolvePath(ctx context.TemplateContext, path string) string {
	isRel := strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../")
	isRelToSource := strings.HasPrefix(path, "~/")
	
	var filePathAbs string

	if isRel {
		absPath, _ := filepath.Abs(ctx.Local.Invoker)

		filePathAbs = GetAbsPathWithSource(path, filepath.Dir(absPath))
	} else if isRelToSource {
		path, _ = strings.CutPrefix(path, "~/")
		path = "./" + path

		filePathAbs = GetAbsPathWithSource(path, ctx.Options.Source)
	} else {
		filePathAbs, _ = filepath.Abs(path)
	}

	return filePathAbs
}

func GetAbsPathWithSource(path, source string) string {
	if !filepath.IsAbs(source) {
		var err error

		source, err = filepath.Abs(source)
		
		if err != nil {
			return ""
		}
	}

	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}

	return filepath.Join(source, path)
}