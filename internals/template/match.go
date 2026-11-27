package template

import (
	"path/filepath"
	"regexp"

	"github.com/codeshelldev/goplater/internals/template/context"
)

func (t *Templater) Match(tmplContext *context.TemplateContext) (bool, error) {
	return matchFile(tmplContext)
}

func matchFile(tmplContext *context.TemplateContext) (bool, error) {
	fileName := filepath.Base(tmplContext.Shared.Path)

	for _, reStr := range tmplContext.Options.Match {
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