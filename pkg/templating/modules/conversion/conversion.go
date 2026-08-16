package conversion

import (
	"fmt"
	"strconv"

	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
)

var Module = modules.NewModule(toStringFunc, toIntFunc, toFloat64Func, toFloat32Func, toBoolFunc)

// Returns a given value as string (using [fmt.Sprint](https://pkg.go.dev/fmt#Sprint)).
//
// @param value string
// @returns string
var toStringFunc = modules.NewFunc("toString", toString)

func toString(_ *templating.Runtime, _ *templating.Context, value any) string  {
	return fmt.Sprint(value)
}

// Attempts to parse a string as int.
//
// @param str string
// @returns int
// @returns error
var toIntFunc = modules.NewFunc("toInt", toInt)

func toInt(_ *templating.Runtime, _ *templating.Context, str string) (int, error)  {
	return strconv.Atoi(str)
}

// Attempts to parse a string as float64.
//
// @param str string
// @returns float64
// @returns error
var toFloat64Func = modules.NewFunc("toFloat64", toFloat64)

func toFloat64(_ *templating.Runtime, _ *templating.Context, str string) (float64, error)  {
	return strconv.ParseFloat(str, 64)
}

// Attempts to parse a string as float32.
//
// @param str string
// @returns float32
// @returns error
var toFloat32Func = modules.NewFunc("toFloat32", toFloat32)

func toFloat32(_ *templating.Runtime, _ *templating.Context, str string) (float32, error)  {
	float, err := strconv.ParseFloat(str, 32)

	return float32(float), err
}

// Attempts to parse a string as a bool (using [strconv.ParseBool](https://pkg.go.dev/strconv#ParseBool)).
//
// @param str string
// @returns bool
// @returns error
var toBoolFunc = modules.NewFunc("toBool", toBool)

func toBool(_ *templating.Runtime, _ *templating.Context, str string) (bool, error)  {
	return strconv.ParseBool(str)
}