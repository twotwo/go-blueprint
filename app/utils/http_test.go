package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const HeaderXTest = "X-Test"

// TestResponse is a sample struct for testing JSON response decoding.
type TestResponse struct {
	Message string `json:"message"`
}

func TestDoHttpPostSuccess(t *testing.T) {
	// Create a test HTTP server with a handler for valid JSON response.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ensure the method is POST.
		if r.Method != http.MethodPost {
			t.Fatalf("Expected POST method, got %s", r.Method)
		}
		// Verify the Content-Type header.
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			// Verify custom header.
			if r.Header.Get(HeaderXTest) != "test-value" {
				t.Fatalf("Expected header %s to be 'test-value', got '%s'", HeaderXTest, r.Header.Get(HeaderXTest))
			}
			t.Fatalf("Expected header X-Test to be 'test-value', got '%s'", r.Header.Get("X-Test"))
		}
		// Return a valid JSON response.
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(TestResponse{Message: "hello"})
	}))
	// Custom header to add.
	headerMap := map[string]string{
		HeaderXTest: "test-value",
	}
	// Call DoHttpPost.
	resp, err := DoHttpPost[TestResponse](ts.URL, headerMap, strings.NewReader(`{"dummy": "data"}`))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.Message != "hello" {
		t.Errorf("Expected Message 'hello', got '%s'", resp.Message)
	}
}

func TestDoHttpPostInvalidJSON(t *testing.T) {
	// Create a test server that returns non-JSON content.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not a json"))
	}))
	defer ts.Close()

	// Call DoHttpPost using TestResponse as the expected type.
	_, err := DoHttpPost[TestResponse](ts.URL, map[string]string{}, strings.NewReader(`{"dummy": "data"}`))
	if err == nil {
		t.Fatal("Expected error due to invalid JSON, but got none")
	}
}

func TestDoHttpPostNon200Status(t *testing.T) {
	// Create a test server that returns a 400 status.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		// Optionally, write a response body.
		w.Write([]byte(`{"error": "bad request"}`))
	}))
	defer ts.Close()

	// Call DoHttpPost.
	_, err := DoHttpPost[TestResponse](ts.URL, map[string]string{}, strings.NewReader(`{"dummy": "data"}`))
	if err == nil {
		t.Fatal("Expected error due to non-200 status code, but got none")
	}
}
