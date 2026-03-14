package strings

import (
	"fmt"
	"strings"

	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
)

var Module = modules.NewModule(
	trimFunc,

	upperFunc, 
	lowerFunc, 

	containsFunc, 
	countFunc, 

	startsWithFunc, 
	endsWithFunc, 

	isEmptyFunc, 

	replaceFunc, 

	splitFunc, 

	afterFunc, 
	beforeFunc, 
	betweenFunc, 

	sliceFunc, 

	joinFunc, 

	repeatFunc, 

	concatFunc, 

	appendFunc, 

	indexOfFunc,
)

var trimFunc = modules.NewFunc("trim", trim)

func trim(_ *templating.Runtime, _ templating.Context, str string) string  {
	return strings.TrimSpace(str)
}

var upperFunc = modules.NewFunc("upper", upper)

func upper(_ *templating.Runtime, _ templating.Context, str string) string  {
	return strings.ToUpper(str)
}

var lowerFunc = modules.NewFunc("lower", lower)

func lower(_ *templating.Runtime, _ templating.Context, str string) string  {
	return strings.ToLower(str)
}

var containsFunc = modules.NewFunc("contains", contains)

func contains(_ *templating.Runtime, _ templating.Context, str string, sub string) bool  {
	return strings.Contains(str, sub)
}

var countFunc = modules.NewFunc("count", count)

func count(_ *templating.Runtime, _ templating.Context, str string, sub string) int  {
	return strings.Count(str, sub)
}

var startsWithFunc = modules.NewFunc("startsWith", startsWith)

func startsWith(_ *templating.Runtime, _ templating.Context, str string, prefix string) bool  {
	return strings.HasPrefix(str, prefix)
}

var endsWithFunc = modules.NewFunc("endsWith", endsWith)

func endsWith(_ *templating.Runtime, _ templating.Context, str string, suffix string) bool  {
	return strings.HasSuffix(str, suffix)
}

var isEmptyFunc = modules.NewFunc("isEmpty", isEmpty)

func isEmpty(_ *templating.Runtime, _ templating.Context, str string) bool  {
	return str == ""
}

var replaceFunc = modules.NewFunc("replace", replace)

func replace(_ *templating.Runtime, _ templating.Context, str string, sub string, replaceWith string) string  {
	return strings.ReplaceAll(str, sub, replaceWith)
}

var splitFunc = modules.NewFunc("split", split)

func split(_ *templating.Runtime, _ templating.Context, str string, sep string) []string  {
	return strings.Split(str, sep)
}

var afterFunc = modules.NewFunc("after", after)

func after(_ *templating.Runtime, _ templating.Context, str string, sub string) string  {
	after, _ := strings.CutPrefix(str, sub)

	return after
}

var beforeFunc = modules.NewFunc("before", before)

func before(_ *templating.Runtime, _ templating.Context, str string, sub string) string  {
	before, _ := strings.CutSuffix(str, sub)

	return before
}

var betweenFunc = modules.NewFunc("between", between)

func between(_ *templating.Runtime, _ templating.Context, str string, startSub, endSub string) string  {
	after, okAfter := strings.CutPrefix(str, startSub)

	if !okAfter {
		return ""
	}

	before, okBefore := strings.CutSuffix(after, endSub)

	if !okBefore {
		return ""
	}

	return before
}

var sliceFunc = modules.NewFunc("slice", slice)

func slice(_ *templating.Runtime, _ templating.Context, str string, start int, end int) string  {
	return str[start:end]
}

var joinFunc = modules.NewFunc("join", join)

func join(_ *templating.Runtime, _ templating.Context, str []any, sep string) string  {
	return strings.Join(toStringSlice(str), sep)
}

var repeatFunc = modules.NewFunc("repeat", repeat)

func repeat(_ *templating.Runtime, _ templating.Context, str string, count int) string  {
	return strings.Repeat(str, count)
}

func toStringSlice(values []any) []string {
    out := make([]string, len(values))

    for i, v := range values {
        out[i] = fmt.Sprint(v)
    }

    return out
}

var concatFunc = modules.NewFunc("concat", concat)

func concat(_ *templating.Runtime, _ templating.Context, args ...any) string  {
	args = modules.UnpackArgs(args...)
	str := toStringSlice(args)

	return strings.Join(str, "")
}

var appendFunc = modules.NewFunc("append", append)

func append(_ *templating.Runtime, _ templating.Context, str string, append string) string  {
	return str + append
}

var indexOfFunc = modules.NewFunc("indexOf", indexOf)

func indexOf(_ *templating.Runtime, _ templating.Context, str string, sub string) int  {
	return strings.Index(str, sub)
}