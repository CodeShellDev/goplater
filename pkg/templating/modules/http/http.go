package http

import (
	"errors"
	"io"
	"net/http"
	urlUtils "net/url"

	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
)

var Module = modules.NewModule(fetchFunc, fetchThrowFunc)

// Fetches a resource from a given url.
// Outputs errors as string.
//
// @param url string
// @returns string
var fetchFunc = modules.NewFunc("fetch", fetch)

func fetch(_ *templating.Runtime, _ *templating.Context, url string) string {
	_, err := urlUtils.Parse(url)

	if err != nil {
		return "invalid url: " + url
	}

	response, err := http.Get(url)
	if err != nil {
		return "remote failed: " + url
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "body malformed: " + url
	}

	return string(body)
}

// Fetches a resource from a given url.
// Throws error instead of outputting it.
//
// @param url string
// @returns string
// @returns error
var fetchThrowFunc = modules.NewFunc("fetchThrow", fetchThrow)

func fetchThrow(_ *templating.Runtime, _ *templating.Context, url string) (string, error) {
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