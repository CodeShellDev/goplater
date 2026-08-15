package funcs

import (
	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
)

var readFunc = modules.NewFunc("read", read)

func read(rt *templating.Runtime, ctx *templating.Context, path string) string {
	str, ctx := readHandler(ctx, path)

	str, err := rt.GetEngine().Execute(":read:" + path, str, nil, rt.GetEngineOptions(), ctx)

	if err != nil {
		panic("could not read " + path + ": " + err.Error())
	}

	return str
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