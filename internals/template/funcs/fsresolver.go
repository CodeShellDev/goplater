package funcs

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/codeshelldev/goplater/internals/template/context"
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
	filePathAbs := resolvePath(*r.ctx, path)

	exists, _ := fsutils.Exists(filePathAbs)

	if !exists {
		filePathAbs = filePathAbs + ".gplt"
	}

	return readFile(filePathAbs)
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
		abs, _ := filepath.Abs(ctx.Invoker)

		filePathAbs = getAbsPathWithSource(path, filepath.Dir(abs))
	} else if isRelToSource {
		path, _ = strings.CutPrefix(path, "~/")
		path = "./" + path

		filePathAbs = getAbsPathWithSource(path, ctx.Options.Source)
	} else {
		filePathAbs, _ = filepath.Abs(path)
	}


	return filePathAbs
}

func readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	
	if err != nil {
		return "", errors.New("file not found: " + path)
	}

	return string(data), nil
}

func getAbsPathWithSource(path, source string) string {
	sourceAbs, _ := filepath.Abs(source)

	fullPath, _ := fsutils.Relative(sourceAbs, path)
	
	fullPath, _ = filepath.Abs(fullPath)

	return fullPath
}