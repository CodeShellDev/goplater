package http

import (
	"errors"
	"net/http"
	urlUtils "net/url"
	"time"

	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
	"github.com/codeshelldev/goplater/pkg/gplt/resolvers"
)

var Module = modules.NewModule(fetchFunc)

// Fetches a resource from a given url.
// Outputs errors as string.
//
// @param url string
// @returns string
var fetchFunc = modules.NewFunc("fetch", fetch)

func fetch(_ *gplt.Runtime, _ *gplt.Context, url string) (string, error) {
	if _, err := urlUtils.Parse(url); err != nil {
		return "", errors.New("invalid URL " + url + ": " + err.Error())
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	return resolvers.FetchURL(client, url)
}