package splitter

import (
	"ecowitt-proxy/local/config"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// mockHTTPClient is a mock implementation of an HTTP client
type mockHTTPClient struct {
	Response *http.Response
	Err      error
	Requests []*http.Request // Capture requests for verification
}

// Do is the mock implementation of the Do method
func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	m.Requests = append(m.Requests, req) // Capture the request
	return m.Response, m.Err
}

func TestSplitter_HandleRequest(t *testing.T) {
	cfg := config.Config{
		Targets: []config.Target{
			{
				Name:     "target1",
				HostAddr: "https://example.com:8220",
				Path:     "/api/webhook",
			},
		},
	}

	// Create a mock response
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body:       io.NopCloser(strings.NewReader("OK")),
	}

	// Create the mock HTTP client
	mockClient := &mockHTTPClient{
		Response: mockResp,
		Err:      nil,
	}

	s := Splitter{
		Config: cfg,
		Client: mockClient,
	}

	// Create a request with real-world GET parameters
	req := httptest.NewRequest("GET", "http://localhost?baromrelin=888.82&humidityin=52&tempinf=74.2", nil)
	req.Header.Add("Content-Type", "application/json")

	s.HandleRequest(req, nil)

	// Verify that the request was forwarded correctly
	if mockClient.Response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", mockClient.Response.StatusCode)
	}

	// Verify that the forwarded request contains the correct parameters
	if len(mockClient.Requests) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(mockClient.Requests))
	}

	forwardedReq := mockClient.Requests[0]
	params := forwardedReq.URL.Query()

	if params.Get("baromrelin") != "888.82" {
		t.Errorf("Expected baromrelin=888.82, got %s", params.Get("baromrelin"))
	}

	if params.Get("humidityin") != "52" {
		t.Errorf("Expected humidityin=52, got %s", params.Get("humidityin"))
	}

	if params.Get("tempinf") != "74.2" {
		t.Errorf("Expected tempinf=74.2, got %s", params.Get("tempinf"))
	}
}

func TestSplitter_forwardRequest(t *testing.T) {
	cfg := config.Config{
		Targets: []config.Target{
			{
				Name:     "target1",
				HostAddr: "http://example.com:8220",
				Path:     "/api/webhook",
			},
		},
	}

	// Create a mock response
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body:       io.NopCloser(strings.NewReader("OK")),
	}

	// Create the mock HTTP client
	mockClient := &mockHTTPClient{
		Response: mockResp,
		Err:      nil,
	}

	s := Splitter{
		Config: cfg,
		Client: mockClient,
	}

	// Create a request with real-world GET parameters
	req := httptest.NewRequest("GET", "http://localhost?baromrelin=888.82&humidityin=52&tempinf=74.2", nil)
	req.Header.Add("Content-Type", "application/json")

	target := cfg.Targets[0]

	err := s.forwardRequest(req, target)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify that the request was forwarded correctly
	if mockClient.Response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", mockClient.Response.StatusCode)
	}

	// Verify that the forwarded request contains the correct parameters
	if len(mockClient.Requests) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(mockClient.Requests))
	}

	forwardedReq := mockClient.Requests[0]
	params := forwardedReq.URL.Query()

	if params.Get("baromrelin") != "888.82" {
		t.Errorf("Expected baromrelin=888.82, got %s", params.Get("baromrelin"))
	}

	if params.Get("humidityin") != "52" {
		t.Errorf("Expected humidityin=52, got %s", params.Get("humidityin"))
	}

	if params.Get("tempinf") != "74.2" {
		t.Errorf("Expected tempinf=74.2, got %s", params.Get("tempinf"))
	}
}
