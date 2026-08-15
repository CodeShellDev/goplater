package core

import (
	"bytes"
	"fmt"

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

func (s *FuncOutputStore) Get(key string) any {
	return s.outputs[key]
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
	out, _ := rt.GetStore(funcOutputsStoreID).Get(callerID).([]any)
	return out
}

type FuncContext struct {
	CallerID string
	Name     string
}

const FuncContextKey templating.ContextKey = "funcContext"

var callFunc = modules.NewFunc("call", callFn)

func callFn(rt *templating.Runtime, ctx *templating.Context, name string, args ...any) any {
	initCallStore(rt)

	args = modules.UnpackArgs(args...)

	output, err := call(rt, ctx, name, args...)

	if err != nil {
		panic("could not call \"" + name + "\": " + err.Error())
	}

	return output
}

func call(rt *templating.Runtime, ctx *templating.Context, name string, args ...any) (any, error) {
	scope, _ := ctx.Get(templating.CurrentTreeKey).(templating.TreeScope)

	tmpl := rt.GetRegistry().Lookup(scope.Name, name)

	if tmpl == nil {
		return nil, fmt.Errorf("function %q is not defined or not accessible", name)
	}

	params, _ := rt.Params(name)
	data := bindParams(params, args)

	callerID := uuid.NewString()

	previous, existsPrevious := ctx.Get(FuncContextKey), ctx.Has(FuncContextKey)

	ctx.Set(FuncContextKey, FuncContext{CallerID: callerID, Name: name})

	defer func() {
		if existsPrevious {
			ctx.Set(FuncContextKey, previous)
		} else {
			ctx.Delete(FuncContextKey)
		}
		
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
		// no return -> just render text
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