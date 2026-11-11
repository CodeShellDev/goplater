package modules

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/codeshelldev/goplater/pkg/gplt/utils/funcutils"
)

type Module struct {
	funcs 	map[string]Func
	id 		string
}

func NewModule(funcs ...Func) Module {
	m := Module{
		funcs: map[string]Func{},
	}

	m.id = ModuleKey(m)
	
	for _, f := range funcs {
		m.funcs[f.Name] = f
	}

	return m
}

func (m Module) GetFuncMap() map[string]Func {
	return m.funcs
}

type Func struct {
	Name string
	handlers []*entry
}

type entry struct {
	handlerFunc reflect.Value
	typ      	reflect.Type
	inputs    	[]reflect.Type
	variadic 	bool
	outputLen 	int
	hasErr   	bool
}

var errType = reflect.TypeFor[error]()

func NewFunc(name string, handlers ...any) Func {
	fn := Func{
		Name:     name,
		handlers: make([]*entry, 0, len(handlers)),
	}

	for _, handler := range handlers {
		h, err := newEntry(handler)
		if err != nil {
			panic(fmt.Sprintf("%s: %v", name, err))
		}
		fn.handlers = append(fn.handlers, h)
	}

	return fn
}

func (fn Func) GetHandler(rt any, ctx any) (any, error) {
	bound := make([]*entry, 0, len(fn.handlers))
	for _, h := range fn.handlers {
		b, err := funcutils.BindContext(h.handlerFunc.Interface(), rt, ctx)
		if err != nil {
			return nil, err
		}

		e, err := newEntry(b)
		if err != nil {
			return nil, err
		}

		bound = append(bound, e)
	}

	newFn := Func{Name: fn.Name, handlers: bound}

	// raw func, without variants or variadics
	if len(bound) == 1 && !bound[0].variadic {
		return bound[0].handlerFunc.Interface(), nil
	}

	// wrap advanced func
	return func(args ...any) (any, error) { return newFn.Call(args) }, nil
}

func newEntry(handler any) (*entry, error) {
	value := reflect.ValueOf(handler)

	if !value.IsValid() {
		return nil, errors.New("got nil")
	}
	if value.Kind() != reflect.Func {
		return nil, fmt.Errorf("not func, got %T", handler)
	}

	valueType := value.Type()

	e := &entry{
		handlerFunc: value,
		typ: valueType,
		variadic: valueType.IsVariadic(),
		outputLen: valueType.NumOut(),
	}

	switch valueType.NumOut() {
	case 0, 1:
	case 2:
		if valueType.Out(1) != errType {
			return nil, errors.New("2nd return must be error")
		}

		e.hasErr = true
	default:
		return nil, errors.New("too many returns")
	}

	for i := 0; i < valueType.NumIn(); i++ {
		e.inputs = append(e.inputs, valueType.In(i))
	}

	return e, nil
}

func (fn Func) Call(args []any) (any, error) {
	var best *entry
	var bestInputs []reflect.Value
	
	// lower is better (best is 0)
	bestScore := -1

	for _, e := range fn.handlers {
		if !e.accepts(len(args)) {
			continue
		}

		score := 0

		// punish variadic funcs
		if e.variadic {
			score++
		}

		inputs := make([]reflect.Value, len(args))

		next := false
		for i, arg := range args {
			s, value, ok := fit(e.param(i), arg)

			if !ok {
				next = true
				break
			}

			score += s
			inputs[i] = value
		}

		if next {
			continue
		}

		if bestScore < 0 || score < bestScore {
			best = e
			bestInputs = inputs
			bestScore = score

			// perfect choice
			if score == 0 {
				break
			}
		}
	}

	if best == nil {
		return nil, fn.noMatch(args)
	}

	out := best.handlerFunc.Call(bestInputs)
	switch {
	case best.outputLen == 0:
		return "", nil
	case best.hasErr && !out[1].IsNil():
		return nil, out[1].Interface().(error)
	}

	return out[0].Interface(), nil
}

func (fn Func) noMatch(args []any) error {
	got := make([]string, len(args))

	for i, a := range args {
		got[i] = fmt.Sprintf("%T", a)
	}

	var available []string
	for _, e := range fn.handlers {
		available = append(available, e.typ.String())
	}

	return fmt.Errorf("%s(%s): no match; have:\n%s", fn.Name, strings.Join(got, ", "), strings.Join(available, "\n  "))
}

// return param at index (last one if variadic)
func (e *entry) param(i int) reflect.Type {
	n := len(e.inputs)

	if e.variadic && i >= n - 1 {
		return e.inputs[n - 1].Elem()
	}

	return e.inputs[i]
}

func (e *entry) accepts(argCount int) bool {
	// allow infinite args at end for variadics
	if e.variadic {
		return argCount >= len(e.inputs) - 1
	}
	return argCount == len(e.inputs)
}

// is type of int / uint (8 through 64) or float (32 / 64)
func isNum(k reflect.Kind) bool {
	return k >= reflect.Int && k <= reflect.Float64
}

// does arg fit paramType or not? (also scores result)
func fit(paramType reflect.Type, arg any) (int, reflect.Value, bool) {
	// handle emtpy args
	if arg == nil {
		switch paramType.Kind() {
		case reflect.Interface, reflect.Pointer, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan:
			return 1, reflect.Zero(paramType), true
		}
		return 0, reflect.Value{}, false
	}

	value := reflect.ValueOf(arg)
	valueType := value.Type()

	switch {
	// perfect match
	case valueType == paramType:
		return 0, value, true

	// can arg be used as paramType
	case valueType.AssignableTo(paramType):
		return 1, value, true

	// special case for numbers
	case isNum(valueType.Kind()) && isNum(paramType.Kind()):
		return 2, value.Convert(paramType), true
	}

	return 0, reflect.Value{}, false
}