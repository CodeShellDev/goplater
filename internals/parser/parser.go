package parser

import (
	"errors"
	"fmt"

	"github.com/codeshelldev/goplater/pkg/gplt/types"
)

type Delims struct {
	Left  string
	Right string
}

type Context struct {
	Delims types.Delims
}

type Processor interface {
	Name() string
	Process(input string, ctx Context) (string, error)
}

type ProcessorFunc struct {
	ProcessorName 	string
	Func           	func(input string, ctx Context) (string, error)
}

func (p ProcessorFunc) Name() string {
	return p.ProcessorName
}

func (p ProcessorFunc) Process(input string, ctx Context) (string, error) {
	return p.Func(input, ctx)
}

type Parser struct {
	delims     types.Delims
	processors []Processor
}

func New() *Parser {
	return &Parser{
		delims: types.Delims{
			Left:  "{{",
			Right: "}}",
		},
	}
}

func (p *Parser) Delims(delims types.Delims) *Parser {
	p.delims = delims
	return p
}

func (p *Parser) Use(processors ...Processor) *Parser {
	p.processors = append(p.processors, processors...)
	return p
}

func (p *Parser) Process(input string) (string, error) {
	if p.delims.Left == "" || p.delims.Right == "" {
		return "", errors.New("missing delimiters")
	}

	ctx := Context{Delims: p.delims}
	output := input

	for _, processor := range p.processors {
		if processor == nil {
			return "", errors.New("nil processor registered")
		}

		name := processor.Name()
		result, err := processor.Process(output, ctx)
		if err != nil {
			return "", fmt.Errorf("processor %q: %w", name, err)
		}

		output = result
	}

	return output, nil
}

func (p *Parser) ProcessorNames() []string {
	names := make([]string, 0, len(p.processors))
	for _, processor := range p.processors {
		if processor == nil {
			names = append(names, "<nil>")
			continue
		}
		names = append(names, processor.Name())
	}
	return names
}