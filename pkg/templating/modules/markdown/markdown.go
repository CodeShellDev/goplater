package markdown

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

var Module = modules.NewModule(
	markdownDecodeFunc,
	markdownEncodeFunc,
	markdownHTMLFunc,

	markdownTextFunc,

	markdownHeadingsFunc,
	markdownLinksFunc,
	markdownImagesFunc,
	markdownParagraphsFunc,
	markdownCodeBlocksFunc,
	markdownBlockquotesFunc,

	markdownLinkURLFunc,
	markdownLinkSetURLFunc,
	markdownLinkTextFunc,

	markdownImageURLFunc,
	markdownImageSetURLFunc,
	markdownImageAltFunc,

	markdownRemoveFunc,
	markdownAppendFunc,
	markdownPrependFunc,

	markdownFindFunc,
	markdownFindAllFunc,
	markdownIsFunc,
)

var markdownParser = goldmark.DefaultParser()

type Document struct {
	Root ast.Node

	Source  []byte
	Sources [][]byte

	textSources map[*ast.Text][]byte

	nodeSources map[ast.Node][]byte
}

type Node struct {
	Document *Document
	Node     ast.Node
}

type Searchable interface {
	searchRoot() (ast.Node, *Document)
}

func (d *Document) searchRoot() (ast.Node, *Document) { 
	return d.Root, d 
}

func (n *Node) searchRoot() (ast.Node, *Document) { 
	return n.Node, n.Document 
}

var markdownDecodeFunc = modules.NewFunc("markdownDecode", markdownDecode)

func markdownDecode(_ *templating.Runtime, _ *templating.Context, str string) *Document {
	doc := newDocument([]byte(str))
	return doc
}

func newDocument(source []byte) *Document {
	doc := &Document{
		Source:      source,
		Sources:     [][]byte{source},
		textSources: make(map[*ast.Text][]byte),
		nodeSources: make(map[ast.Node][]byte),
	}

	doc.Root = markdownParser.Parse(text.NewReader(source))
	doc.registerSources(doc.Root, source)

	return doc
}

func (d *Document) registerSources(root ast.Node, source []byte) {
	ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch n := n.(type) {
		case *ast.Text:
			d.textSources[n] = source

		case *ast.CodeBlock:
			d.nodeSources[n] = source

		case *ast.FencedCodeBlock:
			d.nodeSources[n] = source
		}

		return ast.WalkContinue, nil
	})
}

var markdownEncodeFunc = modules.NewFunc("markdownEncode", markdownEncode)

func markdownEncode(_ *templating.Runtime, _ *templating.Context, doc *Document) string {
	var buf bytes.Buffer

	encodeNode(&buf, doc, doc.Root, 0)

	return buf.String()
}

func encodeNode(
	buf *bytes.Buffer,
	doc *Document,
	node ast.Node,
	depth int,
) {
	switch n := node.(type) {
	case *ast.Document:
		encodeChildren(buf, doc, n, depth)

	case *ast.Heading:
		buf.WriteString(strings.Repeat("#", n.Level))
		buf.WriteByte(' ')
		encodeChildren(buf, doc, n, depth)
		buf.WriteString("\n\n")

	case *ast.Paragraph:
		encodeChildren(buf, doc, n, depth)
		buf.WriteString("\n\n")

	case *ast.Text:
		source := doc.textSources[n]
		if source != nil {
			buf.Write(n.Segment.Value(source))
		}

	case *ast.String:
		buf.Write(n.Value)

	case *ast.Emphasis:
		buf.WriteString(strings.Repeat("*", n.Level))
		encodeChildren(buf, doc, n, depth)
		buf.WriteString(strings.Repeat("*", n.Level))

	case *ast.CodeSpan:
		buf.WriteByte('`')
		encodeChildren(buf, doc, n, depth)
		buf.WriteByte('`')

	case *ast.Link:
		buf.WriteByte('[')
		encodeChildren(buf, doc, n, depth)
		buf.WriteString("](")
		buf.Write(n.Destination)

		if n.Title != nil {
			buf.WriteString(` "`)
			buf.Write(n.Title)
			buf.WriteByte('"')
		}

		buf.WriteByte(')')

	case *ast.Image:
		buf.WriteString("![")
		encodeChildren(buf, doc, n, depth)
		buf.WriteString("](")
		buf.Write(n.Destination)

		if n.Title != nil {
			buf.WriteString(` "`)
			buf.Write(n.Title)
			buf.WriteByte('"')
		}

		buf.WriteByte(')')

	case *ast.Blockquote:
		encodeBlockquote(buf, doc, n, depth)

	case *ast.List:
		encodeList(buf, doc, n, depth)

	case *ast.ListItem:
		encodeListItem(buf, doc, n, depth)

	case *ast.FencedCodeBlock:
		encodeFencedCodeBlock(buf, doc, n)

	case *ast.CodeBlock:
		encodeIndentedCodeBlock(buf, doc, n)

	case *ast.ThematicBreak:
		buf.WriteString("---\n\n")

	case *ast.AutoLink:
		buf.WriteByte('<')
		buf.Write(n.URL(nil))
		buf.WriteByte('>')

	case *ast.RawHTML:
		buf.Write(n.Segments.Value(doc.nodeSources[n]))

	default:
		// unknown node: preserve children just in case
		encodeChildren(buf, doc, n, depth)
	}
}

func encodeChildren(
	buf *bytes.Buffer,
	doc *Document,
	node ast.Node,
	depth int,
) {
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		encodeNode(buf, doc, child, depth)
	}
}

func encodeBlockquote(
	buf *bytes.Buffer,
	doc *Document,
	node *ast.Blockquote,
	depth int,
) {
	var inner bytes.Buffer
	encodeChildren(&inner, doc, node, depth + 1)

	for _, line := range strings.SplitAfter(inner.String(), "\n") {
		if line == "" {
			continue
		}

		buf.WriteString("> ")
		buf.WriteString(line)
	}

	buf.WriteByte('\n')
}

func encodeList(
	buf *bytes.Buffer,
	doc *Document,
	node *ast.List,
	depth int,
) {
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		encodeNode(buf, doc, child, depth)
	}

	buf.WriteByte('\n')
}

func encodeListItem(
	buf *bytes.Buffer,
	doc *Document,
	node *ast.ListItem,
	depth int,
) {
	buf.WriteString("- ")

	first := true

	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if !first {
			buf.WriteString("  ")
		}

		encodeNode(buf, doc, child, depth + 1)
		first = false
	}

	buf.WriteByte('\n')
}

func encodeFencedCodeBlock(
	buf *bytes.Buffer,
	doc *Document,
	node *ast.FencedCodeBlock,
) {
	buf.WriteString("```")

	if node.Info != nil {
		buf.WriteByte(' ')
		buf.Write(node.Lines().Value(doc.nodeSources[node]))
	}

	buf.WriteByte('\n')

	source := doc.nodeSources[node]

	if source != nil {
		lines := node.Lines()

		for i := 0; i < lines.Len(); i++ {
			segment := lines.At(i)
			buf.Write(segment.Value(source))
		}
	}

	if !strings.HasSuffix(buf.String(), "\n") {
		buf.WriteByte('\n')
	}

	buf.WriteString("```\n\n")
}

func encodeIndentedCodeBlock(
	buf *bytes.Buffer,
	doc *Document,
	node *ast.CodeBlock,
) {
	source := doc.nodeSources[node]

	if source == nil {
		return
	}

	lines := node.Lines()

	for i := 0; i < lines.Len(); i++ {
		segment := lines.At(i)
		buf.WriteString("    ")
		buf.Write(segment.Value(source))
	}

	buf.WriteByte('\n')
}

func wrapNode(doc *Document, node ast.Node) *Node {
	return &Node{
		Document: doc,
		Node:     node,
	}
}

// splits "type[attr=value]" into its parts
// if there's no bracket, hasAttr is false and base is the whole selector
func parseSelector(sel string) (base, attr, val string, hasAttr bool) {
	open := strings.Index(sel, "[")
	if open == -1 || !strings.HasSuffix(sel, "]") {
		return sel, "", "", false
	}

	base = sel[:open]
	inner := sel[open + 1 : len(sel) - 1]

	parts := strings.SplitN(inner, "=", 2)
	if len(parts) != 2 {
		return base, "", "", false
	}

	return base, strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), true
}

// checks a node against a selector like "heading" or "heading[level=2]"
func nodeMatches(doc *Document, node ast.Node, selector string) bool {
	base, attr, val, hasAttr := parseSelector(selector)

	if !typeMatches(node, base) {
		return false
	}

	if !hasAttr {
		return true
	}

	return attrMatches(doc, node, attr, val)
}

func typeMatches(node ast.Node, selector string) bool {
	switch selector {
	case "document":
		_, ok := node.(*ast.Document)
		return ok

	case "heading":
		_, ok := node.(*ast.Heading)
		return ok

	case "paragraph":
		_, ok := node.(*ast.Paragraph)
		return ok

	case "link":
		_, ok := node.(*ast.Link)
		return ok

	case "image":
		_, ok := node.(*ast.Image)
		return ok

	case "emphasis":
		_, ok := node.(*ast.Emphasis)
		return ok

	case "codeSpan":
		_, ok := node.(*ast.CodeSpan)
		return ok

	case "codeBlock":
		switch node.(type) {
		case *ast.CodeBlock, *ast.FencedCodeBlock:
			return true
		}

	case "fencedCodeBlock":
		_, ok := node.(*ast.FencedCodeBlock)
		return ok

	case "blockquote":
		_, ok := node.(*ast.Blockquote)
		return ok

	case "list":
		_, ok := node.(*ast.List)
		return ok

	case "listItem":
		_, ok := node.(*ast.ListItem)
		return ok

	case "thematicBreak":
		_, ok := node.(*ast.ThematicBreak)
		return ok

	case "text":
		_, ok := node.(*ast.Text)
		return ok

	case "string":
		_, ok := node.(*ast.String)
		return ok
	}

	return false
}

func attrMatches(doc *Document, node ast.Node, attr, val string) bool {
	switch attr {
	case "level":
		if h, ok := node.(*ast.Heading); ok {
			return strconv.Itoa(h.Level) == val
		}
		return false

	case "lang":
		if fcb, ok := node.(*ast.FencedCodeBlock); ok {
			source := doc.nodeSources[fcb]
			if fcb.Info == nil || source == nil {
				return false
			}
			info := string(fcb.Lines().Value(source))
			// info line may have extra attrs after the language, e.g. "go run"
			lang := strings.Fields(info)
			return len(lang) > 0 && lang[0] == val
		}
		return false
	}

	return false
}

var markdownFindAllFunc = modules.NewFunc("markdownFindAll", markdownFindAll)

func markdownFindAll(_ *templating.Runtime, _ *templating.Context, target Searchable, selector string) []*Node {
	root, doc := target.searchRoot()

	var result []*Node

	ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering || n == root {
			return ast.WalkContinue, nil
		}

		if nodeMatches(doc, n, selector) {
			result = append(result, wrapNode(doc, n))
		}

		return ast.WalkContinue, nil
	})

	return result
}

var markdownFindFunc = modules.NewFunc("markdownFind", markdownFind)

func markdownFind(rt *templating.Runtime, ctx *templating.Context, target Searchable, selector string) *Node {
	all := markdownFindAll(rt, ctx, target, selector)

	if len(all) == 0 {
		return nil
	}

	return all[0]
}

var markdownIsFunc = modules.NewFunc("markdownIs", markdownIs)

func markdownIs(_ *templating.Runtime, _ *templating.Context, node *Node, selector string) bool {
	if node == nil || node.Node == nil {
		return false
	}
	return nodeMatches(node.Document, node.Node, selector)
}

var markdownHTMLFunc = modules.NewFunc("markdownHTML", markdownHTML)

func markdownHTML(_ *templating.Runtime, _ *templating.Context, doc *Document) (string, error) {
	var buf bytes.Buffer

	md := goldmark.New()

	err := md.Renderer().Render(&buf, doc.Source, doc.Root)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

var markdownTextFunc = modules.NewFunc("markdownText", markdownText)

func markdownText(_ *templating.Runtime, _ *templating.Context, node *Node) string {
	return nodeText(node.Document, node.Node)
}

func nodeText(doc *Document, root ast.Node) string {
	var buf bytes.Buffer

	ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch n := n.(type) {
		case *ast.Text:
			source := doc.textSources[n]
			if source != nil {
				buf.Write(n.Segment.Value(source))
			}

		case *ast.String:
			buf.Write(n.Value)
		}

		return ast.WalkContinue, nil
	})

	return buf.String()
}

var markdownHeadingsFunc = modules.NewFunc("markdownHeadings", markdownHeadings)

func markdownHeadings(_ *templating.Runtime, _ *templating.Context, doc *Document) []*Heading {
	var result []*Heading

	ast.Walk(doc.Root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		heading, ok := n.(*ast.Heading)
		if !ok {
			return ast.WalkContinue, nil
		}

		result = append(result, &Heading{
			Node: &Node{
				Document: doc,
				Node:     heading,
			},
			Level: heading.Level,
		})

		return ast.WalkContinue, nil
	})

	return result
}

var markdownLinksFunc = modules.NewFunc("markdownLinks", markdownLinks)

func markdownLinks(_ *templating.Runtime, _ *templating.Context, doc *Document) []*Link {
	var result []*Link

	ast.Walk(doc.Root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		link, ok := n.(*ast.Link)
		if !ok {
			return ast.WalkContinue, nil
		}

		result = append(result, &Link{
			Node: &Node{
				Document: doc,
				Node:     link,
			},
		})

		return ast.WalkContinue, nil
	})

	return result
}

var markdownLinkURLFunc = modules.NewFunc("markdownLinkURL", markdownLinkURL)

func markdownLinkURL(_ *templating.Runtime, _ *templating.Context, link *Link) string {
	return string(link.Node.Node.(*ast.Link).Destination)
}

var markdownLinkSetURLFunc = modules.NewFunc("markdownLinkSetURL", markdownLinkSetURL)

func markdownLinkSetURL(
	_ *templating.Runtime,
	_ *templating.Context,
	link *Link,
	url string,
) {
	link.Node.Node.(*ast.Link).Destination = []byte(url)
}

var markdownLinkTextFunc = modules.NewFunc("markdownLinkText", markdownLinkText)

func markdownLinkText(_ *templating.Runtime, _ *templating.Context, link *Link) string {
	return nodeText(link.Node.Document, link.Node.Node)
}

var markdownImagesFunc = modules.NewFunc("markdownImages", markdownImages)

func markdownImages(_ *templating.Runtime, _ *templating.Context, doc *Document) []*Image {
	var result []*Image

	ast.Walk(doc.Root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		image, ok := n.(*ast.Image)
		if !ok {
			return ast.WalkContinue, nil
		}

		result = append(result, &Image{
			Node: &Node{
				Document: doc,
				Node:     image,
			},
		})

		return ast.WalkContinue, nil
	})

	return result
}

var markdownImageURLFunc = modules.NewFunc("markdownImageURL", markdownImageURL)

func markdownImageURL(_ *templating.Runtime, _ *templating.Context, image *Image) string {
	return string(image.Node.Node.(*ast.Image).Destination)
}

var markdownImageSetURLFunc = modules.NewFunc("markdownImageSetURL", markdownImageSetURL)

func markdownImageSetURL(
	_ *templating.Runtime,
	_ *templating.Context,
	image *Image,
	url string,
) {
	image.Node.Node.(*ast.Image).Destination = []byte(url)
}

var markdownImageAltFunc = modules.NewFunc("markdownImageAlt", markdownImageAlt)

func markdownImageAlt(_ *templating.Runtime, _ *templating.Context, image *Image) string {
	return nodeText(image.Node.Document, image.Node.Node)
}

var markdownParagraphsFunc = modules.NewFunc("markdownParagraphs", markdownParagraphs)

func markdownParagraphs(_ *templating.Runtime, _ *templating.Context, doc *Document) []*Paragraph {
	var result []*Paragraph

	ast.Walk(doc.Root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		p, ok := n.(*ast.Paragraph)
		if !ok {
			return ast.WalkContinue, nil
		}

		result = append(result, &Paragraph{
			Node: &Node{
				Document: doc,
				Node:     p,
			},
		})

		return ast.WalkContinue, nil
	})

	return result
}

var markdownCodeBlocksFunc = modules.NewFunc("markdownCodeBlocks", markdownCodeBlocks)

func markdownCodeBlocks(_ *templating.Runtime, _ *templating.Context, doc *Document) []*CodeBlock {
	var result []*CodeBlock

	ast.Walk(doc.Root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch n.(type) {
		case *ast.CodeBlock, *ast.FencedCodeBlock:
			result = append(result, &CodeBlock{
				Node: &Node{
					Document: doc,
					Node:     n,
				},
			})
		}

		return ast.WalkContinue, nil
	})

	return result
}

var markdownBlockquotesFunc = modules.NewFunc("markdownBlockquotes", markdownBlockquotes)

func markdownBlockquotes(_ *templating.Runtime, _ *templating.Context, doc *Document) []*Blockquote {
	var result []*Blockquote

	ast.Walk(doc.Root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		bq, ok := n.(*ast.Blockquote)
		if !ok {
			return ast.WalkContinue, nil
		}

		result = append(result, &Blockquote{
			Node: &Node{
				Document: doc,
				Node:     bq,
			},
		})

		return ast.WalkContinue, nil
	})

	return result
}

var markdownRemoveFunc = modules.NewFunc("markdownRemove", markdownRemove)

func markdownRemove(_ *templating.Runtime, _ *templating.Context, node *Node) string {
	if node == nil || node.Node == nil {
		return ""
	}

	parent := node.Node.Parent()

	if parent != nil {
		parent.RemoveChild(parent, node.Node)
	}

	return ""
}

var markdownAppendFunc = modules.NewFunc("markdownAppend", markdownAppend)

func markdownAppend(_ *templating.Runtime, _ *templating.Context, doc *Document, str string) *Document {
	source := []byte(str)
	fragment := markdownParser.Parse(text.NewReader(source))

	doc.Sources = append(doc.Sources, source)
	doc.registerSources(fragment, source)

	for child := fragment.FirstChild(); child != nil; {
		next := child.NextSibling()

		doc.Root.AppendChild(doc.Root, child)

		child = next
	}

	return doc
}

var markdownPrependFunc = modules.NewFunc("markdownPrepend", markdownPrepend)

func markdownPrepend(_ *templating.Runtime, _ *templating.Context, doc *Document, str string) *Document {
	source := []byte(str)
	fragment := markdownParser.Parse(text.NewReader(source))

	doc.Sources = append(doc.Sources, source)
	doc.registerSources(fragment, source)

	for child := fragment.LastChild(); child != nil; {
		prev := child.PreviousSibling()

		doc.Root.InsertBefore(doc.Root, doc.Root.FirstChild(), child)

		child = prev
	}

	return doc
}

type Heading struct {
	*Node
	Level int
}

type Link struct {
	*Node
}

type Image struct {
	*Node
}

type Paragraph struct {
	*Node
}

type CodeBlock struct {
	*Node
}

type Blockquote struct {
	*Node
}