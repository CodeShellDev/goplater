package template

import (
	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/internals/template/funcs"
	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/collections"
	"github.com/codeshelldev/goplater/pkg/gplt/resolvers"
	"github.com/codeshelldev/goplater/pkg/gplt/types"
	"github.com/codeshelldev/goplater/pkg/parser"
	"github.com/codeshelldev/goplater/pkg/parser/processors/blanklines"
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
	p := parser.New()

	blankLines := tmplContext.Options.TrimMaxBlankLines
	if (blankLines >= 0) {
		p.Use(blanklines.NewWith(
			blanklines.WithMaxBlankLines(blankLines),
		))
	}

	e := gplt.NewEngine()

	e.Use(funcs.Module)

	e.UseSafe(collections.All...)
	e.SetResolver(
		gplt.NewResolverChain(
			funcs.NewFsResolver(tmplContext),
			resolvers.NewHttpResolver(),
		),
	)

	ctx := &gplt.Context{}
	ctx.Set(context.TemplateContextKey, tmplContext)

	return e.Execute(tmplContext.Path, str, nil, gplt.EngineOptions{
		Delims: types.Delims{
			Left: "+{{", Right: "}}",
		},
		PreparseFunc: p.Process,
	}, ctx)
}