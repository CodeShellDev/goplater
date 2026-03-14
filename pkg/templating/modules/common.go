package modules

import (
	"sort"
	"strings"
)

type Module struct {
	funcs 	map[string]Func
	id 		string
}

func NewModule(funcs ...Func) Module {
	m := Module{
		funcs: map[string]Func{},
	}

	m.id = ModuleKey(m)
	
	for _, f := range funcs {
		m.funcs[f.Name] = f
	}

	return m
}

func (m Module) GetFuncMap() map[string]Func {
	return m.funcs
}

type Func struct {
	Name string
	Handler any
}

func NewFunc(name string, handler any) Func {
	return Func{
		Name: name,
		Handler: handler,
	}
}

func UnpackArgs(args ...any) []any {
	if len(args) == 1 {
		inner, ok := args[0].([]any)

		if ok {
			args = inner
		}
	}

	return args
}

func ModuleKey(m Module) string {
	keys := make([]string, 0, len(m.funcs))

	for k := range m.funcs {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return strings.Join(keys, ",")
}

func UniqueModules(modules []Module) []Module {
	seen := make(map[string]struct{})
	var result []Module

	for _, m := range modules {
		key := ModuleKey(m)

		_, ok := seen[key]
		if !ok {
			seen[key] = struct{}{}
			result = append(result, m)
		}
	}

	return result
}