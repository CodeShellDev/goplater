package funcs

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
	"github.com/codeshelldev/goplater/utils/fsutils"
)

var readFunc = modules.NewFunc("read", read)

func read(rt *templating.Runtime, context templating.Context, path string) string {
	str, ctx := readHandler(context, path)

	str, err := rt.GetEngine().Execute(":read:" + path, str, nil, rt.GetEngineOptions(), ctx)

	if err != nil {
		panic("could not read " + path + ": " + err.Error())
	}

	return str
}

var readArgsFunc = modules.NewFunc("readArgs", readArgs)

func readArgs(rt *templating.Runtime, context templating.Context, path string, args ...any) string {
	args = modules.UnpackArgs(args...)

	str, ctx := readHandler(context, path)

	data := map[string]any{
		"args": args,
	}

	str, err := rt.GetEngine().Execute(":read:" + path, str, data, rt.GetEngineOptions(), ctx)

	if err != nil {
		panic("could not read " + path + ": " + err.Error())
	}

	return str
}

func readHandler(ctx templating.Context, path string) (string, templating.Context) {
	var res string

	tmplContext := ctx.Get(context.TemplateContextKey).(context.TemplateContext)

	isRel := strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../")
	isRelToSource := strings.HasPrefix(path, "~/")
	
	var filePathAbs string

	if isRel {
		abs, _ := filepath.Abs(tmplContext.Invoker)

		filePathAbs = getAbsPathWithSource(path, filepath.Dir(abs))
	} else if isRelToSource {
		path, _ = strings.CutPrefix(path, "~/")
		path = "./" + path

		filePathAbs = getAbsPathWithSource(path, tmplContext.Options.Source)
	} else {
		filePathAbs, _ = filepath.Abs(path)
	}

	res = readFile(filePathAbs)

	tmplContext.Invoker = filePathAbs

	newContext := templating.Context{}

	newContext.Set(context.TemplateContextKey, tmplContext)

	return res, newContext
}

func readFile(path string) string {
	data, err := os.ReadFile(path)
	
	if err != nil {
		return "file not found: " + path
	}

	return string(data)
}

func getAbsPathWithSource(path, source string) string {
	sourceAbs, _ := filepath.Abs(source)

	fullPath := fsutils.Relative(sourceAbs, path)
	
	fullPath, _ = filepath.Abs(fullPath)

	return fullPath
}