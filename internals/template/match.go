package template

import (
	"path/filepath"
	"regexp"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/internals/template/core"
)

func (t *Templater) Match(ctx *context.TemplateContext) bool {
    return matchFile(ctx)
}

var _ core.IMatcher = (*Templater)(nil)

func matchFile(ctx *context.TemplateContext) bool {
	fileName := filepath.Base(ctx.Path)

	for _, reStr := range ctx.Options.Match {
		re, err := regexp.Compile(reStr)

		if err == nil {
			if re.MatchString(fileName) {
				return true
			}
		}
	}

	return false
}