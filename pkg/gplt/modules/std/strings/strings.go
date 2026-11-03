package strings

import (
	"fmt"
	"strings"

	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
)

var Module = modules.NewModule(
	trimFunc,
	trimSpaceFunc,

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
	afterLastFunc,
	beforeFunc,
	beforeLastFunc,
	betweenFunc,

	cutPrefixFunc,
	cutSuffixFunc,

	sliceFunc, 

	joinFunc, 

	repeatFunc, 

	concatFunc, 

	indexOfFunc,
)

// Trims space from a string (leading and trailing).
//
// @param str string
// @returns string
var trimSpaceFunc = modules.NewFunc("trimSpace", trimSpace)

func trimSpace(_ *gplt.Runtime, _ *gplt.Context, str string) string  {
	return strings.TrimSpace(str)
}

// Trims all chars defined in the cutset param from a string (leading and trailing).
//
// @param str string
// @param cutset string
// @returns string
var trimFunc = modules.NewFunc("trim", trim)

func trim(_ *gplt.Runtime, _ *gplt.Context, str string, cutset string) string  {
	return strings.Trim(str, cutset)
}

// Uppercases a string.
//
// @param str string
// @returns string
var upperFunc = modules.NewFunc("upper", upper)

func upper(_ *gplt.Runtime, _ *gplt.Context, str string) string  {
	return strings.ToUpper(str)
}

// Lowercases a string.
//
// @param str string
// @returns string
var lowerFunc = modules.NewFunc("lower", lower)

func lower(_ *gplt.Runtime, _ *gplt.Context, str string) string  {
	return strings.ToLower(str)
}

// Returns whether a string contains the given substring.
//
// @param str string
// @param sub string
// @returns bool
var containsFunc = modules.NewFunc("contains", contains)

func contains(_ *gplt.Runtime, _ *gplt.Context, str string, sub string) bool  {
	return strings.Contains(str, sub)
}

// Counts all occurences of a substring in the given string.
//
// @param str string
// @param sub string
// @returns int
var countFunc = modules.NewFunc("count", count)

func count(_ *gplt.Runtime, _ *gplt.Context, str string, sub string) int  {
	return strings.Count(str, sub)
}

// Returns whether a string starts with the given prefix.
//
// @param str string
// @param prefix string
// @returns bool
var startsWithFunc = modules.NewFunc("startsWith", startsWith)

func startsWith(_ *gplt.Runtime, _ *gplt.Context, str string, prefix string) bool  {
	return strings.HasPrefix(str, prefix)
}

// Returns whether a string ends with the given suffix.
//
// @param str string
// @param suffix string
// @returns bool
var endsWithFunc = modules.NewFunc("endsWith", endsWith)

func endsWith(_ *gplt.Runtime, _ *gplt.Context, str string, suffix string) bool  {
	return strings.HasSuffix(str, suffix)
}

// Returns whether the given string is empty.
//
// @param str string
// @returns bool
var isEmptyFunc = modules.NewFunc("isEmpty", isEmpty)

func isEmpty(_ *gplt.Runtime, _ *gplt.Context, str string) bool  {
	return str == ""
}

// Replaces all occurences of a substring in the given string with the replaceWith param.
//
// @param str string
// @param sub string
// @param replaceWith string
// @returns string
var replaceFunc = modules.NewFunc("replace", replace)

func replace(_ *gplt.Runtime, _ *gplt.Context, str string, sub string, replaceWith string) string  {
	return strings.ReplaceAll(str, sub, replaceWith)
}
// Split a string by the given seperator.
//
// @param str string
// @param sep string
// @returns []string
var splitFunc = modules.NewFunc("split", split)

func split(_ *gplt.Runtime, _ *gplt.Context, str string, sep string) []string  {
	return strings.Split(str, sep)
}

// Returns the string after the given substring in a string.
//
// @param str string
// @param sub string
// @returns string
var afterFunc = modules.NewFunc("after", after)

func after(_ *gplt.Runtime, _ *gplt.Context, str string, sub string) string  {
	_, a, found := strings.Cut(str, sub)
	if !found {
		return str
	}
	return a
}

// Returns the string after the last occurrence of the given substring.
//
// @param str string
// @param sub string
// @returns string
var afterLastFunc = modules.NewFunc("afterLast", afterLast)

func afterLast(_ *gplt.Runtime, _ *gplt.Context, str string, sub string) string {
	i := strings.LastIndex(str, sub)
	if i == -1 {
		return str
	}
	return str[i + len(sub):]
}

// Returns the string before the given substring in a string.
//
// @param str string
// @param sub string
// @returns string
var beforeFunc = modules.NewFunc("before", before)

func before(_ *gplt.Runtime, _ *gplt.Context, str string, sub string) string  {
	b, _, found := strings.Cut(str, sub)
	if !found {
		return str
	}
	return b
}

// Returns the string before the last occurrence of the given substring.
//
// @param str string
// @param sub string
// @returns string
var beforeLastFunc = modules.NewFunc("beforeLast", beforeLast)

func beforeLast(_ *gplt.Runtime, _ *gplt.Context, str string, sub string) string {
	i := strings.LastIndex(str, sub)
	if i == -1 {
		return str
	}
	return str[:i]
}

// Returns the substring between the start and end substring.
//
// @param str string
// @param startSub string
// @param endSub string
// @returns string
var betweenFunc = modules.NewFunc("between", between)

func between(_ *gplt.Runtime, _ *gplt.Context, str string, startSub, endSub string) string  {
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

// Removes a substring at the start of the given string.
//
// @param str string
// @param sub string
// @returns string
var cutPrefixFunc = modules.NewFunc("cutPrefix", cutPrefix)

func cutPrefix(_ *gplt.Runtime, _ *gplt.Context, str string, sub string) string  {
	after, _ := strings.CutPrefix(str, sub)

	return after
}

// Removes a substring at the end of the given string.
//
// @param str string
// @param sub string
// @returns string
var cutSuffixFunc = modules.NewFunc("cutSuffix", cutSuffix)

func cutSuffix(_ *gplt.Runtime, _ *gplt.Context, str string, sub string) string  {
	before, _ := strings.CutSuffix(str, sub)

	return before
}


// Slices a string based on start and end index.
//
// @param str string
// @param start int
// @param end int
// @returns string
var sliceFunc = modules.NewFunc("slice", slice)

func slice(_ *gplt.Runtime, _ *gplt.Context, str string, start int, end int) string  {
	return str[start:end]
}

// Joins multiple strings together by the given separator.
//
// @param sep string
// @param strings []string
// @returns string
var joinFunc = modules.NewFunc("join", join)

func join(_ *gplt.Runtime, _ *gplt.Context, sep string, args ...any) string  {
	if len(args) <= 1 {
		args = modules.UnpackArgs(args...)
	}

	return strings.Join(toStringSlice(args), sep)
}

// Repeat a string n times.
//
// @param str string
// @param count int
// @returns string
var repeatFunc = modules.NewFunc("repeat", repeat)

func repeat(_ *gplt.Runtime, _ *gplt.Context, str string, count int) string  {
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

func concat(_ *gplt.Runtime, _ *gplt.Context, strs ...string) string  {
	return strings.Join(strs, "")
}

var indexOfFunc = modules.NewFunc("indexOf", indexOf)

func indexOf(_ *gplt.Runtime, _ *gplt.Context, str string, sub string) int  {
	return strings.Index(str, sub)
}
