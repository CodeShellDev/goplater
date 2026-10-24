package pretty

import (
	"fmt"
	"io"
	"os"
)

type Color string

const (
	Reset   Color = "\033[0m"
	Black   Color = "\033[30m"
	Red     Color = "\033[31m"
	Green   Color = "\033[32m"
	Yellow  Color = "\033[33m"
	Blue    Color = "\033[34m"
	Magenta Color = "\033[35m"
	Cyan    Color = "\033[36m"
	White   Color = "\033[37m"
	Gray    Color = "\033[90m"
	Bold    Color = "\033[1m"
	Dim     Color = "\033[2m"
)

type Symbols struct {
	Success  string
	Error    string
	Warning  string
	Info     string
	Pending  string
	Branch   string
	Last     string
	Vertical string
	Spinner  []string
}

var UnicodeSymbols = Symbols{
	Success:  "✓",
	Error:    "✗",
	Warning:  "!",
	Info:     "ℹ",
	Pending:  "○",
	Branch:   "├── ",
	Last:     "└── ",
	Vertical: "│   ",
	Spinner:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
}

var ASCIISymbols = Symbols{
	Success:  "OK",
	Error:    "ERROR",
	Warning:  "!",
	Info:     "i",
	Pending:  "...",
	Branch:   "|-- ",
	Last:     "`-- ",
	Vertical: "|   ",
	Spinner:  []string{"-", "\\", "|", "/"},
}

type Style struct {
	Enabled bool
}

func NewStyle(w io.Writer, force bool) Style {
	enabled := force || (IsTerminal(w) && os.Getenv("NO_COLOR") == "")
	return Style{Enabled: enabled}
}

func (s Style) Paint(c Color, text string) string {
	if !s.Enabled || c == "" {
		return text
	}
	return string(c) + text + string(Reset)
}

func (s Style) Success(text string) string  { return s.Paint(Green, text) }
func (s Style) Error(text string) string    { return s.Paint(Red, text) }
func (s Style) Warning(text string) string  { return s.Paint(Yellow, text) }
func (s Style) Info(text string) string     { return s.Paint(Cyan, text) }
func (s Style) Muted(text string) string    { return s.Paint(Gray, text) }
func (s Style) Emphasis(text string) string { return s.Paint(Bold, text) }

func (s Style) Printf(w io.Writer, c Color, format string, args ...any) {
	fmt.Fprint(w, s.Paint(c, fmt.Sprintf(format, args...)))
}

func IsTerminal(w io.Writer) bool {
	if os.Getenv("TERM") == "dumb" {
		return false
	}

	f, ok := w.(*os.File)
	if !ok || (f != os.Stdout && f != os.Stderr) {
		return false
	}

	info, err := f.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

func ColorsEnabled(w io.Writer, force bool) bool {
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		return false
	}
	return force || IsTerminal(w)
}

func DefaultSymbols() Symbols {
	if os.Getenv("GOPLATER_ASCII") != "" {
		return ASCIISymbols
	}
	return UnicodeSymbols
}

func SymbolsFor(w io.Writer) Symbols {
	if os.Getenv("GOPLATER_ASCII") != "" || !IsTerminal(w) {
		return ASCIISymbols
	}
	return UnicodeSymbols
}

func Info(text string, enabled bool) string {
	return NewStyleFromEnabled(enabled).Info(text)
}

func Success(text string, enabled bool) string {
	return NewStyleFromEnabled(enabled).Success(text)
}

func Warning(text string, enabled bool) string {
	return NewStyleFromEnabled(enabled).Warning(text)
}

func Error(text string, enabled bool) string {
	return NewStyleFromEnabled(enabled).Error(text)
}

func Muted(text string, enabled bool) string {
	return NewStyleFromEnabled(enabled).Muted(text)
}

func NewStyleFromEnabled(enabled bool) Style {
	return Style{Enabled: enabled}
}