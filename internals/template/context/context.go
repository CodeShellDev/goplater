package context

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

func New() TemplateContext {
	return TemplateContext{ Options: TemplateOptions{ } }
}