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
	"github.com/codeshelldev/goplater/utils/outputresolver"
	"github.com/codeshelldev/goplater/utils/pretty"
	"github.com/codeshelldev/goplater/utils/prettyerrors"
	"github.com/codeshelldev/goplater/utils/tree"
)

func Run(tmplContext *context.TemplateContext) error {
	ctxPath := tmplContext.Shared.Path
	opts := tmplContext.Options

	isDir := fsutils.IsDir(ctxPath)

	if isDir {
		err := outputresolver.Validate(ctxPath, opts.Output)
		if err != nil {
			return err
		}
	}

	isFile := fsutils.IsFile(ctxPath)
	
	tmplContext.Local.Invoker = ctxPath

	if isDir {
		t := tree.NewTree(os.Stdout, ctxPath)

		t.Start()

		err := filepath.WalkDir(ctxPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// always walk ctx path, even if dotfolder
			if path == ctxPath {
				return nil
			}

			relPath, err := filepath.Rel(ctxPath, path)
			if err != nil {
				return err
			}
			relPath = filepath.ToSlash(relPath)

			// exclude hidden dotfiles and dotfolders
			if !opts.IncludeHidden && isHidden(d.Name()) {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			if isExcluded(opts, relPath) {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			if d.IsDir() {
				if !opts.Recursive {
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

			if !matchRelative(opts, relPath) {
				return nil
			}

			if !opts.IncludeBinary {
				isBin, err := fsutils.IsBinary(path)
				if err != nil {
					return err
				}

				if isBin {
					return nil
				}
			}

			newContext := *tmplContext
			newContext.Local.Invoker = absPath
			newContext.Shared = &context.Shared{Path: absPath}

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

				if !opts.IgnoreErrors {
					return err
				}

				return nil
			}

			if newContext.Shared.OutputPath != "" {
				relOutput, err := fsutils.SafeRelative(ctxPath, newContext.Shared.OutputPath)

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
