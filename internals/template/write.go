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

func templateFile(context context.TemplateContext) (string, error) {
	if !fsutils.IsFile(context.Path) {
		return "", os.ErrNotExist
	}

	data, err := os.ReadFile(context.Path)

	if err != nil {
		return string(data), err
	}

	if data == nil {
		return string(data), errors.New("empty file")
	}

	tmplStr, err := templateContent(string(data), context)

	if err != nil {
		return string(data), err
	}

	return tmplStr, nil
}

func handleFile(context context.TemplateContext) {
	if !matchFile(context) {
		if context.Options.Verbose {
			fmt.Println("skipped", context.Path)
		}

		return
	}

	if context.Options.Verbose {
		fmt.Println("templating", context.Path)
	}

	content, err := templateFile(context)

	if err != nil && !context.Options.Supress {
		fmt.Println("error templating:", err.Error())

		return
	}

	handleFileWrite(content, context)
}

func handleFileWrite(content string, context context.TemplateContext) error {
	fullPath := resolveOutput(context.Path, context.Options.Output, !context.Options.Flatten)

	if context.Options.Verbose {
		fmt.Println("writing to", fullPath)
	}

	dir := filepath.Dir(fullPath)
	err := os.MkdirAll(dir, 0755);

	if err != nil {
		return err
	}

	err = os.WriteFile(fullPath, []byte(content), 0644)

	if err != nil {
		return err
	}

	return nil
}