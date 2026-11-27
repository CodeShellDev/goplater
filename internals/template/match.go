package template

import (
	"fmt"
	"regexp"

	"github.com/codeshelldev/goplater/internals/template/context"
)

func CompilePatterns(in []string) ([]*regexp.Regexp, error) {
	out := make([]*regexp.Regexp, 0, len(in))

	for _, s := range in {
		re, err := regexp.Compile(s)

		if err != nil {
			return nil, fmt.Errorf("invalid pattern %q: %w", s, err)
		}
		out = append(out, re)
	}

	return out, nil
}

// all dotfiles and dotfolders except "." and ".."
func isHidden(name string) bool {
	return len(name) > 1 && name[0] == '.' && name != ".."
}

func anyMatch(res []*regexp.Regexp, s string) bool {
	for _, re := range res {
		if re.MatchString(s) {
			return true
		}
	}
	return false
}

func isExcluded(opts *context.TemplateOptions, rel string) bool {
	return anyMatch(opts.ExcludeRe, rel)
}

func matchRelative(opts *context.TemplateOptions, rel string) bool {
	return anyMatch(opts.MatchRe, rel)
}