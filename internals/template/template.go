package template

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/pkg/gplt/utils/tmplerrors"
	"github.com/codeshelldev/goplater/utils/fsutils"
	"github.com/codeshelldev/goplater/utils/pretty"
	"github.com/codeshelldev/goplater/utils/prettyerrors"
	"github.com/codeshelldev/goplater/utils/tree"
)

type Templater struct{}

func (Templater) Run(tmplContext *context.TemplateContext) error {
	return run(tmplContext)
}

func New() *Templater {
	return &Templater{}
}

func run(tmplContext *context.TemplateContext) error {
	isDir := fsutils.IsDir(tmplContext.Shared.Path)
	isFile := fsutils.IsFile(tmplContext.Shared.Path)
	
	tmplContext.Local.Invoker = tmplContext.Shared.Path

	if isDir {
		t := tree.NewTree(os.Stdout, tmplContext.Shared.Path)

		t.Start()

		err := filepath.WalkDir(tmplContext.Shared.Path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				if path != tmplContext.Shared.Path && !tmplContext.Options.Recursive {
					return filepath.SkipDir
				}
				return nil
			}

			// skip symlinks
			if d.Type()&fs.ModeSymlink != 0 {
				return nil
			}

			absPath, err := filepath.Abs(path)

			if err != nil {
				return err
			}

			newContext := *tmplContext
			newContext.Local.Invoker = absPath
			newContext.Shared = &context.Shared{Path: absPath}

			ok, err := matchFile(&newContext)

			if err != nil {
				return err
			}

			if !ok {
				return nil
			}

			relPath, err := filepath.Rel(tmplContext.Shared.Path, path)
			if err != nil {
				return err
			}

			file := t.Add(relPath)

			err = handleFile(&newContext)
			if err != nil {
				tmplErr, ok := errors.AsType[*tmplerrors.TemplateError](err)
		
				if ok {
					prettyErr := prettyerrors.Convert(tmplErr)
					color := pretty.ColorsEnabled(os.Stdout, false)
					file.Fail(errors.New(tmplErr.CleanContext()), prettyErr.ErrorSnippet(color))
				} else {
					file.Fail(err)
				}

				fmt.Println()

				if !tmplContext.Options.IgnoreErrors {
					return err
				}

				return nil
			}

			if newContext.Shared.OutputPath != "" {
				relOutput, err := fsutils.SafeRelative(tmplContext.Shared.Path, newContext.Shared.OutputPath)

				if err != nil {
					// not under source: show abs path instead
					relOutput = newContext.Shared.OutputPath
				}

				file.SetOutput(relOutput)
			}

			file.Success()
			return nil
		})

		t.Close()

		// should have already errored
		if err != nil {
			return fmt.Errorf("error recursively templating: %w", err)
		}
	} else if isFile {
		err := handleFile(tmplContext)
		if err != nil {
			return err
		}
	}

	return nil
}
