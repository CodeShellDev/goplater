package pathcomponents

import (
	"fmt"
	"regexp"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/internals/template/path/types"
)

func init() {
	Register(types.TemplatePathComponent{
		Name: "http",
		AllowedProtocols: []types.TemplateProtocol{{ Name: "remote" }},
		ApplyFunc: func(key string, context context.TemplateContext) (string, context.TemplateContext) {
			return key, context
		},
		MatchFunc: func(match string) (bool, string) {
			re, err := regexp.Compile(`(https?://([a-zA-Z0-9\-]+\.[a-zA-Z0-9\-\.]+)([a-zA-Z0-9\-\._~!$&'()\*\+,;=:@/%]+)?(\?[A-Za-z0-9\-\._~!$&'()\*\+,;=:@/\?]+)?)`)

			if err != nil {
				fmt.Println("regex error:", err.Error())

				return false, match
			}

			if re.MatchString(match) {
				return true, re.FindString(match)
			}

			return false, match
		},
		Priority: 10,
	})
}