package templating

import (
	"bytes"
	"errors"
	"fmt"
	"maps"
	"strings"
	"text/template"
	"text/template/parse"

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

type Runtime struct {
	engine        *Engine
	engineOptions EngineOptions
	store         map[string]StoreContainer
	registry      *ModuleRegistry
	signatures    map[string][]string
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

func (rt *Runtime) Import(callerScope, path string, aliasOverride string, ctx *Context) error {
	if rt.registry.Loaded(path) {
		rt.registry.AllowReach(callerScope, path)
		return nil
	}

	resolved, err := rt.engine.resolver.Resolve(path)
	if err != nil {
		return errors.New("module " + path + ": " + err.Error())
	}

	name := resolved.Name
	if aliasOverride != "" {
		name = aliasOverride
	}
	
	if name == "" {
		return fmt.Errorf("module %q: resolver returned no alias and none was provided", name)
	}

	funcMap := rt.FuncMap(ctx, resolved.Trusted)

	t := rt.engine.NewTemplate(":import:" + path, rt.engineOptions.Delims)
	t.Funcs(funcMap)

	_, err = t.Parse(resolved.Template)
	if err != nil {
		return errors.New("module " + path + ": parse error: " + err.Error())
	}

	err = validateOnlyDefines(t)
	if err != nil {
		return errors.New("module " + path + ": " + err.Error())
	}

	sigs, err := ExtractSignatures(t, name)
	if err != nil {
		return errors.New("module " + path + ": " + err.Error())
	}

	rt.registry.Register(path, t, resolved.Trusted)
	rt.registry.AllowReach(callerScope, path)

	maps.Copy(rt.signatures, sigs)
	return nil
}

func validateOnlyDefines(t *template.Template) error {
	for _, tmpl := range t.Templates() {
		if tmpl.Name() != t.Name() {
			continue
		}

		err := onlyWhitespace(tmpl.Root)
		if err != nil {
			return errors.New("top-level content outside define blocks: " + err.Error())
		}
	}

	return nil
}

func onlyWhitespace(node *parse.ListNode) error {
	if node == nil {
		return nil
	}

	for _, n := range node.Nodes {
		tn, ok := n.(*parse.TextNode)

		if !ok {
			return errors.New("found " + n.String() + "node")
		}

		if strings.TrimSpace(string(tn.Text)) != "" {
			return errors.New("found non-whitespace text")
		}
	}

	return nil
}

func (rt *Runtime) GetEngine() *Engine {
	return rt.engine
}

func (rt *Runtime) GetRegistry() *ModuleRegistry {
	return rt.registry
}

func (rt *Runtime) GetEngineOptions() EngineOptions {
	return rt.engineOptions
}

func (rt *Runtime) FuncMap(ctx any, includeHost bool) template.FuncMap {
	m := template.FuncMap{}

	for _, mod := range rt.engine.GetSafeModules() {
		for _, f := range mod.GetFuncMap() {
			m[f.Name] = funcutils.BindContext(f.Handler, rt, ctx)
		}
	}

	if includeHost {
		for _, mod := range rt.engine.GetModules() {
			for _, f := range mod.GetFuncMap() {
				m[f.Name] = funcutils.BindContext(f.Handler, rt, ctx)
			}
		}
	}

	return m
}

func (rt *Runtime) Params(name string) ([]string, bool) { 
	p, ok := rt.signatures[name]
	
	return p, ok 
}

func (rt *Runtime) Render(name, body string, data any, ctx *Context) (string, error) {
	ctx.Set(InputContextKey, TemplateInputContext{Data: data, Body: body})
	ctx.Set(CurrentTreeKey, TreeScope{Name: "root", Trusted: true})

	t := rt.engine.NewTemplate(name, rt.engineOptions.Delims)
	t.Funcs(rt.FuncMap(ctx, true))

	_, err := t.Parse(body)
	if err != nil {
		return "", err
	}

	sigs, err := ExtractSignatures(t, "")
	if err != nil {
		return "", err
	}
	rt.signatures = sigs

	rt.registry = NewModuleRegistry(t)

	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	return buf.String(), err
}