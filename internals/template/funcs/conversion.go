package funcs

import (
	"fmt"
	"strconv"

	"github.com/codeshelldev/goplater/internals/template/context"
)

var toStringFunc = TemplateFunc{
	Name: "toString",
	Handler: func(context context.TemplateContext, value any) string {
		return fmt.Sprint(value)
	},
}

var toIntFunc = TemplateFunc{
	Name: "toInt",
	Handler: func(context context.TemplateContext, str string) (int, error) {
		return strconv.Atoi(str)
	},
}

var toFloat64Func = TemplateFunc{
	Name: "toFloat64",
	Handler: func(context context.TemplateContext, str string) (float64, error) {
		return strconv.ParseFloat(str, 64)
	},
}

var toFloat32Func = TemplateFunc{
	Name: "toFloat32",
	Handler: func(context context.TemplateContext, str string) (float32, error) {
		float, err := strconv.ParseFloat(str, 32)

		return float32(float), err
	},
}

var toBoolFunc = TemplateFunc{
	Name: "toBool",
	Handler: func(context context.TemplateContext, str string) (bool, error) {
		return strconv.ParseBool(str)
	},
}

func init() {
	Register(toStringFunc)
	Register(toIntFunc)
	Register(toFloat64Func)
	Register(toFloat32Func)
	Register(toBoolFunc)
}