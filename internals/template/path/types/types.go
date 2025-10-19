package types

import "github.com/codeshelldev/goplater/internals/template/context"

type TemplatePath struct {
	Operator TemplateOperator
	Protocol TemplateProtocol
	PathComponent TemplatePathComponent
}

type TemplateOperator struct {
	Value string
	Name string
	Priority int
	AllowedProtocols []TemplateProtocol
	ApplyFunc func(pathComponent, content string, context context.TemplateContext) (string, context.TemplateContext)
	MatchFunc func(match string) (bool, string)
}

type TemplateProtocol struct {
	Value string
	Name string
	Priority int
	ApplyFunc func(pathComponent string, context context.TemplateContext) (string, context.TemplateContext)
	MatchFunc func(match string) (bool, string)
}

type TemplatePathComponent struct {
	Value string
	Name string
	Priority int
	AllowedProtocols []TemplateProtocol
	ApplyFunc func(key string, context context.TemplateContext) (string, context.TemplateContext)
	MatchFunc func(match string) (bool, string)
}