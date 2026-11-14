package funcs

import (
	"errors"
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

var OutsideOfScopeErr = errors.New("outside of allowed scope")

func outputToFn(rt *gplt.Runtime, ctx *gplt.Context, path string) error {
	tmplContext := ctx.Get(context.TemplateContextKey).(*context.TemplateContext)

	absPath := ResolvePath(*tmplContext, path)

	if !isPathAllowed(absPath, tmplContext.Options.AllowedOutputs) {
		return fmt.Errorf("%q cannot be ouputted to %q: %w", tmplContext.Path, absPath, OutsideOfScopeErr)
	}

	tmplContext.OutputPath = absPath

	return nil
}