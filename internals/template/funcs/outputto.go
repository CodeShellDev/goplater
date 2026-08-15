package funcs

import (
	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
)

var outputToFunc = modules.NewFunc("outputTo", outputToFn)

func outputToFn(rt *templating.Runtime, ctx *templating.Context, path string) string {
	tmplContext := ctx.Get(context.TemplateContextKey).(*context.TemplateContext)
	fullPath := resolvePath(*tmplContext, path)

	if !isPathAllowed(fullPath, tmplContext) {
		panic(tmplContext.Path + " may not be saved to " + fullPath + " as it is not inside of the allowed scope")
	}

	tmplContext.OutputPath = fullPath

	return ""
}