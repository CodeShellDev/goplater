package pathcomponents

import (
	"github.com/codeshelldev/goplater/internals/template/path/types"
)

var registry = map[string]types.TemplatePathComponent{}

func Register(pathComp types.TemplatePathComponent) {
	registry[pathComp.Name] = pathComp
}

func Get(name string) (types.TemplatePathComponent, bool) {
	op, ok := registry[name]
	return op, ok
}

func All() map[string]types.TemplatePathComponent {
	return registry
}
