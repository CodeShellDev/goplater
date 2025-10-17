package cmd

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/codeshelldev/goplater/utils/fsutils"
	"github.com/codeshelldev/goplater/utils/get"
	"github.com/codeshelldev/goplater/utils/templating"
	"github.com/spf13/cobra"
)

var context string = "."
var output string = "."
var localContext string = "."

var whitespace []string
var match []string

var verbose bool
var recursive bool
var preserveStructure bool

var templateCmd = &cobra.Command{
	Use:   "template",
    Short: "Template files",
	Args: validate,
    Long:  `Template files by using local or remote files.`,
	Run: run,
}

func validate(cmd *cobra.Command, args []string) error {
	if len(args) <= 0 && !recursive {
		return errors.New("not enough args")
	} else if len(args) > 0 {
		_, err := os.Stat(args[0])

		if errors.Is(err, os.ErrNotExist) {
			return errors.New("invalid context path")
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(templateCmd)

	templateCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "recusively template files")
	templateCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "print additional information")
	templateCmd.Flags().BoolVarP(&preserveStructure, "preserve-struct", "p", true, "preserves source folder structure in output path")

	templateCmd.Flags().StringVarP(&output, "output", "o", ".", "output path for templated files")
	templateCmd.Flags().StringVarP(&localContext, "source", "s", ".", "source path for local files")

	templateCmd.Flags().StringSliceVarP(&match, "match", "m", []string{".*"}, "regex match for templating")

	templateCmd.Flags().StringSliceVarP(&whitespace, "whitespace", "w", []string{"l","t"}, "remove whitespace from files")
}

func run(cmd *cobra.Command, args []string) {
	context = args[0]

	fullPath, _ := filepath.Abs(context)

	isDir := fsutils.IsDir(fullPath)
	isFile := fsutils.IsFile(fullPath)

	if isDir {
		localContext = context

		filepath.WalkDir(context, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() {
				handleFile(path)
			} else if path != context && !recursive {
				return filepath.SkipDir
			}

			return nil
		})
	} else if isFile {
		handleFile(context)
	}
}

func matchFile(path string) bool {
	fileName := filepath.Base(path)

	for _, reStr := range match {
		re, err := regexp.Compile(reStr)

		if err == nil {
			if re.MatchString(fileName) {
				return true
			}
		}
	}

	return false
}

func handleFile(relativePath string) {
	if !matchFile(relativePath) {
		if verbose {
			fmt.Println("skipped", relativePath)
		}

		return
	}

	if verbose {
		fmt.Println("templating", relativePath)
	}

	err := templateFile(relativePath)

	if err != nil {
		fmt.Println("error templating:", err.Error())
	}
}

func templateFile(relativePath string) error {
	if !fsutils.IsFile(relativePath) {
		return os.ErrNotExist
	}

	data, err := os.ReadFile(relativePath)

	if err != nil {
		return err
	}

	if data == nil {
		return errors.New("empty file")
	}

	normalized := normalize(string(data))

	tmplStr, err := templateStr("main", normalized, nil)

	if err != nil {
		return err
	}

	return handleFileWrite(tmplStr, relativePath)
}

func handleFileWrite(content, relativePath string) error {
	fullPath := fsutils.ResolveOutput(relativePath, output, preserveStructure)

	if verbose {
		fmt.Println("writing to", fullPath)
	}

	dir := filepath.Dir(fullPath)
	err := os.MkdirAll(dir, 0755);

	if err != nil {
		return err
	}

	err = os.WriteFile(fullPath, []byte(content), 0644)

	if err != nil {
		return err
	}

	return nil
}

func normalize(content string) string {
	normalizeLocal, err := templating.ReplaceTemplatePrefix(content, `#://`, "#://")

	if err == nil {
		content = normalizeLocal
	}

	normalizeRemote, err := templating.ReplaceTemplatePrefix(content, `@://`, "@://")

	if err == nil {
		content = normalizeRemote
	}

	return content
}

func removeWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func templateGet(key string) any {
	var res string

	switch(key[:4]) {
	case "@://":
		res = get.Remote(key[4:])
	case "#://":
		res = get.Local(key[4:], localContext)
	}

	if slices.Contains(whitespace, "l") {
		res = strings.TrimLeftFunc(res, removeWhitespace)
	}

	if slices.Contains(whitespace, "t") {
		res = strings.TrimRightFunc(res, removeWhitespace)
	}

	return res
}

func templateStr(name, str string, variables map[string]any) (string, error) {
	tmplStr, err := templating.AddTemplateFunc(str, "get")

	if err != nil {
		return str, err
	}

	templt := templating.CreateTemplateWithFunc(name, template.FuncMap{
		"get": templateGet,
	})

	tmplStr, err = templating.ParseTemplate(templt, tmplStr, variables)

	if err != nil {
		return str, err
	}

	return tmplStr, nil
}