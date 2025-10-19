package template

import (
	"text/template"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/internals/template/core"
	"github.com/codeshelldev/goplater/utils/templateutils"
)

func (t *Templater) Render(content string, context context.TemplateContext) (string, error) {
    return templateContent(content, context)
}

var _ core.IRenderer = (*Templater)(nil)

func templateContent(content string, context context.TemplateContext) (string, error) {
	normalized := normalize(content)

	tmplStr, err := templateStr(normalized, context, nil)

	return tmplStr, err
}

func templateStr(str string, context context.TemplateContext, variables map[string]any) (string, error) {
	tmplStr, err := templateutils.AddTemplateFunc(str, "get")

	if err != nil {
		return str, err
	}

	templt := templateutils.CreateTemplateWithFunc(context.Path, template.FuncMap{
		"get": func (str string) any {
			return templateGet(str, context)
		},
	})

	if context.Options.Supress {
		templt = templateutils.AddTemplateOptions(templt, "missingkey=zero")
	}

	tmplStr, err = templateutils.ParseTemplate(templt, tmplStr, variables)

	if err != nil {
		return str, err
	}

	return tmplStr, nil
}