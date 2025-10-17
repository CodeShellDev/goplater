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

var fileContext string = "."
var outputContext string = "."
var sourceContext string = "."

var whitespace []string
var match []string

var recursive bool
var flatten bool

type TemplateContext struct {
	invoker 	string
	name 		string
	// ...
}

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
	templateCmd.Flags().BoolVarP(&flatten, "flatten", "f", false, "flatten output path: don't preserve source folder structure")

	templateCmd.Flags().StringVarP(&outputContext, "output", "o", ".", "output path for templated files")
	templateCmd.Flags().StringVarP(&sourceContext, "source", "s", ".", "source path for local files")

	templateCmd.Flags().StringSliceVarP(&match, "match", "m", []string{".*"}, "regex match for templating")

	templateCmd.Flags().StringSliceVarP(&whitespace, "whitespace", "w", []string{"l","t"}, "remove whitespace from files")
}

func run(cmd *cobra.Command, args []string) {
	fileContext = args[0]

	fullPath, _ := filepath.Abs(fileContext)

	isDir := fsutils.IsDir(fullPath)
	isFile := fsutils.IsFile(fullPath)

	if isDir {
		filepath.WalkDir(fileContext, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() {
				handleFile(path)
			} else if path != fileContext && !recursive {
				return filepath.SkipDir
			}

			return nil
		})
	} else if isFile {
		handleFile(fileContext)
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

	if err != nil && !supress {
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

	tmplStr, err := templateContent(string(data), TemplateContext{invoker: relativePath, name: relativePath})

	if err != nil {
		return err
	}

	return handleFileWrite(tmplStr, relativePath)
}

func templateContent(content string, context TemplateContext) (string, error) {
	normalized := normalize(content)

	tmplStr, err := templateStr(normalized, context, nil)

	return tmplStr, err
}

func handleFileWrite(content, relativePath string) error {
	fullPath := fsutils.ResolveOutput(relativePath, outputContext, !flatten)

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

func templateGet(key string, context TemplateContext) any {
	var res string

	keyWithoutPrefix := key[4:]

	switch(key[:4]) {
	case "@://":
		res = get.Remote(keyWithoutPrefix)
	case "#://":
		isRel := strings.HasPrefix(keyWithoutPrefix, "./") || strings.HasPrefix(keyWithoutPrefix, "../")

		filePathAbs, _ := filepath.Abs(context.invoker)

		if isRel {
			res = get.Local(keyWithoutPrefix, filepath.Dir(filePathAbs))
		} else {
			res = get.Local(keyWithoutPrefix, sourceContext)
		}

		keyWithoutPrefix = fsutils.Relative(filepath.Dir(filePathAbs), keyWithoutPrefix)

		fmt.Println(keyWithoutPrefix)
	}

	if slices.Contains(whitespace, "l") {
		res = strings.TrimLeftFunc(res, removeWhitespace)
	}

	if slices.Contains(whitespace, "t") {
		res = strings.TrimRightFunc(res, removeWhitespace)
	}

	fileName := filepath.Base(keyWithoutPrefix)

	if matchFile(fileName) {
		res, _ = templateContent(res, TemplateContext{invoker: keyWithoutPrefix, name: keyWithoutPrefix})
	}

	return res
}

func templateStr(str string, context TemplateContext, variables map[string]any) (string, error) {
	tmplStr, err := templating.AddTemplateFunc(str, "get")

	if err != nil {
		return str, err
	}

	templt := templating.CreateTemplateWithFunc(context.name, template.FuncMap{
		"get": func (key string) any {
			return templateGet(key, context)	
		},
	})

	if supress {
		templt = templating.AddTemplateOption(templt, "missingkey=zero")
	}

	tmplStr, err = templating.ParseTemplate(templt, tmplStr, variables)

	if err != nil {
		return str, err
	}

	return tmplStr, nil
}