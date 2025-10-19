package operators

import "github.com/codeshelldev/goplater/internals/template/path/types"

var registry = map[string]types.TemplateOperator{}

func Register(op types.TemplateOperator) {
	registry[op.Name] = op
}

func Get(name string) (types.TemplateOperator, bool) {
	op, ok := registry[name]
	return op, ok
}

func All() map[string]types.TemplateOperator {
	return registry
}