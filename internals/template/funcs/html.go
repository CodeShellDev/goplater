package funcs

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/codeshelldev/goplater/internals/template/context"
)

var htmlDecodeFunc = TemplateFunc{
	Name: "htmlDecode",
	Handler: func(context context.TemplateContext, str string) (*goquery.Document, error) {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(str))

		return doc, err
	},
}

var htmlDocFindFunc = TemplateFunc{
	Name: "htmlDocFind",
	Handler: func(context context.TemplateContext, doc *goquery.Document, selector string) *goquery.Selection {
		return doc.Find(selector).First()
	},
}

var htmlFindFunc = TemplateFunc{
	Name: "htmlFind",
	Handler: func(context context.TemplateContext, el *goquery.Selection, selector string) *goquery.Selection {
		return el.Find(selector).First()
	},
}

var htmlTextFunc = TemplateFunc{
	Name: "htmlText",
	Handler: func(context context.TemplateContext, el *goquery.Selection) string {
		return el.Text()
	},
}

var htmlAttrFunc = TemplateFunc{
	Name: "htmlAttr",
	Handler: func(context context.TemplateContext, el *goquery.Selection, attr string) string {
		val, _ := el.Attr(attr)

		return val
	},
}

var htmlInnerFunc = TemplateFunc{
	Name: "htmlInner",
	Handler: func(context context.TemplateContext, el *goquery.Selection) string {
		val, _ := el.Html()

		return val
	},
}

func init() {
	Register(htmlDecodeFunc)
	Register(htmlDocFindFunc)
	Register(htmlFindFunc)
	Register(htmlTextFunc)
	Register(htmlInnerFunc)
	Register(htmlAttrFunc)
}