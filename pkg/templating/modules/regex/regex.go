package regex

import (
	"regexp"

	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
)

var Module = modules.NewModule(regexMatchFunc, regexFindFunc, regexFindGroupsFunc, regexReplaceFunc)

var regexMatchFunc = modules.NewFunc("regexMatch", regexMatch)

func regexMatch(_ *templating.Runtime, _ templating.Context, regex string, str string) bool  {
	re, err := regexp.Compile(regex)

	if err != nil {
		panic("error parsing regex: " + err.Error())
	}

	return re.MatchString(str)
}

var regexFindFunc = modules.NewFunc("regexFind", regexFind)

func regexFind(_ *templating.Runtime, _ templating.Context, regex string, str string) []string  {
	re, err := regexp.Compile(regex)

	if err != nil {
		panic("error parsing regex: " + err.Error())
	}

	return re.FindAllString(str, -1)
}

var regexFindGroupsFunc = modules.NewFunc("regexFindGroups", regexFindGroups)

func regexFindGroups(_ *templating.Runtime, _ templating.Context, regex string, str string) [][]string  {
	re, err := regexp.Compile(regex)

	if err != nil {
		panic("error parsing regex: " + err.Error())
	}

	return re.FindAllStringSubmatch(str, -1)
}

var regexReplaceFunc = modules.NewFunc("regexReplace", regexReplace)

func regexReplace(_ *templating.Runtime, _ templating.Context, str string, regex string, replaceWith string) string  {
	re, err := regexp.Compile(regex)

	if err != nil {
		panic("error parsing regex: " + err.Error())
	}

	return re.ReplaceAllString(str, replaceWith)
}