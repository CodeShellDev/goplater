package funcs

import (
	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/internals/template/core"
)

var importFunc = TemplateFunc{
	Name: "import",
	Handler: func(context context.TemplateContext, path string) any {
		tmpltBody, err := readHandlerWithOpts(context, path, ReadOptions{
			Recursive: false,
		})

		if err != nil {
			panic("error during import of " + path + ": " + err.Error())
		}

		if tmpltBody != "" {
			_, err := core.Renderer.Render(tmpltBody, context)

			if err != nil {
				panic("error during import of " + path + ": " + err.Error())
			}
		}

		return ""
	},
}

func init() {
	Register(importFunc)
}