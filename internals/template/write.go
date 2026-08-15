package template

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/utils/fsutils"
)

func resolveOutput(source, output string, preserveStruct bool) string {
	if preserveStruct {
		return fsutils.ResolveOutputPreserved(source, output)
	}

	return fsutils.ResolveOutput(source, output)
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

func handleFile(ctx *context.TemplateContext) {
	if !matchFile(ctx) {
		if ctx.Options.Verbose {
			fmt.Println("skipped", ctx.Path)
		}
		return
	}

	if ctx.Options.Verbose {
		fmt.Println("templating", ctx.Path)
	}

	ctx.OutputPath, _= filepath.Abs(resolveOutput(ctx.Path, ctx.Options.Output, !ctx.Options.Flatten))

	content, err := templateFile(ctx)

	if err != nil && !ctx.Options.Supress {
		fmt.Println("error templating:", err.Error())
		return
	}

	handleFileWrite(content, ctx)
}

func handleFileWrite(content string, ctx *context.TemplateContext) error {
	filePathAbs, _ := filepath.Abs(ctx.OutputPath)

	allowed := false
	for _, try := range ctx.Options.AllowedOutputFolders {
		if fsutils.IsInside(filePathAbs, try) {
			allowed = true
		}
	}

	if !allowed {
		absPath, _ := filepath.Abs(ctx.Path)
		fmt.Println("error outputting " + absPath + " to " + filePathAbs)
		return nil
	}

	if ctx.Options.Verbose {
		fmt.Println("writing to", filePathAbs)
	}

	dir := filepath.Dir(filePathAbs)
	err := os.MkdirAll(dir, 0755);

	if err != nil {
		return err
	}

	err = os.WriteFile(filePathAbs, []byte(content), 0644)

	if err != nil {
		return err
	}

	return nil
}