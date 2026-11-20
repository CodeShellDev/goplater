package prettyerrors

import (
	"github.com/codeshelldev/goplater/pkg/gplt/utils/tmplerrors"
	"github.com/codeshelldev/goplater/utils/snippets"
)

type TemplateError struct {
    tmplerrors.TemplateError
}

var _ error = *new(TemplateError)

func Convert(err tmplerrors.TemplateError) TemplateError {
	return TemplateError{err}
}

func (e TemplateError) ErrorSnippet(color bool) string {
	if e.Line <= 0 {
		return ""
	}

	start, source := e.ErrorLines()

	length := 0
	if e.Expression != "" {
		length = len(e.Expression)
	}

	ptr := snippets.SnippetPointer{
		Start: start, Line: 
		e.Line, Column: e.Column,
		Length: length,
		Message: e.ErrorMessage,
	}

	return snippets.PointToPositionInSnippet(source, ptr, color)
}