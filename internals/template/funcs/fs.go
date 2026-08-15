package funcs

import (
	"os"
	"path/filepath"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
	"github.com/codeshelldev/goplater/utils/fsutils"
)

func isPathAllowed(path string, tmplContext *context.TemplateContext) bool {
	allowed := false

	for _, try := range tmplContext.Options.AllowedOutputFolders {
		if fsutils.IsInside(path, try) {
			allowed = true
		}
	}

	return allowed
}

var readFunc = modules.NewFunc("read", read)

func read(rt *templating.Runtime, ctx *templating.Context, path string) string {
	str, ctx := readHandler(ctx, path)

	str, err := rt.GetEngine().Execute(":read:" + path, str, nil, rt.GetEngineOptions(), ctx)

	if err != nil {
		panic("could not read " + path + ": " + err.Error())
	}

	return str
}

var readRawFunc = modules.NewFunc("readRaw", readRaw)

func readRaw(_ *templating.Runtime, ctx *templating.Context, path string) string {
	tmplContext := ctx.Get(context.TemplateContextKey).(*context.TemplateContext)

	filePathAbs := resolvePath(*tmplContext, path)

	res, err := readFile(filePathAbs)

	if err != nil {
		res = err.Error()
	}

	return res
}

var readArgsFunc = modules.NewFunc("readArgs", readArgs)

func readArgs(rt *templating.Runtime, ctx *templating.Context, path string, args ...any) string {
	args = modules.UnpackArgs(args...)

	str, newContext := readHandler(ctx, path)

	data := map[string]any{
		"args": args,
	}

	str, err := rt.GetEngine().Execute(":read:" + path, str, data, rt.GetEngineOptions(), newContext)

	if err != nil {
		panic("could not read " + path + ": " + err.Error())
	}

	return str
}

func readHandler(ctx *templating.Context, path string) (string, *templating.Context) {
	tmplContext := ctx.Get(context.TemplateContextKey).(*context.TemplateContext)

	filePathAbs := resolvePath(*tmplContext, path)

	res, err := readFile(filePathAbs)

	if err != nil {
		res = err.Error()
	}

	tmplContext.Invoker = filePathAbs

	newContext := &templating.Context{}

	newContext.Set(context.TemplateContextKey, tmplContext)

	return res, newContext
}

var writeFunc = modules.NewFunc("write", write)

func write(_ *templating.Runtime, ctx *templating.Context, path string, content string) error {
	tmplContext := ctx.Get(context.TemplateContextKey).(*context.TemplateContext)

	filePathAbs := resolvePath(*tmplContext, path)

	if !isPathAllowed(filePathAbs, tmplContext) {
		panic("writing to " + filePathAbs + " is not allowed as it is not inside of the allowed scope")
	}

	err := os.WriteFile(filePathAbs, []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}

var mkdirFunc = modules.NewFunc("mkdir", mkdir)

func mkdir(_ *templating.Runtime, ctx *templating.Context, path string) error {
	tmplContext := ctx.Get(context.TemplateContextKey).(*context.TemplateContext)

	filePathAbs := resolvePath(*tmplContext, path)

	if !isPathAllowed(filePathAbs, tmplContext) {
		panic("creating folder " + filePathAbs + " is not allowed as it is not inside of the allowed scope")
	}

	err := os.MkdirAll(filePathAbs, 0755)
	if err != nil {
		return err
	}

	return nil
}

var appendFileFunc = modules.NewFunc("appendFile", appendFile)

func appendFile(_ *templating.Runtime, ctx *templating.Context, path string, content string) error {
	tmplContext := ctx.Get(context.TemplateContextKey).(*context.TemplateContext)

	filePathAbs := resolvePath(*tmplContext, path)

	if !isPathAllowed(filePathAbs, tmplContext) {
		panic("appending to " + filePathAbs + " is not allowed as it is not inside of the allowed scope")
	}

	f, err := os.OpenFile(filePathAbs, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

var fsExistsFunc = modules.NewFunc("fsExists", fsExists)

func fsExists(_ *templating.Runtime, ctx *templating.Context, path string) bool {
	tmplContext := ctx.Get(context.TemplateContextKey).(*context.TemplateContext)

	filePathAbs := resolvePath(*tmplContext, path)

	_, err := os.Stat(filePathAbs)
	if err != nil {
		return false
	}

	return true
}

var isFileFunc = modules.NewFunc("isFile", isFile)

func isFile(_ *templating.Runtime, ctx *templating.Context, path string) bool {
	tmplContext := ctx.Get(context.TemplateContextKey).(*context.TemplateContext)

	filePathAbs := resolvePath(*tmplContext, path)

	info, err := os.Stat(filePathAbs)
	if err != nil {
		return false
	}

	return info.Mode().IsRegular()
}

var isDirFunc = modules.NewFunc("isDir", isDir)

func isDir(_ *templating.Runtime, ctx *templating.Context, path string) bool {
	tmplContext := ctx.Get(context.TemplateContextKey).(*context.TemplateContext)

	filePathAbs := resolvePath(*tmplContext, path)

	info, err := os.Stat(filePathAbs)
	if err != nil {
		return false
	}

	return info.IsDir()
}

var listDirFunc = modules.NewFunc("listDir", listDir)

func listDir(_ *templating.Runtime, ctx *templating.Context, path string) ([]string, error) {
	tmplContext := ctx.Get(context.TemplateContextKey).(*context.TemplateContext)

	filePathAbs := resolvePath(*tmplContext, path)

	entries, err := os.ReadDir(filePathAbs)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(entries))

	for _, entry := range entries {
		result = append(result, entry.Name())
	}

	return result, nil
}

var fsRemoveFunc = modules.NewFunc("fsRemove", fsRemove)

func fsRemove(_ *templating.Runtime, ctx *templating.Context, path string) error {
	tmplContext := ctx.Get(context.TemplateContextKey).(*context.TemplateContext)

	filePathAbs := resolvePath(*tmplContext, path)

	if !isPathAllowed(filePathAbs, tmplContext) {
		panic("removing " + filePathAbs + " is not allowed as it is not inside of the allowed scope")
	}

	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	return nil
}

var joinPathFunc = modules.NewFunc("joinPath", joinPath)

func joinPath(_ *templating.Runtime, _ *templating.Context, paths ...string) string {
	return filepath.Join(paths...)
}

var basePathFunc = modules.NewFunc("basePath", basePath)

func basePath(_ *templating.Runtime, _ *templating.Context, path string) string {
	return filepath.Base(path)
}

var pathDirFunc = modules.NewFunc("pathDir", pathDir)

func pathDir(_ *templating.Runtime, _ *templating.Context, path string) string {
	return filepath.Dir(path)
}

var fileExtFunc = modules.NewFunc("fileExt", fileExt)

func fileExt(_ *templating.Runtime, _ *templating.Context, path string) string {
	return filepath.Ext(path)
}

var absPathFunc = modules.NewFunc("absPath", absPath)

func absPath(_ *templating.Runtime, _ *templating.Context, path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return abs, nil
}

var relPathFunc = modules.NewFunc("relPath", relPath)

func relPath(_ *templating.Runtime, _ *templating.Context, base string, target string) (string, error) {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return "", err
	}

	return rel, nil
}