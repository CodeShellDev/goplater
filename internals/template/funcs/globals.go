package funcs

import "github.com/codeshelldev/goplater/internals/template/context"

var globalSetFunc = TemplateFunc{
	Name: "globalSet",
	Handler: func(context context.TemplateContext, key string, value any) any {
		setGlobal(key, value)
		return nil
	},
}

var globalGetFunc = TemplateFunc{
	Name: "globalGet",
	Handler: func(context context.TemplateContext, key string) any {
		return getGlobal(key)
	},
}

func setGlobal(key string, value any) {
	runtime.globals[key] = value
}

func getGlobal(key string) any {
	return runtime.globals[key]
}

func init() {
	Register(globalSetFunc)
	Register(globalGetFunc)
}