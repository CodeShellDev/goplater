package funcs

import (
	"fmt"
	"os"
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
		str, err := readHandlerWithOpts(context, path, ReadOptions{
			Recursive: true,
		})

		if err != nil {
			panic("could not read " + path + ": " + err.Error())
		}

		return str
	},
}

var readOptsFunc = TemplateFunc{
	Name: "readOpts",
	Handler: func(context context.TemplateContext, path string, args []string) string {
		var opts ReadOptions

		flags.ParseArgs(&opts, args)

		str, err := readHandlerWithOpts(context, path, opts)

		if err != nil {
			panic("could not read " + path + ": " + err.Error())
		}

		return str
	},
}


func readHandlerWithOpts(context context.TemplateContext, path string, opts ReadOptions) (string, error) {
	var err error

	res := readHandler(context, path)

	if opts.Recursive {
		res, err = core.Renderer.Render(res, context)
	}

	return res, err
}

func readHandler(context context.TemplateContext, path string) string {
	var res string

	isRel := strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../")

	filePathAbs, _ := filepath.Abs(context.Invoker)
	
	relPathComponent := fsutils.Relative(filepath.Dir(filePathAbs), path)

	context.Invoker = relPathComponent

	os.WriteFile("/tmp/test.txt", fmt.Append(nil, context,), 0644)

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