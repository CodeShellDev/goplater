package encoding

import (
	"encoding/hex"

	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
)

// Encodes a string as hex.
//
// @param str string
// @returns string
var hexEncodeFunc = modules.NewFunc("hexEncode", hexEncode)

func hexEncode(_ *gplt.Runtime, _ *gplt.Context, bytes []byte) string {
	return hex.EncodeToString(bytes)
}

// Decodes a string as hex.
//
// @param str string
// @returns string
// @returns error
var hexDecodeFunc = modules.NewFunc("hexDecode", hexDecode)

func hexDecode(_ *gplt.Runtime, _ *gplt.Context, str string) ([]byte, error) {
	return hex.DecodeString(str)
}