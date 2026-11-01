package resolvers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/codeshelldev/goplater/pkg/gplt"
)

type HttpResolver struct {
	Client *http.Client
}

func NewHttpResolver() *HttpResolver {
	return &HttpResolver{Client: &http.Client{Timeout: 5 * time.Second}}
}

func (r *HttpResolver) Trusted() bool { 
	return false 
}

func (r *HttpResolver) CanResolve(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

func (r *HttpResolver) Resolve(path string) (string, error) {
	path += "." + gplt.IMPORT_FILE_EXTENSION

	return FetchURL(r.Client, path)
}

func FetchURL(client *http.Client, url string) (string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("remote request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("remote returned %q: %s", url, resp.Status)
	}

	const maxResponseSize = 10 << 20 // 10 MB

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize + 1))
	if err != nil {
		return "", fmt.Errorf("could not read response body: %w", err)
	}

	if len(body) > maxResponseSize {
		return "", errors.New("remote response exceeds 10 MB limit")
	}

	return string(body), nil
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