package tests

import (
	"testing"
)

func TestDeleteSlice(t *testing.T) {
	tmpl := `
		{{ $slice := list "Apple" "Banana" "Orange" }}
		{{ $slice = delete $slice 1 }}
		{{ emit "result" $slice }}
	`

	trt := tester.Collect(t, tmpl)
	trt.Equal(t, "result", []any{"Apple", "Orange"})
}

func TestDeleteSliceFirstAndLast(t *testing.T) {
	tmpl := `
		{{ $first := list "Apple" "Banana" "Orange" }}
		{{ $first = delete $first 0 }}
		{{ emit "first" $first }}

		{{ $last := list "Apple" "Banana" "Orange" }}
		{{ $last = delete $last 2 }}
		{{ emit "last" $last }}
	`

	trt := tester.Collect(t, tmpl)
	trt.Equal(t, "first", []any{"Banana", "Orange"})
	trt.Equal(t, "last", []any{"Apple", "Banana"})
}

func TestDeleteMap(t *testing.T) {
	tmpl := `
		{{ $map := dict "name" "Alice" "age" 30 }}
		{{ delete $map "name" }}
		{{ emit "result" $map }}
		{{ emit "has" (hasKey $map "name") }}
	`

	trt := tester.Collect(t, tmpl)
	trt.Equal(t, "has", false)
	trt.Equal(t, "result", map[string]any{"age": 30})
}

func TestDeleteErrors(t *testing.T) {
	tests := []string{
		`{{ delete (list "Apple" "Banana") -1 }}`,
		`{{ delete (list "Apple" "Banana") 2 }}`,
		`{{ delete (list "Apple") "0" }}`,
		`{{ delete "not a container" 0 }}`,
		`{{ delete (dict "name" "Alice") 123 }}`,
	}

	for _, tmpl := range tests {
		tester.ExpectErr(t, tmpl, nil)
	}
}

func TestSetSlice(t *testing.T) {
	tmpl := `
		{{ $slice := list "Apple" "Banana" }}
		{{ set $slice 0 "Strawberry" }}
		{{ emit "result" $slice }}
	`

	trt := tester.Collect(t, tmpl)
	trt.Equal(t, "result", []any{"Strawberry", "Banana"})
}

func TestSetMap(t *testing.T) {
	tmpl := `
		{{ $map := dict "name" "Alice" }}
		{{ set $map "name" "Bob" }}
		{{ set $map "age" 30 }}
		{{ emit "result" $map }}
	`

	trt := tester.Collect(t, tmpl)
	trt.Equal(t, "result", map[string]any{
		"name": "Bob",
		"age":	30,
	})
}

func TestSetErrors(t *testing.T) {
	tests := []string{
		`{{ set (list "Apple") -1 "Orange" }}`,
		`{{ set (list "Apple") 1 "Orange" }}`,
		`{{ set (list "Apple") "0" "Orange" }}`,
		`{{ set (dict "name" "Alice") 123 "Bob" }}`,
		`{{ set "not a container" 0 "value" }}`,
	}

	for _, tmpl := range tests {
		tester.ExpectErr(t, tmpl, nil)
	}
}

func TestHasKeyMap(t *testing.T) {
	tmpl := `
		{{ $map := dict "name" "Alice" }}
		{{ emit "existing" (hasKey $map "name") }}
		{{ emit "missing" (hasKey $map "age") }}
	`

	trt := tester.Collect(t, tmpl)
	trt.Equal(t, "existing", true)
	trt.Equal(t, "missing", false)
}

func TestHasKeySlice(t *testing.T) {
	tmpl := `
		{{ $slice := list "Apple" "Banana" }}
		{{ emit "first" (hasKey $slice 0) }}
		{{ emit "last" (hasKey $slice 1) }}
		{{ emit "missing" (hasKey $slice 2) }}
		{{ emit "negative" (hasKey $slice -1) }}
	`

	trt := tester.Collect(t, tmpl)
	trt.Equal(t, "first", true)
	trt.Equal(t, "last", true)
	trt.Equal(t, "missing", false)
	trt.Equal(t, "negative", false)
}

func TestHasKeyErrors(t *testing.T) {
	tests := []string{
		`{{ hasKey "not a container" "key" }}`,
	}

	for _, tmpl := range tests {
		tester.ExpectErr(t, tmpl, nil)
	}
}

func TestIncludesSlice(t *testing.T) {
	tmpl := `
		{{ $slice := list "Apple" "Banana" "Orange" }}
		{{ emit "existing" (includes $slice "Banana") }}
		{{ emit "missing" (includes $slice "Strawberry") }}
	`

	trt := tester.Collect(t, tmpl)
	trt.Equal(t, "existing", true)
	trt.Equal(t, "missing", false)
}

func TestIncludesMapValues(t *testing.T) {
	tmpl := `
		{{ $map := dict "fruit" "Apple" "color" "Red" }}
		{{ emit "existing" (includes $map "Apple") }}
		{{ emit "missing" (includes $map "Banana") }}
	`

	trt := tester.Collect(t, tmpl)
	trt.Equal(t, "existing", true)
	trt.Equal(t, "missing", false)
}

func TestIncludesErrors(t *testing.T) {
	tmpl := `{{ includes "not a container" "value" }}`
	tester.ExpectErr(t, tmpl, nil)
}

func TestAppend(t *testing.T) {
	tmpl := `
		{{ $slice := list "Apple" "Banana" }}
		{{ $newSlice := append $slice "Orange" }}
		{{ emit "original" $slice }}
		{{ emit "result" $newSlice }}
	`

	trt := tester.Collect(t, tmpl)
	trt.Equal(t, "original", []any{"Apple", "Banana"})
	trt.Equal(t, "result", []any{"Apple", "Banana", "Orange"})
}

func TestAppendErrors(t *testing.T) {
	tests := []string{
		`{{ append "not a slice" "value" }}`,
	}

	for _, tmpl := range tests {
		tester.ExpectErr(t, tmpl, nil)
	}
}

func TestList(t *testing.T) {
	tmpl := `
		{{ emit "empty" (list) }}
		{{ emit "values" (list "Apple" "Banana" 123) }}
	`

	trt := tester.Collect(t, tmpl)
	trt.Equal(t, "empty", []any{})
	trt.Equal(t, "values", []any{"Apple", "Banana", 123})
}

func TestDict(t *testing.T) {
	tmpl := `
		{{ emit "empty" (dict) }}
		{{ emit "values" (dict "name" "Alice" "age" 30) }}
	`

	trt := tester.Collect(t, tmpl)
	trt.Equal(t, "empty", map[string]any{})
	trt.Equal(t, "values", map[string]any{
		"name": "Alice",
		"age":	30,
	})
}

func TestDictErrors(t *testing.T) {
	tests := []string{
		`{{ dict "name" }}`,
		`{{ dict 123 "Alice" }}`,
	}

	for _, tmpl := range tests {
		tester.ExpectErr(t, tmpl, nil)
	}
}
