package template

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/utils/fsutils"
	"github.com/codeshelldev/goplater/utils/printer"
	"github.com/codeshelldev/goplater/utils/tree"
)

type Templater struct{}

func (Templater) Run(ctx *context.TemplateContext) {
	run(ctx)
}

func New() *Templater {
	return &Templater{}
}

func run(ctx *context.TemplateContext) {
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

			absPath, _ := filepath.Abs(path)

			newContext := *ctx
			newContext.Path = absPath
			newContext.Invoker = absPath

			if !matchFile(&newContext) {
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
				return err
			}

			relOutput, err := fsutils.SafeRelative(ctx.Options.Source, newContext.OutputPath)

			if err != nil {
				file.Fail(err)
				return err
			}

			if newContext.OutputPath != "" {
				file.SetOutput(relOutput)
			}

			file.Success()
			return nil
		})

		// should have already errored
		if err != nil {
			printer.Error(err.Error())
		}

		t.Start()
		t.Close()
	} else if isFile {
		err := handleFile(ctx)
		if err != nil {
			printer.Error(err.Error())
		}
	}
}
