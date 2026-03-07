package funcs

import (
	"github.com/codeshelldev/goplater/internals/template/context"
	"go.yaml.in/yaml/v4"
)

var yamlEncodeFunc = TemplateFunc{
	Name: "yamlEncode",
	Handler: func(context context.TemplateContext, obj any) (string, error) {
		bytes, err := yaml.Marshal(obj)

		return string(bytes), err
	},
}

var yamlDecodeFunc = TemplateFunc{
	Name: "yamlDecode",
	Handler: func(context context.TemplateContext, str string) (map[string]any, error) {
		var res map[string]any

		err := yaml.Unmarshal([]byte(str), &res)

		return res, err
	},
}

func init() {
	Register(yamlEncodeFunc)
	Register(yamlDecodeFunc)
}