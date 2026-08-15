package template

import (
	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/internals/template/core"
	"github.com/codeshelldev/goplater/internals/template/funcs"
	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/collections"
	"github.com/codeshelldev/goplater/pkg/templating/resolvers"
	"github.com/codeshelldev/goplater/pkg/templating/types"
)

func (t *Templater) Render(content string, ctx *context.TemplateContext) (string, error) {
    return templateContent(content, ctx)
}

var _ core.IRenderer = (*Templater)(nil)

func templateContent(content string, ctx *context.TemplateContext) (string, error) {
	normalized := content

	tmplStr, err := templateStr(normalized, ctx)

	return tmplStr, err
}

func templateStr(str string, tmplContext *context.TemplateContext) (string, error) {
	e := templating.NewEngine()

	e.Use(funcs.Module)

	e.UseSafeModules(collections.All...)
	e.SetResolver(
		templating.NewResolverChain(
			funcs.NewFsResolver(tmplContext),
			resolvers.NewHttpResolver(),
		),
	)

	ctx := &templating.Context{}
	ctx.Set(context.TemplateContextKey, tmplContext)

	return e.Execute(tmplContext.Path, str, nil, templating.EngineOptions{
		Delims: types.Delims{
			Left: "+{{{", Right: "}}}",
		},
	}, ctx)
}