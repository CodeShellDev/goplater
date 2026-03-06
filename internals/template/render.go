package template

import (
	"text/template"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/internals/template/core"
	"github.com/codeshelldev/goplater/internals/template/funcs"
	"github.com/codeshelldev/gotl/pkg/templating"
)

func (t *Templater) Render(content string, context context.TemplateContext) (string, error) {
    return templateContent(content, context)
}

var _ core.IRenderer = (*Templater)(nil)

func templateContent(content string, context context.TemplateContext) (string, error) {
	normalized := content

	tmplStr, err := templateStr(normalized, context)

	return tmplStr, err
}

func templateStr(str string, context context.TemplateContext) (string, error) {
	templt := template.New(context.Path)
	templt.Delims("${{{", "}}}")
	
	templt.Funcs(funcs.GetFuncMap(context))

	err := templating.ParseTemplate(templt, str)

	if err != nil {
		return str, err
	}

	return templating.ExecuteTemplate(templt, nil)
}