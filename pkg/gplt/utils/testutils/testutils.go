package testutils

import (
	"errors"
	"reflect"
	"testing"

	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/collections"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/test"
	"github.com/codeshelldev/goplater/pkg/gplt/resolvers"
)

type TestRuntime struct {
	store *test.Store
	Output string
}

func render(tmpl string, ctx *gplt.Context, modules []modules.Module, data ...any) (string, *gplt.Runtime, error) {
	e := gplt.NewEngine()

	e.Use(collections.Unsafe...)
	e.UseSafe(collections.Safe...)

	e.Use(test.Module)

	if modules != nil {
		e.Use(modules...)
	}

	e.SetResolver(gplt.NewResolverChain(resolvers.NewHttpResolver()))

	rt := e.NewRuntime(gplt.EngineOptions{})

	input := any(data)
	if len(data) == 1 {
		input = data[0]
	}

	if ctx == nil {
		ctx = &gplt.Context{}
	}

	output, err := rt.Render("test", tmpl, input, ctx)

	if err != nil {
		return "", rt, err
	}

	return output, rt, nil
}

func Collect(t *testing.T, tmpl string, data ...any) TestRuntime {
	t.Helper()

	return CollectWithContext(t, tmpl, nil, nil, data...)
}

func CollectWithContext(t *testing.T, tmpl string, ctx *gplt.Context,  modules []modules.Module, data ...any) TestRuntime {
	t.Helper()

	output, rt, err := render(tmpl, ctx, modules, data...)
	if err != nil {
		t.Fatalf("unexpected error: %v\ntemplate: %s", err, tmpl)
	}

	s, _ := test.GetResults(rt)

	return TestRuntime{
		store:  s,
		Output: output,
	}
}

func ExpectErr(t *testing.T, tmpl string, expectErr error) TestRuntime {
	t.Helper()

	return ExpectErrWithContext(t, tmpl, nil, nil, expectErr)
}

func ExpectErrWithContext(t *testing.T, tmpl string, ctx *gplt.Context, modules []modules.Module, expectErr error) TestRuntime {
	t.Helper()

	output, _, err := render(tmpl, ctx, modules, nil)

	if err == nil {
		t.Fatalf("returned no errors, expected: %v", expectErr)
	}

	if expectErr != nil && !errors.Is(err, expectErr) {
		t.Fatalf("expected: %v, got: %v", expectErr, err)
	}

	return TestRuntime{
		Output: output,
	}
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