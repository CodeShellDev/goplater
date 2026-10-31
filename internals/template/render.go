package template

import (
	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/internals/template/funcs"
	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/collections"
	"github.com/codeshelldev/goplater/pkg/gplt/resolvers"
)

func (t *Templater) Render(content string, ctx *context.TemplateContext) (string, error) {
    return templateContent(content, ctx)
}

func templateContent(content string, ctx *context.TemplateContext) (string, error) {
	normalized := content

	tmplStr, err := templateStr(normalized, ctx)

	return tmplStr, err
}

func templateStr(str string, tmplContext *context.TemplateContext) (string, error) {
	e := gplt.NewEngine()

	e.Use(funcs.Module)

	e.Use(collections.Unsafe...)
	e.UseSafe(collections.Safe...)
	
	e.SetResolver(
		gplt.NewResolverChain(
			funcs.NewFsResolver(tmplContext),
			resolvers.NewHttpResolver(),
		),
	)

	ctx := &gplt.Context{}
	ctx.Set(context.TemplateContextKey, tmplContext)

	return e.Execute(tmplContext.Path, str, nil, gplt.EngineOptions{
		Delims: DELIMS,
		ParseFunc: parseFunc(tmplContext),
		ProcessFunc: transformFunc(tmplContext),
	}, ctx)
}