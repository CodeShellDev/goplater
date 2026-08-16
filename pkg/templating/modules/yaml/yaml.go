package yaml

import (
	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
	yamlUtils "go.yaml.in/yaml/v3"
)

var Module = modules.NewModule(yamlEncodeFunc, yamlDecodeFunc)

var yamlEncodeFunc = modules.NewFunc("yamlEncode", yamlEncode)

// Encodes an object into a yaml string.
//
// @param obj any
// @returns string
// @returns error
func yamlEncode(_ *templating.Runtime, _ *templating.Context, obj any) (string, error) {
	bytes, err := yamlUtils.Marshal(obj)

	return string(bytes), err
}

// Decodes a yaml string into an object.
//
// @param yaml string
// @returns any
// @returns error
var yamlDecodeFunc = modules.NewFunc("yamlDecode", yamlDecode)

func yamlDecode(_ *templating.Runtime, _ *templating.Context, yaml string) (any, error) {
	var obj any

	err := yamlUtils.Unmarshal([]byte(yaml), &obj)

	return obj, err
}