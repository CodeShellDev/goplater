package funcs

import (
	"reflect"

	"github.com/codeshelldev/goplater/internals/template/context"
)

type TemplateFunc struct {
	Name string
	Handler any
}

var funcs = map[string]any{}

func Register(f TemplateFunc) {
	funcs[f.Name] = f.Handler
}

func GetFuncMap(context context.TemplateContext) map[string]any{
	m := make(map[string]any, len(funcs))

	for k, v := range funcs {
		m[k] = bindContext(v, context)
	}

	return m
}

func bindContext(fn any, ctx ...any) any {
	v := reflect.ValueOf(fn)
	t := v.Type()

	if t.Kind() != reflect.Func {
		panic("bindContext: fn is not a function")
	}

	if t.NumIn() < len(ctx) {
		panic("bindContext: function must have at least as many parameters as ctx values")
	}

	newFuncType := reflect.FuncOf(
		func() []reflect.Type {
			ins := []reflect.Type{}
			
			for i := len(ctx); i < t.NumIn(); i++ {
				ins = append(ins, t.In(i))
			}

			return ins
		}(),

		func() []reflect.Type {
			outs := []reflect.Type{}

			for out := range t.Outs() {
				outs = append(outs, out)
			}

			return outs
		}(),

		t.IsVariadic(),
	)

	newFunc := reflect.MakeFunc(newFuncType, func(args []reflect.Value) (results []reflect.Value) {
		ctxValues := make([]reflect.Value, len(ctx))

		for i, c := range ctx {
			ctxValues[i] = reflect.ValueOf(c)
		}

		allArgs := append(ctxValues, args...)
		return v.Call(allArgs)
	})

	return newFunc.Interface()
}