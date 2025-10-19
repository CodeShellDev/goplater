package template

import (
	"path/filepath"
	"regexp"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/internals/template/core"
)

func (t *Templater) Match(context context.TemplateContext) bool {
    return matchFile(context)
}

var _ core.IMatcher = (*Templater)(nil)

func matchFile(context context.TemplateContext) bool {
	fileName := filepath.Base(context.Path)

	for _, reStr := range context.Options.Match {
		re, err := regexp.Compile(reStr)

		if err == nil {
			if re.MatchString(fileName) {
				return true
			}
		}
	}

	return false
}