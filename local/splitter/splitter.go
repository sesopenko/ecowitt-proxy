package splitter

import (
	"context"
	"ecowitt-proxy/local/config"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
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

type OriginalRequest struct {
	Body   string
	Method string
	Header http.Header
}

// HandleRequest forwards the request to multiple targets in goroutines
func (s Splitter) HandleRequest(w http.ResponseWriter, req *http.Request) {
	original, err := buildOriginal(req)
	if err != nil {
		s.log("Error extracting values from request: %s", err)
		return
	}
	var wg sync.WaitGroup
	for _, target := range s.Config.Targets {
		wg.Add(1)
		go func(target config.Target) {
			defer wg.Done()
			err := s.forwardRequest(original, target)
			if err != nil {
				s.log("Error forwarding request to %s: %v", target.Name, err)
			}
		}(target)
	}

	// Respond to the client immediately
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request forwarded to all targets"))

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
	}()
}

func buildOriginal(req *http.Request) (OriginalRequest, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return OriginalRequest{}, err
	}
	bodyRead := string(body)
	original := OriginalRequest{
		Body:   bodyRead,
		Method: req.Method,
		Header: req.Header,
	}
	return original, nil
}

// forwardRequest creates and sends a request to the specified target
func (s Splitter) forwardRequest(original OriginalRequest, target config.Target) error {
	targetUrl := ""
	if proxyURL, err := url.Parse(target.HostAddr); err != nil {
		return err
	} else {
		targetUrl = proxyURL.String()
	}

	// Add original query parameters to the new URL if present

	// Create a new request with the same method and headers

	proxyReq, err := http.NewRequestWithContext(context.TODO(), original.Method, targetUrl, strings.NewReader(original.Body))
	if err != nil {
		return err
	}

	// Copy headers from the original request
	for header, values := range original.Header {
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
