package container

import (
	"reflect"

	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
)

var Module = modules.NewModule(containerDeleteFunc, containerSetFunc, containerHasFunc, containerIncludesFunc, slicePushFunc, sliceCreateFunc)

var containerDeleteFunc = modules.NewFunc("delete", delete)

func delete(_ *templating.Runtime, _ templating.Context, container any, key any) any  {
	return deleteKey(container, key)
}

var containerSetFunc = modules.NewFunc("set", set)

func set(_ *templating.Runtime, _ templating.Context, container any, key any, value any) any  {
	return setKey(container, key, value)
}

var containerHasFunc = modules.NewFunc("has", has)

func has(_ *templating.Runtime, _ templating.Context, container any, key any) bool  {
	return hasKey(container, key)
}

var containerIncludesFunc = modules.NewFunc("includes", includes)

func includes(_ *templating.Runtime, _ templating.Context, container any, value any) bool  {
	return hasValue(container, value)
}

var slicePushFunc = modules.NewFunc("slicePush", slicePush)

func slicePush(_ *templating.Runtime, _ templating.Context, container []any, value any) []any  {
	return append(container, value)
}

var sliceCreateFunc = modules.NewFunc("sliceCreate", sliceCreate)

func sliceCreate(_ *templating.Runtime, _ templating.Context, value ...any) []any  {
	value = modules.UnpackArgs(value...)

	return value
}

func hasKey(data any, key any) bool {
	val := reflect.ValueOf(data)

	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if reflect.DeepEqual(k.Interface(), key) {
				return true
			}
		}

		return false

	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(val.Index(i).Interface(), key) {
				return true
			}
		}

		return false

	default:
		return false
	}
}

func hasValue(data any, value any) bool {
	val := reflect.ValueOf(data)

	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Map:
		iter := val.MapRange()

		for iter.Next() {
			if reflect.DeepEqual(iter.Value().Interface(), value) {
				return true
			}
		}
		return false

	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(val.Index(i).Interface(), value) {
				return true
			}
		}

		return false

	default:
		return false
	}
}

func deleteKey(data any, key any) any {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if reflect.DeepEqual(k.Interface(), key) {
				val.SetMapIndex(k, reflect.Value{})
			}
		}
	case reflect.Slice:
		newLen := 0

		for i := 0; i < val.Len(); i++ {
			if !reflect.DeepEqual(val.Index(i).Interface(), key) {
				val.Index(newLen).Set(val.Index(i))

				newLen++
			}
		}

		val.Set(val.Slice(0, newLen))
	}

	return data
}

func setKey(data any, key any, value any) any {
	val := reflect.ValueOf(data)
	
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Map:
		val.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
	case reflect.Slice:
		i, ok := key.(int)
		
		if ok && i >= 0 {
			for val.Len() <= i {
				val.Set(reflect.Append(val, reflect.Zero(val.Type().Elem())))
			}

			val.Index(i).Set(reflect.ValueOf(value))
		}
	}

	return data
}