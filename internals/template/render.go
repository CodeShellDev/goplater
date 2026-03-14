package template

import (
	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/internals/template/core"
	"github.com/codeshelldev/goplater/internals/template/funcs"
	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/collections"
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

func templateStr(str string, tmplContext context.TemplateContext) (string, error) {
	e := templating.NewEngine()

	e.Use(funcs.Module)

	e.UseModules(collections.All...)

	ctx := templating.Context{}
	ctx.Set(context.TemplateContextKey, tmplContext)

	return e.Execute(tmplContext.Path, str, nil, templating.EngineOptions{
		Delims: templating.Delims{
			Left: "+{{{", Right: "}}}",
		},
		FuncDelims: templating.Delims{
			Left: "{{{", Right: "}}}",
		},
	}, ctx)
}