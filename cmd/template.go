package cmd

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/codeshelldev/goplater/internals/template"
	tmplContext "github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/internals/template/core"
	"github.com/codeshelldev/goplater/utils/printer"

	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template",
    Short: "Template files",
	Args: validate,
    Long:  `Template files by using local or remote files.`,
	Run: run,
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

	templateCmd.Flags().StringSliceVar(&templateOptions.AllowedOutputs, "allowed-outputs", []string{}, "Allowed output directories, including their subdirectories")
	templateCmd.Flags().StringSliceVarP(&templateOptions.Match, "match", "m", []string{".*"}, "Regular expressions used to match files for processing")

	templateCmd.Flags().BoolVar(&templateOptions.DryRun, "dry-run", false, "Show what would be written, change nothing")
	templateCmd.Flags().BoolVarP(&templateOptions.IgnoreErrors, "ignore-errors", "i", false, "Continue processing when errors occur")
}

func validate(cmd *cobra.Command, args []string) error {
	if len(args) <= 0 && !templateOptions.Recursive {
		return errors.New("not enough args")
	} else if len(args) > 0 {
		_, err := os.Stat(args[0])

		if errors.Is(err, os.ErrNotExist) {
			printer.Error("invalid context path: " + args[0])
		}
	}

	return nil
}

func run(cmd *cobra.Command, args []string) {
	templateContext := tmplContext.New(templateOptions)

	allowedOutputFolders := append(templateOptions.AllowedOutputs, templateOptions.Output)

	templateOptions.AllowedOutputs = []string{}

	for _, path := range allowedOutputFolders {
		absPath := path

		if !filepath.IsAbs(path) {
			var err error
			absPath, err = filepath.Abs(path)

			if err != nil {
				printer.Error(err.Error())
			}
		}

		templateOptions.AllowedOutputs = append(templateOptions.AllowedOutputs, absPath)
	}

	outputPathAbs, err := filepath.Abs(templateOptions.Output)

	if err != nil {
		printer.Error(err.Error())
	}

	templateOptions.Output = outputPathAbs

	if templateOptions.Root != "" {
		rootPathAbs, err := filepath.Abs(templateOptions.Root)

		if err != nil {
			printer.Error(err.Error())
		}

		templateOptions.Root = rootPathAbs
	}

	sourcePathAbs, err := filepath.Abs(args[0])

	if err != nil {
		printer.Error(err.Error())
	}

	templateContext.Path = sourcePathAbs

	core.Renderer = &templater
	core.Matcher = &templater

	templater.Run(templateContext)
}