package regex

import (
	"errors"
	"regexp"

	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
)

var Module = modules.NewModule(regexMatchFunc, regexFindFunc, regexFindGroupsFunc, regexFindGroupsIndexFunc, regexReplaceFunc)

// Returns whether a string matches the given regex pattern.
//
// @param regex string
// @params str string
// @returns bool
// @returns error
var regexMatchFunc = modules.NewFunc("regexMatch", regexMatch)

func regexMatch(_ *templating.Runtime, _ *templating.Context, regex string, str string) (bool, error)  {
	re, err := regexp.Compile(regex)

	if err != nil {
		return false, errors.New("error parsing regex: " + err.Error())
	}

	return re.MatchString(str), nil
}

// Returns all matches against the given regex pattern in a string.
//
// @param regex string
// @params str string
// @returns []string
// @returns error
var regexFindFunc = modules.NewFunc("regexFind", regexFind)

func regexFind(_ *templating.Runtime, _ *templating.Context, regex string, str string) ([]string, error)  {
	re, err := regexp.Compile(regex)

	if err != nil {
		return nil, errors.New("error parsing regex: " + err.Error())
	}

	return re.FindAllString(str, -1), nil
}

// Returns all matches and submatches from the given regex pattern in a string.
//
// @param regex string
// @params str string
// @returns [][]string
// @returns error
var regexFindGroupsFunc = modules.NewFunc("regexFindGroups", regexFindGroups)

func regexFindGroups(_ *templating.Runtime, _ *templating.Context, regex string, str string) ([][]string, error)  {
	re, err := regexp.Compile(regex)

	if err != nil {
		return nil, errors.New("error parsing regex: " + err.Error())
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

func regexFindGroupsIndex(_ *templating.Runtime, _ *templating.Context, regex string, str string) ([][]int, error)  {
	re, err := regexp.Compile(regex)

	if err != nil {
		return nil, errors.New("error parsing regex: " + err.Error())
	}

	return re.FindAllStringSubmatchIndex(str, -1), nil
}

// Replaces all matches from the given regex pattern in a string with the replaceWith param (which allows using '$1', etc.).
//
// @param regex string
// @params str string
// @params replaceWith string
// @returns string
// @returns error
var regexReplaceFunc = modules.NewFunc("regexReplace", regexReplace)

func regexReplace(_ *templating.Runtime, _ *templating.Context, regex string, str string, replaceWith string) string  {
	re, err := regexp.Compile(regex)

	if err != nil {
		panic("error parsing regex: " + err.Error())
	}

	return re.ReplaceAllString(str, replaceWith)
}