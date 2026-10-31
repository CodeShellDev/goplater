package core

import (
	"errors"

	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
)

// Tries to import a template file with func definitions.
// Optionally supply an alias for easier access.
// The fs resolver is only available as an internal module and does not ship with the `pkg/gplt/modules` directory.
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

func importFn(rt *gplt.Runtime, ctx *gplt.Context, args ...any) error {
	if len(args) <= 0 {
		return errors.New("wrong number of args for import: want at least 1 got 0")
	}

	scope, _ := ctx.Get(gplt.CurrentTreeKey).(gplt.TreeScope)

	path := args[0].(string)

	var override string

	if len(args) == 2 {
		path = args[1].(string)

		override = args[0].(string)
	}

	err := rt.Import(scope.Name, path, override, ctx)
	if err != nil {
		return errors.New("error importing " + path + ": " + err.Error())
	}

	return nil
}