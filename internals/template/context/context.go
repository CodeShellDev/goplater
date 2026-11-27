package context

import (
	"errors"
	"regexp"

	"github.com/codeshelldev/goplater/pkg/gplt"
)

type TemplateOptions struct {
	Root				string
	Output				string
	Source				string

	AllowedOutputs		[]string
	AllowedInputs		[]string

	Match				[]string
	MatchRe				[]*regexp.Regexp

	Exclude				[]string
	ExcludeRe			[]*regexp.Regexp

	IncludeHidden		bool
	IncludeBinary		bool

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
	Depth		int
}

type Shared struct {
	Path		string
	OutputPath	string
}

const TemplateContextKey gplt.ContextKey = "templateContext"

func New(options *TemplateOptions) *TemplateContext {
	return &TemplateContext{Options: options, Local: Local{}, Shared: &Shared{}}
}

func Extract(ctx *gplt.Context) (*TemplateContext, error) {
	tmplContext, ok := ctx.Get(TemplateContextKey).(*TemplateContext)

	if !ok {
		return nil, errors.New("missing template context")
	}

	return tmplContext, nil
}