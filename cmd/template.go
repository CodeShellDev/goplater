package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/codeshelldev/goplater/internals/template"
	tmplContext "github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/internals/template/core"
	"github.com/codeshelldev/goplater/utils/fsutils"

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

	templateCmd.Flags().BoolVarP(&templateOptions.Recursive, "recursive", "r", false, "recusively template files")
	templateCmd.Flags().BoolVarP(&templateOptions.Flatten, "flatten", "f", false, "flatten output path: don't preserve source folder structure")

	templateCmd.Flags().StringVarP(&templateOptions.Output, "output", "o", ".", "output path for templated files")
	templateCmd.Flags().StringVarP(&templateOptions.Source, "source", "s", ".", "source path for local files")

	templateCmd.Flags().StringSliceVar(&templateOptions.AllowedOutputFolders, "allowed-output-folders", []string{}, "allowed output folders")

	templateCmd.Flags().StringSliceVarP(&templateOptions.Match, "match", "m", []string{".*"}, "regexes used for matching files to process")

	templateCmd.Flags().StringSliceVarP(&templateOptions.Whitespace, "whitespace", "w", []string{"l","t"}, "remove whitespace from files")

	templateCmd.Flags().BoolVarP(&templateOptions.Verbose, "verbose", "v", false, "print additional information")
	templateCmd.Flags().BoolVarP(&templateOptions.Supress, "ignore-errors", "i", false, "ignore / supress errors")
}

func validate(cmd *cobra.Command, args []string) error {
	if len(args) <= 0 && !templateOptions.Recursive {
		return errors.New("not enough args")
	} else if len(args) > 0 {
		_, err := os.Stat(args[0])

		if errors.Is(err, os.ErrNotExist) {
			return errors.New("invalid context path")
		}
	}

	return nil
}

func run(cmd *cobra.Command, args []string) {
	templateContext := tmplContext.New(templateOptions)

	allowedOutputFolders := append(templateOptions.AllowedOutputFolders, templateOptions.Output)

	templateOptions.AllowedOutputFolders = []string{}

	for _, path := range allowedOutputFolders {
		absPath, err := fsutils.Relative(templateOptions.Source, path)

		if err != nil {
			fmt.Println("could not convert " + path +  " to absolute: " + err.Error())
		}

		templateOptions.AllowedOutputFolders = append(templateOptions.AllowedOutputFolders, absPath)
	}

	templateContext.Path = args[0]

	core.Renderer = &templater
	core.Matcher = &templater

	templater.Run(templateContext)
}