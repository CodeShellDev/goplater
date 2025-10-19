package protocols

import (
	"github.com/codeshelldev/goplater/internals/template/path/types"
)

var registry = map[string]types.TemplateProtocol{}

func Register(prot types.TemplateProtocol) {
	registry[prot.Name] = prot
}

func Get(name string) (types.TemplateProtocol, bool) {
	op, ok := registry[name]
	return op, ok
}

func All() map[string]types.TemplateProtocol {
	return registry
}