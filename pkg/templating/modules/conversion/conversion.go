package conversion

import (
	"fmt"
	"strconv"

	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
)

var Module = modules.NewModule(toStringFunc, toIntFunc, toFloat64Func, toFloat32Func, toBoolFunc)

var toStringFunc = modules.NewFunc("toString", toString)

func toString(_ *templating.Runtime, _ templating.Context, value any) string  {
	return fmt.Sprint(value)
}

var toIntFunc = modules.NewFunc("toInt", toInt)

func toInt(_ *templating.Runtime, _ templating.Context, str string) (int, error)  {
	return strconv.Atoi(str)
}

var toFloat64Func = modules.NewFunc("toFloat64", toFloat64)

func toFloat64(_ *templating.Runtime, _ templating.Context, str string) (float64, error)  {
	return strconv.ParseFloat(str, 64)
}

var toFloat32Func = modules.NewFunc("toFloat32", toFloat32)

func toFloat32(_ *templating.Runtime, _ templating.Context, str string) (float32, error)  {
	float, err := strconv.ParseFloat(str, 32)

	return float32(float), err
}

var toBoolFunc = modules.NewFunc("toBool", toBool)

func toBool(_ *templating.Runtime, _ templating.Context, str string) (bool, error)  {
	return strconv.ParseBool(str)
}