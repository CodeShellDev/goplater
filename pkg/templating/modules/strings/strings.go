package strings

import (
	"fmt"
	"strings"

	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
)

var Module = modules.NewModule(
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
	beforeFunc, 
	betweenFunc, 

	cutPrefixFunc,
	cutSuffixFunc,

	sliceFunc, 

	joinFunc, 

	repeatFunc, 

	concatFunc, 

	indexOfFunc,
)

// Trims space from a string.
//
// @param str string
// @returns string
var trimSpaceFunc = modules.NewFunc("trimSpace", trimSpace)

func trimSpace(_ *templating.Runtime, _ *templating.Context, str string) string  {
	return strings.TrimSpace(str)
}

// Uppercases a string.
//
// @param str string
// @returns string
var upperFunc = modules.NewFunc("upper", upper)

func upper(_ *templating.Runtime, _ *templating.Context, str string) string  {
	return strings.ToUpper(str)
}

// Lowercases a string.
//
// @param str string
// @returns string
var lowerFunc = modules.NewFunc("lower", lower)

func lower(_ *templating.Runtime, _ *templating.Context, str string) string  {
	return strings.ToLower(str)
}

// Returns whether a string contains the given substring.
//
// @param str string
// @param sub string
// @returns bool
var containsFunc = modules.NewFunc("contains", contains)

func contains(_ *templating.Runtime, _ *templating.Context, str string, sub string) bool  {
	return strings.Contains(str, sub)
}

// Counts all occurences of a substring in the given string.
//
// @param str string
// @param sub string
// @returns int
var countFunc = modules.NewFunc("count", count)

func count(_ *templating.Runtime, _ *templating.Context, str string, sub string) int  {
	return strings.Count(str, sub)
}

// Returns whether a string starts with the given prefix.
//
// @param str string
// @param prefix string
// @returns bool
var startsWithFunc = modules.NewFunc("startsWith", startsWith)

func startsWith(_ *templating.Runtime, _ *templating.Context, str string, prefix string) bool  {
	return strings.HasPrefix(str, prefix)
}

// Returns whether a string ends with the given suffix.
//
// @param str string
// @param suffix string
// @returns bool
var endsWithFunc = modules.NewFunc("endsWith", endsWith)

func endsWith(_ *templating.Runtime, _ *templating.Context, str string, suffix string) bool  {
	return strings.HasSuffix(str, suffix)
}

// Returns whether the given string is empty.
//
// @param str string
// @returns bool
var isEmptyFunc = modules.NewFunc("isEmpty", isEmpty)

func isEmpty(_ *templating.Runtime, _ *templating.Context, str string) bool  {
	return str == ""
}

// Replaces all occurences of a substring in the given string with the replaceWith param.
//
// @param str string
// @param sub string
// @param replaceWith string
// @returns string
var replaceFunc = modules.NewFunc("replace", replace)

func replace(_ *templating.Runtime, _ *templating.Context, str string, sub string, replaceWith string) string  {
	return strings.ReplaceAll(str, sub, replaceWith)
}
// Split a string by the given seperator.
//
// @param str string
// @param sub string
// @param sep string
// @returns []string
var splitFunc = modules.NewFunc("split", split)

func split(_ *templating.Runtime, _ *templating.Context, str string, sep string) []string  {
	return strings.Split(str, sep)
}

// Returns the string after the given substring in a string.
//
// @param str string
// @param sub string
// @returns string
var afterFunc = modules.NewFunc("after", after)

func after(_ *templating.Runtime, _ *templating.Context, str string, sub string) string  {
	substrings := strings.SplitN(str, sub, 2)

	if len(substrings) < 2 {
		return str
	}

	return substrings[1]
}

// Returns the string before the given substring in a string.
//
// @param str string
// @param sub string
// @returns string
var beforeFunc = modules.NewFunc("before", before)

func before(_ *templating.Runtime, _ *templating.Context, str string, sub string) string  {
	substrings := strings.SplitN(str, sub, 2)

	if len(substrings) < 2 {
		return str
	}

	return substrings[1]
}

// Returns the substring between the start and end substring.
//
// @param str string
// @param startSub string
// @param endSub string
// @returns string
var betweenFunc = modules.NewFunc("between", between)

func between(_ *templating.Runtime, _ *templating.Context, str string, startSub, endSub string) string  {
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

func cutPrefix(_ *templating.Runtime, _ *templating.Context, str string, sub string) string  {
	after, _ := strings.CutPrefix(str, sub)

	return after
}

// Removes a substring at the end of the given string.
//
// @param str string
// @param sub string
// @returns string
var cutSuffixFunc = modules.NewFunc("cutSuffix", cutSuffix)

func cutSuffix(_ *templating.Runtime, _ *templating.Context, str string, sub string) string  {
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

func slice(_ *templating.Runtime, _ *templating.Context, str string, start int, end int) string  {
	return str[start:end]
}

// Joins multiple strings together by the given separator.
//
// @param sep string
// @param strings []string
// @returns string
var joinFunc = modules.NewFunc("join", join)

func join(_ *templating.Runtime, _ *templating.Context, sep string, args ...any) string  {
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

func repeat(_ *templating.Runtime, _ *templating.Context, str string, count int) string  {
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

func concat(_ *templating.Runtime, _ *templating.Context, strs ...string) string  {
	return strings.Join(strs, "")
}

var indexOfFunc = modules.NewFunc("indexOf", indexOf)

func indexOf(_ *templating.Runtime, _ *templating.Context, str string, sub string) int  {
	return strings.Index(str, sub)
}