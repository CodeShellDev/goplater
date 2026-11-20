package gplt

import (
	"strings"
	"text/template"
)

type treeEntry struct {
	body		string
	template	*template.Template
	trusted		bool
	alias		string
}

type rootEntry struct {
	body		string
	template	*template.Template
}

const CurrentTreeKey ContextKey = "currentTree"

type TreeScope struct {
	Name	string
	Trusted bool
}

type ModuleRegistry struct {
	root	  rootEntry
	trees	  map[string]*treeEntry		 // import name -> tree
	reachable map[string]map[string]bool // owner scope -> names it is allowed to call into
}

func NewModuleRegistry(root *template.Template, body string) *ModuleRegistry {
	return &ModuleRegistry{
		root:	   rootEntry{template: root, body: body},
		trees:	   map[string]*treeEntry{},
		reachable: map[string]map[string]bool{"root": {}},
	}
}

func (r *ModuleRegistry) Loaded(name string) bool {
	_, ok := r.trees[name]
	return ok
}

func (r *ModuleRegistry) Register(name string, tmpl *template.Template, body string, trusted bool, alias string) {
	r.trees[name] = &treeEntry{template: tmpl, body: body, trusted: trusted, alias: alias}
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
func (r *ModuleRegistry) Lookup(callerScope, name string) (*template.Template, string) {
	if callerScope == "root" {
		t := r.root.template.Lookup(name)
		if t != nil {
			return t, r.root.body
		}
	}
	
	if callerScope != "root" {
		entry, exists := r.trees[callerScope]

		if exists {
			t := entry.template.Lookup(name)
			if t != nil {
				return t, entry.body
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

		t := entry.template.Lookup(name)
		if t != nil {
			return t, entry.body
		}
	}

	return nil, ""
}

// looks up the tree scope of a template
func (r *ModuleRegistry) ScopeOf(tmpl *template.Template) TreeScope {
	if r.root.template.Lookup(tmpl.Name()) == tmpl {
		return TreeScope{Name: "root", Trusted: true}
	}

	for path, e := range r.trees {
		if e.template.Lookup(tmpl.Name()) == tmpl {
			return TreeScope{Name: path, Trusted: e.trusted}
		}
	}

	return TreeScope{Name: "unknown", Trusted: false}
}