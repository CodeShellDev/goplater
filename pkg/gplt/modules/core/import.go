package core

import (
	"fmt"

	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
)

func parseImportArgs(args []any, requireHash bool) (string, string, string, error) {
	minArgs := 1
	if requireHash {
		minArgs = 2
	}

	err := modules.RequireRangeArgs(minArgs, minArgs + 1, args)
	if err != nil {
		return "", "", "", err
	}

	getString := func(index int) (string, error) {
		value, ok := args[index].(string)
		if !ok {
			return "", fmt.Errorf("argument %d must be a string", index+1)
		}
		return value, nil
	}

	var alias string
	var path string
	var hash string
	if requireHash {
		if len(args) == 3 {
			alias, err = getString(0)
			if err != nil {
				return "", "", "", err
			}
		}

		pathIndex := len(args) - 2
		path, err = getString(pathIndex)
		if err != nil {
			return "", "", "", err
		}

		hash, err = getString(pathIndex + 1)
	} else {
		if len(args) == 2 {
			alias, err = getString(0)
			if err != nil {
				return "", "", "", err
			}
		}

		path, err = getString(len(args) - 1)
	}

	return alias, path, hash, nil
}

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

func importFn(rt *gplt.Runtime, ctx *gplt.Context, args ...any) (string, error) {
	alias, path, _, err := parseImportArgs(args, false)
	if err != nil {
		return "", err
	}

	scope, _ := ctx.Get(gplt.CurrentTreeKey).(gplt.TreeScope)

	err = rt.Import(scope.Name, path, alias, "", ctx)
	if err != nil {
		return "", fmt.Errorf("error importing %q: %w", path, err)
	}

	return "", nil
}

// Same `import`, but only imports if content matches expected SHA-256 hash.
//
// @param alias string?
// @param path string
// @param hash string?
// @returns error
//
// @example
//	 +{{ importExact "funcs" "https://domain.com/functions" "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" }}
// @output
//	 +{{ call "funcs.greet" "John" }}
var importExactFunc = modules.NewFunc("importExact", importExactFn)

func importExactFn(rt *gplt.Runtime, ctx *gplt.Context, args ...any) (string, error) {
	alias, path, hash, err := parseImportArgs(args, true)
	if err != nil {
		return "", err
	}

	scope, _ := ctx.Get(gplt.CurrentTreeKey).(gplt.TreeScope)

	err = rt.Import(scope.Name, path, alias, hash, ctx)
	if err != nil {
		return "", fmt.Errorf("error importing %q: %w", path, err)
	}

	return "", nil
}