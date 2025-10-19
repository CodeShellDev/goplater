package template

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/internals/template/path"
	"github.com/codeshelldev/goplater/utils/templateutils"
)

func normalize(content string) string {
	normalizedContent, err := quoteTemplate(content)

	if err != nil {
		return content
	}

	return normalizedContent
}

func quoteTemplate(content string) (string, error) {
	tmplStr, err := templateutils.TransformTemplateKeys(content, func(re *regexp.Regexp, match string) string {
		varRe, _ := regexp.Compile(path.GetRegex())

		return varRe.ReplaceAllStringFunc(match, func(varMatch string) string {
			path, err := path.GetTemplatePath(varMatch)
			
			if err != nil {
				fmt.Println("syntax error:", err)
				return varMatch
			}

			tmplPath := path.Operator.Value + path.Protocol.Value + path.PathComponent.Value

			return `"` + tmplPath + `"`
		})
	})

	return tmplStr, err
}

func removeWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func cleanOutput(str string, context context.TemplateContext) string {
	var res string = str

	if slices.Contains(context.Options.Whitespace, "l") {
		res = strings.TrimLeftFunc(res, removeWhitespace)
	}

	if slices.Contains(context.Options.Whitespace, "t") {
		res = strings.TrimRightFunc(res, removeWhitespace)
	}

	return res
}