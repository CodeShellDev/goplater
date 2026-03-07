package funcs

import (
	"fmt"
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

var countFunc = TemplateFunc{
	Name: "count",
	Handler: func(context context.TemplateContext, str string, sub string) int {
		return strings.Count(str, sub)
	},
}

var startsWithFunc = TemplateFunc{
	Name: "startsWith",
	Handler: func(context context.TemplateContext, str string, prefix string) bool {
		return strings.HasPrefix(str, prefix)
	},
}

var endsWithFunc = TemplateFunc{
	Name: "endsWith",
	Handler: func(context context.TemplateContext, str string, suffix string) bool {
		return strings.HasSuffix(str, suffix)
	},
}

var isEmptyFunc = TemplateFunc{
	Name: "isEmpty",
	Handler: func(context context.TemplateContext, str string) bool {
		return str == ""
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

var afterFunc = TemplateFunc{
	Name: "after",
	Handler: func(context context.TemplateContext, str string, sub string) string {
		after, _ := strings.CutPrefix(str, sub)

		return after
	},
}

var beforeFunc = TemplateFunc{
	Name: "before",
	Handler: func(context context.TemplateContext, str string, sub string) string {
		before, _ := strings.CutSuffix(str, sub)

		return before
	},
}

var betweenFunc = TemplateFunc{
	Name: "between",
	Handler: func(context context.TemplateContext, str string, startSub, endSub string) string {
		after, okAfter := strings.CutPrefix(str, startSub)

		if !okAfter {
			return ""
		}

		before, okBefore := strings.CutSuffix(after, endSub)

		if !okBefore {
			return ""
		}

		return before
	},
}

var sliceFunc = TemplateFunc{
	Name: "slice",
	Handler: func(context context.TemplateContext, str string, start int, end int) string {
		return str[start:end]
	},
}

var joinFunc = TemplateFunc{
	Name: "join",
	Handler: func(context context.TemplateContext, str []any, sep string) string {
		return strings.Join(toStringSlice(str), sep)
	},
}

var repeatFunc = TemplateFunc{
	Name: "repeat",
	Handler: func(context context.TemplateContext, str string, count int) string {
		return strings.Repeat(str, count)
	},
}

func toStringSlice(values []any) []string {
    out := make([]string, len(values))

    for i, v := range values {
        out[i] = fmt.Sprint(v)
    }

    return out
}

var concatFunc = TemplateFunc{
	Name: "concat",
	Handler: func(context context.TemplateContext, args ...any) string {
		args = unpackArgs(args...)
		str := toStringSlice(args)

		return strings.Join(str, "")
	},
}

var appendFunc = TemplateFunc{
	Name: "append",
	Handler: func(context context.TemplateContext, str string, append string) string {
		return str + append
	},
}

var indexOfFunc = TemplateFunc{
	Name: "indexOf",
	Handler: func(context context.TemplateContext, str string, sub string) int {
		return strings.Index(str, sub)
	},
}

func init() {
	Register(trimFunc)
	Register(upperFunc)
	Register(lowerFunc)
	Register(containsFunc)
	Register(startsWithFunc)
	Register(endsWithFunc)
	Register(isEmptyFunc)
	Register(replaceFunc)
	Register(splitFunc)
	Register(sliceFunc)
	Register(joinFunc)
	Register(concatFunc)
	Register(appendFunc)
	Register(indexOfFunc)
	Register(beforeFunc)
	Register(afterFunc)
	Register(betweenFunc)
	Register(countFunc)
	Register(repeatFunc)
}