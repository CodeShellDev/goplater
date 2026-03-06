package funcs

import (
	"github.com/codeshelldev/goplater/internals/template/context"
	"go.yaml.in/yaml/v4"
)

var yamlFunc = TemplateFunc{
	Name: "yaml",
	Handler: func(context context.TemplateContext, str string) map[string]any {
		var res map[string]any

		err := yaml.Unmarshal([]byte(str), &res)

		if err != nil {
			panic("error parsing yaml: " + err.Error())
		}

		return res
	},
}

func init() {
	Register(yamlFunc)
}