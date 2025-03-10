package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/anilvpatel21/snapurl/internal/ports"
	"github.com/anilvpatel21/snapurl/pkg/csv"
	"github.com/anilvpatel21/snapurl/pkg/downloader"
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

	/* ------ stage 1 done ------- */

	/* ------ stage 2 start ------- */
	// Channels for passing data between stages
	contentChan := make(chan string)

	// Initialize downloader dependencies
	downloader := &downloader.HttpDownloader{
		Timeout: 10,
	}

	// Start goroutine to download URLs
	var downloadWait sync.WaitGroup
	downloadWait.Add(1)
	go func() {
		defer func() {
			downloadWait.Done()
		}()
		downloadURLs(downloader, urlChan, contentChan, &downloadWait)
	}()

	go func() {
		downloadWait.Wait()
		close(contentChan)
	}()

	for content := range contentChan {
		fmt.Println(content)
	}

	/* ------ stage 2 done ------- */
}

// Read URLs from CSV file line by line
func startReading(r ports.Reader, readChan chan<- string) {
	if err := r.ReadLines(readChan); err != nil {
		log.Fatalf("Error reading: %v", err)
	}
}

func downloadURLs(downloader ports.Downloader, urlChan chan string, contentChan chan string, wg *sync.WaitGroup) {
	semaphore := make(chan struct{}, 50)
	for url := range urlChan {
		wg.Add(1)

		go func(url string) {
			defer func() {
				wg.Done()
				<-semaphore
			}()
			semaphore <- struct{}{}

			content, err := downloader.Download(url)
			if err != nil {
				log.Printf("Error downloading %s: %v", url, err)
				return
			}
			contentChan <- content
		}(url)
	}
}
