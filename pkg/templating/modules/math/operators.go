package math

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
)

// Adds numbers together.
//
// @param numbers [](int|float64)
// @returns int|float64
// @returns error
var addFunc = modules.NewFunc("add", add)

func add(_ *templating.Runtime, _ *templating.Context, args ...any) (any, error) {
	if len(args) < 2 {
		return nil, errors.New("wrong number of args for add: want at least 2 got " + strconv.Itoa(len(args)))
	}

	return handleOperator(args, func(a, b int) int { return a + b }, func(a, b float64) float64 { return a + b })
}

// Subtracts numbers from each other.
//
// @param numbers [](int|float64)
// @returns int|float64
// @returns error
var subFunc = modules.NewFunc("sub", sub)

func sub(_ *templating.Runtime, _ *templating.Context, args ...any) (any, error) {
	if len(args) < 2 {
		return nil, errors.New("wrong number of args for sub: want at least 2 got " + strconv.Itoa(len(args)))
	}

	return handleOperator(args, func(a, b int) int { return a - b }, func(a, b float64) float64 { return a - b })
}

// Multiplies numbers together.
//
// @param numbers [](int|float64)
// @returns int|float64
// @returns error
var multFunc = modules.NewFunc("mult", mult)

func mult(_ *templating.Runtime, _ *templating.Context, args ...any) (any, error) {
	if len(args) < 2 {
		return nil, errors.New("wrong number of args for mult: want at least 2 got " + strconv.Itoa(len(args)))
	}

	return handleOperator(args, func(a, b int) int { return a * b }, func(a, b float64) float64 { return a * b })
}

// Divides numbers from each other.
//
// @param numbers [](int|float64)
// @returns int|float64
// @returns error
var divdFunc = modules.NewFunc("divd", divd)

func divd(_ *templating.Runtime, _ *templating.Context, args ...any) (any, error) {
	if len(args) < 2 {
		return nil, errors.New("wrong number of args for divd: want at least 2 got " + strconv.Itoa(len(args)))
	}

	return handleOperator(args, func(a, b int) int { return a / b }, func(a, b float64) float64 { return a / b })
}

// Performs modulo on numbers.
//
// @param numbers [](int|float64)
// @returns int|float64
// @returns error
var modFunc = modules.NewFunc("mod", mod)

func mod(_ *templating.Runtime, _  *templating.Context, args ...any) (any, error) {
	if len(args) < 2 {
		return nil, errors.New("wrong number of args for mod: want at least 2 got " + strconv.Itoa(len(args)))
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
			return nil, errors.New(fmt.Sprint(v) + " is not of a numeric type")
		}
    }

    return out, nil
}

func toIntSlice(values []any) ([]int, error) {
    out := make([]int, len(values))

    for i, v := range values {
		num, ok := v.(int)

		if !ok {
			return nil, errors.New(fmt.Sprint(v) + " is not of type int")
		}

		out[i] = num
    }

    return out, nil
}