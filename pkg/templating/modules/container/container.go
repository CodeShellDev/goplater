package container

import (
	"errors"
	"reflect"
	"strconv"

	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
)

var Module = modules.NewModule(containerDeleteFunc, containerSetFunc, containerHasFunc, containerIncludesFunc, slicePushFunc, sliceCreateFunc, mapCreateFunc)

// Deletes a container entry with the given key.
//
// @param container map[any]any|[]any
// @param key any
// @returns any
// @returns error
// 
// @example
//	 +{{ $newSlice := delete (sliceCreate "Apple" "Banana") 0 }}
//	 +{{ echo $newSlice }}
// @output
//	 [Banana]
var containerDeleteFunc = modules.NewFunc("delete", delete)

func delete(_ *templating.Runtime, _ *templating.Context, container any, key any) (any, error) {
	err := ensureInterfaceIsContainer(container)
	return deleteKey(container, key), err
}

// Set a container entry value with the given key.
//
// @param container map[any]any|[]any
// @param key any
// @param value any
// @returns any
// @returns error
// 
// @example
//	 +{{ $newSlice := set (sliceCreate "Apple" "Banana") 0 "Strawberry" }}
//	 +{{ echo $newSlice }}
// @output
//	 [Strawberry Banana]
var containerSetFunc = modules.NewFunc("set", set)

func set(_ *templating.Runtime, _ *templating.Context, container any, key any, value any) (any, error) {
	err := ensureInterfaceIsContainer(container)
	return setKey(container, key, value), err
}

// Returns if an entry with the given key exists in container.
//
// @param container map[any]any|[]any
// @param key any
// @returns bool
// @returns error
// 
// @example
//	 +{{ if (has (mapCreate "version" "v1") "version") }}
//		Container has "version"!
//	 +{{ end }}
var containerHasFunc = modules.NewFunc("has", has)

func has(_ *templating.Runtime, _ *templating.Context, container any, key any) (bool, error) {
	err := ensureInterfaceIsContainer(container)
	return hasKey(container, key), err
}

var containerIncludesFunc = modules.NewFunc("includes", includes)

// Returns if an entry with the given value exists in container.
//
// @param container map[any]any|[]any
// @param value any
// @returns bool
// @returns error
// 
// @example
//	 +{{ if (includes (mapCreate "version" "v1") "v1") }}
//		Container includes "v1"!
//	 +{{ end }}
func includes(_ *templating.Runtime, _ *templating.Context, container any, value any) (bool, error) {
	err := ensureInterfaceIsContainer(container)
	return hasValue(container, value), err
}

// Appends a value to a slice.
//
// @param container []any
// @param value any
// @returns []any
// 
// @example
//	 +{{ $newSlice := sliceCreate "Apple" "Banana" }}
//	 +{{ $newSlice = slicePush $newSlice "Strawberry" }}
//	 +{{ echo $newSlice }}
// @output
//	 [Apple Banana Strawberry]
var slicePushFunc = modules.NewFunc("slicePush", slicePush)

func slicePush(_ *templating.Runtime, _ *templating.Context, container []any, value any) []any {
	return append(container, value)
}

// Creates a new slice, optionally with values.
//
// @param values []any
// @returns []any
// 
// @example
//	 +{{ $newSlice := sliceCreate "Apple" "Banana" }}
//	 +{{ echo $newSlice }}
// @output
//	 [Apple Banana]
var sliceCreateFunc = modules.NewFunc("sliceCreate", sliceCreate)

func sliceCreate(_ *templating.Runtime, _ *templating.Context, values ...any) []any {
	return values
}

// Creates a new map, optionally with key value pairs.
// Pairs are constructed one by one, a key followed by the given value.
// Currently only supports string keys!
//
// @param entries []any
// @returns map[string]any
// 
// @example
//	 +{{ $newMap := mapCreate "key1" "value1" "key2" "value2" }}
//	 +{{ echo $newMap }}
// @output
//	 map[key1:value1 key2:value2]
var mapCreateFunc = modules.NewFunc("mapCreate", mapCreate)

func mapCreate(_ *templating.Runtime, _ *templating.Context, entries ...any) map[string]any {
	if len(entries) == 0 {
		return map[string]any{}
	}

	if len(entries) % 2 != 0 {
		panic("missing value for key '" + entries[len(entries) - 1].(string) + "'")
	}

	out := map[string]any{}

	for i := 1; i < len(entries); i += 2 {
		key, ok := entries[i - 1].(string)

		if !ok {
			panic("item at index " + strconv.Itoa(i) + " is not of type string, map keys must be strings")
		}

		out[key] = entries[i]
	}

	return out
}

func ensureInterfaceIsContainer(data any) error {
	val := reflect.ValueOf(data)
	
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Map:
		return nil

	case reflect.Slice, reflect.Array:
		return nil
	default:
		return errors.New("object is not of container type")
	}
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