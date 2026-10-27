package ignore

import (
	"errors"
	"strconv"
	"strings"
	"unicode"

	"github.com/codeshelldev/goplater/internals/parser"
)

const (
	IgnoreStartMarker = "GOPLATER_IGNORE"
	IgnoreEndMarker   = "GOPLATER_ENDIGNORE"
)

type IgnoreProcessor struct{}

func New() IgnoreProcessor { return IgnoreProcessor{} }

func (p IgnoreProcessor) Name() string { return "ignore" }

func (p IgnoreProcessor) Process(input string, ctx parser.Context) (string, error) {
	left, right := ctx.Delims.Left, ctx.Delims.Right

	if left == "" || right == "" {
		return "", errors.New("delimiters missing")
	}

	var output strings.Builder
	remaining := input

	for {
		startStart, startEnd := findMarker(remaining, left, right, IgnoreStartMarker)
		endStart, endEnd := findMarker(remaining, left, right, IgnoreEndMarker)

		if startStart == -1 {
			if endStart != -1 {
				return "", errors.New("end marker found without a matching start marker")
			}
			output.WriteString(remaining)
			break
		}

		if endStart != -1 && endStart < startStart {
			return "", errors.New("end marker found before start marker")
		}

		// preserve content before start marker
		output.WriteString(remaining[:startStart])

		if endStart == -1 {
			return "", errors.New("start marker found without a matching end marker")
		}

		// extract content, exclude both markers
		content := remaining[startEnd:endStart]

		output.WriteString(left)
		output.WriteString("print ")
		output.WriteString(strconv.Quote(content))
		output.WriteString(right)

		remaining = remaining[endEnd:]
	}

	return output.String(), nil
}

// returns start and exclusive end offsets of marker
func findMarker(input, left, right, marker string) (int, int) {
	searchFrom := 0

	for searchFrom < len(input) {
		relativeStart := strings.Index(input[searchFrom:], left)
		if relativeStart == -1 {
			return -1, -1
		}

		start := searchFrom + relativeStart
		contentStart := start + len(left)

		relativeRight := strings.Index(input[contentStart:], right)
		if relativeRight == -1 {
			return -1, -1
		}

		contentEnd := contentStart + relativeRight
		content := strings.TrimFunc(input[contentStart:contentEnd], unicode.IsSpace)

		if content == marker {
			return start, contentEnd + len(right)
		}

		searchFrom = start + len(left)
	}

	return -1, -1
}