package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/codeshelldev/goplater/internals/template"
	tmplContext "github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/internals/template/core"

	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template",
    Short: "Template files",
	Args: validate,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveDefault
	},
    Long:  "Template files by using local or remote files.",
	RunE: run,
}

var templateOptions = &tmplContext.TemplateOptions{}

var templater template.Templater = *template.New()

func init() {
	rootCmd.AddCommand(templateCmd)

	templateCmd.Flags().BoolVarP(&templateOptions.Recursive, "recursive", "r", false, "Recursively process files")

	templateCmd.Flags().BoolVarP(&templateOptions.Flatten, "flatten", "f", false, "Flatten output paths without preserving source directory structure")

	templateCmd.Flags().StringVar(&templateOptions.Root, "root", "", "Source root, output keeps paths relative to it (default: inferred)")
	templateCmd.Flags().StringVarP(&templateOptions.Output, "output", "o", ".", "Output directory for templated files")
	templateCmd.Flags().StringVarP(&templateOptions.Source, "source", "s", ".", "Source directory containing files")

	templateCmd.Flags().IntVar(&templateOptions.TrimMaxBlankLines, "trim-max-blank-lines", -1, "experimental: Max blank lines kept around non-output actions")

	templateCmd.Flags().StringSliceVar(&templateOptions.AllowedOutputs, "allowed-outputs", []string{}, "Allowed output directories, including their subdirectories (always includes output path)")
	templateCmd.Flags().StringSliceVar(&templateOptions.AllowedInputs, "allowed-inputs", []string{}, "Allowed input directories, including their subdirectories (always includes source path)")
	templateCmd.Flags().StringSliceVarP(&templateOptions.Match, "match", "m", []string{".*"}, "Regular expressions used to match files for processing")

	templateCmd.Flags().BoolVar(&templateOptions.DryRun, "dry-run", false, "Show what would be written, change nothing")
	templateCmd.Flags().BoolVarP(&templateOptions.IgnoreErrors, "ignore-errors", "i", false, "Continue processing when errors occur")
}

func validate(cmd *cobra.Command, args []string) error {
	err := cobra.ExactArgs(1)(cmd, args)
    if err != nil {
        return err
    }

	_, err = os.Stat(args[0])

	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("context path %q does not exist: %w", args[0], err)
	}

	return nil
}

func makePathsAbsolute(in []string) ([]string, error) {
	out := []string{}

	for _, path := range in {
		absPath := path

		if !filepath.IsAbs(path) {
			var err error
			absPath, err = filepath.Abs(path)

			if err != nil {
				return nil, fmt.Errorf("could not make %q absolute: %w", path, err)
			}
		}

		out = append(out, absPath)
	}

	return out, nil
}

func run(cmd *cobra.Command, args []string) error {
	templateContext := tmplContext.New(templateOptions)

	allowedOutputFolders := append(templateOptions.AllowedOutputs, templateOptions.Output)

	var err error
	templateOptions.AllowedOutputs, err = makePathsAbsolute(allowedOutputFolders)

	if err != nil {
		return err
	}

	allowedInputFolders := append(templateOptions.AllowedInputs, templateOptions.Source)

	templateOptions.AllowedInputs, err = makePathsAbsolute(allowedInputFolders)

	if err != nil {
		return err
	}

	outputPathAbs, err := filepath.Abs(templateOptions.Output)

	if err != nil {
		return fmt.Errorf("could not make output path %q absolute: %w", templateOptions.Output, err)
	}

	templateOptions.Output = outputPathAbs

	if templateOptions.Root != "" {
		rootPathAbs, err := filepath.Abs(templateOptions.Root)

		if err != nil {
			return fmt.Errorf("could not make root path %q absolute: %w", templateOptions.Root, err)
		}

		templateOptions.Root = rootPathAbs
	}

	sourcePathAbs, err := filepath.Abs(args[0])

	if err != nil {
		return fmt.Errorf("could not make source path %q absolute: %w", templateOptions.Source, err)
	}

	templateContext.Path = sourcePathAbs

	core.Renderer = &templater
	core.Matcher = &templater

	templater.Run(templateContext)
	return nil
}