package csv

import (
	"strings"
	"testing"
)

// Test case for reading an file not found (edge case).
func TestCSVReader_EmptyFile(t *testing.T) {
	// Create CSVReader with the empty file path
	csvReader := &CSVReader{FilePath: "./testfiles/empty_file.csv"}
	urlChan := make(chan string)

	// Run the reader in a goroutine
	go func() {
		err := csvReader.ReadLines(urlChan)
		if err != nil {
			t.Errorf("Error reading lines: %v", err)
		}
		close(urlChan)
	}()

	// Capture URLs from the channel (should be empty)
	var urls []string
	for url := range urlChan {
		urls = append(urls, url)
	}

	// Ensure no URLs are read
	if len(urls) > 0 {
		t.Errorf("Expected no URLs, but got %d", len(urls))
	}
}

func TestCSVReader_FileNotFound(t *testing.T) {
	// Create CSVReader with the no file path specified
	csvReader := &CSVReader{FilePath: ""}
	urlChan := make(chan string)

	// Run the reader in a goroutine
	go func() {
		err := csvReader.ReadLines(urlChan)
		if err != nil {
			expectedErr := "open : no such file or directory" // This depends on the Go implementation
			if !strings.Contains(err.Error(), expectedErr) {
				t.Errorf("Expected error to contain '%s', but got '%s'", expectedErr, err.Error())
			}
		}

		close(urlChan)
	}()

	// Capture URLs from the channel (should be empty)
	var urls []string
	for url := range urlChan {
		urls = append(urls, url)
	}

	// Ensure no URLs are read
	if len(urls) > 0 {
		t.Errorf("Expected no URLs, but got %d", len(urls))
	}
}

func TestCSVReader_ReadLines(t *testing.T) {
	// Initialize the CSVReader with the temporary file path
	csvReader := &CSVReader{FilePath: "./testfiles/test_file.csv"}
	urlChan := make(chan string)

	// Run the ReadLines method in a goroutine to avoid blocking
	go func() {
		err := csvReader.ReadLines(urlChan)
		if err != nil {
			t.Errorf("Error reading lines: %v", err)
		}
		close(urlChan)
	}()

	// Capture the URLs from the channel
	var urls []string
	for url := range urlChan {
		urls = append(urls, url)
	}

	// Validate the expected URLs
	expectedURLs := []string{
		"http://example1.com",
		"http://example2.com",
		"http://example3.com",
	}

	if len(urls) != len(expectedURLs) {
		t.Fatalf("Expected %d URLs, but got %d", len(expectedURLs), len(urls))
	}

	for i, url := range urls {
		if url != expectedURLs[i] {
			t.Errorf("Expected URL %s, but got %s", expectedURLs[i], url)
		}
	}
}
