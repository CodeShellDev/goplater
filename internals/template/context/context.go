package context

import "github.com/codeshelldev/goplater/pkg/templating"

type TemplateOptions struct {
	Output string
	Source string

	Whitespace []string
	Match []string

	Recursive bool
	Flatten bool

	Verbose bool
	Supress bool
}

type TemplateContext struct {
	Invoker 	string
	Path 		string
	Options		TemplateOptions
}

const TemplateContextKey templating.ContextKey = "templateContext"

func New() TemplateContext {
	return TemplateContext{ Options: TemplateOptions{ } }
}