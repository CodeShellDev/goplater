package core

import (
	"errors"
	"fmt"
	"slices"

	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
)

func mustFuncContext(ctx *gplt.Context) (FuncContext, error) {
	fc, ok := ctx.Get(FuncContextKey).(FuncContext)
	
	if !ok {
		return FuncContext{}, errors.New("this function can only be used inside a func body invoked via 'call'")
	}

	return fc, nil
}


// Sets an output value at the given index (default: -1 – append).
//
// @param i int?
// @param value any
// @returns error
var returnFunc = modules.NewFunc("return", returnNext, returnAll)

func returnNext(rt *gplt.Runtime, ctx *gplt.Context, value any) (string, error) {
	fc, err := mustFuncContext(ctx)
	if err != nil {
		return "", err
	}

	outputs := slices.Clone(GetOutputs(rt, fc.CallerID))
	outputs = append(outputs, value)

	SetOutput(rt, fc.CallerID, outputs)
	return "", nil
}

func returnAt(rt *gplt.Runtime, ctx *gplt.Context, i int, value any) (string, error) {
	fc, err := mustFuncContext(ctx)
	if err != nil {
		return "", err
	}

	outputs := slices.Clone(GetOutputs(rt, fc.CallerID))
	if i < 0 || i >= len(outputs) {
		return "", fmt.Errorf("could not set output at %d: index out of bounds", i)
	}
	
	outputs[i] = value

	SetOutput(rt, fc.CallerID, outputs)
	return "", nil
}

// Sets the output object as a whole.
//
// @param value []any
// @returns error
//
// @example
//	 +{{ returnAll "1" "2" "3" }}
// @example
//	 +{{ returnAll (sliceCreate "1" "2" "3") }}
var returnAllFunc = modules.NewFunc("returnAll", returnAll, returnAllSlice)

func returnAll(rt *gplt.Runtime, ctx *gplt.Context, values ...any) (string, error) {
	funcContext, err := mustFuncContext(ctx)

	if err != nil {
		return "", err
	}

	SetOutput(rt, funcContext.CallerID, values)
	return "", nil
}

func returnAllSlice(rt *gplt.Runtime, ctx *gplt.Context, values []any) (string, error) {
	return returnAll(rt, ctx, values...)
}

// Returns current outputs.
//
// @returns []any
var getOutputsFunc = modules.NewFunc("getOutputs", getOutputs)

func getOutputs(rt *gplt.Runtime, ctx *gplt.Context) ([]any, error) {
	funcContext, err := mustFuncContext(ctx)

	if err != nil {
		return nil, err
	}

	outputs := GetOutputs(rt, funcContext.CallerID)

	return outputs, nil
}