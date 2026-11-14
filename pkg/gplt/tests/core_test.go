package tests

import (
	"testing"

	"github.com/codeshelldev/goplater/pkg/gplt/modules/core"
	"github.com/codeshelldev/goplater/pkg/gplt/utils/testutils"
)

func TestGlobals(t *testing.T) {
	tmpl := `
	{{ globalSet "abc" 123 }}
	{{ emit "has" (globalHas "abc") }}
	{{ emit "get" (globalGet "abc") }}
	`

	trt, err := testutils.Collect(t, tmpl)

	if err != nil {
		t.Fatalf("unexpected error: %v\ntemplate: %s", err, tmpl)
	}

	trt.Equal(t, "has", true)
	trt.Equal(t, "get", 123)
}

func TestFuncs(t *testing.T) {
	tmpl := `
		{{ define "simple(a, b)" }}
			{{ return ( printf "%d, %d" .a .b ) }}
		{{ end }}

		{{ define "output(a, b)" }}
			{{- .a }}, {{ .b -}}
		{{ end }}

		{{ define "multiple(a, b)" }}
			{{ return .a .b }}
		{{ end }}
		
		{{ define "nested(a, b)" }}
			{{ return (call "simple" .a .b) }}
		{{ end }}

		{{ emit "simple" (call "simple" 1 2) }}
		{{ emit "output" (call "output" 1 2) }}
		{{ emit "multiple" (call "multiple" 1 2) }}
		{{ emit "nested" (call "nested" 1 2) }}
	`

	trt, err := testutils.Collect(t, tmpl)

	if err != nil {
		t.Fatalf("unexpected error: %v\ntemplate: %s", err, tmpl)
	}

	trt.Equal(t, "simple", "1, 2")
	trt.Equal(t, "output", "1, 2")
	trt.Equal(t, "nested", "1, 2")
	trt.Equal(t, "multiple", []any{1, 2})
}

func TestReturns(t *testing.T) {
	tmpl := `
		{{ define "simple" }}
			{{ return "doing simple stuff" }}
		{{ end }}

		{{ define "earlyExit" }}
			{{ return "doing simple stuff" }}
			{{ emit "exit" "" }}
		{{ end }}

		{{ emit "simple" (call "simple") }}
		{{ emit "earlyExit" (call "earlyExit") }}
	`

	trt, err := testutils.Collect(t, tmpl)

	if err != nil {
		t.Fatalf("unexpected error: %v\ntemplate: %s", err, tmpl)
	}

	trt.Equal(t, "simple", "doing simple stuff")
	trt.Equal(t, "earlyExit", "doing simple stuff")
	trt.Unexpect(t, "exit")

	tmpl = `{{ return "naughty return" }}`

	_, err = testutils.ExpectErr(t, tmpl, core.OutsideFuncErr)

	if err != nil {
		t.Fatalf("unexpected error: %v\ntemplate: %s", err, tmpl)
	}
}