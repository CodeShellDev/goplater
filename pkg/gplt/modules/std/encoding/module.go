package encoding

import "github.com/codeshelldev/goplater/pkg/gplt/modules"

var Module = modules.NewModule(
	base64EncodeFunc, 
	base64DecodeFunc,

	hexEncodeFunc,
	hexDecodeFunc,
)
