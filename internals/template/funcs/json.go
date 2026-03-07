package funcs

import (
	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/gotl/pkg/jsonutils"
)

var jsonEncodeFunc = TemplateFunc{
	Name: "jsonEncode",
	Handler: func(context context.TemplateContext, obj any) (string, error) {
		return jsonutils.ToJsonSafe(obj)
	},
}

var jsonDecodeFunc = TemplateFunc{
	Name: "jsonDecode",
	Handler: func(context context.TemplateContext, str string) (map[string]any, error) {
		return jsonutils.GetJsonSafe[map[string]any](str)
	},
}

func init() {
	Register(jsonEncodeFunc)
	Register(jsonDecodeFunc)
}