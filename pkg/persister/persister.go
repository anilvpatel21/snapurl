package persister

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// FilePersister implements the Persister interface for saving files locally with a folder structure
type FilePersister struct {
	BaseDir string // Base directory to save files
}

func (f *FilePersister) Persist(content string) (string, error) {
	// Create a random filename
	rand.Seed(time.Now().UnixNano())
	filename := fmt.Sprintf("%d.txt", rand.Int())

	// Generate folder structure based on domain and current date
	dateStr := time.Now().Format("2006-01-02")
	folderPath := filepath.Join(f.BaseDir, dateStr)

	// Ensure the folder exists
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	// Create file in the structured folder path
	filePath := filepath.Join(folderPath, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
