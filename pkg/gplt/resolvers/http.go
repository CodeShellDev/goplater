package resolvers

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/utils/httpclient"
)

type HttpResolver struct {
	Client *http.Client
}

// imports: https only
var ImportPolicy = httpclient.Policy{
	AllowedSchemes:	[]string{"https"},
	MaxRedirects:	5,
	Timeout:		5 * time.Second,
}

func NewHttpResolver() *HttpResolver {
	return &HttpResolver{Client: httpclient.NewClient(ImportPolicy)}
}

func (r *HttpResolver) Trusted() bool { 
	return false 
}

func (r *HttpResolver) CanResolve(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

func (r *HttpResolver) Resolve(ctx *gplt.Context, path string) (string, error) {
	u, err := url.Parse(path)
	if err != nil {
		return "", fmt.Errorf("malformed url %q: %w", path, err)
	}

	if u.Scheme != "https" {
		return "", fmt.Errorf("import %q: only https is allowed", path)
	}

	// default file ext
	u.Path += "." + gplt.IMPORT_FILE_EXTENSION
	u.RawPath = ""

	return httpclient.FetchURL(r.Client, u.String())
}

func (r *HttpResolver) DeriveName(path string) (string, error) {
	URL, err := url.Parse(path)

	if err != nil {
		return "", fmt.Errorf("malformed path: %q cannot derive import name", path)
	}

	name := filepath.Base(strings.TrimRight(URL.Path, "/"))
	name = strings.TrimSuffix(name, filepath.Ext(name))

	return name, nil
}

func (r *HttpResolver) CleanPath(ctx *gplt.Context, path string) (string, error) {
	URL, err := url.Parse(path)

	if err != nil {
		return "", fmt.Errorf("malformed path: %q cannot derive import name", path)
	}

	return URL.String(), nil
}