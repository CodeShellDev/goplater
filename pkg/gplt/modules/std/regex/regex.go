package regex

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
)

var cache sync.Map

func compile(expr string) (*regexp.Regexp, error) {
	re, ok := cache.Load(expr)
	if ok {
		return re.(*regexp.Regexp), nil
	}

	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, fmt.Errorf("error parsing regex: %w", err)
	}

	cache.Store(expr, re)
	return re.(*regexp.Regexp), nil
}

var Module = modules.NewModule(regexMatchFunc, regexFindAllFunc, regexFindGroupsFunc, regexFindGroupsIndexFunc, regexReplaceAllFunc)

// Returns whether a string matches the given regex pattern.
//
// @param regex string
// @params str string
// @returns bool
// @returns error
var regexMatchFunc = modules.NewFunc("regexMatch", regexMatch)

func regexMatch(_ *gplt.Runtime, _ *gplt.Context, regex string, str string) (bool, error)  {
	re, err := compile(regex)

	if err != nil {
		return false, fmt.Errorf("error parsing regex: %w", err)
	}

	return re.MatchString(str), nil
}

// Returns all (optionally up to n) matches against the given regex pattern in a string.
//
// @param regex string
// @params str string
// @params n int?
// @returns []string
// @returns error
var regexFindAllFunc = modules.NewFunc("regexFindAll", regexFindAll, regexFindAllN)

func regexFindAll(_ *gplt.Runtime, _ *gplt.Context, regex, str string) ([]string, error) {
	return regexFindAllN(nil, nil, regex, str, -1)
}

func regexFindAllN(_ *gplt.Runtime, _ *gplt.Context, regex, str string, n int) ([]string, error) {
	re, err := compile(regex)
	if err != nil {
		return nil, err
	}
	return re.FindAllString(str, n), nil
}

// Returns all matches and submatches from the given regex pattern in a string.
//
// @param regex string
// @params str string
// @returns [][]string
// @returns error
var regexFindGroupsFunc = modules.NewFunc("regexFindGroups", regexFindGroups)

func regexFindGroups(_ *gplt.Runtime, _ *gplt.Context, regex string, str string) ([][]string, error)  {
	re, err := compile(regex)

	if err != nil {
		return nil, fmt.Errorf("error parsing regex: %w", err)
	}

	return re.FindAllStringSubmatch(str, -1), nil
}

// Returns all match and submatch positions from the given regex pattern in a string.
//
// @param regex string
// @params str string
// @returns [][]int
// @returns error
var regexFindGroupsIndexFunc = modules.NewFunc("regexFindGroupsIndex", regexFindGroupsIndex)

func regexFindGroupsIndex(_ *gplt.Runtime, _ *gplt.Context, regex string, str string) ([][]int, error)  {
	re, err := compile(regex)

	if err != nil {
		return nil, fmt.Errorf("error parsing regex: %w", err)
	}

	return re.FindAllStringSubmatchIndex(str, -1), nil
}

// Replaces all matches from the given regex pattern in a string with new (which for allows using '$1', etc.).
//
// @params str string
// @param regex string
// @params new string
// @returns string
// @returns error
var regexReplaceAllFunc = modules.NewFunc("regexReplaceAll", regexReplaceAll)

func regexReplaceAll(_ *gplt.Runtime, _ *gplt.Context, str string, regex string, new string) (string, error)  {
	re, err := compile(regex)

	if err != nil {
		return "", fmt.Errorf("error parsing regex: %w", err)
	}

	return re.ReplaceAllString(str, new), nil
}