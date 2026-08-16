package math

import (
	"errors"
	"math"

	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
)

// Returns the abs value of a number.
//
// @param numbers int|float64
// @returns int|float64
// @returns error
var absFunc = modules.NewFunc("abs", abs)

func abs(_ *templating.Runtime, _ *templating.Context, number any) (any, error) {
	switch n := number.(type) {
	case int:
		if n < 0 {
			return -n, nil
		}

		return n, nil
	case float32:
		return math.Abs(float64(n)), nil
	case float64:
		return math.Abs(n), nil
	default:
		return nil, errors.New("abs requires an int or float")
	}
}

// Returns the smaller value of two numbers.
//
// @param a int|float64
// @param b int|float64
// @returns int|float64
// @returns error
var minFunc = modules.NewFunc("min", min)

func min(_ *templating.Runtime, _ *templating.Context, a, b any) (any, error) {
	return handleOperator([]any{a, b},
		func(a, b int) int {
			if b < a {
				return b
			}

			return a
		},
		func(a, b float64) float64 {
			if b < a {
				return b
			}

			return a
		},
	)
}

// Returns the bigger value of two numbers.
//
// @param a int|float64
// @param b int|float64
// @returns int|float64
// @returns error
var maxFunc = modules.NewFunc("max", max)

func max(_ *templating.Runtime, _ *templating.Context, a, b any) (any, error) {
	return handleOperator([]any{a, b},
		func(a, b int) int {
			if b > a {
				return b
			}

			return a
		},
		func(a, b float64) float64 {
			if b > a {
				return b
			}

			return a
		},
	)
}

// Clamps a value between a minimum and maximum value.
//
// @param value int|float64
// @param min int|float64
// @param max int|float64
// @returns int|float64
// @returns error
var clampFunc = modules.NewFunc("clamp", clamp)

func clamp(_ *templating.Runtime, _ *templating.Context, value, min, max any) (any, error) {
	switch v := value.(type) {
	case int:
		minValue, ok := min.(int)
		if !ok {
			return nil, errors.New("clamp requires all arguments to have the same type")
		}

		maxValue, ok := max.(int)
		if !ok {
			return nil, errors.New("clamp requires all arguments to have the same type")
		}

		if minValue > maxValue {
			return nil, errors.New("clamp minimum cannot be greater than maximum")
		}

		if v < minValue {
			return minValue, nil
		}

		if v > maxValue {
			return maxValue, nil
		}

		return v, nil

	case float64:
		minValue, ok := min.(float64)
		if !ok {
			return nil, errors.New("clamp requires all arguments to have the same type")
		}

		maxValue, ok := max.(float64)
		if !ok {
			return nil, errors.New("clamp requires all arguments to have the same type")
		}

		if minValue > maxValue {
			return nil, errors.New("clamp minimum cannot be greater than maximum")
		}

		if v < minValue {
			return minValue, nil
		}

		if v > maxValue {
			return maxValue, nil
		}

		return v, nil

	default:
		return nil, errors.New("clamp requires int or float64 arguments")
	}
}

// Performs pow on base with exponent.
//
// @param base int|float64
// @param exponent int|float64
// @returns float64
// @returns error
var powFunc = modules.NewFunc("pow", pow)

func pow(_ *templating.Runtime, _ *templating.Context, base, exponent any) (float64, error) {
	var baseFloat float64
	var exponentFloat float64

	switch v := base.(type) {
	case int:
		baseFloat = float64(v)
	case float32:
		exponentFloat = float64(v)
	case float64:
		baseFloat = v
	default:
		return 0, errors.New("pow base must be an int or float")
	}

	switch v := exponent.(type) {
	case int:
		exponentFloat = float64(v)
	case float32:
		exponentFloat = float64(v)
	case float64:
		exponentFloat = v
	default:
		return 0, errors.New("pow exponent must be an int or float")
	}

	return math.Pow(baseFloat, exponentFloat), nil
}

// Returns the square root of a number.
//
// @param number int|float64
// @returns float64
// @returns error
var sqrtFunc = modules.NewFunc("sqrt", sqrt)

func sqrt(_ *templating.Runtime, _ *templating.Context, number any) (float64, error) {
	var n float64

	switch v := number.(type) {
	case int:
		n = float64(v)
	case float32:
		n = float64(v)
	case float64:
		n = v
	default:
		return 0, errors.New("sqrt requires an int or float")
	}

	if n < 0 {
		return 0, errors.New("sqrt requires a non-negative number")
	}

	return math.Sqrt(n), nil
}

// Rounds a number to the nearest integer.
//
// @param number int|float64
// @returns int
// @returns error
var roundFunc = modules.NewFunc("round", round)

func round(_ *templating.Runtime, _ *templating.Context, number any) (int, error) {
	var n float64

	switch v := number.(type) {
	case int:
		return v, nil
	case float32:
		n = float64(v)
	case float64:
		n = v
	default:
		return 0, errors.New("round requires an int or float")
	}

	return int(math.Round(n)), nil
}

// Returns the greatest integer value less than or equal to number.
//
// @param number int|float64
// @returns int
// @returns error
var floorFunc = modules.NewFunc("floor", floor)

func floor(_ *templating.Runtime, _ *templating.Context, number any) (int, error) {
	var n float64

	switch v := number.(type) {
	case int:
		return v, nil
	case float32:
		n = float64(v)
	case float64:
		n = v
	default:
		return 0, errors.New("floor requires an int or float")
	}

	return int(math.Floor(n)), nil
}

// Returns the least integer value greater than or equal to number.
//
// @param number int|float64
// @returns int
// @returns error
var ceilFunc = modules.NewFunc("ceil", ceil)

func ceil(_ *templating.Runtime, _ *templating.Context, number any) (int, error) {
	var n float64

	switch v := number.(type) {
	case int:
		return v, nil
	case float32:
		n = float64(v)
	case float64:
		n = v
	default:
		return 0, errors.New("ceil requires an int or float")
	}

	return int(math.Ceil(n)), nil
}