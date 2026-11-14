package tests

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/internals/template/funcs"
	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
	"github.com/codeshelldev/goplater/pkg/gplt/utils/testutils"
)

func newFSTestContext(t *testing.T) (*gplt.Context, string, string) {
	t.Helper()

	root := t.TempDir()
	source := filepath.Join(root, "source")
	output := filepath.Join(root, "output")

	err := os.MkdirAll(source, 0755)
	if err != nil {
		t.Fatal(err)
	}

	err = os.MkdirAll(output, 0755)
	if err != nil {
		t.Fatal(err)
	}

	tmplContext := &context.TemplateContext{
		Invoker: filepath.Join(source, "main.gplt"),
		Options: &context.TemplateOptions{
			Root:           root,
			Source:         source,
			Output:         output,
			AllowedInputs:  []string{source},
			AllowedOutputs: []string{output},
		},
	}

	ctx := &gplt.Context{}
	ctx.Set(context.TemplateContextKey, tmplContext)

	return ctx, source, output
}

func collectFS(t *testing.T, tmpl string, ctx *gplt.Context) testutils.TestRuntime {
	t.Helper()

	return testutils.CollectWithContext(
		t,
		tmpl,
		ctx,
		[]modules.Module{funcs.Module},
	)
}

func expectErrFS(t *testing.T, tmpl string, ctx *gplt.Context, expectErr error) testutils.TestRuntime {
	t.Helper()

	return testutils.ExpectErrWithContext(
		t,
		tmpl,
		ctx,
		[]modules.Module{funcs.Module},
		expectErr,
	)
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()

	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}
}


func TestReadRaw(t *testing.T) {
	ctx, source, _ := newFSTestContext(t)
	writeTestFile(t, filepath.Join(source, "hello.txt"), "Hello, World!")

	trt := collectFS(t, `{{ emit "result" (readRaw "./hello.txt") }}`, ctx)

	trt.Equal(t, "result", "Hello, World!")
}

func TestReadRawMissingFile(t *testing.T) {
	ctx, _, _ := newFSTestContext(t)

	expectErrFS(t, `{{ readRaw "./missing.txt" }}`, ctx, nil)
}

func TestReadRawOutsideAllowedInputs(t *testing.T) {
	ctx, _, output := newFSTestContext(t)

	path := filepath.Join(output, "secret.txt")
	writeTestFile(t, path, "secret")

	expectErrFS(t, `{{ readRaw ` + strconv.Quote(path) + ` }}`, ctx, nil)
}


func TestRead(t *testing.T) {
	ctx, source, _ := newFSTestContext(t)

	want := "Hello from nested template"

	writeTestFile(t, filepath.Join(source, "hello.txt"), `{{ print ` + strconv.Quote(want) + ` }}`)

	trt := collectFS(t, `{{ read "./hello.txt" }}`, ctx)

	if trt.Output != want {
		t.Fatalf("file templating failed: got %q, want %q", trt.Output, want)
	}
}

func TestReadMissingFile(t *testing.T) {
	ctx, _, _ := newFSTestContext(t)

	expectErrFS(t, `{{ read "./missing.txt" }}`, ctx, nil)
}


func TestWrite(t *testing.T) {
	ctx, _, output := newFSTestContext(t)
	path := filepath.Join(output, "hello.txt")

	want := "Hello, World!"

	collectFS(t, `{{ write ` + strconv.Quote(path) + ` ` + strconv.Quote(want) + ` }}`, ctx)

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestWriteOutsideAllowedOutputs(t *testing.T) {
	ctx, source, _ := newFSTestContext(t)
	path := filepath.Join(source, "forbidden.txt")

	expectErrFS(t, `{{ write ` + strconv.Quote(path) + ` "Hello" }}`, ctx, nil)

	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		t.Fatal("file should not have been created")
	}
}


func TestMkdir(t *testing.T) {
	ctx, _, output := newFSTestContext(t)
	path := filepath.Join(output, "nested", "folder")

	collectFS(t, `{{ mkdir ` + strconv.Quote(path) + ` }}`, ctx)

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	if !info.IsDir() {
		t.Fatal("expected a directory")
	}
}

func TestMkdirOutsideAllowedOutputs(t *testing.T) {
	ctx, source, _ := newFSTestContext(t)
	path := filepath.Join(source, "forbidden")

	expectErrFS(t, `{{ mkdir ` + strconv.Quote(path) + ` }}`, ctx, nil)
}


func TestAppendFile(t *testing.T) {
	ctx, _, output := newFSTestContext(t)
	path := filepath.Join(output, "hello.txt")
	writeTestFile(t, path, "Hello")

	collectFS(t, `{{ appendFile ` + strconv.Quote(path) + ` ", World!" }}`, ctx)

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != "Hello, World!" {
		t.Fatalf("got %q, want %q", got, "Hello, World!")
	}
}

func TestAppendFileCreatesFile(t *testing.T) {
	ctx, _, output := newFSTestContext(t)
	path := filepath.Join(output, "new.txt")

	collectFS(t, `{{ appendFile ` + strconv.Quote(path) + ` "Hello" }}`, ctx)

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != "Hello" {
		t.Fatalf("got %q, want %q", got, "Hello")
	}
}


func TestFsExists(t *testing.T) {
	ctx, source, _ := newFSTestContext(t)
	writeTestFile(t, filepath.Join(source, "file.txt"), "content")

	err := os.Mkdir(filepath.Join(source, "folder"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{"file", "./file.txt", true},
		{"directory", "./folder", true},
		{"missing", "./missing", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trt := collectFS(t, `{{ emit "result" (fsExists ` + strconv.Quote(tt.path) + `) }}`, ctx)

			trt.Equal(t, "result", tt.want)
		})
	}
}


func TestIsFile(t *testing.T) {
	ctx, source, _ := newFSTestContext(t)
	writeTestFile(t, filepath.Join(source, "file.txt"), "content")

	err := os.Mkdir(filepath.Join(source, "folder"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{"file", "./file.txt", true},
		{"directory", "./folder", false},
		{"missing", "./missing", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trt := collectFS(t, `{{ emit "result" (isFile ` + strconv.Quote(tt.path) + `) }}`, ctx)

			trt.Equal(t, "result", tt.want)
		})
	}
}


func TestIsDir(t *testing.T) {
	ctx, source, _ := newFSTestContext(t)
	writeTestFile(t, filepath.Join(source, "file.txt"), "content")

	err := os.Mkdir(filepath.Join(source, "folder"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{"directory", "./folder", true},
		{"file", "./file.txt", false},
		{"missing", "./missing", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trt := collectFS(t, `{{ emit "result" (isDir ` + strconv.Quote(tt.path) + `) }}`, ctx)

			trt.Equal(t, "result", tt.want)
		})
	}
}


func TestListDir(t *testing.T) {
	ctx, source, _ := newFSTestContext(t)
	writeTestFile(t, filepath.Join(source, "b.txt"), "B")
	writeTestFile(t, filepath.Join(source, "a.txt"), "A")

	err := os.Mkdir(filepath.Join(source, "folder"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	trt := collectFS(t, `{{ emit "result" (listDir "./") }}`, ctx)

	trt.Equal(t, "result", []string{"a.txt", "b.txt", "folder"})
}

func TestListDirMissingDirectory(t *testing.T) {
	ctx, _, _ := newFSTestContext(t)

	expectErrFS(t, `{{ listDir "./missing" }}`, ctx, nil)
}


func TestWalkDir(t *testing.T) {
	ctx, source, _ := newFSTestContext(t)
	writeTestFile(t, filepath.Join(source, "file.txt"), "content")
	writeTestFile(t, filepath.Join(source, "nested", "inner.txt"), "inner")

	trt := collectFS(t, `{{ emit "result" (walkDir "./") }}`, ctx)

	trt.Equal(t, "result", []string{
		source,
		filepath.Join(source, "file.txt"),
		filepath.Join(source, "nested"),
		filepath.Join(source, "nested", "inner.txt"),
	})
}

func TestWalkDirMissingDirectory(t *testing.T) {
	ctx, _, _ := newFSTestContext(t)

	expectErrFS(t, `{{ walkDir "./missing" }}`, ctx, nil)
}


func TestFsRemoveFile(t *testing.T) {
	ctx, _, output := newFSTestContext(t)
	path := filepath.Join(output, "file.txt")
	writeTestFile(t, path, "content")

	collectFS(t, `{{ fsRemove ` + strconv.Quote(path) + ` }}`, ctx)

	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		t.Fatal("file should have been removed")
	}
}

func TestFsRemoveDirectory(t *testing.T) {
	ctx, _, output := newFSTestContext(t)
	dir := filepath.Join(output, "folder")
	writeTestFile(t, filepath.Join(dir, "nested.txt"), "content")

	collectFS(t, `{{ fsRemove ` + strconv.Quote(dir) + ` }}`, ctx)

	_, err := os.Stat(dir)
	if !os.IsNotExist(err) {
		t.Fatal("directory should have been removed")
	}
}

func TestFsRemoveOutsideAllowedOutputs(t *testing.T) {
	ctx, source, _ := newFSTestContext(t)
	path := filepath.Join(source, "protected.txt")
	writeTestFile(t, path, "content")

	expectErrFS(t, `{{ fsRemove ` + strconv.Quote(path) + ` }}`, ctx, nil)

	_, err := os.Stat(path)
	if err != nil {
		t.Fatal("file should still exist")
	}
}


func TestJoinPath(t *testing.T) {
	trt := collectFS(t, `{{ emit "result" (joinPath "assets" "images" "logo.png") }}`, nil)

	trt.Equal(t, "result", filepath.Join("assets", "images", "logo.png"))
}

func TestBasePath(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"assets/images/logo.png", "logo.png"},
		{"assets/images", "images"},
		{"logo.png", "logo.png"},
		{".", "."},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			trt := collectFS(t, `{{ emit "result" (basePath ` + strconv.Quote(tt.path) + `) }}`, nil)

			trt.Equal(t, "result", tt.want)
		})
	}
}

func TestPathDir(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"assets/images/logo.png", filepath.Join("assets", "images")},
		{"assets/logo.png", "assets"},
		{"logo.png", "."},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			trt := collectFS(t, `{{ emit "result" (pathDir ` + strconv.Quote(tt.path) + `) }}`, nil)

			trt.Equal(t, "result", tt.want)
		})
	}
}

func TestFileExt(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"logo.png", ".png"},
		{"archive.tar.gz", ".gz"},
		{"README", ""},
		{"folder/file.txt", ".txt"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			trt := collectFS(t, `{{ emit "result" (fileExt ` + strconv.Quote(tt.path) + `) }}`, nil)

			trt.Equal(t, "result", tt.want)
		})
	}
}

func TestAbsPath(t *testing.T) {
	trt := collectFS(t, `{{ emit "result" (absPath ".") }}`, nil)

	want, err := filepath.Abs(".")
	if err != nil {
		t.Fatal(err)
	}
	trt.Equal(t, "result", want)
}

func TestRelPath(t *testing.T) {
	trt := collectFS(t, `{{ emit "result" (relPath "/project/src" "/project/assets/logo.png") }}`, nil)

	trt.Equal(t, "result", filepath.Join("..", "assets", "logo.png"))
}
