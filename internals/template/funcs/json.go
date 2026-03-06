package funcs

import (
	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/gotl/pkg/jsonutils"
)

var jsonFunc = TemplateFunc{
	Name: "json",
	Handler: func(context context.TemplateContext, str string) map[string]any {
		res, err := jsonutils.GetJsonSafe[map[string]any](str)

		if err != nil {
			panic("error parsing json: " + err.Error())
		}

		return res
	},
}

func init() {
	Register(jsonFunc)
}