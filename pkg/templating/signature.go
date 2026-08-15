package templating

import (
	"fmt"
	"strings"
	"text/template"
)

// ExtractSignatures walks all defines already parsed onto t and, for any whose
// name looks like "name(params)", re-registers that tree under the plain name.
// If prefix is non-empty, every define (signature or not) is also namespaced.
func ExtractSignatures(t *template.Template, prefix string) (map[string][]string, error) {
	signatures := map[string][]string{}

	for _, tmpl := range t.Templates() {
		if tmpl.Name() == t.Name() {
			continue // this is the top-level body itself, not a define
		}

		base, params, hasSig := parseSignature(tmpl.Name())

		fullName := tmpl.Name()
		if hasSig {
			fullName = base
		}
		if prefix != "" {
			fullName = prefix + "." + fullName
		}

		if fullName != tmpl.Name() {
			if _, err := t.AddParseTree(fullName, tmpl.Tree); err != nil {
				return nil, fmt.Errorf("define %q: %w", tmpl.Name(), err)
			}
		}

		if hasSig {
			signatures[fullName] = params
		}
	}

	return signatures, nil
}

func parseSignature(name string) (base string, params []string, ok bool) {
	open := strings.Index(name, "(")
	if open == -1 || !strings.HasSuffix(name, ")") {
		return "", nil, false
	}
	base = name[:open]
	return base, splitParams(name[open+1 : len(name)-1]), true
}

func splitParams(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	params := make([]string, 0, len(parts))
	for _, p := range parts {
		params = append(params, strings.TrimSpace(p))
	}
	return params
}