package math

import (
	"errors"
	"math"

	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
)

// Adds numbers together.
//
// @param numbers [](int|float64)
// @returns int|float64
// @returns error
var addFunc = modules.NewFunc("add", add)

func add(_ *gplt.Runtime, _ *gplt.Context, a, b any, rest ...any) (any, error) {
	return reduce(
		append([]any{a, b}, rest...),
		func(x, y int) (int, error) { return x + y, nil },
		func(x, y float64) (float64, error) { return x + y, nil })
}

// Subtracts numbers from each other.
//
// @param numbers [](int|float64)
// @returns int|float64
// @returns error
var subFunc = modules.NewFunc("sub", sub)

func sub(_ *gplt.Runtime, _ *gplt.Context, a, b any, rest ...any) (any, error) {
	return reduce(
		append([]any{a, b}, rest...),
		func(x, y int) (int, error) { return x - y, nil },
		func(x, y float64) (float64, error) { return x - y, nil })
}

// Multiplies numbers together.
//
// @param numbers [](int|float64)
// @returns int|float64
// @returns error
var mulFunc = modules.NewFunc("mul", mul)

func mul(_ *gplt.Runtime, _ *gplt.Context, a, b any, rest ...any) (any, error) {
	return reduce(
		append([]any{a, b}, rest...),
		func(x, y int) (int, error) { return x * y, nil },
		func(x, y float64) (float64, error) { return x * y, nil })
}

// Divides numbers from each other.
//
// @param numbers [](int|float64)
// @returns int|float64
// @returns error
var divFunc = modules.NewFunc("div", div)

func div(_ *gplt.Runtime, _ *gplt.Context, a, b any, rest ...any) (any, error) {
	return reduce(
		append([]any{a, b}, rest...),
		func(x, y int) (int, error) {
			if y == 0 {
				return 0, errors.New("division by zero")
			}
			return x / y, nil
		},
		func(x, y float64) (float64, error) {
			if y == 0 {
				return 0, errors.New("division by zero")
			}
			return x / y, nil
		})
}
// Performs modulo on numbers.
//
// @param numbers [](int|float64)
// @returns int|float64
// @returns error
var modFunc = modules.NewFunc("mod", mod)

func mod(_ *gplt.Runtime, _ *gplt.Context, a, b any, rest ...any) (any, error) {
	return reduce(
		append([]any{a, b}, rest...),
		func(x, y int) (int, error) {
			if y == 0 {
				return 0, errors.New("modulo by zero")
			}
			return x % y, nil
		},
		func(x, y float64) (float64, error) { return math.Mod(x, y), nil })
}