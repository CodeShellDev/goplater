package http

import (
	"errors"
	"io"
	"net/http"
	urlUtils "net/url"

	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
)

var Module = modules.NewModule(fetchFunc)

// Fetches a resource from a given url.
// Outputs errors as string.
//
// @param url string
// @returns string
var fetchFunc = modules.NewFunc("fetch", fetch)

func fetch(_ *gplt.Runtime, _ *gplt.Context, url string) (string, error) {
	_, err := urlUtils.Parse(url)

	if err != nil {
		return "", errors.New("invalid url: " + url)
	}

	response, err := http.Get(url)
	if err != nil {
		return "", errors.New("remote failed: " + url)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("body malformed: " + url)
	}

	return string(body), nil
}