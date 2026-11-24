package template

import (
	"fmt"
	"os"
	"path/filepath"

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

func handleFileWrite(content string, ctx *context.TemplateContext) error {
	filePathAbs, err := filepath.Abs(ctx.Shared.OutputPath)

	if err != nil {
		return err
	}

	if !ctx.Options.DryRun {
		var allowed bool
		for _, path := range ctx.Options.AllowedOutputs {
			if filePathAbs == path || fsutils.IsInside(filePathAbs, path) {
				allowed = true

				err := os.MkdirAll(filepath.Dir(filePathAbs), 0755)

				if err != nil {
					return err
				}

				// explicitly allowed file => direct write
				if filePathAbs == path {
					return os.WriteFile(filePathAbs, []byte(content), 0644)
				}
			}
		}

		if !allowed {
			return fmt.Errorf("target %q is not allowed", filePathAbs)
		}
	}

	return fsutils.WithAllowedRoot(filePathAbs, ctx.Options.AllowedOutputs, func(root *os.Root, relative string) error {
		if ctx.Options.DryRun {
			return nil
		}

		return root.WriteFile(relative, []byte(content), 0644)
	})
}