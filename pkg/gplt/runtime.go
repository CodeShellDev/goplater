package gplt

import (
	"bytes"
	"errors"
	"fmt"
	"maps"
	"strings"
	"text/template"
	"text/template/parse"

	"github.com/codeshelldev/goplater/pkg/gplt/utils/funcutils"
)

type Runtime struct {
	engine        *Engine
	engineOptions EngineOptions
	store         map[string]StoreContainer
	registry      *ModuleRegistry
	signatures    map[string][]string
}

type StoreContainer interface {
	Set(key string, value any)
	Get(key string) (any, error)
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
		return fmt.Errorf("store with %q already registered", id)
	}

	rt.store[id] = store

	return nil
}

func (rt *Runtime) UnegisterStore(id string) error {
	_, exists := rt.store[id]

	if !exists {
		return fmt.Errorf("no store with %q found", id)
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
		return fmt.Errorf("module %q: %w", path, err)
	}

	name := resolved.Name
	if aliasOverride != "" {
		name = aliasOverride
	}
	
	if name == "" {
		return fmt.Errorf("module %q: resolver returned no alias and none was provided", path)
	}

	funcMap, err := rt.FuncMap(ctx, resolved.Trusted)

	if err != nil {
		return err
	}

	t := rt.engine.NewTemplate(":import:" + path, rt.engineOptions.Delims)
	t.Funcs(funcMap)

	_, err = t.Parse(resolved.Template)
	if err != nil {
		return fmt.Errorf("module %q: parse error: %w", path, err)
	}

	err = validateOnlyDefines(t)
	if err != nil {
		return fmt.Errorf("module %q: %w", path, err)
	}

	sigs, err := ExtractSignatures(t, name)
	if err != nil {
		return fmt.Errorf("module %q: %w", path, err)
	}

	rt.registry.Register(path, t, resolved.Trusted, name)
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
			return fmt.Errorf("top-level content outside define blocks: %w", err)
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
			return fmt.Errorf("found %q node", node)
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

func (rt *Runtime) FuncMap(ctx any, includeUnsafe bool) (template.FuncMap, error) {
	m := template.FuncMap{}

	for _, mod := range rt.engine.GetSafeModules() {
		for _, f := range mod.GetFuncMap() {
			var err error

			m[f.Name], err = funcutils.BindContext(f.Handler, rt, ctx)

			if err != nil {
				return nil, err
			}
		}
	}

	if includeUnsafe {
		for _, mod := range rt.engine.GetModules() {
			for _, f := range mod.GetFuncMap() {
				var err error

				m[f.Name], err = funcutils.BindContext(f.Handler, rt, ctx)

				if err != nil {
					return nil, err
				}
			}
		}
	}

	return m, nil
}

func (rt *Runtime) Params(name string) ([]string, bool) { 
	p, ok := rt.signatures[name]
	
	return p, ok 
}

// Render runs gplt engine on a template body with optional data and context.
// Be sure to create a new runtime for every Render(), a runtime is expected to only be a single-use object
func (rt *Runtime) Render(name, body string, data any, ctx *Context) (string, error) {
	t, err := rt.PrepareRender(name, body, data, ctx)
	
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

func (rt *Runtime) PrepareRender(name, body string, data any, ctx *Context) (*template.Template, error) {
	ctx.Set(InputContextKey, TemplateInputContext{Data: data, Body: body})
	ctx.Set(CurrentTreeKey, TreeScope{Name: "root", Trusted: true})

	t := rt.engine.NewTemplate(name, rt.engineOptions.Delims)

	funcMap, err := rt.FuncMap(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("could no create func map: %w", err)
	}

	t.Funcs(funcMap)

	if rt.engineOptions.ParseFunc == nil {
		rt.engineOptions.ParseFunc = func(s string, t *template.Template) (error) {
			_, err := t.Parse(s)
			return err
		}
	}

	err = rt.engineOptions.ParseFunc(body, t)

	if err != nil {
		return nil, err
	}

	if rt.engineOptions.ProcessFunc != nil {
		err = rt.engineOptions.ProcessFunc(t)

		if err != nil {
			return nil, err
		}
	}

	return t, nil
}