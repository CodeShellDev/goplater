package core

import (
	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
)

var FuncModule = modules.NewModule(returnFunc, returnNextFunc, returnAllFunc, returnOutputsFunc, getOutputsFunc)

var returnFunc = modules.NewFunc("return", returnFn)

func returnFn(rt *templating.Runtime, ctx templating.Context, i int, value any) any  {
	funcContext := ctx.Get(FuncContextKey).(FuncContext)

	outputs := GetOutputs(rt, funcContext.CallerID)

	for len(outputs) <= i {
		outputs = append(outputs, nil)
	}

	outputs[i] = value

	SetOutput(rt, funcContext.CallerID, outputs)
	return nil
}

var returnNextFunc = modules.NewFunc("returnNext", returnNext)

func returnNext(rt *templating.Runtime, ctx templating.Context, value any) any  {
	funcContext := ctx.Get(FuncContextKey).(FuncContext)

	outputs := GetOutputs(rt, funcContext.CallerID)

	outputs = append(outputs, value)

	SetOutput(rt, funcContext.CallerID, outputs)
	return nil
}

var returnAllFunc = modules.NewFunc("returnAll", returnAll)

func returnAll(rt *templating.Runtime, ctx templating.Context, values ...any) any  {
	funcContext := ctx.Get(FuncContextKey).(FuncContext)

	values = modules.UnpackArgs(values...)

	SetOutput(rt, funcContext.CallerID, values)
	return nil
}

var returnOutputsFunc = modules.NewFunc("returnOutputs", returnOutputs)

func returnOutputs(rt *templating.Runtime, ctx templating.Context, value []any) any  {
	funcContext := ctx.Get(FuncContextKey).(FuncContext)

	SetOutput(rt, funcContext.CallerID, value)
	return nil
}

var getOutputsFunc = modules.NewFunc("getOutputs", getOutputs)

func getOutputs(rt *templating.Runtime, ctx templating.Context) []any  {
	funcContext := ctx.Get(FuncContextKey).(FuncContext)

	return GetOutputs(rt, funcContext.CallerID)
}