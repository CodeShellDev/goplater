package http

import (
	"fmt"
	urlUtils "net/url"
	"time"

	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
	"github.com/codeshelldev/goplater/pkg/gplt/utils/httpclient"
)

// fetch: http(s) only
var FetchPolicy = httpclient.Policy{
	AllowedSchemes:	[]string{"https", "http"},
	MaxRedirects:	5,
	Timeout:		15 * time.Second,
}

var Module = modules.NewModule(fetchFunc)

// Fetches a resource from a given url.
//
// @param url string
// @returns string
// @returns error
var fetchFunc = modules.NewFunc("fetch", fetch)

func fetch(_ *gplt.Runtime, _ *gplt.Context, url string) (string, error) {
	_, err := urlUtils.Parse(url)
	if err != nil {
		return "", fmt.Errorf("invalid URL %q: %w", url, err)
	}

	return httpclient.FetchURL(httpclient.NewClient(FetchPolicy), url)
}