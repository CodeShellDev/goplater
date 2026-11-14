package tests

import (
	"testing"

	"github.com/codeshelldev/goplater/pkg/gplt/utils/testutils"
)

func TestJSONEncode(t *testing.T) {
	tmpl := `
		{{ emit "string" (jsonEncode "hello") }}
		{{ emit "number" (jsonEncode 123) }}
		{{ emit "bool" (jsonEncode true) }}
		{{ emit "array" (jsonEncode (list 1 2 3)) }}
		{{ emit "object" (jsonEncode (dict "name" "Alice" "age" 30)) }}
		{{ emit "nil" (jsonEncode nil) }}
	`

	trt := testutils.Collect(t, tmpl)
	trt.Equal(t, "string", `"hello"`)
	trt.Equal(t, "number", "123")
	trt.Equal(t, "bool", "true")
	trt.Equal(t, "array", `[1,2,3]`)
	trt.Equal(t, "object", `{"age":30,"name":"Alice"}`)
	trt.Equal(t, "nil", "null")
}

func TestJSONDecode(t *testing.T) {
	tmpl := `
		{{ emit "string" (jsonDecode "\"hello\"") }}
		{{ emit "number" (jsonDecode "123") }}
		{{ emit "float" (jsonDecode "12.5") }}
		{{ emit "bool" (jsonDecode "true") }}
		{{ emit "array" (jsonDecode "[1,2,3]") }}
		{{ emit "object" (jsonDecode "{\"name\":\"Alice\",\"age\":30}") }}
		{{ emit "null" (jsonDecode "null") }}
	`

	trt := testutils.Collect(t, tmpl)
	trt.Equal(t, "string", "hello")
	trt.Equal(t, "number", float64(123))
	trt.Equal(t, "float", float64(12.5))
	trt.Equal(t, "bool", true)
	trt.Equal(t, "array", []any{float64(1), float64(2), float64(3)})
	trt.Equal(t, "object", map[string]any{
		"name": "Alice",
		"age":  float64(30),
	})
	trt.Equal(t, "null", nil)
}

func TestJSONDecodeInvalidJSON(t *testing.T) {
	tmpl := `{{ jsonDecode "{invalid}" }}`

	testutils.ExpectErr(t, tmpl, nil)
}