package operators

import (
	"strings"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/internals/template/path/types"
)

func init() {
	Register(types.TemplateOperator{
		Name: "supress",
		AllowedProtocols: []types.TemplateProtocol{{ Name: "local" },{ Name: "remote" }},
		ApplyFunc: func(pathComponent, content string, context context.TemplateContext) (string, context.TemplateContext) {
			return content, context
		},
		MatchFunc: func(match string) (bool, string) {
			return strings.HasPrefix(match, "!"), "!"
		},
		Priority: 10,
	})
}