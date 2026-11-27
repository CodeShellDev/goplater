package template

import (
	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/internals/template/funcs"
	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/ext/html"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/ext/markdown"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/ext/yaml"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/std"
	"github.com/codeshelldev/goplater/pkg/gplt/resolvers"
)

// share resolver for shared cache
var httpResolver = resolvers.NewHttpResolver()

func Render(content string, tmplContext *context.TemplateContext) (string, error) {
	normalized := content

	tmplStr, err := templateStr(normalized, tmplContext)

	return tmplStr, err
}

func templateStr(input string, tmplContext *context.TemplateContext) (string, error) {
	e := gplt.NewEngine()

	e.Use(funcs.Module)
	e.Use(yaml.Module, html.Module, markdown.Module)

	e.Use(std.Collections.Unsafe...)
	e.UseSafe(std.Collections.Safe...)
	
	e.SetResolver(
		gplt.NewResolverChain(
			funcs.NewFsResolver(),
			httpResolver,
		),
	)

	ctx := &gplt.Context{}
	ctx.Set(context.TemplateContextKey, tmplContext)

	output, err := e.Execute(tmplContext.Shared.Path, input, nil, gplt.EngineOptions{
		Delims: DELIMS,
		ParseFunc: parseFunc(tmplContext),
		ProcessFunc: transformFunc(tmplContext),
	}, ctx)

	return output, err
}