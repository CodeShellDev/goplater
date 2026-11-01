package template

import (
	"path/filepath"
	"regexp"

	"github.com/codeshelldev/goplater/internals/template/context"
)

func (t *Templater) Match(ctx *context.TemplateContext) (bool, error) {
    return matchFile(ctx)
}

func matchFile(ctx *context.TemplateContext) (bool, error) {
	fileName := filepath.Base(ctx.Path)

	for _, reStr := range ctx.Options.Match {
		re, err := regexp.Compile(reStr)

		if err != nil {
			return false, err
		}

		if re.MatchString(fileName) {
			return true, nil
		}
	}

	return false, nil
}