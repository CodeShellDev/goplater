package core

import (
	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
)

// Tries to import a template file with func definitions.
// Optionally supply an alias for easier access.
// The fs resolver is only available as an internal module and does not ship with the `pkg/templating/modules` directory.
//
// @param alias string?
// @param path string
// @returns error
//
// @example
//	 +{{ import "https://domain.com/functions" }}
// @output
//	 +{{ call "functions.greet" "John" }}
//
// @example
//	 +{{ import "funcs" "https://domain.com/functions" }}
// @output
//	 +{{ call "funcs.greet" "John" }}
var importFunc = modules.NewFunc("import", importFn)

func importFn(rt *templating.Runtime, ctx *templating.Context, args ...any) string {
	if len(args) <= 0 {
		panic("wrong number of args for import: want at least 1 got 0")
	}

	scope, _ := ctx.Get(templating.CurrentTreeKey).(templating.TreeScope)

	path := args[0].(string)

	var override string

	if len(args) == 2 {
		path = args[1].(string)

		override = args[0].(string)
	}

	err := rt.Import(scope.Name, path, override, ctx)
	if err != nil {
		panic(err.Error())
	}

	return ""
}