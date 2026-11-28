package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/codeshelldev/goplater/internals/template"
	tmplContext "github.com/codeshelldev/goplater/internals/template/context"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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

var opts = &tmplContext.TemplateOptions{}

func init() {
	rootCmd.AddCommand(templateCmd)

	templateCmd.Flags().BoolVarP(&opts.Recursive, "recursive", "r", false, "Recursively process files")

	templateCmd.Flags().BoolVarP(&opts.Flatten, "flatten", "f", false, "Flatten output paths without preserving source directory structure")

	templateCmd.Flags().StringVar(&opts.Root, "root", "", "Source root, output keeps paths relative to it (default: inferred)")
	templateCmd.Flags().StringVarP(&opts.Output, "output", "o", ".", "Output directory for templated files")
	templateCmd.Flags().StringVarP(&opts.Source, "source", "s", ".", "Source directory, ~/ is relative to it")

	templateCmd.Flags().IntVar(&opts.TrimMaxBlankLines, "trim-max-blank-lines", -1, "Experimental: Max blank lines kept around non-output actions")

	templateCmd.Flags().StringSliceVar(&opts.AllowedOutputs, "allowed-outputs", []string{}, "Allowed output directories, including their subdirectories (always includes output path)")
	templateCmd.Flags().StringSliceVar(&opts.AllowedInputs, "allowed-inputs", []string{}, "Allowed input directories, including their subdirectories (always includes context path)")

	templateCmd.Flags().StringArrayVarP(&opts.Match, "match", "m", []string{".*"}, "Only template files/folders whose relative path matches regex")
	templateCmd.Flags().StringArrayVar(&opts.Exclude, "exclude", []string{}, "Skip files/folders whose relative path matches regex")
	templateCmd.Flags().BoolVar(&opts.IncludeHidden, "hidden", false, "Include dotfiles and dotfolders")
	templateCmd.Flags().BoolVar(&opts.IncludeBinary, "binary", false, "Include binary files")

	templateCmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Show what would be written, change nothing")
	templateCmd.Flags().BoolVarP(&opts.IgnoreErrors, "ignore-errors", "i", false, "Continue processing when errors occur")

	templateCmd.Flags().SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		switch name {
		case "allowed-output":
			name = "allowed-outputs"
		case "allowed-input":
			name = "allowed-inputs"
		}
		return pflag.NormalizedName(name)
	})
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

func makePathAbsolute(path string) (string, error) {
	if !filepath.IsAbs(path) {
		var err error
		path, err = filepath.Abs(path)

		if err != nil {
			return "", fmt.Errorf("could not make %q absolute: %w", path, err)
		}
	}

	return path, nil
}

func makePathsAbsolute(in []string) ([]string, error) {
	out := []string{}

	for _, path := range in {
		absPath, err := makePathAbsolute(path)

		if err != nil {
			return nil, err
		}

		out = append(out, absPath)
	}

	return out, nil
}

func parseStringSlice(slice []string) []string {
	for i, item := range slice {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			slice[i] = trimmed
		}
	}
	return slice
}

func run(cmd *cobra.Command, args []string) error {
	tmplContext := tmplContext.New(opts)

	// ctx path
	ctxPathAbs, err := filepath.Abs(args[0])

	if err != nil {
		return fmt.Errorf("could not make context path %q absolute: %w", args[0], err)
	}

	tmplContext.Shared.Path = ctxPathAbs

	// allowed dirs
	allowedOutputFolders := append(parseStringSlice(opts.AllowedOutputs), opts.Output)

	opts.AllowedOutputs, err = makePathsAbsolute(allowedOutputFolders)

	if err != nil {
		return err
	}

	allowedInputFolders := append(parseStringSlice(opts.AllowedInputs), ctxPathAbs)

	opts.AllowedInputs, err = makePathsAbsolute(allowedInputFolders)

	if err != nil {
		return err
	}

	// output
	outputPathAbs, err := filepath.Abs(opts.Output)

	if err != nil {
		return fmt.Errorf("could not make output path %q absolute: %w", opts.Output, err)
	}

	opts.Output = outputPathAbs

	// root
	if opts.Root != "" {
		rootPathAbs, err := filepath.Abs(opts.Root)

		if err != nil {
			return fmt.Errorf("could not make root path %q absolute: %w", opts.Root, err)
		}

		opts.Root = rootPathAbs
	}

	opts.MatchRe, err = template.CompilePatterns(opts.Match)
	if err != nil {
		return err
	}

	opts.ExcludeRe, err = template.CompilePatterns(opts.Exclude)
	if err != nil {
		return err
	}

	return template.Run(tmplContext)
}