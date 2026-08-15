package funcs

import (
	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
	"github.com/codeshelldev/goplater/utils/fsutils"
)

var outputToFunc = modules.NewFunc("outputTo", outputToFn)

func outputToFn(rt *templating.Runtime, ctx *templating.Context, path string) string {
	tmplContext := ctx.Get(context.TemplateContextKey).(*context.TemplateContext)
	fullPath := resolvePath(*tmplContext, path)

	allowed := false

	for _, try := range tmplContext.Options.AllowedOutputFolders {
		if fsutils.IsInside(fullPath, try) {
			allowed = true
		}
	}

	if !allowed {
		panic(tmplContext.Path + " may not be saved to " + fullPath + " as it is not inside of the allowed scope")
	}

	tmplContext.OutputPath = fullPath

	return ""
}