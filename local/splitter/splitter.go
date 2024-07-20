package splitter

import (
	"ecowitt-proxy/local/config"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
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
func (s Splitter) HandleRequest(w http.ResponseWriter, req *http.Request) {
	for _, target := range s.Config.Targets {
		err := s.forwardRequest(req, target)
		if err != nil {
			log.Printf("Error forwarding request to %s: %v", target.Name, err)
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request forwarded to all targets"))
}

// forwardRequest creates and sends a request to the specified target
func (s Splitter) forwardRequest(req *http.Request, target config.Target) error {
	proxyURL, err := url.Parse(target.HostAddr)
	if err != nil {
		return err
	}

	// Add original query parameters to the new URL if present
	query := proxyURL.Query()
	for key, values := range req.URL.Query() {
		for _, value := range values {
			query.Add(key, value)
		}
	}
	proxyURL.RawQuery = query.Encode()

	// Create a new request with the same method and headers
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	bodyRead := string(body)
	s.log("body: %s", bodyRead)

	req.Body = io.NopCloser(strings.NewReader(bodyRead)) // Reset the body for reuse

	proxyReq, err := http.NewRequestWithContext(req.Context(), req.Method, proxyURL.String(), io.NopCloser(strings.NewReader(string(body))))
	if err != nil {
		return err
	}

	// Copy headers from the original request
	for header, values := range req.Header {
		for _, value := range values {
			proxyReq.Header.Add(header, value)
		}
	}
	proxyReq.Header.Add("X-ECOWITT-PROXY-TARGET", target.Name)
	s.log("request url: %s", proxyReq.URL.String())

	resp, err := s.Client.Do(proxyReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	s.log("Forwarded request to %s, received response: %s", target.Name, resp.Status)

	return nil
}

func (s Splitter) log(format string, v ...any) {
	if s.Config.Server.Verbose {
		log.Printf(format, v...)
	}
}
