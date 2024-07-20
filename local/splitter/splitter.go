package splitter

import (
	"ecowitt-proxy/local/config"
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
	"net/url"
)

// HTTPClient interface to allow injecting a mock client
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Splitter handles splitting requests to multiple targets
type Splitter struct {
	Config config.Config
	Client HTTPClient
}

// HandleRequest forwards the request to multiple targets
func (s Splitter) HandleRequest(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	for _, target := range s.Config.Targets {
		err := s.forwardRequest(req, target)
		if err != nil {
			log.Printf("Error forwarding request to %s: %v", target.Name, err)
		}
	}
	return req, nil
}

// forwardRequest creates and sends a request to the specified target
func (s Splitter) forwardRequest(req *http.Request, target config.Target) error {
	proxyURL, err := url.Parse(target.HostAddr)
	if err != nil {
		return err
	}

	// Add original query parameters to the new URL
	query := proxyURL.Query()
	for key, values := range req.URL.Query() {
		for _, value := range values {
			query.Add(key, value)
		}
	}
	proxyURL.RawQuery = query.Encode()
	proxyURL.Path = target.Path

	proxyReq, err := http.NewRequestWithContext(req.Context(), req.Method, proxyURL.String(), req.Body)
	if err != nil {
		return err
	}

	// Copy headers from the original request
	for header, values := range req.Header {
		for _, value := range values {
			proxyReq.Header.Add(header, value)
		}
	}

	resp, err := s.Client.Do(proxyReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Printf("Forwarded request to %s, received response: %s", target.Name, resp.Status)
	return nil
}
