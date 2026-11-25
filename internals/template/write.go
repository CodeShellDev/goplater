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

func templateFile(ctx *context.TemplateContext) (string, error) {
	root, err := os.OpenRoot(filepath.Dir(ctx.Shared.Path))
	if err != nil {
		return "", err
	}
	defer root.Close()

	name := filepath.Base(ctx.Shared.Path)

	info, err := root.Lstat(name)
	if err != nil {
		return "", err
	}

	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("%q is not a regular file", ctx.Shared.Path)
	}

	data, err := root.ReadFile(name)
	if err != nil {
		return "", err
	}

	return templateContent(string(data), ctx)
}

func handleFile(ctx *context.TemplateContext) error {
	output, err := resolveOutput(ctx.Options.Root, ctx.Shared.Path, ctx.Options.Output, !ctx.Options.Flatten)

	if err != nil {
		return err
	}

	ctx.Shared.OutputPath, err = filepath.Abs(output)

	if err != nil {
		return err
	}

	content, err := templateFile(ctx)

	if err != nil {
		return err
	}

	return handleFileWrite(content, ctx)
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