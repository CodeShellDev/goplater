package context

import "github.com/codeshelldev/goplater/pkg/templating"

type TemplateOptions struct {
	Output string
	Source string

	AllowedOutputFolders []string

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
	OutputPath	string
	Options		*TemplateOptions
}

const TemplateContextKey templating.ContextKey = "templateContext"

func New(options *TemplateOptions) *TemplateContext {
	return &TemplateContext{ Options: options }
}