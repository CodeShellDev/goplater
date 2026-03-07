package funcs

import (
	"encoding/base64"

	"github.com/codeshelldev/goplater/internals/template/context"
)

var base64EncodeFunc = TemplateFunc{
	Name: "base64Encode",
	Handler: func(context context.TemplateContext, decoded string) string {
		encoded := base64.StdEncoding.EncodeToString([]byte(decoded))

		return encoded
	},
}

var base64DecodeFunc = TemplateFunc{
	Name: "base64Decode",
	Handler: func(context context.TemplateContext, encoded string) (string, error) {
		decoded, err := base64.StdEncoding.DecodeString(encoded)

		return string(decoded), err
	},
}

func init() {
	Register(base64EncodeFunc)
	Register(base64DecodeFunc)
}