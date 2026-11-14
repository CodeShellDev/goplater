package tests

import (
	"testing"

	"github.com/codeshelldev/goplater/pkg/gplt/utils/testutils"
)

func TestAdd(t *testing.T) {
	tmpl := `
		{{ emit "int" (add 1 2 3 4) }}
		{{ emit "float" (add 1.5 2.5 3.0) }}
	`

	trt := testutils.Collect(t, tmpl)
	trt.Equal(t, "int", 10)
	trt.Equal(t, "float", 7.0)
}

func TestSub(t *testing.T) {
	tmpl := `
		{{ emit "int" (sub 20 5 3) }}
		{{ emit "float" (sub 10.5 2.5 1.0) }}
	`

	trt := testutils.Collect(t, tmpl)
	trt.Equal(t, "int", 12)
	trt.Equal(t, "float", 7.0)
}

func TestMul(t *testing.T) {
	tmpl := `
		{{ emit "int" (mul 2 3 4) }}
		{{ emit "float" (mul 1.5 2.0 3.0) }}
	`

	trt := testutils.Collect(t, tmpl)
	trt.Equal(t, "int", 24)
	trt.Equal(t, "float", 9.0)
}

func TestDiv(t *testing.T) {
	tmpl := `
		{{ emit "int" (div 20 2 2) }}
		{{ emit "float" (div 12.0 2.0 3.0) }}
	`

	trt := testutils.Collect(t, tmpl)
	trt.Equal(t, "int", 5)
	trt.Equal(t, "float", 2.0)
}

func TestMod(t *testing.T) {
	tmpl := `
		{{ emit "int" (mod 17 5) }}
		{{ emit "float" (mod 17.5 5.0) }}
	`

	trt := testutils.Collect(t, tmpl)
	trt.Equal(t, "int", 2)
	trt.Equal(t, "float", 2.5)
}

func TestDivByZero(t *testing.T) {
	tmpl := `{{ div 10 0 }}`

	testutils.ExpectErr(t, tmpl, nil)
}

func TestModByZero(t *testing.T) {
	tmpl := `{{ mod 10 0 }}`

	testutils.ExpectErr(t, tmpl, nil)
}

func TestAbs(t *testing.T) {
	tmpl := `
		{{ emit "negativeInt" (abs -10) }}
		{{ emit "positiveInt" (abs 10) }}
		{{ emit "negativeFloat" (abs -3.14) }}
		{{ emit "positiveFloat" (abs 3.14) }}
	`

	trt := testutils.Collect(t, tmpl)
	trt.Equal(t, "negativeInt", 10)
	trt.Equal(t, "positiveInt", 10)
	trt.Equal(t, "negativeFloat", 3.14)
	trt.Equal(t, "positiveFloat", 3.14)
}

func TestMin(t *testing.T) {
	tmpl := `
		{{ emit "int" (min 5 2 8 1) }}
		{{ emit "float" (min 5.5 2.2 8.8) }}
	`

	trt := testutils.Collect(t, tmpl)
	trt.Equal(t, "int", 1)
	trt.Equal(t, "float", 2.2)
}

func TestMax(t *testing.T) {
	tmpl := `
		{{ emit "int" (max 5 2 8 1) }}
		{{ emit "float" (max 5.5 2.2 8.8) }}
	`

	trt := testutils.Collect(t, tmpl)
	trt.Equal(t, "int", 8)
	trt.Equal(t, "float", 8.8)
}

func TestClamp(t *testing.T) {
	tmpl := `
		{{ emit "below" (clamp 1 5 10) }}
		{{ emit "inside" (clamp 7 5 10) }}
		{{ emit "above" (clamp 15 5 10) }}
		{{ emit "float" (clamp 3.5 1.5 2.5) }}
	`

	trt := testutils.Collect(t, tmpl)
	trt.Equal(t, "below", 5)
	trt.Equal(t, "inside", 7)
	trt.Equal(t, "above", 10)
	trt.Equal(t, "float", 2.5)
}

func TestClampInvalidRange(t *testing.T) {
	tmpl := `{{ clamp 5 10 1 }}`

	testutils.ExpectErr(t, tmpl, nil)
}

func TestPow(t *testing.T) {
	tmpl := `
		{{ emit "int" (pow 2 3) }}
		{{ emit "float" (pow 2.5 2) }}
	`

	trt := testutils.Collect(t, tmpl)
	trt.Equal(t, "int", 8.0)
	trt.Equal(t, "float", 6.25)
}

func TestSqrt(t *testing.T) {
	tmpl := `
		{{ emit "perfect" (sqrt 25) }}
		{{ emit "float" (sqrt 2.25) }}
		{{ emit "zero" (sqrt 0) }}
	`

	trt := testutils.Collect(t, tmpl)
	trt.Equal(t, "perfect", 5.0)
	trt.Equal(t, "float", 1.5)
	trt.Equal(t, "zero", 0.0)
}

func TestSqrtNegative(t *testing.T) {
	tmpl := `{{ sqrt -1 }}`

	testutils.ExpectErr(t, tmpl, nil)
}

func TestRound(t *testing.T) {
	tmpl := `
		{{ emit "positive" (round 3.7) }}
		{{ emit "negative" (round -3.7) }}
		{{ emit "integer" (round 5) }}
		{{ emit "positiveDigits" (round 3.14159 2) }}
		{{ emit "negativeDigits" (round 3.14159 0) }}
		{{ emit "negativePrecision" (round 1234.0 -2) }}
	`

	trt := testutils.Collect(t, tmpl)
	trt.Equal(t, "positive", 4)
	trt.Equal(t, "negative", -4)
	trt.Equal(t, "integer", 5)
	trt.Equal(t, "positiveDigits", 3.14)
	trt.Equal(t, "negativeDigits", 3.0)
	trt.Equal(t, "negativePrecision", 1200.0)
}

func TestFloor(t *testing.T) {
	tmpl := `
		{{ emit "positive" (floor 3.7) }}
		{{ emit "negative" (floor -3.2) }}
		{{ emit "integer" (floor 5) }}
	`

	trt := testutils.Collect(t, tmpl)
	trt.Equal(t, "positive", 3)
	trt.Equal(t, "negative", -4)
	trt.Equal(t, "integer", 5)
}

func TestCeil(t *testing.T) {
	tmpl := `
		{{ emit "positive" (ceil 3.2) }}
		{{ emit "negative" (ceil -3.7) }}
		{{ emit "integer" (ceil 5) }}
	`

	trt := testutils.Collect(t, tmpl)
	trt.Equal(t, "positive", 4)
	trt.Equal(t, "negative", -3)
	trt.Equal(t, "integer", 5)
}