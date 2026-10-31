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

type FsResolver struct {
	ctx *context.TemplateContext
}

func NewFsResolver(ctx *context.TemplateContext) *FsResolver {
	return &FsResolver{ctx: ctx}
}

func (r *FsResolver) Trusted() bool {
	return true
}

func (r *FsResolver) CanResolve(path string) bool {
	return strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../") || strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "/")
}

func (r *FsResolver) Resolve(path string) (string, error) {
	filePathAbs := resolvePath(*r.ctx, path + "." + gplt.IMPORT_FILE_EXTENSION)

	info, err := os.Stat(filePathAbs)

	if err != nil {
		return "", errors.New("file does not exist")
	}

	if info.IsDir() {
		return "", errors.New("import is not a valid file")
	}

	return safeRead(filePathAbs, r.ctx)
}

func safeRead(path string, tmplContext *context.TemplateContext) (string, error) {
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

func resolvePath(ctx context.TemplateContext, path string) string {
	isRel := strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../")
	isRelToSource := strings.HasPrefix(path, "~/")
	
	var filePathAbs string

	if isRel {
		absPath, _ := filepath.Abs(ctx.Invoker)

		filePathAbs = getAbsPathWithSource(path, filepath.Dir(absPath))
	} else if isRelToSource {
		path, _ = strings.CutPrefix(path, "~/")
		path = "./" + path

		filePathAbs = getAbsPathWithSource(path, ctx.Options.Source)
	} else {
		filePathAbs, _ = filepath.Abs(path)
	}

	return filePathAbs
}

func getAbsPathWithSource(path, source string) string {
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