package funcs

import (
	"bytes"
	"errors"
	"maps"
	"text/template"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/gotl/pkg/templating"
	"github.com/google/uuid"
)

var functionFuncs = map[string]any{}

const functionOutputsKey = "out"

func setLocal(callerID string, key string, value any) {
	_, exists := runtime.locals[callerID]

	if !exists {
		runtime.locals[callerID] = map[string]any{}
	}

	runtime.locals[callerID][key] = value
}

func getLocal(callerID string, key string) any {
	scope, exists := runtime.locals[callerID]

	if exists {
		return scope[key]
	}

	return nil
}

func RegisterFunction(f TemplateFunc) {
	functionFuncs[f.Name] = f.Handler
}

func GetFunctionFuncMap(context context.TemplateContext, callerID string) map[string]any{
	m := make(map[string]any, len(functionFuncs))

	for k, v := range functionFuncs {
		m[k] = bindContext(v, context, callerID)
	}

	return m
}

var funcDefineFunc = TemplateFunc{
	Name: "funcDefine",
	Handler: func(context context.TemplateContext, name string, tmplBody string) any {
		funcCreate(context, name, tmplBody)
		return ""
	},
}

var funcCallFunc = TemplateFunc{
	Name: "funcCall",
	Handler: func(context context.TemplateContext, name string) []any {
		outputs, err := funcCall(context, name, nil)

		if err != nil {
			panic("could not call func: " + err.Error())
		}

		return outputs
	},
}

var funcCallWithArgsFunc = TemplateFunc{
	Name: "funcCallArgs",
	Handler: func(context context.TemplateContext, name string, args ...any) []any {
		args = unpackArgs(args...)

		outputs, err := funcCall(context, name, args...)

		if err != nil {
			panic("could not call func: " + err.Error())
		}

		return outputs
	},
}

func unpackArgs(args ...any) []any {
	if len(args) == 1 {
		inner, ok := args[0].([]any)

		if ok {
			args = inner
		}
	}

	return args
}

func funcCreate(_ context.TemplateContext, name, tmplBody string) any {
    runtime.funcs[name] = tmplBody

    return ""
}

func funcCall(context context.TemplateContext, name string, args ...any) ([]any, error) {
	tmplBody, ok := runtime.funcs[name]
	
	if !ok {
		return nil, errors.New("function \"" + name + "\" not defined")
	}

	var buf bytes.Buffer

	data := map[string]any{
		"args": args,
	}

	callerID := uuid.NewString()

	templt := template.New(context.Path + ":" + "func:" + name)
	templt.Delims("{{{", "}}}")

	funcMap := GetFuncMap(context)
	addMap := GetFunctionFuncMap(context, callerID)

	maps.Copy(funcMap, addMap)

	templt.Funcs(funcMap)

	err := templating.ParseTemplate(templt, tmplBody)

	if err != nil {
		return nil, err
	}

	err = templt.Execute(&buf, data)

	if err != nil {
		return nil, err
	}

	out, ok := getLocal(callerID, functionOutputsKey).([]any)

	if ok {
		return out, nil
	}

	return nil, nil
}

var returnFunc = TemplateFunc{
	Name: "return",
	Handler: func(context context.TemplateContext, callerID string, i int, value any) any {
		out, ok := getLocal(callerID, functionOutputsKey).([]any)

		if !ok {
			out = []any{}
		}

		for len(out) <= i {
			out = append(out, nil)
		}

		out[i] = value

		setLocal(callerID, functionOutputsKey, out)
		return ""
	},
}

var returnNextFunc = TemplateFunc{
	Name: "returnNext",
	Handler: func(context context.TemplateContext, callerID string, value any) any {
		out, ok := getLocal(callerID, functionOutputsKey).([]any)

		if !ok {
			out = []any{}
		}

		out = append(out, value)

		setLocal(callerID, functionOutputsKey, out)
		return ""
	},
}

var returnAllFunc = TemplateFunc{
	Name: "returnAll",
	Handler: func(context context.TemplateContext, callerID string, values ...any) any {
		values = unpackArgs(values...)

		setLocal(callerID, functionOutputsKey, values)
		return ""
	},
}

var returnOutputsFunc = TemplateFunc{
	Name: "returnOutputs",
	Handler: func(context context.TemplateContext, callerID string, value []any) any {
		setLocal(callerID, functionOutputsKey, value[0])

		return ""
	},
}

var outputsGetFunc = TemplateFunc{
	Name: "outputsGet",
	Handler: func(context context.TemplateContext, callerID string) []any {
		out, ok := getLocal(callerID, functionOutputsKey).([]any)

		if !ok {
			out = []any{}
		}

		return out
	},
}

func init() {
	Register(funcDefineFunc)
	Register(funcCallWithArgsFunc)
	Register(funcCallFunc)

	RegisterFunction(returnFunc)
	RegisterFunction(returnNextFunc)
	RegisterFunction(returnAllFunc)
	RegisterFunction(returnOutputsFunc)

	RegisterFunction(outputsGetFunc)
}