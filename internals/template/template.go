package template

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/utils/fsutils"
	"github.com/codeshelldev/goplater/utils/tree"
)

type Templater struct{}

func (Templater) Run(ctx *context.TemplateContext) error {
	return run(ctx)
}

func New() *Templater {
	return &Templater{}
}

func run(ctx *context.TemplateContext) error {
	isDir := fsutils.IsDir(ctx.Path)
	isFile := fsutils.IsFile(ctx.Path)
	
	ctx.Invoker = ctx.Path

	if isDir {
		t := tree.NewTree(os.Stdout, ctx.Path)

		err := filepath.WalkDir(ctx.Path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				if path != ctx.Path && !ctx.Options.Recursive {
					return filepath.SkipDir
				}
				return nil
			}

			absPath, err := filepath.Abs(path)

			if err != nil {
				return err
			}

			newContext := *ctx
			newContext.Path = absPath
			newContext.Invoker = absPath

			ok, err := matchFile(&newContext)

			if err != nil {
				return err
			}

			if !ok {
				return nil
			}

			relPath, err := filepath.Rel(ctx.Path, path)
			if err != nil {
				return err
			}

			file := t.Add(relPath)

			err = handleFile(&newContext)

			if err != nil {
				file.Fail(err)

				if !ctx.Options.IgnoreErrors {
					os.Exit(1)
					return err
				}
			}

			relOutput, err := fsutils.SafeRelative(ctx.Options.Source, newContext.OutputPath)

			if err != nil {
				file.Fail(err)
				os.Exit(1)
				return err
			}

			if newContext.OutputPath != "" {
				file.SetOutput(relOutput)
			}

			file.Success()
			return nil
		})

		t.Start()
		t.Close()

		// should have already errored
		if err != nil {
			return fmt.Errorf("error recursively templating: %w", err)
		}
	} else if isFile {
		err := handleFile(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
