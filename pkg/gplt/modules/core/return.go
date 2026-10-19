package core

import (
	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
)

func mustFuncContext(ctx *gplt.Context) FuncContext {
	fc, ok := ctx.Get(FuncContextKey).(FuncContext)
	
	if !ok {
		panic("this function can only be used inside a func body invoked via \"call\"")
	}

	return fc
}


// Sets an output value at the given index.
//
// @param i int
// @param value any
// @returns error
var returnFunc = modules.NewFunc("return", returnFn)

func returnFn(rt *gplt.Runtime, ctx *gplt.Context, i int, value any) string {
	funcContext := mustFuncContext(ctx)

	outputs := GetOutputs(rt, funcContext.CallerID)

	for len(outputs) <= i {
		panic("index out of bounds")
	}
	outputs[i] = value

	SetOutput(rt, funcContext.CallerID, outputs)
	return ""
}

// Appends a value to the output.
//
// @param value any
var returnNextFunc = modules.NewFunc("returnNext", returnNext)

func returnNext(rt *gplt.Runtime, ctx *gplt.Context, value any) string {
	funcContext := mustFuncContext(ctx)

	outputs := GetOutputs(rt, funcContext.CallerID)

	outputs = append(outputs, value)

	SetOutput(rt, funcContext.CallerID, outputs)
	return ""
}

// Sets the output object as a whole.
//
// @param value []any
//
// @example
//	 +{{ returnAll "1" "2" "3" }}
// @example
//	 +{{ returnAll (sliceCreate "1" "2" "3") }}
var returnAllFunc = modules.NewFunc("returnAll", returnAll)

func returnAll(rt *gplt.Runtime, ctx *gplt.Context, values ...any) string {
	funcContext := mustFuncContext(ctx)

	SetOutput(rt, funcContext.CallerID, values)
	return ""
}

// Returns current outputs.
//
// @returns []any
var getOutputsFunc = modules.NewFunc("getOutputs", getOutputs)

func getOutputs(rt *gplt.Runtime, ctx *gplt.Context) []any {
	funcContext := mustFuncContext(ctx)

	outputs := GetOutputs(rt, funcContext.CallerID)

	return outputs
}