package core

import (
	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
)

func mustFuncContext(ctx *templating.Context) FuncContext {
	fc, ok := ctx.Get(FuncContextKey).(FuncContext)
	
	if !ok {
		panic("this function can only be used inside a func body invoked via \"call\"")
	}

	return fc
}

var returnFunc = modules.NewFunc("return", returnFn)

func returnFn(rt *templating.Runtime, ctx *templating.Context, i int, value any) any {
	funcContext := mustFuncContext(ctx)

	outputs := GetOutputs(rt, funcContext.CallerID)
	for len(outputs) <= i {
		outputs = append(outputs, nil)
	}
	outputs[i] = value

	SetOutput(rt, funcContext.CallerID, outputs)
	return nil
}

var returnNextFunc = modules.NewFunc("returnNext", returnNext)

func returnNext(rt *templating.Runtime, ctx *templating.Context, value any) any {
	funcContext := mustFuncContext(ctx)

	outputs := GetOutputs(rt, funcContext.CallerID)
	outputs = append(outputs, value)

	SetOutput(rt, funcContext.CallerID, outputs)
	return nil
}

var returnAllFunc = modules.NewFunc("returnAll", returnAll)

func returnAll(rt *templating.Runtime, ctx *templating.Context, values ...any) any {
	funcContext := mustFuncContext(ctx)

	values = modules.UnpackArgs(values...)
	SetOutput(rt, funcContext.CallerID, values)
	return nil
}

var returnOutputsFunc = modules.NewFunc("returnOutputs", returnOutputs)

func returnOutputs(rt *templating.Runtime, ctx *templating.Context, value []any) any {
	funcContext := mustFuncContext(ctx)

	SetOutput(rt, funcContext.CallerID, value)
	return nil
}

var getOutputsFunc = modules.NewFunc("getOutputs", getOutputs)

func getOutputs(rt *templating.Runtime, ctx *templating.Context) []any {
	funcContext := mustFuncContext(ctx)

	return GetOutputs(rt, funcContext.CallerID)
}