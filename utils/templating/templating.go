package templating

import (
	"bytes"
	"regexp"
	"text/template"
)

func ParseTemplate(templt *template.Template, tmplStr string, variables any) (string, error) {
	tmpl, err := templt.Parse(tmplStr)

	if err != nil {
		return "", err
	}
	var buf bytes.Buffer

	err = tmpl.Execute(&buf, variables)

	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func CreateTemplateWithFunc(name string, funcMap template.FuncMap) *template.Template {
	templt := template.New(name).Funcs(funcMap)

	templt = AddTemplateDelim(templt, "{{{", "}}}")

	return templt
}

func AddTemplateOptions(templt *template.Template, options ...string) *template.Template {
	return templt.Option(options...)
}

func AddTemplateDelim(templt *template.Template, left, right string) *template.Template {
	return templt.Delims(left, right)
}

func TransformTemplateKeys(tmplStr string, prefix string, transform func(varRegex *regexp.Regexp, m string) string) (string, error) {
	re, err := regexp.Compile(`{{{[^{}]+}}}`)

	if err != nil {
		return tmplStr, err
	}

	varRe, err := regexp.Compile(`"`+prefix+`([^"\n]+)"|`+prefix+`([a-zA-Z0-9_.]+)`)

	if err != nil {
		return tmplStr, err
	}

	tmplStr = re.ReplaceAllStringFunc(tmplStr, func(match string) string {
		return transform(varRe, match)
	})

	return tmplStr, nil
}

func AddTemplateFunc(tmplStr string, funcName string) (string, error) {
	return TransformTemplateKeys(tmplStr, ``, func(re *regexp.Regexp, match string) string {
		return re.ReplaceAllStringFunc(match, func(varMatch string) string {
			return "("+funcName+" "+varMatch+")"
		})
	})
}

func ReplaceTemplatePrefix(content string, fromPrefix, toPrefix string) (string, error) {
	tmplStr, err := TransformTemplateKeys(content, fromPrefix, func(re *regexp.Regexp, match string) string {
		re, _ = regexp.Compile(`"`+fromPrefix+`([^"\n]+)"|`+fromPrefix+`([^"\s{}]+)`)

		return re.ReplaceAllStringFunc(match, func(varMatch string) string {
			var varName string

			submatches := re.FindStringSubmatch(varMatch)

			if submatches[1] != "" {
				varName = submatches[1]
			} else {
				varName = submatches[2]
			}

			return `"` + toPrefix + varName + `"`
		})
	})

	return tmplStr, err
}