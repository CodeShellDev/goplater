package funcs

import (
	"path/filepath"
	"strings"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/internals/template/core"
	"github.com/codeshelldev/goplater/utils/fetchutils"
	"github.com/codeshelldev/goplater/utils/fsutils"

	"github.com/jessevdk/go-flags"
)

type ReadOptions struct {
	Recursive	bool	`short:"r" long:"recursive"`
}

var readFunc = TemplateFunc{
	Name: "read",
	Handler: func(context context.TemplateContext, path string) string {
		return readHandlerWithOpts(context, path, ReadOptions{
			Recursive: true,
		})
	},
}

var readOptsFunc = TemplateFunc{
	Name: "readOpts",
	Handler: func(context context.TemplateContext, path string, args []string) string {
		var opts ReadOptions

		flags.ParseArgs(&opts, args)

		return readHandlerWithOpts(context, path, opts)
	},
}


func readHandlerWithOpts(context context.TemplateContext, path string, opts ReadOptions) string {
	res := readHandler(context, path)

	if opts.Recursive {
		res, _ = core.Renderer.Render(res, context)
	}

	return res
}

func readHandler(context context.TemplateContext, path string) string {
	var res string

	isRel := strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../")

	filePathAbs, _ := filepath.Abs(context.Invoker)
	
	relPathComponent := fsutils.Relative(filepath.Dir(filePathAbs), path)

	context.Invoker = relPathComponent

	if isRel {
		res = fetchutils.Local(path, filepath.Dir(filePathAbs))
	} else {
		res = fetchutils.Local(path, context.Options.Source)
	}

	res, _ = core.Renderer.Render(res, context)

	return res
}

func init() {
	Register(readFunc)
	Register(readOptsFunc)
}