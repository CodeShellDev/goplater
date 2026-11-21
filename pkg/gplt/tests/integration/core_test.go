package tests

import (
	"testing"

	"github.com/codeshelldev/goplater/pkg/gplt/modules/core"
)

func TestGlobals(t *testing.T) {
	tmpl := `
	{{ globalSet "abc" 123 }}
	{{ emit "has" (globalHas "abc") }}
	{{ emit "get" (globalGet "abc") }}
	`

	trt := tester.Collect(t, tmpl)

	trt.Equal(t, "has", true)
	trt.Equal(t, "get", 123)
}

func TestFuncs(t *testing.T) {
	tmpl := `
		{{ define "empty()" }}
			{{ return 1 }}
		{{ end }}

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

		{{ emit "empty" (call "empty") }}
		{{ emit "simple" (call "simple" 1 2) }}
		{{ emit "output" (call "output" 1 2) }}
		{{ emit "multiple" (call "multiple" 1 2) }}
		{{ emit "nested" (call "nested" 1 2) }}
	`

	trt := tester.Collect(t, tmpl)

	trt.Equal(t, "empty", 1)
	trt.Equal(t, "simple", "1, 2")
	trt.Equal(t, "output", "1, 2")
	trt.Equal(t, "nested", "1, 2")
	trt.Equal(t, "multiple", []any{1, 2})
}

func TestOptionalFunc(t *testing.T) {
	tmpl := `
		{{ define "multiple(a, b?)" }}
			{{ return .a .b }}
		{{ end }}

		{{ emit "multipleA" (call "multiple" 1) }}
		{{ emit "multipleAB" (call "multiple" 1 2) }}
	`

	trt := tester.Collect(t, tmpl)

	trt.Equal(t, "multipleA", []any{1, nil})
	trt.Equal(t, "multipleAB", []any{1, 2})
}

func TestVariadicFunc(t *testing.T) {
	tmpl := `
		{{ define "simple(...a)" }}
			{{ return .a }}
		{{ end }}

		{{ emit "simpleEmpty" (call "simple") }}
		{{ emit "simple1" (call "simple" 1) }}
		{{ emit "simple2" (call "simple" 1 2) }}
	`

	trt := tester.Collect(t, tmpl)

	trt.Equal(t, "simpleEmpty", []any{})
	trt.Equal(t, "simple1", []any{1})
	trt.Equal(t, "simple2", []any{1, 2})
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

	trt := tester.Collect(t, tmpl)

	trt.Equal(t, "simple", "doing simple stuff")
	trt.Equal(t, "earlyExit", "doing simple stuff")
	trt.Unexpect(t, "exit")

	tmpl = `{{ return "naughty return" }}`

	tester.ExpectErr(t, tmpl, core.OutsideFuncErr)
}