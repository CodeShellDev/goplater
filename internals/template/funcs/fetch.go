package funcs

import (
	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/utils/fetchutils"
)

var fetchFunc = TemplateFunc{
	Name: "fetch",
	Handler: func(context context.TemplateContext, url string) string {
		return fetchutils.Remote(url)
	},
}

func init() {
	Register(fetchFunc)
}