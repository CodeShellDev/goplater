package template

import (
	"errors"
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
	if !fsutils.IsFile(ctx.Path) {
		return "", os.ErrNotExist
	}

	data, err := os.ReadFile(ctx.Path)

	if err != nil {
		return string(data), err
	}

	if data == nil {
		return string(data), errors.New("empty file")
	}

	tmplStr, err := templateContent(string(data), ctx)

	if err != nil {
		return string(data), err
	}

	return tmplStr, nil
}

func handleFile(ctx *context.TemplateContext) error {
	output, err := resolveOutput(ctx.Options.Root, ctx.Path, ctx.Options.Output, !ctx.Options.Flatten)

	if err != nil {
		return err
	}

	ctx.OutputPath, err = filepath.Abs(output)

	if err != nil {
		return err
	}

	content, err := templateFile(ctx)

	if err != nil {
		if !ctx.Options.IgnoreErrors {
			return err
		}
		return nil
	}

	return handleFileWrite(content, ctx)
}

func handleFileWrite(content string, ctx *context.TemplateContext) error {
	filePathAbs, _ := filepath.Abs(ctx.OutputPath)

	allowed := false
	for _, try := range ctx.Options.AllowedOutputs {
		if fsutils.IsInside(filePathAbs, try) {
			allowed = true
		}
	}

	if !allowed {
		return errors.New("not allowed to output to " + filePathAbs)
	}

	if ctx.Options.DryRun {
		return nil
	}
	
	dir := filepath.Dir(filePathAbs)
	err := os.MkdirAll(dir, 0755)

	if err != nil {
		return err
	}

	return os.WriteFile(filePathAbs, []byte(content), 0644)
}