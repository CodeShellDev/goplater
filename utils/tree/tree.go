package tree

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/codeshelldev/goplater/utils/pretty"
)

type Status uint8

const (
	Waiting Status = iota
	Running
	Succeeded
	Failed
	Skipped
)

func (s Status) String() string {
	switch s {
	case Running:
		return "processing"
	case Succeeded:
		return "done"
	case Failed:
		return "failed"
	case Skipped:
		return "skipped"
	default:
		return "waiting"
	}
}

type TreeOptions struct {
	Root         string
	Writer       io.Writer
	Style        pretty.Style
	Symbols      pretty.Symbols
	Live         bool
	SpinnerEvery time.Duration
}

type Tree struct {
	mu              sync.Mutex
	opts            TreeOptions
	root            *Node
	stop            chan struct{}
	done            chan struct{}
	ticker          *time.Ticker
	started, closed bool
	frame, rendered int
}

type Node struct {
	tree                  *Tree
	parent                *Node
	children              []*Node
	Source, Output        string
	Status                Status
	Err                   error
	StartedAt, FinishedAt time.Time
}

func NewTree(w io.Writer, root string) *Tree {
	if w == nil {
		w = os.Stdout
	}

	return NewTreeWithOptions(TreeOptions{
		Root:         root,
		Writer:       w,
		Style:        pretty.NewStyle(w, false),
		Symbols:      pretty.SymbolsFor(w),
		Live:         pretty.IsTerminal(w),
		SpinnerEvery: 90 * time.Millisecond,
	})
}

func NewTreeWithOptions(opts TreeOptions) *Tree {
	if opts.Writer == nil {
		opts.Writer = os.Stdout
	}

	if opts.Root == "" {
		opts.Root = "."
	}

	if opts.SpinnerEvery <= 0 {
		opts.SpinnerEvery = 100 * time.Millisecond
	}

	if opts.Symbols.Success == "" {
		opts.Symbols = pretty.SymbolsFor(opts.Writer)
	}

	if opts.Style == (pretty.Style{}) {
		opts.Style = pretty.NewStyle(opts.Writer, false)
	}

	root := &Node{Source: opts.Root, Status: Waiting}
	return &Tree{opts: opts, root: root}
}

func (t *Tree) Root() *Node {
	return t.root
}

func (t *Tree) Add(source string) *Node {
	t.mu.Lock()
	defer t.mu.Unlock()

	clean := filepath.ToSlash(filepath.Clean(source))
	clean = strings.TrimPrefix(clean, "./")

	if clean == "" || clean == "." {
		return t.root
	}

	parts := strings.Split(clean, "/")
	cur := t.root

	for i, part := range parts {
		if part == "" || part == "." {
			continue
		}

		var child *Node
		for _, c := range cur.children {
			if c.Source == part && (i < len(parts) - 1 || c.Output == "") {
				child = c
				break
			}
		}

		if child == nil {
			child = &Node{
				tree:   t,
				parent: cur,
				Source: part,
				Status: Waiting,
			}
			cur.children = append(cur.children, child)
		}

		cur = child
		cur.tree = t
	}

	return cur
}

func (n *Node) SetOutput(output string) *Node {
	if n == nil || n.tree == nil {
		return n
	}

	t := n.tree
	t.mu.Lock()
	defer t.mu.Unlock()

	n.Output = output
	if t.started {
		t.renderLocked()
	}

	return n
}

func (n *Node) AddChild(label string) *Node {
	if n == nil || n.tree == nil {
		return nil
	}

	t := n.tree
	t.mu.Lock()
	defer t.mu.Unlock()

	child := &Node{
		tree:   t,
		parent: n,
		Source: label,
		Status: Waiting,
	}
	n.children = append(n.children, child)

	if t.started {
		t.renderLocked()
	}

	return child
}

func (t *Tree) Start() {
	t.mu.Lock()
	if t.started || t.closed {
		t.mu.Unlock()
		return
	}

	t.started = true
	t.renderLocked()

	if !t.opts.Live {
		t.mu.Unlock()
		return
	}

	t.stop = make(chan struct{})
	t.done = make(chan struct{})
	t.ticker = time.NewTicker(t.opts.SpinnerEvery)
	stop, done, ticker := t.stop, t.done, t.ticker

	t.mu.Unlock()

	go func() {
		defer close(done)

		for {
			select {
			case <-ticker.C:
				t.mu.Lock()
				
				if !t.closed {
					t.frame++
					t.renderLocked()
				}

				t.mu.Unlock()
			case <-stop:
				return
			}
		}
	}()
}

func (t *Tree) Close() {
	t.mu.Lock()
	if t.closed {
		t.mu.Unlock()
		return
	}

	t.closed = true
	if t.ticker != nil {
		t.ticker.Stop()
	}

	stop, done := t.stop, t.done
	if stop != nil {
		close(stop)
	}
	t.mu.Unlock()

	if done != nil {
		<-done
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if t.opts.Live && t.rendered > 0 {
		fmt.Fprintf(t.opts.Writer, "\033[%dA\033[J", t.rendered)
		t.rendered = 0
	}

	t.renderLocked()
}

func (t *Tree) Snapshot() string {
	t.mu.Lock()
	defer t.mu.Unlock()

	var b strings.Builder
	t.writeNode(&b, t.root, nil, true, true)
	return b.String()
}

func (n *Node) SetStatus(s Status) {
	if n != nil && n.tree != nil {
		n.tree.setStatus(n, s, nil)
	}
}

func (n *Node) Start() {
	n.SetStatus(Running)
}

func (n *Node) Success() {
	n.SetStatus(Succeeded)
}

func (n *Node) Skip() {
	n.SetStatus(Skipped)
}

func (n *Node) Fail(err error) {
	if n != nil && n.tree != nil {
		n.tree.setStatus(n, Failed, err)
	}
}

func (t *Tree) setStatus(n *Node, s Status, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	n.Status, n.Err = s, err
	now := time.Now()

	switch s {
	case Running:
		n.StartedAt = now
	case Succeeded, Failed, Skipped:
		n.FinishedAt = now
	}

	if t.started {
		t.renderLocked()
	}
}

func (t *Tree) renderLocked() {
	var b strings.Builder
	t.writeNode(&b, t.root, nil, true, true)

	text := strings.TrimRight(b.String(), "\n") + "\n"
	if t.opts.Live && t.rendered > 0 {
		fmt.Fprintf(t.opts.Writer, "\033[%dA\033[J", t.rendered)
	}

	io.WriteString(t.opts.Writer, text)
	t.rendered = strings.Count(text, "\n")
}

func (t *Tree) writeNode(b *strings.Builder, n *Node, ancestors []bool, last, root bool) {
	if root {
		b.WriteString(n.Source)
		b.WriteByte('\n')
	} else {
		for _, hasNext := range ancestors {
			if hasNext {
				b.WriteString(t.opts.Symbols.Vertical)
			} else {
				b.WriteString("    ")
			}
		}

		if last {
			b.WriteString(t.opts.Symbols.Last)
		} else {
			b.WriteString(t.opts.Symbols.Branch)
		}

		b.WriteString(t.statusGlyph(n.Status))
		b.WriteByte(' ')
		b.WriteString(n.Source)

		if n.Output != "" {
			b.WriteString(" → ")
			b.WriteString(n.Output)
		}

		if n.Status == Failed && n.Err != nil {
			fmt.Fprintf(b, " (%v)", n.Err)
		}

		b.WriteByte('\n')
	}

	next := ancestors
	if !root {
		next = append(append([]bool(nil), ancestors...), !last)
	}

	for i, child := range n.children {
		t.writeNode(b, child, next, i == len(n.children) - 1, false)
	}
}

func (t *Tree) statusGlyph(s Status) string {
	sym := t.opts.Symbols

	switch s {
	case Running:
		glyph := sym.Pending
		if len(sym.Spinner) > 0 {
			glyph = sym.Spinner[t.frame%len(sym.Spinner)]
		}
		return t.opts.Style.Paint(pretty.Cyan, glyph)

	case Succeeded:
		return t.opts.Style.Paint(pretty.Green, sym.Success)

	case Failed:
		return t.opts.Style.Paint(pretty.Red, sym.Error)

	case Skipped:
		return t.opts.Style.Paint(pretty.Yellow, "–")

	default:
		return t.opts.Style.Paint(pretty.Gray, sym.Pending)
	}
}