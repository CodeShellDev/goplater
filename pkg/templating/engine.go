package templating

import (
	"text/template"

	"github.com/codeshelldev/goplater/pkg/templating/modules"
	"github.com/codeshelldev/goplater/pkg/templating/types"
)

type Engine struct {
	modules 	[]modules.Module
	safeModules []modules.Module // modules that are safe to be called from imported resources
	resolver 	*ResolverChain
}

type EngineOptions struct {
	Delims types.Delims
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

func (e *Engine) UseSafe(m modules.Module) { 
	e.safeModules = append(e.safeModules, m) 
}

func (e *Engine) UseSafeModules(m ...modules.Module) { 
	e.safeModules = append(e.safeModules, m...) 
}

func (e *Engine) GetSafeModules() []modules.Module { 
	return modules.UniqueModules(e.safeModules) 
}


func (e *Engine) SetResolver(r *ResolverChain) {
	e.resolver = r
}

func (e *Engine) NewTemplate(name string, delims types.Delims) *template.Template {
	t := template.New(name)

	t.Delims(delims.Left, delims.Right)

	return t
}

func (e *Engine) Execute(name, body string, data any, options EngineOptions, ctx *Context) (string, error) {
	rt := &Runtime{engine: e, engineOptions: options, store: map[string]StoreContainer{}}
	
	return rt.Render(name, body, data, ctx)
}