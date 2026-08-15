package resolvers

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
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
	resp, err := r.Client.Get(path)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("http resolver: " + path + " returned " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (r *HttpResolver) DeriveName(path string) (string, error) {
	URL, err := url.Parse(path)

	if err != nil {
		return "", errors.New("malformed path: " + path + ", cannot derive import name")
	}

	name := filepath.Base(strings.TrimRight(URL.Path, "/"))
	name = strings.TrimSuffix(name, filepath.Ext(name))

	return name, nil
}