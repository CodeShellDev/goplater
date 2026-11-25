package template

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/utils/fsutils"
	"github.com/codeshelldev/goplater/utils/outputresolver"
)

func resolveOutput(root, source, output string, preserveStruct bool) (string, error) {
	if preserveStruct {
		if root != "" {
			return outputresolver.ResolveOutputPreserved(root, source, output)
		}

		return outputresolver.ResolveOutputInferred(source, output)
	}

	return outputresolver.ResolveOutput(source, output), nil
}

func templateFile(tmplContext *context.TemplateContext) (string, error) {
	root, err := os.OpenRoot(filepath.Dir(tmplContext.Shared.Path))
	if err != nil {
		return "", err
	}
	defer root.Close()

	name := filepath.Base(tmplContext.Shared.Path)

	info, err := root.Lstat(name)
	if err != nil {
		return "", err
	}

	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("%q is not a regular file", tmplContext.Shared.Path)
	}

	data, err := root.ReadFile(name)
	if err != nil {
		return "", err
	}

	return templateContent(string(data), tmplContext)
}

func handleFile(tmplContext *context.TemplateContext) error {
	output, err := resolveOutput(tmplContext.Options.Root, tmplContext.Shared.Path, tmplContext.Options.Output, !tmplContext.Options.Flatten)

	if err != nil {
		return err
	}

	tmplContext.Shared.OutputPath, err = filepath.Abs(output)

	if err != nil {
		return err
	}

	content, err := templateFile(tmplContext)

	if err != nil {
		return err
	}

	return handleFileWrite(content, tmplContext)
}

func handleFileWrite(content string, tmplContext *context.TemplateContext) error {
	filePathAbs, err := filepath.Abs(tmplContext.Shared.OutputPath)

	if err != nil {
		return err
	}

	// case 1: explicitly allowed file => direct write
	if slices.Contains(tmplContext.Options.AllowedOutputs, filePathAbs) {
		if tmplContext.Options.DryRun {
			return nil
		}

		err := os.MkdirAll(filepath.Dir(filePathAbs), 0755)
		if err != nil {
			return err
		}

		return os.WriteFile(filePathAbs, []byte(content), 0644)
	}

	// case 2: allowed dir containing target?
	var allowedDir string

	for _, p := range tmplContext.Options.AllowedOutputs {
		if fsutils.IsInside(filePathAbs, p) {
			allowedDir = p
			break
		}
	}

	if allowedDir == "" {
		return fmt.Errorf("target %q is not allowed", filePathAbs)
	}

	if tmplContext.Options.DryRun {
		return nil
	}

	err = os.MkdirAll(allowedDir, 0755)
	if err != nil {
		return err
	}

	return fsutils.WithAllowedRoot(filePathAbs, []string{allowedDir},
		func(root *os.Root, relative string) error {
			dir := filepath.Dir(relative)

			err := root.MkdirAll(dir, 0755)
			if err != nil {
				return fmt.Errorf("could not create output %q: %w", dir, err)
			}

			return root.WriteFile(relative, []byte(content), 0644)
		},
	)
}