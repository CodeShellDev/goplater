package fsutils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func Exists(path string) (bool, os.FileInfo) {
	info, err := os.Stat(path)

	return !errors.Is(err, os.ErrNotExist), info
}

func IsDir(path string) bool {
	exists, info := Exists(path)

	if exists {
		return info.IsDir()
	}

	return false
}

func IsFile(path string) bool {
	exists, info := Exists(path)

	if exists {
		return !info.IsDir()
	}

	return false
}

// Resolves Path based on relative Source and relative Output Path.
// Use `preserverStruct` to preserve folder structure
func ResolveOutput(source, output string, preserveStruct bool) string {
	if preserveStruct {
		return resolveOutputPreserved(source, output)
	}

	return resolveOutput(source, output)
}

func resolveOutputPreserved(source, output string) string {
	isDirOutput := strings.HasSuffix(output, "/")

	if isDirOutput && output != "." {
		// Add output directory and source directory starting at first unequal folder
		sourceComponents := strings.SplitAfter(source, string(filepath.Separator))
		outputComponents := strings.SplitAfter(output, string(filepath.Separator))

		for i, component := range sourceComponents {
			if component != outputComponents[i] {
				sourceComponents = sourceComponents[i+1:]
				break
			}
		}

		source = strings.Join(sourceComponents, "")

		return output + source
	}

	if output == "." {
		return filepath.Join(output, filepath.Base(source))
	}

	// Rename or overwrite file
	return output
}

func resolveOutput(source, output string) string {
	isDirSource := strings.HasSuffix(source, "/")
	isDirOutput := strings.HasSuffix(output, "/")

	// Directory source
	if isDirSource {
		// Rename or overwrite
		return output
	}

	// File source
	if isDirOutput || output == "." {
		// Place inside output directory
		return filepath.Join(output, filepath.Base(source))
	}

	// Rename or overwrite file
	return output
}