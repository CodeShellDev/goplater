package gplt

import (
	"fmt"
	"strings"
	"text/template"
)

// walks all defines already parsed, searches for "name(params)" pattern, re-registers that tree under plan name.
// if prefix is not empty, every define is also namespaced
func ExtractSignatures(tmpl *template.Template, prefix string) (map[string][]string, error) {
	signatures := map[string][]string{}

	for _, t := range tmpl.Templates() {
		if t.Name() == tmpl.Name() {
			continue
		}

		base, params, hasSig := parseSignature(t.Name())
		if !hasSig {
			base = t.Name()
		}

		// always register the clean bare name (signature-stripped, for self-reference)
		if base != t.Name() {
			_, err := tmpl.AddParseTree(base, t.Tree)
			if err != nil {
				return nil, fmt.Errorf("define %q: %w", t.Name(), err)
			}
		}

		if hasSig {
			signatures[base] = params
		}

		// additionally register the prefixed name (for external callers)
		if prefix != "" {
			prefixed := prefix + "." + base

			_, err := tmpl.AddParseTree(prefixed, t.Tree)
			if err != nil {
				return nil, fmt.Errorf("define %q: %w", t.Name(), err)
			}

			if hasSig {
				signatures[prefixed] = params
			}
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
	return base, splitParams(name[open + 1:len(name) - 1]), true
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