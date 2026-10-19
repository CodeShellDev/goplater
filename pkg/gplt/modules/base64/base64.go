package base64

import (
	"encoding/base64"

	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
)

var Module = modules.NewModule(base64EncodeFunc, base64DecodeFunc)

// Encodes strings with base64.
//
// @param encode string
// @returns string
var base64EncodeFunc = modules.NewFunc("base64Encode", base64Encode)

func base64Encode(_ *gplt.Runtime, _ *gplt.Context, str string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(str))

	return encoded
}

// Decodes a base64 string.
//
// @param decode string
// @returns string
// @returns error
var base64DecodeFunc = modules.NewFunc("base64Decode", base64Decode)

func base64Decode(_ *gplt.Runtime, _ *gplt.Context, encoded string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)

	return string(decoded), err
}