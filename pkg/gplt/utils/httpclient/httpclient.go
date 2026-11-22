package httpclient

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"slices"
	"syscall"
	"time"
)
type Policy struct {
	AllowedSchemes	[]string
	AllowPrivate	bool // allow loopback, private, link-local targets
	MaxRedirects	int
	Timeout			time.Duration
}

var blockedPrefixes = []netip.Prefix{
	netip.MustParsePrefix("100.64.0.0/10"), // CGNAT
	netip.MustParsePrefix("192.0.0.0/24"),
	netip.MustParsePrefix("198.18.0.0/15"),
	netip.MustParsePrefix("240.0.0.0/4"),
	netip.MustParsePrefix("64:ff9b::/96"), // NAT64
}

func isBlocked(ip netip.Addr) bool {
	ip = ip.Unmap() // ipv6 to ipv4

	if ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() ||
		ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() ||
		ip.IsInterfaceLocalMulticast() || ip.IsMulticast() {
		return true
	}

	for _, p := range blockedPrefixes {
		if p.Contains(ip) {
			return true
		}
	}

	return false
}

// runs after dns resolve, on every connect attempt => stops dns rebinding
func guardControl(network, address string, _ syscall.RawConn) error {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}

	ip, err := netip.ParseAddr(host)
	if err != nil {
		return err
	}

	if isBlocked(ip) {
		return fmt.Errorf("blocked address %s", ip)
	}

	return nil
}

// checks scheme + userinfo on every request, redirects included
type guardTransport struct {
	base	http.RoundTripper
	schemes	[]string
}

func (t *guardTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if !slices.Contains(t.schemes, req.URL.Scheme) {
		return nil, fmt.Errorf("scheme %q not allowed", req.URL.Scheme)
	}

	if req.URL.User != nil {
		return nil, errors.New("credentials in url not allowed")
	}

	if req.URL.Hostname() == "" {
		return nil, errors.New("url has no host")
	}

	return t.base.RoundTrip(req)
}

func NewClient(p Policy) *http.Client {
	dialer := &net.Dialer{Timeout: 5 * time.Second}

	if !p.AllowPrivate {
		dialer.Control = guardControl
	}

	base := &http.Transport{
		Proxy:					nil, // ignore *_PROXY env
		DialContext:			dialer.DialContext,
		DisableKeepAlives:		true,
		TLSHandshakeTimeout:	5 * time.Second,
		ResponseHeaderTimeout:	p.Timeout,
	}

	return &http.Client{
		Timeout:   p.Timeout,
		Transport: &guardTransport{base: base, schemes: p.AllowedSchemes},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= p.MaxRedirects {
				return errors.New("too many redirects")
			}

			// no https -> http downgrade
			if via[len(via) - 1].URL.Scheme == "https" && req.URL.Scheme != "https" {
				return errors.New("redirect downgrades https")
			}

			return nil
		},
	}
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

	if resp.ContentLength > maxResponseSize {
		return "", errors.New("remote response exceeds 10 MB limit")
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize + 1))
	if err != nil {
		return "", fmt.Errorf("could not read response body: %w", err)
	}

	if len(body) > maxResponseSize {
		return "", errors.New("remote response exceeds 10 MB limit")
	}

	return string(body), nil
}
