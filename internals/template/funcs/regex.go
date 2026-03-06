package funcs

import (
	"regexp"

	"github.com/codeshelldev/goplater/internals/template/context"
)

var regexMatchFunc = TemplateFunc{
	Name: "regexMatch",
	Handler: func(context context.TemplateContext, regex string, str string) bool {
		re, err := regexp.Compile(regex)

		if err != nil {
			panic("error parsing regex: " + err.Error())
		}

		return re.MatchString(str)
	},
}

var regexFindFunc = TemplateFunc{
	Name: "regexFind",
	Handler: func(context context.TemplateContext, regex string, str string) []string {
		re, err := regexp.Compile(regex)

		if err != nil {
			panic("error parsing regex: " + err.Error())
		}

		return re.FindAllString(str, -1)
	},
}

var regexFindGroupsFunc = TemplateFunc{
	Name: "regexFindGroups",
	Handler: func(context context.TemplateContext, regex string, str string) [][]string {
		re, err := regexp.Compile(regex)

		if err != nil {
			panic("error parsing regex: " + err.Error())
		}

		return re.FindAllStringSubmatch(str, -1)
	},
}

var regexReplaceFunc = TemplateFunc{
	Name: "regexReplace",
	Handler: func(context context.TemplateContext, str string, regex string, replaceWith string) string {
		re, err := regexp.Compile(regex)

		if err != nil {
			panic("error parsing regex: " + err.Error())
		}

		return re.ReplaceAllString(str, replaceWith)
	},
}

func init() {
	Register(regexMatchFunc)
	Register(regexFindFunc)
	Register(regexFindGroupsFunc)
	Register(regexReplaceFunc)
}
