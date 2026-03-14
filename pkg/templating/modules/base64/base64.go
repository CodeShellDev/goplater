package base64

import (
	"encoding/base64"

	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
)

var Module = modules.NewModule(base64EncodeFunc, base64DecodeFunc)

var base64EncodeFunc = modules.NewFunc("base64Encode", base64Encode)

func base64Encode(_ *templating.Runtime, _ templating.Context, decoded string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(decoded))

	return encoded
}

var base64DecodeFunc = modules.NewFunc("base64Decode", base64Decode)

func base64Decode(_ *templating.Runtime, _ templating.Context, encoded string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)

	return string(decoded), err
}