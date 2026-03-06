package funcs

import (
	"reflect"

	"github.com/codeshelldev/goplater/internals/template/context"
)

var deleteFunc = TemplateFunc{
	Name: "delete",
	Handler: func(context context.TemplateContext, container any, key any) any {
		return deleteKey(container, key)
	},
}

func deleteKey(data any, key any) any {
	val := reflect.ValueOf(data)

	switch val.Kind() {
	case reflect.Map:
		newMap := reflect.MakeMap(val.Type())

		for _, k := range val.MapKeys() {
			if !reflect.DeepEqual(k.Interface(), key) {
				newMap.SetMapIndex(k, val.MapIndex(k))
			}
		}

		return newMap.Interface()

	case reflect.Slice:
		newSlice := reflect.MakeSlice(val.Type(), 0, val.Len())

		for i := 0; i < val.Len(); i++ {
			if !reflect.DeepEqual(val.Index(i).Interface(), key) {
				newSlice = reflect.Append(newSlice, val.Index(i))
			}

		}
		return newSlice.Interface()

	default:
		return data
	}
}

func init() {
	Register(deleteFunc)
}