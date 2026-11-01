package math

import (
	"fmt"

	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
)

// Adds numbers together.
//
// @param numbers [](int|float64)
// @returns int|float64
// @returns error
var addFunc = modules.NewFunc("add", add)

func add(_ *gplt.Runtime, _ *gplt.Context, args ...any) (any, error) {
	err := modules.RequireArgs(2, args)

	if err != nil {
		return nil, err
	}

	return handleOperator(args, func(a, b int) int { return a + b }, func(a, b float64) float64 { return a + b })
}

// Subtracts numbers from each other.
//
// @param numbers [](int|float64)
// @returns int|float64
// @returns error
var subFunc = modules.NewFunc("sub", sub)

func sub(_ *gplt.Runtime, _ *gplt.Context, args ...any) (any, error) {
	err := modules.RequireArgs(2, args)

	if err != nil {
		return nil, err
	}

	return handleOperator(args, func(a, b int) int { return a - b }, func(a, b float64) float64 { return a - b })
}

// Multiplies numbers together.
//
// @param numbers [](int|float64)
// @returns int|float64
// @returns error
var multFunc = modules.NewFunc("mult", mult)

func mult(_ *gplt.Runtime, _ *gplt.Context, args ...any) (any, error) {
	err := modules.RequireArgs(2, args)

	if err != nil {
		return nil, err
	}

	return handleOperator(args, func(a, b int) int { return a * b }, func(a, b float64) float64 { return a * b })
}

// Divides numbers from each other.
//
// @param numbers [](int|float64)
// @returns int|float64
// @returns error
var divdFunc = modules.NewFunc("divd", divd)

func divd(_ *gplt.Runtime, _ *gplt.Context, args ...any) (any, error) {
	err := modules.RequireArgs(2, args)

	if err != nil {
		return nil, err
	}

	return handleOperator(args, func(a, b int) int { return a / b }, func(a, b float64) float64 { return a / b })
}

// Performs modulo on numbers.
//
// @param numbers [](int|float64)
// @returns int|float64
// @returns error
var modFunc = modules.NewFunc("mod", mod)

func mod(_ *gplt.Runtime, _  *gplt.Context, args ...any) (any, error) {
	err := modules.RequireArgs(2, args)

	if err != nil {
		return nil, err
	}

	return handleOperator(args, func(a, b int) int { return a / b }, func(a, b float64) float64 { return a / b })
}

func handleOperator(values []any, intOp func(a, b int) int, floatOp func(a, b float64) float64) (any, error) {
	if containsFloat(values) {
		numbers, err := toFloatSlice(values)

		if err != nil {
			return nil, err
		}

		var res float64

		res = numbers[0]

		for _, num := range numbers[1:] {
			res = floatOp(res, num)
		}

		return res, nil
	} else {
		numbers, err := toIntSlice(values)

		if err != nil {
			return nil, err
		}

		var res int

		res = numbers[0]

		for _, num := range numbers[1:] {
			res = intOp(res, num)
		}

		return res, nil
	}
}

func containsFloat(values []any) bool {
	for _, v := range values {
        switch v.(type) {
		case float32:
			return true
		case float64:
			return true
		}
    }

	return false
}

func toFloatSlice(values []any) ([]float64, error) {
    out := make([]float64, len(values))

    for i, v := range values {
        switch v := v.(type) {
		case int:
			out[i] = float64(v)
		case float32:
			out[i] = float64(v)
		case float64:
			out[i] = v
		default:
			return nil, fmt.Errorf("%v is not of a numeric type (%T)", v, v)
		}
    }

    return out, nil
}

func toIntSlice(values []any) ([]int, error) {
    out := make([]int, len(values))

    for i, v := range values {
		num, ok := v.(int)

		if !ok {
			return nil, fmt.Errorf("%v is not of type int (%T)", v, v)
		}

		out[i] = num
    }

    return out, nil
}