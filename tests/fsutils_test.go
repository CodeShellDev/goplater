package tests

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/codeshelldev/goplater/utils/fsutils"
)

func expectMkdir(t *testing.T, path string) {
	t.Helper()

	err := os.MkdirAll(path, 0755)
	if err != nil {
		t.Fatal(err)
	}
}

func expectSymlink(t *testing.T, target, link string) {
	t.Helper()

	err := os.Symlink(target, link)
	if err != nil {
		t.Skipf("symlinks not supported: %v", err)
	}
}

func TestWithAllowedRootBroadestWins(t *testing.T) {
	base := t.TempDir()
	sub := filepath.Join(base, "sub")
	expectMkdir(t, sub)

	target := filepath.Join(sub, "f.txt")

	orders := map[string][]string{
		"narrow first": {sub, base},
		"broad first":  {base, sub},
	}

	for name, allowed := range orders {
		t.Run(name, func(t *testing.T) {
			calls := 0

			err := fsutils.WithAllowedRoot(target, allowed, func(root *os.Root, relative string) error {
				calls++

				if root.Name() != base {
					t.Errorf("root: got %q, want %q", root.Name(), base)
				}

				if relative != filepath.Join("sub", "f.txt") {
					t.Errorf("relative: got %q", relative)
				}

				return nil
			})

			if err != nil {
				t.Fatal(err)
			}

			if calls != 1 {
				t.Fatalf("fn called %d times, want 1", calls)
			}
		})
	}
}

func TestWithAllowedRootNoRetryOnFnError(t *testing.T) {
	base := t.TempDir()
	sub := filepath.Join(base, "sub")
	expectMkdir(t, sub)

	errBoom := errors.New("boom")
	calls := 0

	err := fsutils.WithAllowedRoot(filepath.Join(sub, "f.txt"), []string{sub, base}, func(_ *os.Root, _ string) error {
		calls++
		return errBoom
	})

	if !errors.Is(err, errBoom) {
		t.Fatalf("got %v, want %v", err, errBoom)
	}

	if calls != 1 {
		t.Fatalf("fn called %d times, want 1 (no retry)", calls)
	}
}

func TestWithAllowedRootOutside(t *testing.T) {
	allowed := t.TempDir()
	other := t.TempDir()

	calls := 0

	err := fsutils.WithAllowedRoot(filepath.Join(other, "f.txt"), []string{allowed}, func(_ *os.Root, _ string) error {
		calls++
		return nil
	})

	if err == nil {
		t.Fatal("expected error for path outside allowed dirs")
	}

	if calls != 0 {
		t.Fatalf("fn called %d times, want 0", calls)
	}
}

func TestWithAllowedRootEmptyAllowed(t *testing.T) {
	err := fsutils.WithAllowedRoot(filepath.Join(t.TempDir(), "f.txt"), nil, func(_ *os.Root, _ string) error {
		t.Fatal("fn shouldn't be called")
		return nil
	})

	if err == nil {
		t.Fatal("expected error with no allowed dirs")
	}
}

func TestWithAllowedRootSiblingNotAllowed(t *testing.T) {
	base := t.TempDir()
	a := filepath.Join(base, "a")
	ab := filepath.Join(base, "ab")
	expectMkdir(t, a)
	expectMkdir(t, ab)

	err := fsutils.WithAllowedRoot(filepath.Join(ab, "f.txt"), []string{a}, func(_ *os.Root, _ string) error {
		t.Fatal("fn shouldn't be called")
		return nil
	})

	if err == nil {
		t.Fatal("siblings shouldn't be allowed")
	}
}

func TestWithAllowedRootFileEntry(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "main.txt")
	sibling := filepath.Join(dir, "other.txt")

	writeTestFile(t, file, "main")
	writeTestFile(t, sibling, "other")

	t.Run("exact file allowed", func(t *testing.T) {
		calls := 0

		err := fsutils.WithAllowedRoot(file, []string{file}, func(root *os.Root, relative string) error {
			calls++

			if root.Name() != dir {
				t.Errorf("root: got %q, want %q", root.Name(), dir)
			}

			if relative != "main.txt" {
				t.Errorf("relative: got %q", relative)
			}

			return nil
		})

		if err != nil {
			t.Fatal(err)
		}

		if calls != 1 {
			t.Fatalf("fn called %d times, want 1", calls)
		}
	})

	t.Run("sibling denied", func(t *testing.T) {
		err := fsutils.WithAllowedRoot(sibling, []string{file}, func(_ *os.Root, _ string) error {
			t.Fatal("fn shouldn't called")
			return nil
		})

		if err == nil {
			t.Fatal("sibling of allowed file shouldn't be allowed")
		}
	})
}

func TestWithAllowedRootSymlinkEscapeDenied(t *testing.T) {
	base := t.TempDir()
	allowed := filepath.Join(base, "allowed")
	outside := filepath.Join(base, "outside")
	expectMkdir(t, allowed)
	expectMkdir(t, outside)

	writeTestFile(t, filepath.Join(outside, "secret.txt"), "secret")
	expectSymlink(t, outside, filepath.Join(allowed, "link"))

	err := fsutils.WithAllowedRoot(filepath.Join(allowed, "link", "secret.txt"), []string{allowed}, func(root *os.Root, relative string) error {
		_, err := root.ReadFile(relative)
		return err
	})

	if err == nil {
		t.Fatal("symlink escaping allowed root must fail")
	}
}

func TestWithAllowedRootSymlinkInsideBroadRoot(t *testing.T) {
	base := t.TempDir()
	sub := filepath.Join(base, "sub")
	other := filepath.Join(base, "other")
	expectMkdir(t, sub)
	expectMkdir(t, other)

	writeTestFile(t, filepath.Join(other, "f.txt"), "ok")
	expectSymlink(t, filepath.Join("..", "other"), filepath.Join(sub, "link"))

	// symlink leaves "sub" but not "base", so it should be allowed
	for _, allowed := range [][]string{{sub, base}, {base, sub}} {
		var got []byte

		err := fsutils.WithAllowedRoot(filepath.Join(sub, "link", "f.txt"), allowed, func(root *os.Root, relative string) error {
			data, err := root.ReadFile(relative)
			got = data
			return err
		})

		if err != nil {
			t.Fatalf("allowed=%v: %v", allowed, err)
		}

		if string(got) != "ok" {
			t.Fatalf("allowed=%v: got %q", allowed, got)
		}
	}
}