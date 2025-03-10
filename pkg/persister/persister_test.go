package persister

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// case where content is successfully persisted to a file
func TestPersist_ValidContent(t *testing.T) {
	// TempDir is a temporary directory that gets cleaned up automatically
	baseDir := t.TempDir()
	filePersister := &FilePersister{BaseDir: baseDir}

	// Content to persist
	content := "This is some test content."

	// Call the Persist method
	filePath, err := filePersister.Persist(content)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	// Ensure the file is created at the correct path
	dateStr := time.Now().Format("2006-01-02")
	expectedFolderPath := filepath.Join(baseDir, dateStr)
	if _, err := os.Stat(expectedFolderPath); os.IsNotExist(err) {
		t.Errorf("Expected folder '%s' to exist, but it doesn't", expectedFolderPath)
	}

	// Check if the file is created with the expected content
	if !strings.HasPrefix(filePath, expectedFolderPath) {
		t.Errorf("Expected file path to start with '%s', but got '%s'", expectedFolderPath, filePath)
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(fileContent) != content {
		t.Errorf("Expected file content '%s', but got '%s'", content, string(fileContent))
	}
}

// case where creating a folder fails
func TestPersist_ErrorCreatingFolder(t *testing.T) {
	invalidDir := "/invalid_permission_directory"
	filePersister := &FilePersister{BaseDir: invalidDir}

	// Content to persist
	content := "This content won't be saved."

	// Call the Persist method, expecting an error because the directory is invalid
	_, err := filePersister.Persist(content)
	if err == nil {
		t.Fatalf("Expected error, but got nil")
	}

	// Assert that the error contains the expected message (e.g., permission denied)
	if !strings.Contains(err.Error(), "mkdir /invalid_permission_directory: read-only file system") {
		t.Errorf("Expected error message to contain 'permission denied', but got: %v", err)
	}
}

// case where creating a file fails
func TestPersist_ErrorCreatingFile(t *testing.T) {
	baseDir := t.TempDir()
	// Change permissions of baseDir to be read-only
	err := os.Chmod(baseDir, 0444)
	if err != nil {
		t.Fatalf("Failed to change directory permissions: %v", err)
	}

	// Initialize FilePersister
	filePersister := &FilePersister{BaseDir: baseDir}

	// Content to persist
	content := "This content won't be saved."

	// Call the Persist method, expecting an error
	_, err = filePersister.Persist(content)
	if err == nil {
		t.Fatalf("Expected error, but got nil")
	}

	// Assert that the error contains the expected message (e.g., file creation failure)
	if !strings.Contains(err.Error(), "permission denied") {
		t.Errorf("Expected error message to contain 'permission denied', but got: %v", err)
	}

	// Cleanup: Restore directory permissions so the temp directory can be removed
	err = os.Chmod(baseDir, 0755)
	if err != nil {
		t.Fatalf("Failed to restore directory permissions: %v", err)
	}
}

// case where empty content is persisted to a file
func TestPersist_EmptyContent(t *testing.T) {
	baseDir := t.TempDir()
	filePersister := &FilePersister{BaseDir: baseDir}

	// Content is empty
	content := ""

	// Call the Persist method
	filePath, err := filePersister.Persist(content)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	// Ensure the file is created at the correct path
	dateStr := time.Now().Format("2006-01-02")
	expectedFolderPath := filepath.Join(baseDir, dateStr)
	if _, err := os.Stat(expectedFolderPath); os.IsNotExist(err) {
		t.Errorf("Expected folder '%s' to exist, but it doesn't", expectedFolderPath)
	}

	// Check if the file is created with empty content
	if !strings.HasPrefix(filePath, expectedFolderPath) {
		t.Errorf("Expected file path to start with '%s', but got '%s'", expectedFolderPath, filePath)
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(fileContent) != content {
		t.Errorf("Expected file content '%s', but got '%s'", content, string(fileContent))
	}
}
