package tmplerrors

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

type ParsedTemplate struct {
	Raw		string
	Template	*template.Template
}

type TemplateError struct {
	raw				error
	template		ParsedTemplate

	ParseName		string
	Message			string
	Line			int
	Column			int
	Expression		string

	ErrorContext	string
	ErrorMessage	string
}

var _ error = *new(TemplateError)


var (
	parseErrorRe = regexp.MustCompile(
		`^template: (?P<name>.*):(?P<line>\d+)(?::(?P<column>\d+))?:\s*(?P<message>.*)$`,
	)

	executionErrorRe = regexp.MustCompile(
		`^template: (?P<name>.*):(?P<line>\d+):(?P<column>\d+): executing "(?P<template>[^"]+)"(?: at <(?P<expression>.*?)>)?: (?P<message>.*)$`,
	)
)

func TryWrapTemplateErrorOnce(tmpl *template.Template, body string, err error) (TemplateError, bool) {
	var existing TemplateError

	// prevent rewrapping
	if !errors.As(err, &existing) {
		return Parse(ParseTemplate(tmpl, body), err)
	}
	
	return TemplateError{}, false
}

func ParseTemplate(tmpl *template.Template, body string) ParsedTemplate {
	return ParsedTemplate{Template: tmpl, Raw: body}
}

func Parse(t ParsedTemplate, err error) (TemplateError, bool) {
	if err == nil {
		return TemplateError{}, false
	}

	res := TemplateError{
		raw:	  err,
		template: t,
		Column:	  -1,
	}

	raw := err.Error()

	matches := executionErrorRe.FindStringSubmatch(raw)
	if matches != nil {
		return parseExecutionError(res, matches), true
	}

	matches = parseErrorRe.FindStringSubmatch(raw)
	if matches != nil {
		return parseParseError(res, matches), true
	}

	return TemplateError{}, false
}

func parseExecutionError(res TemplateError, matches []string) TemplateError {
	res.ParseName = matches[1]
	res.Line, _ = strconv.Atoi(matches[2])
	res.Column, _ = strconv.Atoi(matches[3])
	res.Expression = matches[5]
	res.ErrorMessage = matches[6]
	res.Message = res.ErrorMessage

	res.ErrorContext = fmt.Sprintf(
		"template: %s:%d:%d: executing %q",
		res.ParseName,
		res.Line,
		res.Column,
		matches[4],
	)

	if res.Expression != "" {
		res.ErrorContext += fmt.Sprintf(" at <%s>", res.Expression)
	}

	res.ErrorContext += ":"

	return res
}

func parseParseError(res TemplateError, matches []string) TemplateError {
	res.ParseName = matches[1]
	res.Line, _ = strconv.Atoi(matches[2])

	if matches[3] != "" {
		res.Column, _ = strconv.Atoi(matches[3])
	}

	res.ErrorMessage = matches[4]
	res.Message = res.ErrorMessage

	res.ErrorContext = fmt.Sprintf("template: %s:%d", res.ParseName, res.Line)

	if res.Column >= 0 {
		res.ErrorContext += fmt.Sprintf(":%d", res.Column)
	}

	res.ErrorContext += ":"

	return res
}

func (e TemplateError) Error() string {
	if e.raw != nil {
		return e.raw.Error()
	}
	return e.ErrorMessage
}

func (e TemplateError) ErrorLines() (int, string) {
	if e.Line <= 0 {
		return 0, ""
	}

	lines := strings.Split(e.template.Raw, "\n")

	from := max(e.Line - 2, 0)
	to := min(e.Line + 1, len(lines))

	lines = lines[from:to]

	// leading emty lines
	for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
		from++
	}

	// trailing empty lines
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines) - 1]) == "" {
		lines = lines[:len(lines) - 1]
	}

	return from, strings.Join(lines, "\n")
}

func (e TemplateError) CleanContext() string {
	ctx := e.ErrorContext

	ctx = strings.TrimPrefix(ctx, "template: ")
	ctx = strings.TrimSuffix(ctx, ":")

	i := strings.Index(ctx, ": executing ")
	if i != -1 {
		ctx = ctx[:i]
	}

	return strings.TrimSpace(ctx)
}

func (e TemplateError) CleanError() string {
	return strings.TrimSpace(e.ErrorMessage)
}