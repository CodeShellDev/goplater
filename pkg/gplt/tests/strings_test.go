package tests

import (
	"testing"

	"github.com/codeshelldev/goplater/pkg/gplt/utils/testutils"
)

func TestTrim(t *testing.T) {
	tmpl := `
		{{ emit "trim" (trim "--hello--" "-") }}
		{{ emit "trimSpace" (trimSpace "  hello	 ") }}
		{{ emit "trimPrefix" (trimPrefix "hello world" "hello ") }}
		{{ emit "trimSuffix" (trimSuffix "hello world" " world") }}
	`

	trt := testutils.Collect(t, tmpl)
	trt.Equal(t, "trim", "hello")
	trt.Equal(t, "trimSpace", "hello")
	trt.Equal(t, "trimPrefix", "world")
	trt.Equal(t, "trimSuffix", "hello")
}

func TestCase(t *testing.T) {
	tmpl := `
		{{ emit "upper" (upper "hello World") }}
		{{ emit "lower" (lower "HELLO WORLD") }}
	`

	trt := testutils.Collect(t, tmpl)
	trt.Equal(t, "upper", "HELLO WORLD")
	trt.Equal(t, "lower", "hello world")
}

func TestContains(t *testing.T) {
	tmpl := `
		{{ emit "contains" (contains "hello world" "world") }}
		{{ emit "count" (count "banana" "a") }}
		{{ emit "hasPrefix" (hasPrefix "hello world" "hello") }}
		{{ emit "hasSuffix" (hasSuffix "hello world" "world") }}
		{{ emit "isEmpty" (isEmpty "") }}
	`

	trt := testutils.Collect(t, tmpl)
	trt.Equal(t, "contains", true)
	trt.Equal(t, "count", 3)
	trt.Equal(t, "hasPrefix", true)
	trt.Equal(t, "hasSuffix", true)
	trt.Equal(t, "isEmpty", true)
}

func TestReplace(t *testing.T) {
	tmpl := `
		{{ emit "replace" (replace "hello world" "world" "Go") }}
	`

	trt := testutils.Collect(t, tmpl)
	trt.Equal(t, "replace", "hello Go")
}

func TestSplit(t *testing.T) {
	tmpl := `
		{{ emit "split" (split "one,two,three" ",") }}
	`

	trt := testutils.Collect(t, tmpl)
	trt.Equal(t, "split", []string{"one", "two", "three"})
}

func TestAfterBefore(t *testing.T) {
	tmpl := `
		{{ emit "after" (after "one/two/three" "/") }}
		{{ emit "afterLast" (afterLast "one/two/three" "/") }}
		{{ emit "before" (before "one/two/three" "/") }}
		{{ emit "beforeLast" (beforeLast "one/two/three" "/") }}
		{{ emit "between" (between "[hello]" "[" "]") }}
	`

	trt := testutils.Collect(t, tmpl)
	trt.Equal(t, "after", "two/three")
	trt.Equal(t, "afterLast", "three")
	trt.Equal(t, "before", "one")
	trt.Equal(t, "beforeLast", "one/two")
	trt.Equal(t, "between", "hello")
}

func TestSlice(t *testing.T) {
	tmpl := `
	{{ emit "substr" (substr "hello world" 0 5) }}
	`

	trt := testutils.Collect(t, tmpl)
	trt.Equal(t, "substr", "hello")
}

func TestJoin(t *testing.T) {
	tmpl := `
		{{ emit "joinStrings" (join "," (list "one" "two" "three")) }}
		{{ emit "joinAny" (join "-" (list "one" 2 true)) }}
		{{ emit "joinVariadic" (join "/" "one" 2 true) }}
	`

	trt := testutils.Collect(t, tmpl)
	trt.Equal(t, "joinStrings", "one,two,three")
	trt.Equal(t, "joinAny", "one-2-true")
	trt.Equal(t, "joinVariadic", "one/2/true")
}

func TestRepeat(t *testing.T) {
	tmpl := `
	{{ emit "repeat" (repeat "ha" 3) }}
	`

	trt := testutils.Collect(t, tmpl)
	trt.Equal(t, "repeat", "hahaha")
}

func TestConcat(t *testing.T) {
	tmpl := `
		{{ emit "concatStrings" (concat (list "one" "two" "three")) }}
		{{ emit "concatAny" (concat (list "one" 2 true)) }}
		{{ emit "concatVariadic" (concat "one" 2 true) }}
	`

	trt := testutils.Collect(t, tmpl)
	trt.Equal(t, "concatStrings", "onetwothree")
	trt.Equal(t, "concatAny", "one2true")
	trt.Equal(t, "concatVariadic", "one2true")
}

func TestIndexOf(t *testing.T) {
	tmpl := `
		{{ emit "found" (indexOf "hello world" "world") }}
		{{ emit "missing" (indexOf "hello world" "xyz") }}
	`

	trt := testutils.Collect(t, tmpl)
	trt.Equal(t, "found", 6)
	trt.Equal(t, "missing", -1)
}
