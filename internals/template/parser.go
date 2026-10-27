package template

import (
	"text/template"

	"github.com/codeshelldev/goplater/internals/parser"
	"github.com/codeshelldev/goplater/internals/parser/processors/ignore"
	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/pkg/gplt/types"
)

var DELIMS = types.Delims{
	Left: "+{{", 
	Right: "}}",
}

func parseFunc(_ *context.TemplateContext) func(s string, t *template.Template) error {
	return func (s string, t *template.Template) error {
		p := parser.New().Delims(DELIMS)

		p.Use(ignore.New())

		parsed, err := p.Process(s)

		if err != nil {
			return err
		}

		_, err = t.Parse(parsed)

		return err
	}
}