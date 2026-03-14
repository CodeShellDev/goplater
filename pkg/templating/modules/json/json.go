package json

import (
	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
	"github.com/codeshelldev/gotl/pkg/jsonutils"
)

var Module = modules.NewModule(jsonEncodeFunc, jsonDecodeFunc)

var jsonEncodeFunc = modules.NewFunc("jsonEncode", jsonEncode)

func jsonEncode(_ *templating.Runtime, _ templating.Context, obj any) (string, error)  {
	return jsonutils.ToJsonSafe(obj)
}

var jsonDecodeFunc = modules.NewFunc("jsonDecode", jsonDecode)

func jsonDecode(_ *templating.Runtime, _ templating.Context, str string) (map[string]any, error)  {
	return jsonutils.GetJsonSafe[map[string]any](str)
}