package hash

import (
	md5Utils "crypto/md5"
	sha1Utils "crypto/sha1"
	sha256Utils "crypto/sha256"
	sha512Utils "crypto/sha512"
	"encoding/hex"

	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
)

var Module = modules.NewModule(
	hashFunc,
	sha1Func,
	sha512Func,
	md5Func,
)

// Hashes a given string with SHA-256.
//
// @param str string
var hashFunc = modules.NewFunc("hash", hash)

func hash(_ *gplt.Runtime, _ *gplt.Context, str string) string {
	sum := sha256Utils.Sum256([]byte(str))
	return hex.EncodeToString(sum[:])
}

// Hashes a given string with SHA-1.
//
// @param str string
var sha1Func = modules.NewFunc("sha1", sha1)

func sha1(_ *gplt.Runtime, _ *gplt.Context, str string) string {
    sum := sha1Utils.Sum([]byte(str))
    return hex.EncodeToString(sum[:])
}

// Hashes a given string with SHA-512.
//
// @param str string
var sha512Func = modules.NewFunc("sha512", sha512)

func sha512(_ *gplt.Runtime, _ *gplt.Context, str string) string {
    sum := sha512Utils.Sum512([]byte(str))
    return hex.EncodeToString(sum[:])
}

// Hashes a given string with MD-5.
//
// @param str string
var md5Func = modules.NewFunc("md5", md5)

func md5(_ *gplt.Runtime, _ *gplt.Context, str string) string {
    sum := md5Utils.Sum([]byte(str))
    return hex.EncodeToString(sum[:])
}