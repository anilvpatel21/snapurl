package downloader

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestDownload_Success(t *testing.T) {
	// Create a mock server that returns a successful response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Success"))
	}))
	defer server.Close()

	downloader := &HttpDownloader{
		Timeout: 5, // 5 seconds timeout
	}

	result, err := downloader.Download(server.URL)
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}

	expected := "Success"
	if result != expected {
		t.Errorf("expected %s, but got %s", expected, result)
	}
}

func TestDownload_Failure_NetworkError(t *testing.T) {
	// Test case to simulate a network failure
	downloader := &HttpDownloader{
		Timeout: 5,
	}

	// Invalid URL to simulate network error
	_, err := downloader.Download("http://nonexistentserver")
	if err == nil {
		t.Errorf("expected an error, but got none")
	}
}

func TestDownload_Failure_Timeout(t *testing.T) {
	// Create a mock server that simulates a long delay (timeout)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second) // Simulate long delay
		w.Write([]byte("Delayed"))
	}))
	defer server.Close()

	downloader := &HttpDownloader{
		Timeout: 1, // 1 second timeout (should trigger a timeout error)
	}

	_, err := downloader.Download(server.URL)
	if err == nil {
		t.Errorf("expected timeout error, but got none")
	} else {
		var urlError *url.Error
		if !errors.As(err, &urlError) {
			t.Errorf("expected timeout error, but got %v", err)
		}
	}
}

func TestDownload_Failure_InvalidResponse(t *testing.T) {
	// Create a mock server that simulates a server error (HTTP 500)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	defer server.Close()

	// Create the downloader instance
	downloader := &HttpDownloader{
		Timeout: 5,
	}

	_, err := downloader.Download(server.URL)
	if err == nil {
		t.Errorf("expected an error, but got none")
	} else {
		expectedErr := fmt.Sprintf("url %s responded with status code: %d", server.URL, http.StatusInternalServerError)
		if err.Error() != expectedErr {
			t.Errorf("expected error string, but got %v", err)
		}
	}
}
