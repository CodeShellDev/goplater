package math

import (
	"errors"
	"fmt"
	"math"
	"reflect"
)

type number struct {
	Int       int
	Float       float64
	isFloat bool
}

func ToFloat(v any) (float64, error) {
	n, err := ToNumber(v)
	if err != nil {
		return 0, err
	}
	return n.Float, nil
}

func ToNumber(v any) (number, error) {
	value := reflect.ValueOf(v)
	switch value.Kind() {
	
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return number{Int: int(value.Int()), Float: float64(value.Int())}, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u := value.Uint()
		if u > math.MaxInt64 {
			return number{Float: float64(u), isFloat: true}, nil // promote to float (prevent overflow)
		}
		return number{Int: int(u), Float: float64(u)}, nil

	case reflect.Float32, reflect.Float64:
		return number{Float: value.Float(), isFloat: true}, nil

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
		n, err := ToNumber(a)

		if err != nil {
			return nil, err
		}

		nums[i] = n
		anyFloat = anyFloat || n.isFloat
	}

	var err error
	if anyFloat {
		res := nums[0].Float
		for _, n := range nums[1:] {
			res, err = floatOp(res, n.Float)

			if err != nil {
				return nil, err
			}
		}
		return res, nil
	}

	res := nums[0].Int
	for _, n := range nums[1:] {
		res, err = intOp(res, n.Int)

		if err != nil {
			return nil, err
		}
	}
	return res, nil
}