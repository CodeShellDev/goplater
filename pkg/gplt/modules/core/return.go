package core

import (
	"errors"

	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
)

var OutsideFuncErr = errors.New(`only callable within a "call"-invoked function body`)

func mustFuncContext(ctx *gplt.Context) (FuncContext, error) {
	fc, ok := ctx.Get(FuncContextKey).(FuncContext)
	
	if !ok {
		return FuncContext{}, OutsideFuncErr
	}

	return fc, nil
}

var ReturnExitErr = errors.New("return exit (you should not see this message)")

// Sets output value(s).
//
// @param values []any
//
// @example
//	 +{{ return "a" "b" 1 2 }}
// @example
//	 +{{ return (list "3" "4" "5") }}
var returnFunc = modules.NewFunc("return", returnFn)

func returnFn(rt *gplt.Runtime, ctx *gplt.Context, values ...any) (string, error) {
	funcContext, err := mustFuncContext(ctx)

	if err != nil {
		return "", err
	}

	SetOutput(rt, funcContext.CallerID, values)
	return "", ReturnExitErr
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