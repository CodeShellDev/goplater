package template

import (
	"io/fs"
	"path/filepath"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/utils/fsutils"
)

type Templater struct{}

func (Templater) Run(context context.TemplateContext) {
	run(context)
}

func New() *Templater {
	return &Templater{}
}

func run(context context.TemplateContext) {
	fullPath, _ := filepath.Abs(context.Path)
	
	isDir := fsutils.IsDir(fullPath)
	isFile := fsutils.IsFile(fullPath)

	context.Invoker = context.Path

	if isDir {
		filepath.WalkDir(context.Path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() {
				newContext := context
				newContext.Path = path
				newContext.Invoker = path

				handleFile(newContext)
			} else if path != context.Path && !context.Options.Recursive {
				return filepath.SkipDir
			}

			return nil
		})
	} else if isFile {
		handleFile(context)
	}
}