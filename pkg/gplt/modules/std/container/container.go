package container

import (
	"fmt"
	"reflect"

	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
)

var Module = modules.NewModule(deleteFunc, setFunc, hasKeyFunc, includesFunc, appendFunc, listFunc, dictFunc)

// Deletes a container entry with the given key (mutates maps, but must return modified slice).
//
// @param container map|slice
// @param key any
// @returns any
// @returns error
// 
// @example
//	 +{{ $slice := list "Apple" "Banana" }}
//	 +{{ $newSlice := delete $slice 0 }}
//	 +{{ echo $newSlice }}
// @output
//	 [Banana]
// @example
//	 +{{ $map := dict "key1" "value1" "key2" "value2" }}
//	 +{{ delete $map "key1" }}
//	 +{{ echo $map }}
// @output
//	 map[key2:value2]
var deleteFunc = modules.NewFunc("delete", delete)

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

		value.SetMapIndex(k, reflect.Value{})
		return value.Interface(), nil
	}

	// case 2: slice
	i, ok := key.(int)
	if !ok || i < 0 || i >= value.Len() {
		return nil, fmt.Errorf("bad slice index %v", key)
	}

	last := value.Len() - 1

	// shift elements left
	reflect.Copy(value.Slice(i, last), value.Slice(i + 1, value.Len()))

	// clear the last element to release references
	value.Index(last).Set(reflect.Zero(value.Type().Elem()))

	// return shortened slice
	return value.Slice(0, last).Interface(), nil
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
//	 +{{ $slice := list "Apple" "Banana" }}
//	 +{{ set $slice 0 "Strawberry" }}
//	 +{{ echo $slice }}
// @output
//	 [Strawberry Banana]
var setFunc = modules.NewFunc("set", set)

func set(_ *gplt.Runtime, _ *gplt.Context, container, key, value any) (string, error) {
	v, err := ensureContainer(container)
	if err != nil {
		return "", err
	}

	// case 1: map
	if v.Kind() == reflect.Map {
		k, ok := convertTo(key, v.Type().Key())
		if !ok {
			return "", fmt.Errorf("key type %T does not fit map", key)
		}

		val, ok := convertTo(value, v.Type().Elem())
		if !ok {
			return "", fmt.Errorf("value type %T does not fit container", value)
		}

		v.SetMapIndex(k, val)
		return "", nil
	}

	// case 2: slice
	i, ok := key.(int)
	if !ok || i < 0 || i >= v.Len() {
		return "", fmt.Errorf("bad slice index %v", key)
	}

	val, ok := convertTo(value, v.Type().Elem())
	if !ok {
		return "", fmt.Errorf("value type %T does not fit container", value)
	}

	v.Index(i).Set(val)
	return "", nil
}


// Returns if an entry with the given key exists in container.
//
// @param container map|slice
// @param key any
// @returns bool
// @returns error
// 
// @example
//	 +{{ if (hasKey (list "version" "v1") "version") }}
//		Container has "version"!
//	 +{{ end }}
var hasKeyFunc = modules.NewFunc("hasKey", has)

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
//	 +{{ if (includes (list "version" "v1") "v1") }}
//		Container includes "v1"!
//	 +{{ end }}
var includesFunc = modules.NewFunc("includes", includes)

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

// Appends a value to a slice (returns a new copy).
//
// @param container []any
// @param value any
// @returns []any
// 
// @example
//	 +{{ $slice := list "Apple" "Banana" }}
//	 +{{ $newSlice = append $slice "Strawberry" }}
//	 +{{ echo $newSlice }}
// @output
//	 [Apple Banana Strawberry]
var appendFunc = modules.NewFunc("append", append)

func append(_ *gplt.Runtime, _ *gplt.Context, container any, value any) (any, error) {
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
//	 +{{ $slice := list "Apple" "Banana" }}
//	 +{{ echo $slice }}
// @output
//	 [Apple Banana]
var listFunc = modules.NewFunc("list", list)

func list(_ *gplt.Runtime, _ *gplt.Context, values ...any) []any {
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
//	 +{{ $map := dict "key1" "value1" "key2" "value2" }}
//	 +{{ echo $map }}
// @output
//	 map[key1:value1 key2:value2]
var dictFunc = modules.NewFunc("dict", dict)

func dict(_ *gplt.Runtime, _ *gplt.Context, entries ...any) (map[string]any, error) {
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