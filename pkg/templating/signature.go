package templating

import (
	"errors"
	"strings"
	"text/template"
)

// walks all defines already parsed, searches for "name(params)" pattern, re-registers that tree under plan name.
// if prefix is not empty, every define is also namespaced
func ExtractSignatures(t *template.Template, prefix string) (map[string][]string, error) {
	signatures := map[string][]string{}

	for _, tmpl := range t.Templates() {
		if tmpl.Name() == t.Name() {
			continue
		}

		base, params, hasSig := parseSignature(tmpl.Name())
		if !hasSig {
			base = tmpl.Name()
		}

		// always register the clean bare name (signature-stripped, for self-reference)
		if base != tmpl.Name() {
			_, err := t.AddParseTree(base, tmpl.Tree)
			if err != nil {
				return nil, errors.New("define " + tmpl.Name() + ": " + err.Error())
			}
		}

		if hasSig {
			signatures[base] = params
		}

		// additionally register the prefixed name (for external callers)
		if prefix != "" {
			prefixed := prefix + "." + base

			_, err := t.AddParseTree(prefixed, tmpl.Tree)
			if err != nil {
				return nil, errors.New("define " + tmpl.Name() + ": " + err.Error())
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