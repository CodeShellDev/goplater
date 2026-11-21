package core

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
	"github.com/codeshelldev/goplater/pkg/gplt/utils/tmplerrors"
	"github.com/google/uuid"
)

type FuncOutputStore struct { 
	outputs map[string][]any 
}

const funcOutputsStoreID = "funcOutputsStore"

func NewCallOutputStore() *FuncOutputStore { 
	return &FuncOutputStore{outputs: map[string][]any{}}
}

func (s *FuncOutputStore) Set(key string, value any) {
	outputs, ok := value.([]any)
	if !ok {
		panic("'value' is expected to be []any")
	}

	if s.outputs == nil {
		s.outputs = make(map[string][]any)
	}

	s.outputs[key] = outputs
}

func (s *FuncOutputStore) Get(key string) (any, error) {
	outputs, ok := s.outputs[key]
	if !ok {
		return nil, fmt.Errorf("outputs for func %q do not exist", key)
	}

	return outputs, nil
}

func (s *FuncOutputStore) Delete(key string) bool {
	_, ok := s.outputs[key]
	if !ok {
		return false
	}

	delete(s.outputs, key)
	return true
}

func (s *FuncOutputStore) Has(key string) bool {
	_, ok := s.outputs[key]
	return ok
}

func (s *FuncOutputStore) Keys() []string {
	keys := make([]string, 0, len(s.outputs))

	for key := range s.outputs {
		keys = append(keys, key)
	}

	return keys
}

func initCallStore(rt *gplt.Runtime) error {
	if !rt.HasStore(funcOutputsStoreID) {
		err := rt.RegisterStore(funcOutputsStoreID, NewCallOutputStore())
		
		if err != nil {
			return fmt.Errorf("error registering call outputs store: %w", err)
		}
	}

	return nil
}

func SetOutput(rt *gplt.Runtime, callerID string, value []any) {
	rt.GetStore(funcOutputsStoreID).Set(callerID, value)
}

func GetOutputs(rt *gplt.Runtime, callerID string) []any {
	raw, _ := rt.GetStore(funcOutputsStoreID).Get(callerID)
	out, _ := raw.([]any)

	return out
}

type FuncContext struct {
	CallerID string
	Name	 string
}

const FuncContextKey gplt.ContextKey = "funcContext"

// Calls a func by name with args, optionally with a namespace handle.
// Fails if func is not accessible, undefined or has no signature.
// Args are bound to the params of the func signature:
//	 - `name`: required, call fails if missing
//	 - `name?`: optional (end only), `nil` if not given
//	 - `...name`: variadic (last only), a slice of all remaining args
//
// Too many args fail unless the last param is variadic.
//
// Result:
//	 - `return` not called (or called with no values): captured template output
//	 - `return` with one value: that value
//	 - `return` with multiple values: a slice of all values
//
// @param name string
// @param args []any
// @returns any
// @returns error
//
// @example
//	 +{{ define "greet(name)" }}
//		 +{{ echo (printf "Hello %s!" .name) }}
//	 +{{ end }}
//	 
//	 +{{ call "greet" "John" }}
// @output
//	 Hello John!
//
// @example
//	 +{{ define "greet(name)" }}
//		 Hello +{{ .name }}!
//	 +{{ end }}
//	 
//	 +{{ output := call "greet" "John" }}
//	 +{{ echo $output }}
// @output
//	 Hello John!
//
// @example
//	 +{{ define "greet(name)" }}
//		 +{{ return (printf "Hello %s!" .name) }}
//	 +{{ end }}
//	 
//	 +{{ output := call "greet" "John" }}
//	 +{{ echo $output }}
// @output
//	 Hello John!
var callFunc = modules.NewFunc("call", callFn)

func callFn(rt *gplt.Runtime, ctx *gplt.Context, name string, args ...any) (any, error) {
	err := initCallStore(rt)

	if err != nil {
		return nil, err
	}

	output, err := call(rt, ctx, name, args...)

	if err != nil {
		return nil, fmt.Errorf("could not call %q: %w", name, err)
	}

	return output, nil
}

func call(rt *gplt.Runtime, ctx *gplt.Context, name string, args ...any) (any, error) {
	scope, _ := ctx.Get(gplt.CurrentTreeKey).(gplt.TreeScope)

	tmpl, body := rt.GetRegistry().Lookup(scope.Name, name)

	if tmpl == nil {
		return nil, fmt.Errorf("function %q is not defined or not accessible", name)
	}

	sig, _ := rt.Signature(name)
	data, err := bindParams(sig, args)
	if err != nil {
		return nil, err
	}

	callerID := uuid.NewString()

	previousFunc, existsPreviousFunc := ctx.Get(FuncContextKey), ctx.Has(FuncContextKey)
	previousScope := scope

	ctx.Set(FuncContextKey, FuncContext{CallerID: callerID, Name: name})
	ctx.Set(gplt.CurrentTreeKey, rt.GetRegistry().ScopeOf(tmpl))

	defer func() {
		if existsPreviousFunc {
			ctx.Set(FuncContextKey, previousFunc)
		} else {
			ctx.Delete(FuncContextKey)
		}

		ctx.Set(gplt.CurrentTreeKey, previousScope)

		rt.GetStore(funcOutputsStoreID).Delete(callerID)
	}()

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, data)

	// the ReturnExitErr exists to make the func stop early
	if err != nil && !errors.Is(err, ReturnExitErr) {
		tmplErr, ok := tmplerrors.TryWrapTemplateErrorOnce(tmpl, body, err)

		if ok {
			err = tmplErr
		}
		
		return nil, err
	}
	
	outputs := GetOutputs(rt, callerID)

	switch len(outputs) {
	case 0:
		return buf.String(), nil
	case 1:
		return outputs[0], nil
	default:
		return outputs, nil
	}
}

func bindParams(sig gplt.Signature, args []any) (map[string]any, error) {
	if len(args) < sig.Required() {
		return nil, fmt.Errorf("missing argument %q", sig.Params[len(args)].Name)
	}

	max := sig.Max()
	if max >= 0 && len(args) > max {
		return nil, fmt.Errorf("too many arguments: takes %d, got %d", max, len(args))
	}

	data := make(map[string]any, len(sig.Params))

	for i, p := range sig.Params {
		switch {
		case p.Variadic:
			rest := []any{}
			if i < len(args) {
				rest = append(rest, args[i:]...)
			}
			data[p.Name] = rest
		case i < len(args):
			data[p.Name] = args[i]
		default:
			// always present as a key, but nil (allows for checking if set)
			data[p.Name] = nil
		}
	}

	return data, nil
}