package funcs

import (
	"fmt"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
)

// Changes the output path for the current file.
// Errors if path is not allowed.
//
// @param path string
// @returns error
var outputToFunc = modules.NewFunc("outputTo", outputToFn)

func outputToFn(rt *gplt.Runtime, ctx *gplt.Context, path string) error {
	tmplContext := ctx.Get(context.TemplateContextKey).(*context.TemplateContext)

	absPath := resolvePath(*tmplContext, path)

	if !isPathAllowed(absPath, tmplContext.Options.AllowedOutputs) {
		return fmt.Errorf("%q may not be ouputted to %q as it is not inside of the allowed scope", tmplContext.Path, absPath)
	}

	tmplContext.OutputPath = absPath

	return nil
}