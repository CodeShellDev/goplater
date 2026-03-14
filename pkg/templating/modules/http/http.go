package http

import (
	"io"
	"net/http"
	"net/url"

	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
)

var Module = modules.NewModule(fetchFunc)

var fetchFunc = modules.NewFunc("fetch", fetch)

func fetch(_ *templating.Runtime, _ templating.Context, urlStr string) string {
	_, err := url.Parse(urlStr)

	if err != nil {
		return "invalid url: " + urlStr
	}

	response, err := http.Get(urlStr)
	if err != nil {
		return "remote failed: " + urlStr
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "body malformed: " + urlStr
	}

	return string(body)
}