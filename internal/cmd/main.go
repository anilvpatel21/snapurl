package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/anilvpatel21/snapurl/internal/ports"
	"github.com/anilvpatel21/snapurl/pkg/csv"
)

func main() {
	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// Read input file path from CLI
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <file_path>")
	}
	filePath := os.Args[1]

	/* ------ stage 1 start------- */
	// Create the URL channel
	urlChan := make(chan string)

	// Initialize reader dependencies
	csvReader := &csv.CSVReader{
		FilePath: filePath,
	}

	// Start goroutine to reading file
	var readingWait sync.WaitGroup
	readingWait.Add(1)
	go func() {
		defer func() {
			readingWait.Done()
		}()
		startReading(csvReader, urlChan)
	}()

	go func() {
		readingWait.Wait()
		close(urlChan)
	}()

	for url := range urlChan {
		fmt.Println(url)
	}

	/* ------ stage 1 done ------- */
}

// Read URLs from CSV file line by line
func startReading(r ports.Reader, readChan chan<- string) {
	if err := r.ReadLines(readChan); err != nil {
		log.Fatalf("Error reading: %v", err)
	}
}
