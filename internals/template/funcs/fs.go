package funcs

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/codeshelldev/goplater/internals/template/context"
	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
	"github.com/codeshelldev/goplater/utils/fsutils"
)

func isPathAllowed(path string, tmplContext *context.TemplateContext) bool {
	allowed := false

	for _, try := range tmplContext.Options.AllowedOutputFolders {
		if fsutils.IsInside(path, try) {
			allowed = true
		}
	}

	return allowed
}

// Reads and templates file.
//
// @param path string
// @returns string
//
// @example
//   +{{ read "/path/to/file" }}
var readFunc = modules.NewFunc("read", read)

func read(rt *templating.Runtime, ctx *templating.Context, path string) string {
	str, ctx := readHandler(ctx, path)

	str, err := rt.GetEngine().Execute(":read:" + path, str, nil, rt.GetEngineOptions(), ctx)

	if err != nil {
		panic("could not read " + path + ": " + err.Error())
	}

	return str
}

// Reads file (without templating).
//
// @param path string
// @returns string
//
// @example
//   +{{ readRaw "/path/to/file" }}
var readRawFunc = modules.NewFunc("readRaw", readRaw)

func readRaw(_ *templating.Runtime, ctx *templating.Context, path string) string {
	tmplContext := ctx.Get(context.TemplateContextKey).(*context.TemplateContext)

	filePathAbs := resolvePath(*tmplContext, path)

	res, err := readFile(filePathAbs)

	if err != nil {
		res = err.Error()
	}

	return res
}

// Reads file and passes arguments to it.
//
// @param path string
// @param args any
// @returns string
//
// @example
//   +{{ readArgs "/path/to/file" (sliceCreate "my" "args") }}
var readArgsFunc = modules.NewFunc("readArgs", readArgs)

func readArgs(rt *templating.Runtime, ctx *templating.Context, path string, args any) string {
	str, newContext := readHandler(ctx, path)

	data := map[string]any{
		"args": args,
	}

	str, err := rt.GetEngine().Execute(":read:" + path, str, data, rt.GetEngineOptions(), newContext)

	if err != nil {
		panic("could not read " + path + ": " + err.Error())
	}

	return str
}

func readHandler(ctx *templating.Context, path string) (string, *templating.Context) {
	tmplContext := ctx.Get(context.TemplateContextKey).(*context.TemplateContext)

	filePathAbs := resolvePath(*tmplContext, path)

	res, err := readFile(filePathAbs)

	if err != nil {
		res = err.Error()
	}

	tmplContext.Invoker = filePathAbs

	newContext := &templating.Context{}

	newContext.Set(context.TemplateContextKey, tmplContext)

	return res, newContext
}

// Writes to a file path.
//
// @param path string
// @param content string
// @returns error
//
// @example
//   +{{ write "/path/to/file" "Hello" }}
var writeFunc = modules.NewFunc("write", write)

func write(_ *templating.Runtime, ctx *templating.Context, path string, content string) error {
	tmplContext := ctx.Get(context.TemplateContextKey).(*context.TemplateContext)

	filePathAbs := resolvePath(*tmplContext, path)

	if !isPathAllowed(filePathAbs, tmplContext) {
		panic("writing to " + filePathAbs + " is not allowed as it is not inside of the allowed scope")
	}

	err := os.WriteFile(filePathAbs, []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}

// Creates directory at path.
//
// @param path string
// @returns error
//
// @example
//   +{{ mkdir "/path/to/somewhere/" }}
var mkdirFunc = modules.NewFunc("mkdir", mkdir)

func mkdir(_ *templating.Runtime, ctx *templating.Context, path string) string {
	tmplContext := ctx.Get(context.TemplateContextKey).(*context.TemplateContext)

	folderPathAbs := resolvePath(*tmplContext, path)

	if !isPathAllowed(folderPathAbs, tmplContext) {
		panic("creating folder " + folderPathAbs + " is not allowed as it is not inside of the allowed scope")
	}

	err := os.MkdirAll(folderPathAbs, 0755)
	if err != nil {
		panic(err.Error())
	}

	return ""
}

// Appends to an existing file.
//
// @param path string
// @param content string
// @returns error
//
// @example
//   +{{ appendFile "/path/to/somewhere/" "Goodbye" }}
var appendFileFunc = modules.NewFunc("appendFile", appendFile)

func appendFile(_ *templating.Runtime, ctx *templating.Context, path string, content string) string {
	tmplContext := ctx.Get(context.TemplateContextKey).(*context.TemplateContext)

	filePathAbs := resolvePath(*tmplContext, path)

	if !isPathAllowed(filePathAbs, tmplContext) {
		panic("appending to " + filePathAbs + " is not allowed as it is not inside of the allowed scope")
	}

	f, err := os.OpenFile(filePathAbs, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		panic(err.Error())
	}

	return ""
}

// Returns if a file or folder exists at a given path.
//
// @param path string
// @returns bool
//
// @example
//   +{{ if (fsExists "/path/to/somewhere") }}
//		Path exists!
//	 +{{ end }}
var fsExistsFunc = modules.NewFunc("fsExists", fsExists)

func fsExists(_ *templating.Runtime, ctx *templating.Context, path string) bool {
	tmplContext := ctx.Get(context.TemplateContextKey).(*context.TemplateContext)

	filePathAbs := resolvePath(*tmplContext, path)

	_, err := os.Stat(filePathAbs)
	if err != nil {
		return false
	}

	return true
}

// Returns if a file exists at a given path.
//
// @param path string
// @returns bool
//
// @example
//   +{{ if (isFile "/path/to/file") }}
//		Path is a file!
//	 +{{ end }}
var isFileFunc = modules.NewFunc("isFile", isFile)

func isFile(_ *templating.Runtime, ctx *templating.Context, path string) bool {
	tmplContext := ctx.Get(context.TemplateContextKey).(*context.TemplateContext)

	filePathAbs := resolvePath(*tmplContext, path)

	info, err := os.Stat(filePathAbs)
	if err != nil {
		return false
	}

	return info.Mode().IsRegular()
}

// Returns if a folder exists at a given path.
//
// @param path string
// @returns bool
//
// @example
//   +{{ if (isDir "/path/to/somewhere") }}
//		Path is a folder!
//	 +{{ end }}
var isDirFunc = modules.NewFunc("isDir", isDir)

func isDir(_ *templating.Runtime, ctx *templating.Context, path string) bool {
	tmplContext := ctx.Get(context.TemplateContextKey).(*context.TemplateContext)

	filePathAbs := resolvePath(*tmplContext, path)

	info, err := os.Stat(filePathAbs)
	if err != nil {
		return false
	}

	return info.IsDir()
}

// Returns all files and folders inside of a given directory.
//
// @param path string
// @returns []string
// @returns error
//
// @example
//   +{{ join ", " (listDir "/path/to") }}
// @output
//	 /path/to/file1.txt, /path/to/file2.txt, /path/to/folder
var listDirFunc = modules.NewFunc("listDir", listDir)

func listDir(_ *templating.Runtime, ctx *templating.Context, path string) ([]string, error) {
	tmplContext := ctx.Get(context.TemplateContextKey).(*context.TemplateContext)

	folderPathAbs := resolvePath(*tmplContext, path)

	entries, err := os.ReadDir(folderPathAbs)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(entries))

	for _, entry := range entries {
		result = append(result, entry.Name())
	}

	return result, nil
}

// Returns all files and folders recursively under a given directory.
//
// @param path string
// @returns []string
// @returns error
//
// @example
//   +{{ join ", " (walkDir "/path/to") }}
// @output
//	 /path/to/file1.txt, /path/to/file2.txt, /path/to/folder, /path/to/folder/file3.txt
var walkDirFunc = modules.NewFunc("walkDir", walkDir)

func walkDir(_ *templating.Runtime, ctx *templating.Context, path string) ([]string, error) {
	tmplContext := ctx.Get(context.TemplateContextKey).(*context.TemplateContext)

	folderPathAbs := resolvePath(*tmplContext, path)

	paths := []string{}

	err := filepath.WalkDir(folderPathAbs, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		paths = append(paths, path)
		return nil
	})

	return paths, err
}


// Removes a file or folder at a given path.
//
// @param path string
// @returns error
//
// @example
//   +{{ fsRemove "/path/to/file/or/folder" }}
var fsRemoveFunc = modules.NewFunc("fsRemove", fsRemove)

func fsRemove(_ *templating.Runtime, ctx *templating.Context, path string) string {
	tmplContext := ctx.Get(context.TemplateContextKey).(*context.TemplateContext)

	filePathAbs := resolvePath(*tmplContext, path)

	if !isPathAllowed(filePathAbs, tmplContext) {
		panic("removing " + filePathAbs + " is not allowed as it is not inside of the allowed scope")
	}

	err := os.RemoveAll(path)
	if err != nil {
		panic(err.Error())
	}

	return ""
}

// Joins paths together.
//
// @param paths []string
// @returns string
var joinPathFunc = modules.NewFunc("joinPath", joinPath)

func joinPath(_ *templating.Runtime, _ *templating.Context, paths ...string) string {
	return filepath.Join(paths...)
}

// Returns the last element of a path.
//
// @param path string
// @returns string
var basePathFunc = modules.NewFunc("basePath", basePath)

func basePath(_ *templating.Runtime, _ *templating.Context, path string) string {
	return filepath.Base(path)
}

// Returns all but the last element of a path.
//
// @param path string
// @returns string
var pathDirFunc = modules.NewFunc("pathDir", pathDir)

func pathDir(_ *templating.Runtime, _ *templating.Context, path string) string {
	return filepath.Dir(path)
}

// Returns the file name extension of a file path.
//
// @param path string
// @returns string
var fileExtFunc = modules.NewFunc("fileExt", fileExt)

func fileExt(_ *templating.Runtime, _ *templating.Context, path string) string {
	return filepath.Ext(path)
}

// Returns an absolute representation of a given path. 
// If the path is not absolute it will be joined with the current working directory to turn it into an absolute path.
//
// @param path string
// @returns string
// @returns error
var absPathFunc = modules.NewFunc("absPath", absPath)

func absPath(_ *templating.Runtime, _ *templating.Context, path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return abs, nil
}

// Returns a relative path that is lexically equivalent to targetPath when joined to basePath with an intervening separator.
// The returned path will always be relative to basePath, even if basePath and targetPath share no elements.
//
// @param basePath string
// @param targetPath string
// @returns string
// @returns error
var relPathFunc = modules.NewFunc("relPath", relPath)

func relPath(_ *templating.Runtime, _ *templating.Context, basePath string, targetPath string) (string, error) {
	rel, err := filepath.Rel(basePath, targetPath)
	if err != nil {
		return "", err
	}

	return rel, nil
}