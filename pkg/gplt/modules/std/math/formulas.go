package math

import (
	"errors"
	"math"

	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
)

// Returns the abs value of a number.
//
// @param numbers int|float64
// @returns int|float64
// @returns error
var absFunc = modules.NewFunc("abs", abs)

func abs(_ *gplt.Runtime, _ *gplt.Context, number any) (any, error) {
	n, err := toNumber(number)
	if err != nil {
		return nil, err
	}
	
	if n.isFloat {
		return math.Abs(n.f), nil
	}

	if n.i < 0 {
		return -n.i, nil
	}

	return n.i, nil
}

// Returns the smaller value of two numbers.
//
// @param a int|float64
// @param b int|float64
// @returns int|float64
// @returns error
var minFunc = modules.NewFunc("min", minFn)

func minFn(_ *gplt.Runtime, _ *gplt.Context, a, b any, rest ...any) (any, error) {
	return reduce(
		append([]any{a, b}, rest...),
		func(x, y int) (int, error) { return min(x, y), nil },
		func(x, y float64) (float64, error) { return math.Min(x, y), nil })
}

// Returns the bigger value of two numbers.
//
// @param a int|float64
// @param b int|float64
// @returns int|float64
// @returns error
var maxFunc = modules.NewFunc("max", maxFn)

func maxFn(_ *gplt.Runtime, _ *gplt.Context, a, b any, rest ...any) (any, error) {
	return reduce(
		append([]any{a, b}, rest...),
		func(x, y int) (int, error) { return max(x, y), nil },
		func(x, y float64) (float64, error) { return math.Max(x, y), nil })
}

// Clamps a value between a minimum and maximum value.
//
// @param value int|float64
// @param min int|float64
// @param max int|float64
// @returns int|float64
// @returns error
var clampFunc = modules.NewFunc("clamp", clamp)

func clamp(rt *gplt.Runtime, ctx *gplt.Context, value, min, max any) (any, error) {
	l, err := toNumber(min)
	if err != nil {
		return nil, err
	}

	h, err := toNumber(max)
	if err != nil {
		return nil, err
	}

	if h.f < l.f {
		return nil, errors.New("clamp minimum cannot be greater than maximum")
	}

	m, err := maxFn(rt, ctx, value, min)
	if err != nil {
		return nil, err
	}

	return minFn(rt, ctx, m, max)
}

// Performs pow on base with exponent.
//
// @param base int|float64
// @param exponent int|float64
// @returns float64
// @returns error
var powFunc = modules.NewFunc("pow", pow)

func pow(_ *gplt.Runtime, _ *gplt.Context, base, exp any) (float64, error) {
	b, err := toFloat(base)
	if err != nil {
		return 0, err
	}

	e, err := toFloat(exp)
	if err != nil {
		return 0, err
	}

	return math.Pow(b, e), nil
}

// Returns the square root of a number.
//
// @param number int|float64
// @returns float64
// @returns error
var sqrtFunc = modules.NewFunc("sqrt", sqrt)

func sqrt(_ *gplt.Runtime, _ *gplt.Context, number any) (float64, error) {
	n, err := toFloat(number)
	if err != nil {
		return 0, err
	}

	if n < 0 {
		return 0, errors.New("sqrt requires a non-negative number")
	}

	return math.Sqrt(n), nil
}

// Rounds a number to the specified number of decimal digits (default: 0).
//
// @param number int|float64
// @param digits int?
// @returns int
// @returns error
var roundFunc = modules.NewFunc("round", round, roundDigits)

func round(_ *gplt.Runtime, _ *gplt.Context, number any) (int, error) {
	return rounder(number, math.Round)
}

func roundDigits(_ *gplt.Runtime, _ *gplt.Context, number any, digits int) (float64, error) {
	n, err := toFloat(number)
	if err != nil {
		return 0, err
	}

	if math.IsNaN(n) || math.IsInf(n, 0) {
		return n, nil
	}

	p := math.Pow10(digits)
	return math.Round(n * p) / p, nil
}

func rounder(v any, op func(float64) float64) (int, error) {
	n, err := toNumber(v)
	if err != nil {
		return 0, err
	}

	if !n.isFloat {
		return n.i, nil
	}

	return int(op(n.f)), nil
}

// Returns the greatest integer value less than or equal to number.
//
// @param number int|float64
// @returns int
// @returns error
var floorFunc = modules.NewFunc("floor", floor)

func floor(_ *gplt.Runtime, _ *gplt.Context, number any) (int, error) {
	return rounder(number, math.Floor)
}

// Returns the least integer value greater than or equal to number.
//
// @param number int|float64
// @returns int
// @returns error
var ceilFunc = modules.NewFunc("ceil", ceil)

func ceil(_ *gplt.Runtime, _ *gplt.Context, number any) (int, error) {
	return rounder(number, math.Ceil)
}