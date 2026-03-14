package templating

import (
	"text/template"

	"github.com/codeshelldev/goplater/pkg/templating/modules"
)

type Engine struct {
	modules []modules.Module
	template *template.Template
}

type EngineOptions struct {
	Delims 		Delims
	FuncDelims 	Delims
}

type Delims struct {
	Left string
	Right string
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Use(m modules.Module) {
	e.modules = append(e.modules, m)
}

func (e *Engine) UseModules(m ...modules.Module) {
	e.modules = append(e.modules, m...)
}

func (e *Engine) GetModules() []modules.Module {
	return modules.UniqueModules(e.modules)
}

func (e *Engine) NewTemplate(name string, delims Delims) *template.Template {
	if e.template != nil {
		name = e.template.Name() + name
	}

	t := template.New(name)
	t.Delims(delims.Left, delims.Right)

	return t
}

func (e *Engine) Execute(name, body string, data any, options EngineOptions, context Context) (string, error) {
	rt := &Runtime{
		engine: e,
		engineOptions: options,
		store: map[string]StoreContainer{},
	}

	return rt.Render(name, body, data, options.Delims, context)
}

func (e *Engine) ExecuteWithRuntime(name, body string, data any, delims Delims, context Context, rt *Runtime) (string, error) {
	rt.engine = e

	return rt.Render(name, body, data, delims, context)
}