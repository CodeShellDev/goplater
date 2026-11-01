package transformer

import (
	"fmt"
	"text/template"
	"text/template/parse"
)

type Document struct {
	Name  string
	Trees map[string]*parse.Tree
}

type Processor interface {
	Name() string
	Process(doc *Document) error
}

type Transformer struct {
	processors []Processor
}

func New() *Transformer { return &Transformer{} }

func (p *Transformer) Use(processors ...Processor) *Transformer {
	for _, proc := range processors {
		if proc != nil {
			p.processors = append(p.processors, proc)
		}
	}

	return p
}

func (p *Transformer) Process(t *template.Template) error {
	if len(p.processors) == 0 {
		return nil
	}

	doc := newDocument(t)

	for _, proc := range p.processors {
		err := proc.Process(doc)
		if err != nil {
			return fmt.Errorf("processor %q: %w", proc.Name(), err)
		}
	}

	return nil
}

func newDocument(t *template.Template) *Document {
	doc := &Document{
		Name:  t.Name(),
		Trees: make(map[string]*parse.Tree),
	}

	for _, tmpl := range t.Templates() {
		if tmpl.Tree != nil && tmpl.Tree.Root != nil {
			doc.Trees[tmpl.Name()] = tmpl.Tree
		}
	}

	return doc
}