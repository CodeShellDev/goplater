package math

import (
	"errors"
	"fmt"
	"math"
	"reflect"
)

type number struct {
	i       int
	f       float64
	isFloat bool
}

func toFloat(v any) (float64, error) {
	n, err := toNumber(v)
	if err != nil {
		return 0, err
	}
	return n.f, nil
}

func toNumber(v any) (number, error) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return number{i: int(rv.Int()), f: float64(rv.Int())}, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u := rv.Uint()
		if u > math.MaxInt64 {
			return number{f: float64(u), isFloat: true}, nil // promote to float (prevent overflow)
		}
		return number{i: int(u), f: float64(u)}, nil

	case reflect.Float32, reflect.Float64:
		return number{f: rv.Float(), isFloat: true}, nil

	}
	return number{}, fmt.Errorf("%v is not numeric (%T)", v, v)
}

func reduce(args []any, intOp func(a, b int) (int, error), floatOp func(a, b float64) (float64, error)) (any, error) {
	if len(args) == 0 {
		return nil, errors.New("no arguments")
	}

	nums := make([]number, len(args))
	anyFloat := false

	for i, a := range args {
		n, err := toNumber(a)

		if err != nil {
			return nil, err
		}

		nums[i] = n
		anyFloat = anyFloat || n.isFloat
	}

	var err error
	if anyFloat {
		res := nums[0].f
		for _, n := range nums[1:] {
			res, err = floatOp(res, n.f)

			if err != nil {
				return nil, err
			}
		}
		return res, nil
	}

	res := nums[0].i
	for _, n := range nums[1:] {
		res, err = intOp(res, n.i)

		if err != nil {
			return nil, err
		}
	}
	return res, nil
}