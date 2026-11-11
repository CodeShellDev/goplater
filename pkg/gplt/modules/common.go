package modules

import (
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"
)

func RequireMinArgs(count int, args []any) error {
	if len(args) < count {
		return fmt.Errorf("wrong number of args: want at least %d got %d", count, len(args))
	}
	return nil
}

func RequireRangeArgs(min, max int, args []any) error {
	if len(args) < min || len(args) > max {
		return fmt.Errorf("wrong number of args: want %d to %d got %d", min, max, len(args))
	}
	return nil
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

func MergeModules(modules ...Module) Module {
	funcs := []Func{}

	for _, m := range modules {
		funcs = slices.AppendSeq(funcs, maps.Values(m.funcs))
	}

	return NewModule(funcs...)
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