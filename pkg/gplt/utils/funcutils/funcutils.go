package funcutils

import "reflect"

func BindContext(fn any, ctx ...any) any {
	v := reflect.ValueOf(fn)
	t := v.Type()

	if t.Kind() != reflect.Func {
		panic("bindContext: fn is not a function")
	}

	if t.NumIn() < len(ctx) {
		panic("bindContext: function must have at least as many parameters as context values")
	}

	ins := make([]reflect.Type, t.NumIn() - len(ctx))
	for i := range ins {
		ins[i] = t.In(i + len(ctx))
	}

	outs := make([]reflect.Type, t.NumOut())
	for i := range outs {
		outs[i] = t.Out(i)
	}

	newFuncType := reflect.FuncOf(ins, outs, t.IsVariadic())

	ctxValues := make([]reflect.Value, len(ctx))
	for i, c := range ctx {
		ctxValues[i] = reflect.ValueOf(c)
	}

	return reflect.MakeFunc(newFuncType, func(args []reflect.Value) []reflect.Value {
		if !t.IsVariadic() {
			allArgs := make([]reflect.Value, 0, len(ctxValues)+len(args))
			allArgs = append(allArgs, ctxValues...)
			allArgs = append(allArgs, args...)

			return v.Call(allArgs)
		}

		// for a variadic MakeFunc, args contains the variadic args as the final []T value
		allArgs := make([]reflect.Value, 0, len(ctxValues) + len(args))
		allArgs = append(allArgs, ctxValues...)
		allArgs = append(allArgs, args...)

		return v.CallSlice(allArgs)
	}).Interface()
}