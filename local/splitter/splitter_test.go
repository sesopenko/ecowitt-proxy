package splitter

import (
	"ecowitt-proxy/local/config"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"
)

// mockHTTPClient is a mock implementation of an HTTP client
type mockHTTPClient struct {
	Response *http.Response
	Err      error
	Requests []*http.Request // Capture requests for verification
	mu       sync.Mutex      // Ensure thread-safe access to Requests slice
}

// Do is the mock implementation of the Do method
func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Requests = append(m.Requests, req) // Capture the request
	return m.Response, m.Err
}

func TestSplitter_HandleRequest(t *testing.T) {
	cfg := config.Config{
		Targets: []config.Target{
			{
				Name:     "target1",
				HostAddr: "https://example.com:8220",
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

	// Create a request with real-world POST parameters
	formData := url.Values{
		"baromrelin": {"888.82"},
		"humidityin": {"52"},
		"tempinf":    {"74.2"},
	}
	req := httptest.NewRequest("POST", "http://localhost", strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.HandleRequest)
	handler.ServeHTTP(rr, req)

	// Verify that the request was forwarded correctly
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rr.Code)
	}

	// Wait for all requests to complete
	time.Sleep(1 * time.Second)

	// Verify that the forwarded request contains the correct body
	if len(mockClient.Requests) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(mockClient.Requests))
	}

	mockClient.mu.Lock()
	defer mockClient.mu.Unlock()

	forwardedReq := mockClient.Requests[0]
	body, err := io.ReadAll(forwardedReq.Body)
	if err != nil {
		t.Fatalf("Error reading forwarded request body: %v", err)
	}

	forwardedReq.Body.Close()
	forwardedParams, err := url.ParseQuery(string(body))
	if err != nil {
		t.Fatalf("Error parsing forwarded request body: %v", err)
	}

	if forwardedParams.Get("baromrelin") != "888.82" {
		t.Errorf("Expected baromrelin=888.82, got %s", forwardedParams.Get("baromrelin"))
	}

	if forwardedParams.Get("humidityin") != "52" {
		t.Errorf("Expected humidityin=52, got %s", forwardedParams.Get("humidityin"))
	}

	if forwardedParams.Get("tempinf") != "74.2" {
		t.Errorf("Expected tempinf=74.2, got %s", forwardedParams.Get("tempinf"))
	}
}

func TestSplitter_forwardRequest(t *testing.T) {
	cfg := config.Config{
		Targets: []config.Target{
			{
				Name:     "target1",
				HostAddr: "http://example.com:8220",
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

	// Create a request with real-world POST parameters
	formData := url.Values{
		"baromrelin": {"888.82"},
		"humidityin": {"52"},
		"tempinf":    {"74.2"},
	}
	req := httptest.NewRequest("POST", "http://localhost", strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	target := cfg.Targets[0]

	o, err := buildOriginal(req)
	if err != nil {
		t.Errorf("Failed to build original from req: %s", err)
	}

	err = s.forwardRequest(o, target)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify that the request was forwarded correctly
	if mockClient.Response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", mockClient.Response.StatusCode)
	}

	// Verify that the forwarded request contains the correct body
	if len(mockClient.Requests) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(mockClient.Requests))
	}

	mockClient.mu.Lock()
	defer mockClient.mu.Unlock()

	forwardedReq := mockClient.Requests[0]
	body, err := io.ReadAll(forwardedReq.Body)
	if err != nil {
		t.Fatalf("Error reading forwarded request body: %v", err)
	}

	forwardedReq.Body.Close()
	forwardedParams, err := url.ParseQuery(string(body))
	if err != nil {
		t.Fatalf("Error parsing forwarded request body: %v", err)
	}

	if forwardedParams.Get("baromrelin") != "888.82" {
		t.Errorf("Expected baromrelin=888.82, got %s", forwardedParams.Get("baromrelin"))
	}

	if forwardedParams.Get("humidityin") != "52" {
		t.Errorf("Expected humidityin=52, got %s", forwardedParams.Get("humidityin"))
	}

	if forwardedParams.Get("tempinf") != "74.2" {
		t.Errorf("Expected tempinf=74.2, got %s", forwardedParams.Get("tempinf"))
	}
}
