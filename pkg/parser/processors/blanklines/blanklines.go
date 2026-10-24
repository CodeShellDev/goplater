package blanklines

import (
	"math"
	"strings"
	"text/template/parse"

	"github.com/codeshelldev/goplater/pkg/parser"
)

type config struct {
	maxBlank int
}

type Option func(*config)

func WithMaxBlankLines(n int) Option {
	return func(c *config) {
		if n < 0 {
			n = 0
		}
		c.maxBlank = n
	}
}

type Processor struct {
	cfg config
}

func New() *Processor { return NewWith() }

func NewWith(opts ...Option) *Processor {
	cfg := config{maxBlank: 1}

	for _, opt := range opts {
		opt(&cfg)
	}

	return &Processor{cfg: cfg}
}

func (p *Processor) Name() string { return "blanklines" }

func (p *Processor) Process(doc *parser.Document) error {
	for name, tree := range doc.Trees {
		isMain := name == doc.Name
		root := tree.Root

		parser.WalkLists(root, func(list *parse.ListNode) {
			top := isMain && list == root
			p.processList(list, top, top)
		})
	}

	return nil
}

// should process node
func (p *Processor) eligible(n parse.Node) bool {
	a, ok := n.(*parse.ActionNode)
	return ok && a.Pipe != nil && len(a.Pipe.Decl) > 0
}

func (p *Processor) processList(list *parse.ListNode, atStart, atEnd bool) {
	normalize(list)

	nodes := list.Nodes
	out := make([]parse.Node, 0, len(nodes))

	for i := 0; i < len(nodes); {
		if !p.eligible(nodes[i]) {
			out = append(out, nodes[i])
			i++
			continue
		}

		j := p.runEnd(nodes, i)

		if p.collapse(nodes, out, j, atStart, atEnd) {
			// keep actions, drop whitespace between them
			for k := i; k <= j; k++ {
				if p.eligible(nodes[k]) {
					out = append(out, nodes[k])
				}
			}
		} else {
			out = append(out, nodes[i:j + 1]...)
		}

		i = j + 1
	}

	list.Nodes = out
	normalize(list)
}

// find end of eligible actions (i.e. whitespace and eligible actions), return index
func (p *Processor) runEnd(nodes []parse.Node, i int) int {
	j := i

	for j+1 < len(nodes) {
		switch {
		case p.eligible(nodes[j + 1]):
			j++
		case j + 2 < len(nodes) && isBlankText(nodes[j + 1]) && p.eligible(nodes[j + 2]):
			j += 2
		default:
			return j
		}
	}

	return j
}

// checks if nodes can be collapsed and collapses them
func (p *Processor) collapse(nodes, out []parse.Node, j int, atStart, atEnd bool) bool {
	var prev, next *parse.TextNode
	var ok bool

	// check if node before group is text
	if len(out) > 0 {
		prev, ok = out[len(out) - 1].(*parse.TextNode)
		if !ok {
			return false
		}
	}

	// check if node after group is text
	if j + 1 < len(nodes) {
		next, ok = nodes[j + 1].(*parse.TextNode)
		if !ok {
			return false
		}
	}

	// check if file end/start
	startsFile := atStart && len(out) <= 1
	endsFile := atEnd && j + 1 >= len(nodes) - 1

	before, after := "", ""
	if prev != nil {
		before = string(prev.Text)
	}
	if next != nil {
		after = string(next.Text)
	}

	// reject if preceding content => inline
	before = strings.TrimRight(before, " \t")
	if !strings.HasSuffix(before, "\n") && !(before == "" && startsFile) {
		return false
	}

	// reject if following content exists => inline 
	after = strings.TrimLeft(after, " \t")
	switch {
	case strings.HasPrefix(after, "\r\n"):
		after = after[2:]
	case strings.HasPrefix(after, "\n"):
		after = after[1:]
	case after == "" && endsFile:
	default:
		return false
	}

	trailCount := countTrailingBlank(before, startsFile)
	leadCount := countLeadingBlank(after)

	limit := p.cfg.maxBlank

	bAll, _ := trimTrailingBlank(before, math.MaxInt, startsFile)
	aAll, _ := trimLeadingBlank(after, math.MaxInt)

	// do not preserve lines at start/end of file
	if (startsFile && strings.TrimSpace(bAll) == "") || (endsFile && strings.TrimSpace(aAll) == "") {
		limit = 0
	}

	// trim excess
	excess := trailCount + leadCount - limit
	if excess > 0 {
		var removed int

		after, removed = trimLeadingBlank(after, excess)
		before, _ = trimTrailingBlank(before, excess-removed, startsFile)
	}

	if prev != nil {
		prev.Text = []byte(before)
	}
	if next != nil {
		next.Text = []byte(after)
	}

	return true
}

func isBlankLine(s string) bool { return strings.Trim(s, " \t\r") == "" }

func isBlankText(n parse.Node) bool {
	t, ok := n.(*parse.TextNode)
	return ok && strings.TrimSpace(string(t.Text)) == ""
}

// removes up to n blank lines at end of str,
// if atStart: str begins at file start => counts as real line
func trimTrailingBlank(str string, n int, atStart bool) (string, int) {
	removed := 0

	for removed < n && strings.HasSuffix(str, "\n") {
		body := str[:len(str) - 1]
		idx := strings.LastIndexByte(body, '\n')

		if idx < 0 {
			if atStart && isBlankLine(body) {
				return "", removed + 1
			}
			break
		}

		if !isBlankLine(body[idx + 1:]) {
			break
		}

		str = str[:idx + 1]
		removed++
	}

	return str, removed
}

// removes up to n blank lines from start of str
func trimLeadingBlank(str string, n int) (string, int) {
	removed := 0

	for removed < n {
		idx := strings.IndexByte(str, '\n')
		if idx < 0 || !isBlankLine(str[:idx]) {
			break
		}

		str = str[idx + 1:]
		removed++
	}

	return str, removed
}

func countTrailingBlank(s string, atStart bool) int {
	_, k := trimTrailingBlank(s, math.MaxInt, atStart)
	return k
}

func countLeadingBlank(s string) int {
	_, k := trimLeadingBlank(s, math.MaxInt)
	return k
}

// drop empty text nodes, merge neighbouring text nodes
func normalize(list *parse.ListNode) {
	out := make([]parse.Node, 0, len(list.Nodes))

	for _, n := range list.Nodes {
		t, ok := n.(*parse.TextNode)
		if !ok {
			out = append(out, n)
			continue
		}

		if len(t.Text) == 0 {
			continue
		}

		if len(out) > 0 {
			prev, ok := out[len(out) - 1].(*parse.TextNode)
			if ok {
				prev.Text = append(prev.Text, t.Text...)
				continue
			}
		}

		out = append(out, n)
	}

	list.Nodes = out
}