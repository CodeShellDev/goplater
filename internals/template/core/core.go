package core

import "github.com/codeshelldev/goplater/internals/template/context"

var Renderer IRenderer
var Matcher IMatcher

type IRenderer interface {
    Render(content string, context context.TemplateContext) (string, error)
}

type IMatcher interface {
	Match(context context.TemplateContext) bool
}