package snippets

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/codeshelldev/goplater/utils/pretty"
)

type SnippetPointer struct {
	Start	int
	Line	int
	Column	int
	Length	int
	Message	string
	Color	bool
}

type SnippetLines struct {
	Before		[]SnippetLine
	After		[]SnippetLine
	PointerLine	string
	MessageLine	string
}

type SnippetLine struct {
	Number	int
	Content	string
}

func (s SnippetLines) ToString(color bool) string {
	var lines []string

	for _, line := range s.Before {
		lines = append(lines, formatSnippetLine(line, color))
	}

	if s.PointerLine != "" {
		pointer := pretty.Error(s.PointerLine, color)
		lines = append(lines, pointer)
	}

	if s.MessageLine != "" {
		message := pretty.Error(s.MessageLine, color)
		lines = append(lines, message)
	}

	for _, line := range s.After {
		lines = append(lines, formatSnippetLine(line, color))
	}

	return strings.Join(lines, "\n")
}

func (s SnippetLines) String() string {
	return s.ToString(false)
}

func formatSnippetLine(line SnippetLine, color bool) string {
	number := pretty.Muted(formatNumber(line.Number, line.Number), color)
	content := pretty.NewStyleFromEnabled(color).Paint(pretty.Cyan, line.Content)

	return fmt.Sprintf("%s | %s", number, content)
}

func PointToPositionInSnippet(snippet string, ptr SnippetPointer, color bool) string {
	return BuildSnippetPointer(snippet, ptr).ToString(color)
}

func BuildSnippetPointer(snippet string, ptr SnippetPointer) SnippetLines {
	lines := ParseCodeSnippet(snippet, ptr.Start)

	lineIndex := ptr.Line - ptr.Start - 1
	if lineIndex < 0 || lineIndex >= len(lines) {
		return SnippetLines{Before: lines}
	}

	sourceLine := lines[lineIndex].Content

	column := ptr.Column + 1
	length := ptr.Length

	// entire line fallback
	if column <= 1 {
		column = 1
		length = len(sourceLine)
	}

	// dynamic length
	if length <= 0 {
		length = len(sourceLine) - column + 1
	}

	if length < 1 {
		length = 1
	}

	lineNumberWidth := countDigits(ptr.Start + len(lines))
	gutterWidth := lineNumberWidth + 3 // " | "

	padding := strings.Repeat(" ", gutterWidth + column - 1)

	pointerLine := padding + strings.Repeat("^", length)

	var messageLine string
	if ptr.Message != "" {
		messageLine = padding + ptr.Message
	}

	return SnippetLines{
		Before: lines[:lineIndex + 1],
		After: lines[lineIndex + 1:],
		PointerLine: pointerLine,
		MessageLine: messageLine,
	}
}

func ParseCodeSnippet(snippet string, start int) []SnippetLine {
	lines := strings.Split(snippet, "\n")
	result := make([]SnippetLine, len(lines))

	for i, line := range lines {
		result[i] = SnippetLine{
			Number:  start + i + 1,
			Content: line,
		}
	}

	return result
}

func FormatCodeSnippet(snippet string, start int) string {
	lines := ParseCodeSnippet(snippet, start)

	maxLine := start + len(lines)
	formatted := make([]string, len(lines))

	for i, line := range lines {
		formatted[i] = fmt.Sprintf(
			"%s | %s",
			formatNumber(line.Number, maxLine),
			line.Content,
		)
	}

	return strings.Join(formatted, "\n")
}

func countDigits(max int) int {
	return len(strconv.Itoa(max))
}

func formatNumber(num, max int) string {
	digits := countDigits(max)
	digitVerb := "%" + strconv.Itoa(digits) + "d"

	return fmt.Sprintf(digitVerb, num)
}