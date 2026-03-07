package funcs

import (
	"strings"

	"github.com/codeshelldev/goplater/internals/template/context"
)

var trimFunc = TemplateFunc{
	Name: "trim",
	Handler: func(context context.TemplateContext, str string) string {
		return strings.TrimSpace(str)
	},
}

var upperFunc = TemplateFunc{
	Name: "upper",
	Handler: func(context context.TemplateContext, str string) string {
		return strings.ToUpper(str)
	},
}

var lowerFunc = TemplateFunc{
	Name: "lower",
	Handler: func(context context.TemplateContext, str string) string {
		return strings.ToLower(str)
	},
}

var containsFunc = TemplateFunc{
	Name: "contains",
	Handler: func(context context.TemplateContext, str string, sub string) bool {
		return strings.Contains(str, sub)
	},
}

var replaceFunc = TemplateFunc{
	Name: "replace",
	Handler: func(context context.TemplateContext, str string, sub string, replaceWith string) string {
		return strings.ReplaceAll(str, sub, replaceWith)
	},
}

var splitFunc = TemplateFunc{
	Name: "split",
	Handler: func(context context.TemplateContext, str string, sep string) []string {
		return strings.Split(str, sep)
	},
}

var joinFunc = TemplateFunc{
	Name: "join",
	Handler: func(context context.TemplateContext, str []string, sep string) string {
		return strings.Join(str, sep)
	},
}

var appendFunc = TemplateFunc{
	Name: "append",
	Handler: func(context context.TemplateContext, str string, append string) string {
		return str + append
	},
}

func init() {
	Register(trimFunc)
	Register(upperFunc)
	Register(lowerFunc)
	Register(containsFunc)
	Register(replaceFunc)
	Register(splitFunc)
	Register(joinFunc)
	Register(appendFunc)
}