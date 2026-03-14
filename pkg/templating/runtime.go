package templating

import (
	"bytes"
	"errors"
	"text/template"

	"github.com/codeshelldev/goplater/utils/funcutils"
)

type Runtimes struct {
	lookup map[uint64]*Runtime
	nextID uint64
}

func (rts *Runtimes) Lookup(id uint64) *Runtime {
	return rts.lookup[id]
}

func (rts *Runtimes) Register(rt *Runtime) uint64 {
   	id := rts.nextID
    rts.nextID++
    rts.lookup[id] = rt

	return id
}

func (rts *Runtimes) Unregister(id uint64) {
    delete(rts.lookup, id)
}

type Runtime struct{
	engine 		*Engine
	engineOptions EngineOptions
	store 		map[string]StoreContainer
}

type StoreContainer interface {
	Set(key string, value any)
	Get(key string) any
	Delete(key string) bool
	Has(key string) bool
	Keys() []string
}

func (rt *Runtime) GetStore(id string) StoreContainer {
	return rt.store[id]
}

func (rt *Runtime) HasStore(id string) bool {
	_, exists := rt.store[id]
	return exists
}

func (rt *Runtime) RegisterStore(id string, store StoreContainer) error {
	_, exists := rt.store[id]

	if exists {
		return errors.New("store with " + id + " already registered")
	}

	rt.store[id] = store

	return nil
}

func (rt *Runtime) UnegisterStore(id string) error {
	_, exists := rt.store[id]

	if !exists {
		return errors.New("no store with " + id + " found")
	}

	delete(rt.store, id)

	return nil
}

func (rt *Runtime) GetEngine() *Engine {
	return rt.engine
}

func (rt *Runtime) GetEngineOptions() EngineOptions {
	return rt.engineOptions
}

func (rt *Runtime) FuncMap(context any) template.FuncMap {
	m := template.FuncMap{}

	for _, mod := range rt.engine.modules {
		for _, f := range mod.GetFuncMap() {
			m[f.Name] = funcutils.BindContext(f.Handler, rt, context)
		}
	}

	return m
}

func (rt *Runtime) Render(name, body string, data any, delims Delims, context Context) (string, error) {
	t := rt.engine.NewTemplate(name, delims)

	context.Set(InputContextKey, TemplateInputContext{
		Data: data,
		Body: body,
	})

	t.Funcs(rt.FuncMap(context))

	_, err := t.Parse(body)
	if err != nil {
		return "", err
	}

	rt.engine.template = t

	var buf bytes.Buffer
	err = t.Execute(&buf, data)

	return buf.String(), err
}