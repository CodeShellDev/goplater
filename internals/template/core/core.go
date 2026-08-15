package core

import (
	"github.com/codeshelldev/goplater/internals/template/context"
)

var Renderer IRenderer
var Matcher IMatcher

type IRenderer interface {
    Render(content string, ctx *context.TemplateContext) (string, error)
}

type IMatcher interface {
	Match(ctx *context.TemplateContext) bool
}