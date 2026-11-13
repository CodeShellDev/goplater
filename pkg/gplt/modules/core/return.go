package core

import (
	"errors"

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


// Sets output value(s).
//
// @param values []any
// @returns error
//
// @example
//	 +{{ return "a" "b" 1 2 }}
// @example
//	 +{{ return (sliceCreate "3" "4" "5") }}
var returnFunc = modules.NewFunc("return", returnAll)

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