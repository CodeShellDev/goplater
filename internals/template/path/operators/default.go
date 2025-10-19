package operators

import (
	"path/filepath"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/internals/template/core"
	"github.com/codeshelldev/goplater/internals/template/path/types"
)

func init() {
	Register(types.TemplateOperator{
		Name: "default",
		AllowedProtocols: []types.TemplateProtocol{{ Name: "local" },{ Name: "remote" }},
		ApplyFunc: func(pathComponent, content string, context context.TemplateContext) (string, context.TemplateContext) {
			newContext := context
			newContext.Path = filepath.Base(context.Path)
			newContext.Invoker = context.Path

			if core.Matcher.Match(newContext) {
				content, _ = core.Renderer.Render(content, newContext)
			}

			return content, context
		},
		MatchFunc: func(match string) (bool, string) {
			return true, match
		},
	})
}