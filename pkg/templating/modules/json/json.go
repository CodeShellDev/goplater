package json

import (
	jsonUtils "encoding/json"

	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
)

var Module = modules.NewModule(jsonEncodeFunc, jsonDecodeFunc)

// Encodes an object into a json string.
//
// @param obj any
// @returns string
// @returns error
var jsonEncodeFunc = modules.NewFunc("jsonEncode", jsonEncode)

func jsonEncode(_ *templating.Runtime, _ *templating.Context, obj any) (string, error)  {
	json, err := jsonUtils.Marshal(obj)

	return string(json), err
}

// Decodes a json string into an object.
//
// @param json string
// @returns any
// @returns error
var jsonDecodeFunc = modules.NewFunc("jsonDecode", jsonDecode)

func jsonDecode(_ *templating.Runtime, _ *templating.Context, json string) (any, error)  {
	var obj any

	err := jsonUtils.Unmarshal([]byte(json), &obj)

	return obj, err
}