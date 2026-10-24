package printer

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/codeshelldev/goplater/utils/pretty"
)

type Printer struct {
	Out        io.Writer
	Err        io.Writer
	OutColors  bool
	ErrColors  bool
	OutSymbols pretty.Symbols
	ErrSymbols pretty.Symbols
}

var Default = New()

func New() *Printer {
	return NewWithWriters(os.Stdout, os.Stderr)
}

func NewWithWriters(out, errOut io.Writer) *Printer {
	if out == nil {
		out = io.Discard
	}
	if errOut == nil {
		errOut = io.Discard
	}

	return &Printer{
		Out:        out,
		Err:        errOut,
		OutColors:  pretty.ColorsEnabled(out, false),
		ErrColors:  pretty.ColorsEnabled(errOut, false),
		OutSymbols: pretty.SymbolsFor(out),
		ErrSymbols: pretty.SymbolsFor(errOut),
	}
}

func SetDefault(p *Printer) {
	if p == nil {
		p = New()
	}
	Default = p
}

func (p *Printer) Info(args ...any) {
	p.write(p.Out, pretty.Info(
		p.SymbolsMessage(p.OutSymbols.Info, args...),
		p.OutColors,
	))
}

func (p *Printer) Success(args ...any) {
	p.write(p.Out, pretty.Success(
		p.SymbolsMessage(p.OutSymbols.Success, args...),
		p.OutColors,
	))
}

func (p *Printer) Warn(args ...any) {
	p.write(p.Err, pretty.Warning(
		p.SymbolsMessage(p.ErrSymbols.Warning, args...),
		p.ErrColors,
	))
}

func (p *Printer) Error(args ...any) {
	p.write(p.Err, pretty.Error(
		p.SymbolsMessage(p.ErrSymbols.Error, args...),
		p.ErrColors,
	))
}

func (p *Printer) Debug(args ...any) {
	p.write(p.Err, pretty.Muted(
		p.SymbolsMessage("DEBUG", args...),
		p.ErrColors,
	))
}

func (p *Printer) SymbolsMessage(symbol string, args ...any) string {
	message := strings.TrimSuffix(fmt.Sprintln(args...), "\n")
	if symbol == "" {
		return message
	}
	return symbol + " " + message
}

func (p *Printer) Blank() {
	fmt.Fprintln(p.Out)
}


func Info(args ...any) {
	Default.Info(args...)
}

func Success(args ...any) {
	Default.Success(args...)
}

func Warn(args ...any) {
	Default.Warn(args...)
}

func Error(args ...any) {
	Default.Error(args...)
}

func Debug(args ...any) {
	Default.Debug(args...)
}

func Blank() {
	Default.Blank()
}

func (p *Printer) write(w io.Writer, message string) {
	fmt.Fprintln(w, message)
}