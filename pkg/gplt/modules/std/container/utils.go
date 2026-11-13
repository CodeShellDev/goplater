package container

import (
	"fmt"
	"reflect"

	"github.com/codeshelldev/goplater/pkg/gplt/modules/std/math"
)

// access underlying value of data (by following pointer refs)
func deref(data any) reflect.Value {
	value := reflect.ValueOf(data)
	for value.Kind() == reflect.Pointer && !value.IsNil() {
		value = value.Elem()
	}
	return value
}

func ensureContainer(data any) (reflect.Value, error) {
	value := deref(data)
	if value.Kind() != reflect.Map && value.Kind() != reflect.Slice {
		return value, fmt.Errorf("not a map or slice, got %T", data)
	}
	return value, nil
}

// can x be stored in type t?
func convertTo(x any, t reflect.Type) (reflect.Value, bool) {
	if x == nil {
		switch t.Kind() {
		case reflect.Interface, reflect.Pointer, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan:
			return reflect.Zero(t), true
		}
		return reflect.Value{}, false
	}

	v := reflect.ValueOf(x)
	if v.Type().AssignableTo(t) {
		return v, true
	}
	
	return reflect.Value{}, false
}

func equal(a, b any) bool {
	na, errA := math.ToNumber(a)
	nb, errB := math.ToNumber(b)

	if errA == nil && errB == nil {
		return na.Float == nb.Float
	}

	return reflect.DeepEqual(a, b)
}