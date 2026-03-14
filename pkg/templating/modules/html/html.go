package html

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
)

var Module = modules.NewModule(htmlDecodeFunc, htmlDocFindFunc, htmlFindFunc, htmlTextFunc, htmlAttrFunc, htmlInnerFunc)

var htmlDecodeFunc = modules.NewFunc("htmlDecode", htmlDecode)

func htmlDecode(_ *templating.Runtime, _ templating.Context, str string) (*goquery.Document, error)  {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(str))

	return doc, err
}

var htmlDocFindFunc = modules.NewFunc("htmlDocFind", htmlDocFind)

func htmlDocFind(_ *templating.Runtime, _ templating.Context, doc *goquery.Document, selector string) *goquery.Selection  {
	return doc.Find(selector).First()
}

var htmlFindFunc = modules.NewFunc("htmlFind", htmlFind)

func htmlFind(_ *templating.Runtime, _ templating.Context, el *goquery.Selection, selector string) *goquery.Selection  {
	return el.Find(selector).First()
}

var htmlTextFunc = modules.NewFunc("htmlText", htmlText)

func htmlText(_ *templating.Runtime, _ templating.Context, el *goquery.Selection) string  {
	return el.Text()
}

var htmlAttrFunc = modules.NewFunc("htmlAttr", htmlAttr)

func htmlAttr(_ *templating.Runtime, _ templating.Context, el *goquery.Selection, attr string) string  {
	val, _ := el.Attr(attr)

	return val
}

var htmlInnerFunc = modules.NewFunc("htmlInner", htmlInner)

func htmlInner(_ *templating.Runtime, _ templating.Context, el *goquery.Selection) string  {
	val, _ := el.Html()

	return val
}