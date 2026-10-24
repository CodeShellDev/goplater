package tests

import (
	"path/filepath"
	"testing"

	"github.com/codeshelldev/goplater/utils/outputresolver"
)

type outputPathCase struct {
	name       	string
	root 		string
	source     	string
	output     	string
	expect     	string
	wantErr    	bool
}

func relToCwd(absPath string) (string, error) {
	cwd, err := filepath.Abs(".")
	if err != nil {
		return "", err
	}

	rel, err := filepath.Rel(cwd, absPath)
	if err != nil {
		return "", err
	}

	return filepath.ToSlash(rel), nil
}

func resolvePreserved(root, source, output string) (string, error) {
	abs, err := outputresolver.ResolveOutputPreserved(root, source, output)
	if err != nil {
		return "", err
	}

	return relToCwd(abs)
}

func resolveInferred(source, output string) (string, error) {
	abs, err := outputresolver.ResolveOutputInferred(source, output)
	if err != nil {
		return "", err
	}

	return relToCwd(abs)
}

func TestOutputPathFlatten(t *testing.T) {
	tests := []outputPathCase{
		{source: "source.txt", output: "source.txt", expect: "source.txt"},
		{source: "source.txt", output: ".", expect: "source.txt"},
		{source: "source.txt", output: "other.txt", expect: "other.txt"},
		{source: "source/source.txt", output: "output/other.txt", expect: "output/other.txt"},
		{source: "source/source.txt", output: "output/", expect: "output/source.txt"},
		{source: "source/source.txt", output: "output/sub/", expect: "output/sub/source.txt"},
		{source: "source/sub/", output: "output/", expect: "output/sub"},
		{source: "source/sub/folder/source.txt", output: "output/", expect: "output/source.txt"},
		{source: "fs/source/sub/folder/source.txt", output: "fs/output/", expect: "fs/output/source.txt"},
	}

	for _, tc := range tests {
		t.Run(tc.source + " -> " + tc.output, func(t *testing.T) {
			got := filepath.ToSlash(outputresolver.ResolveOutput(tc.source, tc.output))
			if got != tc.expect {
				t.Error(
					"ResolveOutput(" + tc.source + ", " + tc.output + ")\n got: " + got + "\nwant: " + tc.expect,
				)
			}
		})
	}
}

func TestOutputPathPreserve(t *testing.T) {
	tests := []outputPathCase{
		{root: ".", source: "source.txt", output: "source.txt", expect: "source.txt"},
		{root: ".", source: "source.txt", output: ".", expect: "source.txt"},
		{root: ".", source: "source.txt", output: "other.txt", expect: "other.txt"},
		{root: ".", source: "source/source.txt", output: "output/other.txt", expect: "output/other.txt"},

		{root: "source", source: "source/source.txt", output: ".", expect: "source.txt"},
		{root: "source", source: "source/source.txt", output: "output/", expect: "output/source.txt"},
		{root: "source", source: "source/source.txt", output: "output/sub/", expect: "output/sub/source.txt"},
		{root: "source", source: "source/sub/", output: "output/", expect: "output/sub"},
		{root: "source", source: "source/sub/folder/source.txt", output: "output/", expect: "output/sub/folder/source.txt"},

		{root: "fs/source", source: "fs/source/sub/folder/source.txt", output: "fs/output/", expect: "fs/output/sub/folder/source.txt"},
		{root: "fs/source", source: "fs/source/sub/folder/", output: "fs/output/", expect: "fs/output/sub/folder"},

		{root: ".fs", source: ".fs/test.md", output: ".", expect: "test.md"},
		{root: ".fs", source: ".fs/test/test.md", output: ".", expect: "test/test.md"},

		// root as non-clean path
		{name: "root with ./", root: "./source", source: "source/sub/x.txt", output: "output/", expect: "output/sub/x.txt"},

		// errors
		{name: "source outside root", root: "source", source: "other/x.txt", output: "output/", wantErr: true},
		{name: "dir source + file output", root: "source", source: "source/sub/", output: "out.txt", wantErr: true},
	}

	for _, tc := range tests {
		name := tc.name
		if name == "" {
			name = tc.root + " | " + tc.source + " -> " + tc.output
		}

		t.Run(name, func(t *testing.T) {
			got, err := resolvePreserved(tc.root, tc.source, tc.output)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got " + got)
				}
				return
			}

			if err != nil {
				t.Fatal("unexpected error: " + err.Error())
			}

			if got != tc.expect {
				t.Error(
					"ResolveOutputPreserved(" + tc.root + ", " + tc.source + ", " + tc.output + ")\n got: " + got + "\nwant: " + tc.expect,
				)
			}
		})
	}
}

func TestOutputPathInferred(t *testing.T) {
	tests := []outputPathCase{
		{source: "source/source.txt", output: ".", expect: "source.txt"},
		{source: "source.txt", output: "output.txt", expect: "output.txt"},
		{source: "source/source.txt", output: "output/output.txt", expect: "output/output.txt"},
		{source: "source/source.txt", output: "output/", expect: "output/source.txt"},
		{source: "source/source.txt", output: "output/sub/", expect: "output/sub/source.txt"},
		{source: "source/sub/", output: "output/", expect: "output/sub"},
		{source: "source/sub/folder/source.txt", output: "output/", expect: "output/sub/folder/source.txt"},
		{source: "fs/source/sub/folder/source.txt", output: "fs/output/", expect: "fs/output/sub/folder/source.txt"},
		{source: "fs/source/sub/folder/", output: "fs/output/", expect: "fs/output/sub/folder"},

		// single segment below common parent: keep
		{name: "single segment", source: "source.txt", output: ".", expect: "source.txt"},

		// output beside source in workspace
		{name: "workspace", source: "ws/src/x.txt", output: "ws/build/", expect: "ws/build/x.txt"},

		// same path, different text
		{name: "leading ./", source: "./source/sub/x.txt", output: "./output/", expect: "output/sub/x.txt"},

		// errors
		{name: "no shared parent", source: "/a/b/x.txt", output: "/c/d/", wantErr: true},
		{name: "dir source + file output", source: "source/sub/", output: "out.txt", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := resolveInferred(tc.source, tc.output)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got " + got)
				}
				return
			}

			if err != nil {
				t.Fatal("unexpected error: " + err.Error())
			}

			if got != tc.expect {
				t.Error(
					"ResolveOutputInferred(" + tc.source + ", " + tc.output + ")\n got: " + got + "\nwant: " + tc.expect,
				)
			}
		})
	}
}