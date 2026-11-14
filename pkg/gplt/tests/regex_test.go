package tests

import (
	"testing"

	"github.com/codeshelldev/goplater/pkg/gplt/utils/testutils"
)

func TestRegexMatch(t *testing.T) {
	tmpl := `
		{{ emit "match" (regexMatch "^hello" "hello world") }}
		{{ emit "noMatch" (regexMatch "^world" "hello world") }}
	`

	trt, err := testutils.Collect(t, tmpl)
	if err != nil {
		t.Fatalf("unexpected error: %v\ntemplate: %s", err, tmpl)
	}

	trt.Equal(t, "match", true)
	trt.Equal(t, "noMatch", false)
}

func TestRegexFindAll(t *testing.T) {
	tmpl := `
		{{ emit "all" (regexFindAll "[a-z]+" "hello 123 world") }}
		{{ emit "limited" (regexFindAll "[a-z]+" "hello 123 world" 1) }}
	`

	trt, err := testutils.Collect(t, tmpl)
	if err != nil {
		t.Fatalf("unexpected error: %v\ntemplate: %s", err, tmpl)
	}

	trt.Equal(t, "all", []string{"hello", "world"})
	trt.Equal(t, "limited", []string{"hello"})
}

func TestRegexFindGroups(t *testing.T) {
	tmpl := `
		{{ emit "groups" (regexFindGroups "([a-z]+)([0-9]+)" "abc123 def456") }}
	`

	trt, err := testutils.Collect(t, tmpl)
	if err != nil {
		t.Fatalf("unexpected error: %v\ntemplate: %s", err, tmpl)
	}

	trt.Equal(t, "groups", [][]string{
		{"abc123", "abc", "123"},
		{"def456", "def", "456"},
	})
}

func TestRegexFindGroupsIndex(t *testing.T) {
	tmpl := `
	{{ emit "indices" (regexFindGroupsIndex "([a-z]+)([0-9]+)" "abc123") }}
	`

	trt, err := testutils.Collect(t, tmpl)
	if err != nil {
		t.Fatalf("unexpected error: %v\ntemplate: %s", err, tmpl)
	}

	trt.Equal(t, "indices", [][]int{
		{0, 6, 0, 3, 3, 6},
	})
}

func TestRegexReplaceAll(t *testing.T) {
	tmpl := `
	{{ emit "replace" (regexReplaceAll "hello 123 world 456" "[0-9]+" "NUM") }}
	{{ emit "groups" (regexReplaceAll "John Doe" "([A-Za-z]+) ([A-Za-z]+)" "$2, $1") }}
	`

	trt, err := testutils.Collect(t, tmpl)
	if err != nil {
		t.Fatalf("unexpected error: %v\ntemplate: %s", err, tmpl)
	}

	trt.Equal(t, "replace", "hello NUM world NUM")
	trt.Equal(t, "groups", "Doe, John")
}

func TestRegexErrors(t *testing.T) {
	tests := []struct {
		name string
		tmpl string
	}{
		{
			name: "match",
			tmpl: `{{ regexMatch "[" "hello" }}`,
		},
		{
			name: "findAll",
			tmpl: `{{ regexFindAll "[" "hello" }}`,
		},
		{
			name: "findGroups",
			tmpl: `{{ regexFindGroups "[" "hello" }}`,
		},
		{
			name: "findGroupsIndex",
			tmpl: `{{ regexFindGroupsIndex "[" "hello" }}`,
		},
		{
			name: "replaceAll",
			tmpl: `{{ regexReplaceAll "hello" "[" "world" }}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := testutils.ExpectErr(t, tt.tmpl, nil)
			if err != nil {
				t.Fatal(err.Error())
			}
		})
	}
}
