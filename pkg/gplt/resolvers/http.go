package resolvers

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/utils/httpclient"
)

type HttpResolver struct {
	Client *http.Client

	mu		sync.Mutex
	cache	map[string]string
}

const MAX_CACHE_ENTRIES = 64

// imports: https only
var ImportPolicy = httpclient.Policy{
	AllowedSchemes:	[]string{"https"},
	MaxRedirects:	5,
	Timeout:		5 * time.Second,
}

func NewHttpResolver() *HttpResolver {
	return &HttpResolver{
		Client: httpclient.NewClient(ImportPolicy),
		cache: map[string]string{},
	}
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
	u.Fragment = ""

	target := u.String()

	r.mu.Lock()
	body, ok := r.cache[target]
	r.mu.Unlock()

	if ok {
		return body, nil
	}

	body, err = httpclient.FetchURL(r.Client, target)
	if err != nil {
		return "", err
	}

	r.mu.Lock()
	if len(r.cache) >= MAX_CACHE_ENTRIES {
		clear(r.cache)
	}
	r.cache[target] = body
	r.mu.Unlock()

	return body, nil
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