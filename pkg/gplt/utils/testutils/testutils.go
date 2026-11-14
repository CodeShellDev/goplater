package testutils

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/collections"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/test"
	"github.com/codeshelldev/goplater/pkg/gplt/resolvers"
)

type TestRuntime struct {
	store *test.Store
	Output string
}

func render(tmpl string, data ...any) (string, *gplt.Runtime, error) {
	e := gplt.NewEngine()

	e.Use(collections.Unsafe...)
	e.UseSafe(collections.Safe...)

	e.Use(test.Module)

	e.SetResolver(gplt.NewResolverChain(resolvers.NewHttpResolver()))

	rt := e.NewRuntime(gplt.EngineOptions{})

	input := any(data)
	if len(data) == 1 {
		input = data[0]
	}

	output, err := rt.Render("test", tmpl, input, &gplt.Context{})

	if err != nil {
		return "", rt, err
	}

	return output, rt, nil
}

func Collect(t *testing.T, tmpl string, data ...any) (TestRuntime, error) {
	t.Helper()

	output, rt, err := render(tmpl, data...)

	if err != nil {
		return TestRuntime{}, err
	}

	s, _ := test.GetResults(rt)
	
	return TestRuntime{
		store: s,
		Output: output,
	}, nil
}

func ExpectErr(t *testing.T, tmpl string, expectErr error) (TestRuntime, error) {
	t.Helper()

	output, _, err := render(tmpl)

	if err == nil {
		return TestRuntime{}, fmt.Errorf("returned no errors, expected: %v", expectErr)
	}

	if expectErr != nil && !errors.Is(err, expectErr) {
		return TestRuntime{}, fmt.Errorf("expected: %v, got: %v", expectErr, err)
	}

	return TestRuntime{
		Output: output,
	}, nil
}

func (trt TestRuntime) Expect(t *testing.T, want ...any) {
	t.Helper()

	if trt.store == nil {
		t.Fatal("did not collect any emits")
	}

	got := trt.store.Entries()

	if len(got) != len(want) {
		t.Fatalf("got %d emits, want %d: %#v", len(got), len(want), got)
	}

	for i, w := range want {
		if !reflect.DeepEqual(got[i].Value, w) {
			t.Errorf("emit %d: got %#v (%T), want %#v (%T)", i, got[i].Value, got[i].Value, w, w)
		}
	}
}

func (trt TestRuntime) Unexpect(t *testing.T, label string) {
	t.Helper()

	if trt.store == nil {
		return
	}

	exists := trt.store.Has(label)

	if exists {
		t.Errorf("unexpected label %q exists", label)
	}
}

func (trt TestRuntime) Equal(t *testing.T, label string, want any) {
	t.Helper()

	if trt.store == nil {
		t.Fatal("did not collect any emits")
	}

	got, err := trt.store.Get(label)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s: got %#v (%T), want %#v (%T)", label, got, got, want, want)
	}
}