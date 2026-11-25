package funcs

import (
	"errors"
	"fmt"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
	"github.com/codeshelldev/goplater/utils/fsutils"
)

// Changes the output path for the current file.
// Errors if path is not allowed.
//
// @param path string
// @returns error
var outputToFunc = modules.NewFunc("outputTo", outputToFn)

var OutsideOfScopeErr = errors.New("outside of allowed scope")

func outputToFn(rt *gplt.Runtime, ctx *gplt.Context, path string) (string, error) {
	tmplContext, err := context.Extract(ctx)

	if err != nil {
		return "", err
	}

	absPath := ResolvePath(tmplContext, path)

	// only lexically, does not prevent symlink bypass (will get caught at final write)
	allowed := false

	for _, dir := range tmplContext.Options.AllowedOutputs {
		if fsutils.IsInside(absPath, dir) {
			allowed = true
			break
		}
	}

	if !allowed {
		return "", fmt.Errorf("%q cannot be outputted to %q: %w", tmplContext.Shared.Path, absPath, OutsideOfScopeErr)
	}

	tmplContext.Shared.OutputPath = absPath
	return "", nil
}