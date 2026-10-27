package template

import (
	"text/template"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/pkg/transformer"
	"github.com/codeshelldev/goplater/pkg/transformer/processors/blanklines"
)

func transformFunc(tmplContext *context.TemplateContext) func(t *template.Template) error {
	return func(t *template.Template) error {
		tf := transformer.New()

		blankLines := tmplContext.Options.TrimMaxBlankLines
		if (blankLines >= 0) {
			tf.Use(blanklines.NewWith(
				blanklines.WithMaxBlankLines(blankLines),
			))
		}

		return tf.Process(t)
	}
}