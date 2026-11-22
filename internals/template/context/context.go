package context

import "github.com/codeshelldev/goplater/pkg/gplt"

type TemplateOptions struct {
	Root				string
	Output				string
	Source				string

	AllowedOutputs		[]string
	AllowedInputs		[]string

	Match				[]string

	Recursive			bool
	Flatten				bool

	TrimMaxBlankLines	int

	DryRun				bool
	IgnoreErrors		bool
}

type TemplateContext struct {
	Options		*TemplateOptions	// read-only
	Local		Local				// copied each read/import
	Shared		*Shared				// shared across top-level file run (entry file)
}

type Local struct {
	Invoker		string
}

type Shared struct {
	Path		string
	OutputPath	string
}

const TemplateContextKey gplt.ContextKey = "templateContext"

func New(options *TemplateOptions) *TemplateContext {
	return &TemplateContext{ Options: options }
}