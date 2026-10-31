package core

import (
	"errors"
	"strconv"

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


// Sets an output value at the given index.
//
// @param i int
// @param value any
// @returns error
var returnFunc = modules.NewFunc("return", returnFn)

func returnFn(rt *gplt.Runtime, ctx *gplt.Context, i int, value any) (string, error) {
	funcContext, err := mustFuncContext(ctx)

	if err != nil {
		return "", err
	}

	outputs := GetOutputs(rt, funcContext.CallerID)

	for len(outputs) <= i {
		return "", errors.New("could not set output at " + strconv.Itoa(i) + ": index out of bounds")
	}
	outputs[i] = value

	SetOutput(rt, funcContext.CallerID, outputs)
	return "", nil
}

// Appends a value to the output.
//
// @param value any
// @returns error
var returnNextFunc = modules.NewFunc("returnNext", returnNext)

func returnNext(rt *gplt.Runtime, ctx *gplt.Context, value any) (string, error) {
	funcContext, err := mustFuncContext(ctx)

	if err != nil {
		return "", err
	}

	outputs := GetOutputs(rt, funcContext.CallerID)

	outputs = append(outputs, value)

	SetOutput(rt, funcContext.CallerID, outputs)
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
var returnAllFunc = modules.NewFunc("returnAll", returnAll)

func returnAll(rt *gplt.Runtime, ctx *gplt.Context, values ...any) (string, error) {
	funcContext, err := mustFuncContext(ctx)

	if err != nil {
		return "", err
	}

	SetOutput(rt, funcContext.CallerID, values)
	return "", nil
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