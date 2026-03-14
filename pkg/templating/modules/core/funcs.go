package core

import (
	"errors"

	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
	"github.com/google/uuid"
)

type FuncStore struct {
	funcTmpls 	map[string]string
}

type FuncOutputStore struct {
	funcOutput 	map[string][]any
}

const funcStoreID = "funcStore"
const funcOutputsStoreID = "funcOutputsStore"

func NewFuncStore() *FuncStore {
	return &FuncStore{
		funcTmpls: map[string]string{},
	}
}

func NewFuncOutputsStore() *FuncOutputStore {
	return &FuncOutputStore{
		funcOutput: map[string][]any{},
	}
}

func (s *FuncStore) Set(key string, value any) {
	s.funcTmpls[key] = value.(string)
}

func (s *FuncStore) Get(key string) any {
	return s.funcTmpls[key]
}

func (s *FuncStore) Delete(key string) bool {
	delete(s.funcTmpls, key)

	return true
}

func (s *FuncStore) Has(key string) bool {
	_, exists := s.funcTmpls[key]
	
	return exists
}

func (s *FuncStore) Keys() []string {
	keys := make([]string, 0, len(s.funcTmpls))

	for k := range s.funcTmpls {
		keys = append(keys, k)
	}

	return keys
}

func (s *FuncOutputStore) Set(key string, value any) {
	s.funcOutput[key] = value.([]any)
}

func (s *FuncOutputStore) Get(key string) any {
	return s.funcOutput[key]
}

func (s *FuncOutputStore) Delete(key string) bool {
	delete(s.funcOutput, key)

	return true
}

func (s *FuncOutputStore) Has(key string) bool {
	_, exists := s.funcOutput[key]
	
	return exists
}

func (s *FuncOutputStore) Keys() []string {
	keys := make([]string, 0, len(s.funcOutput))

	for k := range s.funcOutput {
		keys = append(keys, k)
	}

	return keys
}

func SetOutput(rt *templating.Runtime, callerID string, value []any) {
	s := rt.GetStore(funcOutputsStoreID)

	s.Set(callerID, value)
}

func GetOutputs(rt *templating.Runtime, callerID string) []any {
	s := rt.GetStore(funcOutputsStoreID)

	return s.Get(callerID).([]any)
}

func InitStores(rt *templating.Runtime) {
	if !rt.HasStore(funcStoreID) {
		err := rt.RegisterStore(funcStoreID, NewFuncStore())

		if err != nil {
			panic("error registering func store: " + err.Error())
		}
	}

	if !rt.HasStore(funcOutputsStoreID) {
		err := rt.RegisterStore(funcOutputsStoreID, NewFuncOutputsStore())

		if err != nil {
			panic("error registering func outputs store: " + err.Error())
		}
	}
}

type FuncContext struct {
	CallerID	string
	Name		string
}

const FuncContextKey templating.ContextKey = "funcContext"

var funcDefineFunc = modules.NewFunc("funcDefine", funcDefine)

func funcDefine(rt *templating.Runtime, _ templating.Context, name string, tmplBody string) any {
	InitStores(rt)

	createFunc(rt, name, tmplBody)

	return nil
}

var funcCallFunc = modules.NewFunc("funcCall", funcCall)

func funcCall(rt *templating.Runtime, ctx templating.Context, name string) any {
	InitStores(rt)

	outputs, err := callFunc(rt, ctx, name, nil)

	if err != nil {
		panic("could not call func: " + err.Error())
	}

	return outputs
}

var funcCallWithArgsFunc = modules.NewFunc("funcCallArgs", funcCallArgs)

func funcCallArgs(rt *templating.Runtime, ctx templating.Context, name string, args ...any) any {
	InitStores(rt)

	args = modules.UnpackArgs(args...)

	outputs, err := callFunc(rt, ctx, name, args...)

	if err != nil {
		panic("could not call func: " + err.Error())
	}

	return outputs
}

func createFunc(rt *templating.Runtime, name, tmplBody string) any {
	s := rt.GetStore(funcStoreID)

    s.Set(name, tmplBody)

    return ""
}

func callFunc(rt *templating.Runtime, ctx templating.Context, name string, args ...any) (any, error) {
	s := rt.GetStore(funcStoreID)
	
	exists := s.Has(name)
	
	if !exists {
		return nil, errors.New("function \"" + name + "\" not defined")
	}

	tmplBody := s.Get(name).(string)

	data := map[string]any{
		"args": args,
	}

	callerID := uuid.NewString()

	newEngine := templating.NewEngine()

	newEngine.Use(FuncModule)

	newEngine.UseModules(rt.GetEngine().GetModules()...)

	var newContext templating.Context

	funcContext := FuncContext{
		CallerID: callerID,
		Name: name,
	}
	
	ctx.Copy(&newContext)

	newContext.Set(FuncContextKey, funcContext)

	_, err := newEngine.ExecuteWithRuntime(":func:" + name, tmplBody, data, rt.GetEngineOptions().FuncDelims, newContext, rt)

	if err != nil {
		return nil, err
	}

	outputs := GetOutputs(rt, callerID)

	s.Delete(callerID)

	if len(outputs) == 1 {
		return outputs[0], nil
	}

	return outputs, nil
}