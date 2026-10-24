package context

import "github.com/codeshelldev/goplater/pkg/gplt"

type TemplateOptions struct {
	Root 			string
	Output 			string
	Source 			string

	AllowedOutputs 	[]string

	Match 			[]string

	Recursive 		bool
	Flatten 		bool

	DryRun 			bool
	IgnoreErrors 	bool
}

type TemplateContext struct {
	Invoker 	string
	Path 		string
	OutputPath	string
	Options		*TemplateOptions
}

const TemplateContextKey gplt.ContextKey = "templateContext"

func New(options *TemplateOptions) *TemplateContext {
	return &TemplateContext{ Options: options }
}