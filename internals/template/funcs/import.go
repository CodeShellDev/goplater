package funcs

import (
	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
)

var importFunc = modules.NewFunc("import", importFn)

func importFn(rt *templating.Runtime, context templating.Context, path string) any {
	str, ctx := readHandler(context, path)

	_, err := rt.GetEngine().ExecuteWithRuntime(":import:" + path, str, nil, rt.GetEngineOptions().Delims, ctx, rt)

	if err != nil {
		panic("error during import of " + path + ": " + err.Error())
	}

	return ""
}