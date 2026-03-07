package funcs

import (
	"reflect"

	"github.com/codeshelldev/goplater/internals/template/context"
)

var containerDeleteFunc = TemplateFunc{
	Name: "delete",
	Handler: func(context context.TemplateContext, container any, key any) any {
		return deleteKey(container, key)
	},
}

var containerSetFunc = TemplateFunc{
	Name: "set",
	Handler: func(context context.TemplateContext, container any, key any, value any) any {
		return setKey(container, key, value)
	},
}

var slicePushFunc = TemplateFunc{
	Name: "push",
	Handler: func(context context.TemplateContext, container []any, value any) any {
		i := len(container) + 1

		return setKey(container, i, value)
	},
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

func init() {
	Register(containerDeleteFunc)
	Register(containerSetFunc)
	
	Register(slicePushFunc)
}