package html

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
)

var Module = modules.NewModule(htmlDecodeFunc, htmlDocFindAllFunc, htmlDocFindFunc, htmlFindFunc, htmlFindAllFunc, htmlTextFunc, htmlAttrFunc, htmlInnerFunc)

// Decodes html into a document.
//
// @param html string
// @returns *goquery.Document
// @returns error
var htmlDecodeFunc = modules.NewFunc("htmlDecode", htmlDecode)

func htmlDecode(_ *templating.Runtime, _ *templating.Context, str string) (*goquery.Document, error)  {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(str))

	return doc, err
}

// Finds elements in document by selector (using [goquery.Find](https://pkg.go.dev/github.com/PuerkitoBio/goquery#Selection.Find)).
//
// @param doc *goquery.Document
// @param selector string
// @returns []*goquery.Selection
// @returns error
var htmlDocFindAllFunc = modules.NewFunc("htmlDocFindAll", htmlDocFindAll)

func htmlDocFindAll(_ *templating.Runtime, _ *templating.Context, doc *goquery.Document, selector string) []*goquery.Selection  {
	all := []*goquery.Selection{}

	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		all = append(all, s)
	})

	return all
}

// Finds element in document by selector (using [goquery.Find](https://pkg.go.dev/github.com/PuerkitoBio/goquery#Selection.Find)).
//
// @param doc *goquery.Document
// @param selector string
// @returns *goquery.Selection
// @returns error
var htmlDocFindFunc = modules.NewFunc("htmlDocFind", htmlDocFind)

func htmlDocFind(_ *templating.Runtime, _ *templating.Context, doc *goquery.Document, selector string) *goquery.Selection  {
	return doc.Find(selector).First()
}

// Finds elements in element by selector (using [goquery.Find](https://pkg.go.dev/github.com/PuerkitoBio/goquery#Selection.Find)).
//
// @param doc *goquery.Document
// @param selector string
// @returns []*goquery.Selection
// @returns error
var htmlFindAllFunc = modules.NewFunc("htmlFindAll", htmlFindAll)

func htmlFindAll(_ *templating.Runtime, _ *templating.Context, el *goquery.Selection, selector string) []*goquery.Selection  {
	all := []*goquery.Selection{}

	el.Find(selector).Each(func(i int, s *goquery.Selection) {
		all = append(all, s)
	})

	return all
}


// Finds element in element by selector (using [goquery.Find](https://pkg.go.dev/github.com/PuerkitoBio/goquery#Selection.Find)).
//
// @param doc *goquery.Document
// @param selector string
// @returns *goquery.Selection
// @returns error
var htmlFindFunc = modules.NewFunc("htmlFind", htmlFind)

func htmlFind(_ *templating.Runtime, _ *templating.Context, el *goquery.Selection, selector string) *goquery.Selection  {
	return el.Find(selector).First()
}

// Returns the elements text content.
//
// @param el *goquery.Selection
// @returns string
var htmlTextFunc = modules.NewFunc("htmlText", htmlText)

func htmlText(_ *templating.Runtime, _ *templating.Context, el *goquery.Selection) string  {
	return el.Text()
}

// Returns an elements attribute content.
//
// @param el *goquery.Selection
// @param attr string
// @returns string
var htmlAttrFunc = modules.NewFunc("htmlAttr", htmlAttr)

func htmlAttr(_ *templating.Runtime, _ *templating.Context, el *goquery.Selection, attr string) string  {
	val, _ := el.Attr(attr)

	return val
}

// Returns the elements inner html.
//
// @param el *goquery.Selection
// @returns string
var htmlInnerFunc = modules.NewFunc("htmlInner", htmlInner)

func htmlInner(_ *templating.Runtime, _ *templating.Context, el *goquery.Selection) string  {
	val, _ := el.Html()

	return val
}