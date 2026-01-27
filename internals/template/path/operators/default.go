package operators

import (
	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/internals/template/core"
	"github.com/codeshelldev/goplater/internals/template/path/types"
)

func init() {
	Register(types.TemplateOperator{
		Name: "default",
		AllowedProtocols: []types.TemplateProtocol{{ Name: "local" },{ Name: "remote" }},
		ApplyFunc: func(pathComponent, content string, context context.TemplateContext) (string, context.TemplateContext) {
			if core.Matcher.Match(context) {
				content, _ = core.Renderer.Render(content, context)
			}

			return content, context
		},
		MatchFunc: func(match string) (bool, string) {
			return true, match
		},
	})
}