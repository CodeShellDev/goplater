package container

import (
	"fmt"
	"reflect"

	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
)

var Module = modules.NewModule(containerDeleteFunc, containerSetFunc, containerHasFunc, containerIncludesFunc, slicePushFunc, sliceCreateFunc, mapCreateFunc)

// Deletes a container entry with the given key.
//
// @param container map|slice
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

func delete(_ *gplt.Runtime, _ *gplt.Context, container, key any) (any, error) {
	value, err := ensureContainer(container)
	if err != nil {
		return nil, err
	}

	// case 1: map
	if value.Kind() == reflect.Map {
		k, ok := convertTo(key, value.Type().Key())
		
		if !ok {
			return nil, fmt.Errorf("key type %T does not fit map", key)
		}

		out := copyContainer(value, 0)
		out.SetMapIndex(k, reflect.Value{})

		return out.Interface(), nil
	}

	// case 2: slice
	i, ok := key.(int)
	if !ok || i < 0 || i >= value.Len() {
		return nil, fmt.Errorf("bad slice index %v", key)
	}

	// rebuild slice
	out := reflect.MakeSlice(value.Type(), 0, value.Len() - 1)
	out = reflect.AppendSlice(out, value.Slice(0, i))
	out = reflect.AppendSlice(out, value.Slice(i + 1, value.Len()))

	return out.Interface(), nil
}

// Set a container entry value with the given key.
//
// @param container map|slice
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

func set(_ *gplt.Runtime, _ *gplt.Context, container, key, value any) (any, error) {
	v, err := ensureContainer(container)
	if err != nil {
		return nil, err
	}

	val, ok := convertTo(value, v.Type().Elem())
	if !ok {
		return nil, fmt.Errorf("value type %T does not fit container", value)
	}

	// case 1: map
	if v.Kind() == reflect.Map {
		k, ok := convertTo(key, v.Type().Key())
		if !ok {
			return nil, fmt.Errorf("key type %T does not fit map", key)
		}

		out := copyContainer(v, 0)
		out.SetMapIndex(k, val)

		return out.Interface(), nil
	}

	// case 2: slice
	i, ok := key.(int)
	if !ok || i < 0 {
		return nil, fmt.Errorf("bad slice index %v", key)
	}

	// increase size
	n := max(v.Len(), i + 1)
	out := reflect.MakeSlice(v.Type(), n, n)

	// ... and insert new value
	reflect.Copy(out, v)
	out.Index(i).Set(val)

	return out.Interface(), nil
}

// Returns if an entry with the given key exists in container.
//
// @param container map|slice
// @param key any
// @returns bool
// @returns error
// 
// @example
//	 +{{ if (has (mapCreate "version" "v1") "version") }}
//		Container has "version"!
//	 +{{ end }}
var containerHasFunc = modules.NewFunc("has", has)

func has(_ *gplt.Runtime, _ *gplt.Context, container, key any) (bool, error) {
	value, err := ensureContainer(container)
	if err != nil {
		return false, err
	}

	if value.Kind() == reflect.Map {
		k, ok := convertTo(key, value.Type().Key())

		return ok && value.MapIndex(k).IsValid(), nil
	}

	i, ok := key.(int)

	return ok && i >= 0 && i < value.Len(), nil
}

// Returns if an entry with the given value exists in container.
//
// @param container map|slice
// @param value any
// @returns bool
// @returns error
// 
// @example
//	 +{{ if (includes (mapCreate "version" "v1") "v1") }}
//		Container includes "v1"!
//	 +{{ end }}
var containerIncludesFunc = modules.NewFunc("includes", includes)

func includes(_ *gplt.Runtime, _ *gplt.Context, container, value any) (bool, error) {
	v, err := ensureContainer(container)
	if err != nil {
		return false, err
	}

	// case 1: map
	if v.Kind() == reflect.Map {
		for it := v.MapRange(); it.Next(); {
			if equal(it.Value().Interface(), value) {
				return true, nil
			}
		}
		return false, nil
	}

	// case 2: slice
	for i := 0; i < v.Len(); i++ {
		if equal(v.Index(i).Interface(), value) {
			return true, nil
		}
	}
	return false, nil
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

func slicePush(_ *gplt.Runtime, _ *gplt.Context, container any, value any) (any, error) {
	// get element type
	v := deref(container)

	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("not a slice, got %T", container)
	}

	val, ok := convertTo(value, v.Type().Elem())
	if !ok {
		return nil, fmt.Errorf("value type %T does not fit slice", value)
	}

	// copy and increase size
	out := copyContainer(v, 1)
	return reflect.Append(out, val).Interface(), nil
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

func sliceCreate(_ *gplt.Runtime, _ *gplt.Context, values ...any) []any {
	return values
}

// Creates a new map, optionally with key value pairs.
// Pairs are constructed one by one, a key followed by the given value.
// Currently only supports string keys!
//
// @param entries []any
// @returns map[string]any
// @returns errors
// 
// @example
//	 +{{ $newMap := mapCreate "key1" "value1" "key2" "value2" }}
//	 +{{ echo $newMap }}
// @output
//	 map[key1:value1 key2:value2]
var mapCreateFunc = modules.NewFunc("mapCreate", mapCreate)

func mapCreate(_ *gplt.Runtime, _ *gplt.Context, entries ...any) (map[string]any, error) {
	if len(entries) == 0 {
		return map[string]any{}, nil
	}

	if len(entries) % 2 != 0 {
		return nil, fmt.Errorf("missing value for key %q", entries[len(entries) - 1])
	}

	out := map[string]any{}

	for i := 1; i < len(entries); i += 2 {
		key, ok := entries[i - 1].(string)

		if !ok {
			return nil, fmt.Errorf("item at index %d is not of type string, map keys must be strings", i)
		}

		out[key] = entries[i]
	}

	return out, nil
}