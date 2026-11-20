package testutils

import (
	"errors"
	"reflect"
	"testing"

	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
	"github.com/codeshelldev/goplater/pkg/gplt/modules/test"
	"github.com/codeshelldev/goplater/pkg/gplt/resolvers"
)

type TestRuntime struct {
	store *test.Store
	Output string
}

type Tester struct {
	safeModules []modules.Module
	modules []modules.Module
}

func NewTester(modules []modules.Module, safeModules []modules.Module) Tester {
	return Tester{modules: modules, safeModules: safeModules}
}

func render(tmpl string, ctx *gplt.Context, modules []modules.Module, safeModules []modules.Module, data ...any) (string, *gplt.Runtime, error) {
	e := gplt.NewEngine()

	e.Use(modules...)
	e.UseSafe(safeModules...)

	e.Use(test.Module)

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

func (tester Tester) Collect(t *testing.T, tmpl string, data ...any) TestRuntime {
	t.Helper()

	return tester.CollectWithContext(t, tmpl, nil, nil, data...)
}

func (tester Tester) CollectWithContext(t *testing.T, tmpl string, ctx *gplt.Context, modules []modules.Module, data ...any) TestRuntime {
	t.Helper()

	output, rt, err := render(tmpl, ctx, tester.modules, append(tester.safeModules, modules...), data...)
	if err != nil {
		t.Fatalf("unexpected error: %v\ntemplate: %s", err, tmpl)
	}

	s, _ := test.GetResults(rt)

	return TestRuntime{
		store:	s,
		Output: output,
	}
}

func (tester Tester) ExpectErr(t *testing.T, tmpl string, expectErr error) TestRuntime {
	t.Helper()

	return tester.ExpectErrWithContext(t, tmpl, nil, nil, expectErr)
}

func (tester Tester) ExpectErrWithContext(t *testing.T, tmpl string, ctx *gplt.Context, modules []modules.Module, expectErr error) TestRuntime {
	t.Helper()

	output, _, err := render(tmpl, ctx, tester.modules, append(tester.safeModules, modules...), nil)

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