package yaml

import (
	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
	"go.yaml.in/yaml/v4"
)

var Module = modules.NewModule(yamlEncodeFunc, yamlDecodeFunc)

var yamlEncodeFunc = modules.NewFunc("yamlEncode", yamlEncode)

func yamlEncode(_ *templating.Runtime, _ templating.Context, obj any) (string, error) {
	bytes, err := yaml.Marshal(obj)

	return string(bytes), err
}

var yamlDecodeFunc = modules.NewFunc("yamlDecode", yamlDecode)

func yamlDecode(_ *templating.Runtime, _ templating.Context, str string) (map[string]any, error) {
	var res map[string]any

	err := yaml.Unmarshal([]byte(str), &res)

	return res, err
}