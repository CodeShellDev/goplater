package protocols

import (
	"regexp"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/internals/template/path/types"
	"github.com/codeshelldev/goplater/utils/fetchutils"
)

func init() {
	Register(types.TemplateProtocol{
		Name: "remote",
		ApplyFunc: func(pathComponent string, context context.TemplateContext) (string, context.TemplateContext) {
			return fetchutils.Remote(pathComponent), context
		},
		MatchFunc: func(match string) (bool, string) {
			re, err := regexp.Compile(`@://`)

			if err != nil {
				return false, match
			}

			return re.MatchString(match), re.FindString(match)
		},
	})
}