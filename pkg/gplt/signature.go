package gplt

import (
	"fmt"
	"strings"
	"text/template"
)

type Param struct {
	Name		string
	Optional	bool	// "opts?"
	Variadic	bool	// "...parts"
}

type Signature struct {
	Name	string
	Params	[]Param
}

func (s Signature) Required() int {
	n := 0
	for _, p := range s.Params {
		if !p.Optional && !p.Variadic {
			n++
		}
	}
	return n
}

func (s Signature) Variadic() bool {
	return len(s.Params) > 0 && s.Params[len(s.Params) - 1].Variadic
}

// max positional args, -1 if variadic
func (s Signature) Max() int {
	if s.Variadic() {
		return -1
	}
	return len(s.Params)
}

// walks all defines already parsed, searches for "name(params)" pattern, re-registers that tree under plan name.
// if prefix is not empty, every define is also namespaced.
// only defines with a signature are returned (if not it is just a plain sub-template)
func ExtractSignatures(tmpl *template.Template, prefix string) (map[string]Signature, error) {
	signatures := map[string]Signature{}

	for _, t := range tmpl.Templates() {
		if t.Name() == tmpl.Name() {
			continue
		}

		base, params, ok, err := parseSignature(t.Name())
		if err != nil {
			return nil, fmt.Errorf("define %q: %w", t.Name(), err)
		}

		if !ok {
			base = t.Name()
		}

		// register under clean name (self-reference, no namespacing)
		if base != t.Name() {
			_, err := tmpl.AddParseTree(base, t.Tree)
			if err != nil {
				return nil, fmt.Errorf("define %q: %w", t.Name(), err)
			}
		}

		if ok {
			signatures[base] = Signature{Name: base, Params: params}
		}

		// additionally register namespaced name (external callers)
		if prefix != "" {
			prefixed := prefix + "." + base

			_, err := tmpl.AddParseTree(prefixed, t.Tree)
			if err != nil {
				return nil, fmt.Errorf("define %q: %w", t.Name(), err)
			}

			if ok {
				signatures[prefixed] = Signature{Name: prefixed, Params: params}
			}
		}
	}

	return signatures, nil
}

func parseSignature(name string) (base string, params []Param, ok bool, err error) {
	open := strings.Index(name, "(")
	if open == -1 || !strings.HasSuffix(name, ")") {
		return "", nil, false, nil
	}

	base = name[:open]
	params, err = splitParams(name[open + 1 : len(name) - 1])
	if err != nil {
		return "", nil, false, err
	}

	return base, params, true, nil
}

func splitParams(raw string) ([]Param, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}

	parts := strings.Split(raw, ",")
	params := make([]Param, 0, len(parts))
	seenOptional := false

	for i, p := range parts {
		p = strings.TrimSpace(p)
		param := Param{}

		// variadic
		if strings.HasPrefix(p, "...") {
			param.Variadic = true
			p = strings.TrimPrefix(p, "...")
			if i != len(parts) - 1 {
				return nil, fmt.Errorf("variadic param %q must be last", p)
			}
		}

		// optional
		if strings.HasSuffix(p, "?") {
			param.Optional = true
			p = strings.TrimSuffix(p, "?")
		}

		if p == "" {
			return nil, fmt.Errorf("empty param name")
		}

		if param.Variadic && param.Optional {
			return nil, fmt.Errorf("variadic param %q is already optional", p)
		}

		if param.Optional {
			seenOptional = true
		} else if seenOptional && !param.Variadic {
			return nil, fmt.Errorf("required param %q after optional", p)
		}

		param.Name = p
		params = append(params, param)
	}

	return params, nil
}