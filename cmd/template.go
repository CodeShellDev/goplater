package cmd

import (
	"errors"
	"os"

	"github.com/codeshelldev/goplater/internals/template"
	tmplContext "github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/internals/template/core"

	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template",
    Short: "Template files",
	Args: validate,
    Long:  `Template files by using local or remote files.`,
	Run: run,
}

var templateContext tmplContext.TemplateContext = tmplContext.New()

var templater template.Templater = *template.New()

func init() {
	rootCmd.AddCommand(templateCmd)

	templateCmd.Flags().BoolVarP(&templateContext.Options.Recursive, "recursive", "r", false, "recusively template files")
	templateCmd.Flags().BoolVarP(&templateContext.Options.Flatten, "flatten", "f", false, "flatten output path: don't preserve source folder structure")

	templateCmd.Flags().StringVarP(&templateContext.Options.Output, "output", "o", ".", "output path for templated files")
	templateCmd.Flags().StringVarP(&templateContext.Options.Source, "source", "s", ".", "source path for local files")

	templateCmd.Flags().StringSliceVarP(&templateContext.Options.Match, "match", "m", []string{".*"}, "regex match for templateutils")

	templateCmd.Flags().StringSliceVarP(&templateContext.Options.Whitespace, "whitespace", "w", []string{"l","t"}, "remove whitespace from files")

	templateCmd.Flags().BoolVarP(&templateContext.Options.Verbose, "verbose", "v", false, "print additional information")
	templateCmd.Flags().BoolVarP(&templateContext.Options.Supress, "ignore-errors", "i", false, "ignore / supress errors")
}

func validate(cmd *cobra.Command, args []string) error {
	if len(args) <= 0 && !templateContext.Options.Recursive {
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
	templateContext.Path = args[0]

	core.Renderer = &templater
	core.Matcher = &templater

	templater.Run(templateContext)
}