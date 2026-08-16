package core

import (
	"bytes"
	"errors"

	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
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
	s.outputs[key] = value.([]any)
}

func (s *FuncOutputStore) Get(key string) (any, error) {
	outputs, exists := s.outputs[key]
	
	if !exists {
		return nil, errors.New("outputs for func '" + key + "' do not exist")
	}
	
	return outputs, nil
}

func (s *FuncOutputStore) Delete(key string) bool {
	delete(s.outputs, key)

	return true
}

func (s *FuncOutputStore) Has(key string) bool {
	_, exists := s.outputs[key]
	
	return exists
}

func (s *FuncOutputStore) Keys() []string {
	keys := make([]string, 0, len(s.outputs))

	for k := range s.outputs {
		keys = append(keys, k)
	}

	return keys
}

func initCallStore(rt *templating.Runtime) {
	if !rt.HasStore(funcOutputsStoreID) {
		err := rt.RegisterStore(funcOutputsStoreID, NewCallOutputStore())
		
		if err != nil {
			panic("error registering call outputs store: " + err.Error())
		}
	}
}

func SetOutput(rt *templating.Runtime, callerID string, value []any) {
	rt.GetStore(funcOutputsStoreID).Set(callerID, value)
}

func GetOutputs(rt *templating.Runtime, callerID string) []any {
	raw, _ := rt.GetStore(funcOutputsStoreID).Get(callerID)
	out, _ := raw.([]any)

	return out
}

type FuncContext struct {
	CallerID string
	Name     string
}

const FuncContextKey templating.ContextKey = "funcContext"

var callFunc = modules.NewFunc("call", callFn)

// Tries to call a func by name with args, optionally with a namespace handle.
// Fails if func is not accessible or undefined.
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
func callFn(rt *templating.Runtime, ctx *templating.Context, name string, args ...any) (any, error) {
	initCallStore(rt)

	output, err := call(rt, ctx, name, args...)

	if err != nil {
		return nil, errors.New("could not call \"" + name + "\": " + err.Error())
	}

	return output, nil
}

func call(rt *templating.Runtime, ctx *templating.Context, name string, args ...any) (any, error) {
	scope, _ := ctx.Get(templating.CurrentTreeKey).(templating.TreeScope)

	tmpl := rt.GetRegistry().Lookup(scope.Name, name)

	if tmpl == nil {
		return nil, errors.New("function " + name + " is not defined or not accessible")
	}

	params, _ := rt.Params(name)
	data := bindParams(params, args)

	callerID := uuid.NewString()

	previousFunc, existsPreviousFunc := ctx.Get(FuncContextKey), ctx.Has(FuncContextKey)
	previousScope := scope

	ctx.Set(FuncContextKey, FuncContext{CallerID: callerID, Name: name})
	ctx.Set(templating.CurrentTreeKey, rt.GetRegistry().ScopeOf(tmpl))

	defer func() {
		if existsPreviousFunc {
			ctx.Set(FuncContextKey, previousFunc)
		} else {
			ctx.Delete(FuncContextKey)
		}

		ctx.Set(templating.CurrentTreeKey, previousScope)

		rt.GetStore(funcOutputsStoreID).Delete(callerID)
	}()

	var buf bytes.Buffer

	err := tmpl.Execute(&buf, data)
	if err != nil {
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

func bindParams(params []string, args []any) map[string]any {
	data := map[string]any{}
	
	for i, p := range params {
		if i < len(args) {
			data[p] = args[i]
		}
	}

	return data
}