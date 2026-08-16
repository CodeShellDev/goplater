package templating

import (
	"strings"
	"text/template"
)

type treeEntry struct {
	tmpl    *template.Template
	trusted bool
	alias   string
}

const CurrentTreeKey ContextKey = "currentTree"

type TreeScope struct {
	Name    string
	Trusted bool
}

type ModuleRegistry struct {
	root      *template.Template
	trees     map[string]*treeEntry      // import name -> tree
	reachable map[string]map[string]bool // owner scope -> names it is allowed to call into
}

func NewModuleRegistry(root *template.Template) *ModuleRegistry {
	return &ModuleRegistry{
		root:      root,
		trees:     map[string]*treeEntry{},
		reachable: map[string]map[string]bool{"root": {}},
	}
}

func (r *ModuleRegistry) Loaded(name string) bool {
	_, ok := r.trees[name]
	return ok
}

func (r *ModuleRegistry) Register(name string, tmpl *template.Template, trusted bool, alias string) {
	r.trees[name] = &treeEntry{tmpl: tmpl, trusted: trusted, alias: alias}
}

func (r *ModuleRegistry) Trusted(name string) bool {
	e, ok := r.trees[name]
	return ok && e.trusted
}

func (r *ModuleRegistry) AllowReach(callerScope, importName string) {
	if r.reachable[callerScope] == nil {
		r.reachable[callerScope] = map[string]bool{}
	}

	r.reachable[callerScope][importName] = true
}

func (r *ModuleRegistry) CanReach(callerScope, target string) bool {
	return r.reachable[callerScope][target]
}

// finds name, restriced by what callerScope is allowed to see
func (r *ModuleRegistry) Lookup(callerScope, name string) *template.Template {
	if callerScope == "root" {
		tmplt := r.root.Lookup(name)
		if tmplt != nil {
			return tmplt
		}
	}
	
	if callerScope != "root" {
		entry, exists := r.trees[callerScope]

		if exists {
			tmplt := entry.tmpl.Lookup(name)
			if tmplt != nil {
				return tmplt
			}
		}
	}

	for imported := range r.reachable[callerScope] {
		entry, exists := r.trees[imported]
		if !exists {
			continue
		}

		// external must use fq "alias.name" form
		// prevents bare-name collisions between unrelated imports
		if !strings.HasPrefix(name, entry.alias + ".") {
			continue
		}

		t := entry.tmpl.Lookup(name)
		if t != nil {
			return t
		}
	}

	return nil
}

// looks up the tree scope of a template
func (r *ModuleRegistry) ScopeOf(tmpl *template.Template) TreeScope {
	if r.root.Lookup(tmpl.Name()) == tmpl {
		return TreeScope{Name: "root", Trusted: true}
	}

	for path, e := range r.trees {
		if e.tmpl.Lookup(tmpl.Name()) == tmpl {
			return TreeScope{Name: path, Trusted: e.trusted}
		}
	}

	return TreeScope{Name: "unknown", Trusted: false}
}