package core

import (
	"errors"
	"reflect"

	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
)

// Returns the type of a value as a string (using [reflect.TypeOf](https://pkg.go.dev/reflect#TypeOf)).
//
// @param value any
// @returns string
// @returns error
var typeOfFunc = modules.NewFunc("typeOf", typeOf)

func typeOf(_ *gplt.Runtime, _ *gplt.Context, value any) (string, error) {
    if value == nil {
        return "", errors.New("cannot get type of nil")
    }

    return reflect.TypeOf(value).String(), nil
}

// Reports whether value is of type (using [reflect.TypeOf](https://pkg.go.dev/reflect#TypeOf)).
//
// @param value any
// @param typeName string
// @returns bool
var isTypeFunc = modules.NewFunc("isType", isType)

func isType(_ *gplt.Runtime, _ *gplt.Context, value any, typeName string) bool {
    if value == nil {
        return false
    }

    return reflect.TypeOf(value).String() == typeName
}