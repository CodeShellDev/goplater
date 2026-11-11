package core

import (
	"fmt"

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
var importFunc = modules.NewFunc("import", importPath, importAlias)

func importPath(rt *gplt.Runtime, ctx *gplt.Context, path string) (string, error) {
	return doImport(rt, ctx, "", path, "")
}

func importAlias(rt *gplt.Runtime, ctx *gplt.Context, alias, path string) (string, error) {
	return doImport(rt, ctx, alias, path, "")
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
var importExactFunc = modules.NewFunc("importExact", importExactPath, importExactAlias)

func importExactPath(rt *gplt.Runtime, ctx *gplt.Context, path, hash string) (string, error) {
	return doImport(rt, ctx, "", path, hash)
}

func importExactAlias(rt *gplt.Runtime, ctx *gplt.Context, alias, path, hash string) (string, error) {
	return doImport(rt, ctx, alias, path, hash)
}

func doImport(rt *gplt.Runtime, ctx *gplt.Context, alias, path, hash string) (string, error) {
	scope := ctx.Get(gplt.CurrentTreeKey).(gplt.TreeScope)

	err := rt.Import(scope.Name, path, alias, hash, ctx)
	if err != nil {
		return "", fmt.Errorf("error importing %q: %w", path, err)
	}
	return "", nil
}