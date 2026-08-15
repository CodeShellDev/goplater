package template

import (
	"io/fs"
	"path/filepath"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/utils/fsutils"
)

type Templater struct{}

func (Templater) Run(ctx *context.TemplateContext) {
	run(ctx)
}

func New() *Templater {
	return &Templater{}
}

func run(ctx *context.TemplateContext) {
	fullPath, _ := filepath.Abs(ctx.Path)
	
	isDir := fsutils.IsDir(fullPath)
	isFile := fsutils.IsFile(fullPath)

	ctx.Invoker = ctx.Path

	if isDir {
		filepath.WalkDir(ctx.Path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() {
				newContext := ctx
				newContext.Path = path
				newContext.Invoker = path

				handleFile(newContext)
			} else if path != ctx.Path && !ctx.Options.Recursive {
				return filepath.SkipDir
			}

			return nil
		})
	} else if isFile {
		handleFile(ctx)
	}
}